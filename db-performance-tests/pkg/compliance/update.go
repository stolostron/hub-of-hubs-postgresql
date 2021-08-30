package compliance

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	set "github.com/deckarep/golang-set"
	"github.com/gofrs/uuid"
	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type policyClusterTuple struct {
	PolicyID    uuid.UUID
	ClusterName string
}

func RunUpdate(ctx context.Context, dbConnectionPool *pgxpool.Pool, leafsNumber, startLeafHubIndex int) error {
	entry := time.Now()

	rand.Seed(entry.Unix())

	defer func() {
		now := time.Now()
		elapsed := now.Sub(entry)
		log.Printf("compliance RunUpdate of %d leaf hubs: elapsed %v\n", leafsNumber, elapsed)
	}()

	var wg sync.WaitGroup

	c := make(chan int, leafsNumber)

	for i := 0; i < leafsNumber; i++ {
		wg.Add(1)

		go updateRows(ctx, dbConnectionPool, c, &wg)
	}

	for i := 0; i < leafsNumber; i++ {
		c <- startLeafHubIndex + i
	}
	close(c)

	wg.Wait()

	return nil
}

func updateRows(ctx context.Context, dbConnectionPool *pgxpool.Pool, c chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for leafHubIndex := range c {
		if err := updateRowsForLeafHub(ctx, dbConnectionPool, leafHubIndex); err != nil {
			log.Printf("failed to update rows: %v\n", err)
			break
		}
	}
}

func updateRowsForLeafHub(ctx context.Context, dbConnectionPool *pgxpool.Pool, leafHubIndex int) error {
	leafHubName := fmt.Sprintf("hub%d", leafHubIndex)

	nonCompliantRows, err := dbConnectionPool.Query(ctx, fmt.Sprintf(`SELECT policy_id, cluster_name FROM status.compliance
		WHERE leaf_hub_name = '%s' AND compliance = 'non_compliant'`, leafHubName))
	if err != nil {
		return fmt.Errorf("error in getting non_compliant clusters: %w", err)
	}

	nonCompliantSet, err := getPolicyClusterTupleSet(nonCompliantRows)
	if err != nil {
		return fmt.Errorf("error in scanning non_compliant clusters: %w", err)
	}

	newNonCompliantSet := generatePolicyClusterSet(leafHubIndex,
		policiesNumber*clustersPerLeafHub/compliantToNonCompliantRatio)
	tuplesToBecomeCompliant := nonCompliantSet.Difference(newNonCompliantSet)
	newTuplesToBecomeNonCompliant := newNonCompliantSet.Difference(nonCompliantSet)

	batch := &pgx.Batch{}

	updateCompliance(tuplesToBecomeCompliant, leafHubName, true, batch)
	upsert(newTuplesToBecomeNonCompliant, leafHubName, false, batch)

	batchResult := dbConnectionPool.SendBatch(context.Background(), batch)
	defer batchResult.Close()

	for i := 0; i < batch.Len(); i++ {
		_, err = batchResult.Exec()
		if err != nil {
			return fmt.Errorf("failed to execute batch item %d: %w", i, err)
		}
	}

	return nil
}

func getPolicyClusterTupleSet(rows pgx.Rows) (set.Set, error) {
	s := set.NewSet()

	for rows.Next() {
		var (
			policyID    uuid.UUID
			clusterName string
		)

		err := rows.Scan(&policyID, &clusterName)
		if err != nil {
			return nil, fmt.Errorf("error in scanning policy-cluster rows: %w", err)
		}

		s.Add(policyClusterTuple{PolicyID: policyID, ClusterName: clusterName})
	}

	return s, nil
}

func updateCompliance(policyClusterTuples set.Set, leafHubName string, compliant bool, batch *pgx.Batch) {
	if policyClusterTuples.Cardinality() < 1 {
		return
	}

	compliance := compliantString
	if !compliant {
		compliance = nonCompliantString
	}

	batch.Queue(fmt.Sprintf(`UPDATE status.compliance SET compliance = '%s'
		WHERE leaf_hub_name = '%s' AND (%s)`, compliance, leafHubName, sqlConditionFromTuples(policyClusterTuples)))
}

func upsert(policyClusterTuples set.Set, leafHubName string, compliant bool, batch *pgx.Batch) {
	if policyClusterTuples.Cardinality() < 1 {
		return
	}

	rows := generateRowsFromTuples(policyClusterTuples, leafHubName, compliant)
	sb := generateInsertByMultipleValues(len(rows) / columnSize)

	compliance := compliantString
	if !compliant {
		compliance = nonCompliantString
	}

	sb.WriteString(" ON CONFLICT(policy_id, cluster_name, leaf_hub_name) DO UPDATE SET compliance = '")
	sb.WriteString(compliance)
	sb.WriteString("'")

	batch.Queue(sb.String(), rows...)
}

func sqlConditionFromTuples(policyClusterTuples set.Set) string {
	if policyClusterTuples.Cardinality() < 1 {
		return "FALSE"
	}

	var sb strings.Builder

	policyClusterTuples.Each(func(tuple interface{}) bool {
		pct, correctType := tuple.(policyClusterTuple)
		if !correctType {
			panic("policyClusterTuples contains a member of a wrong type")
		}
		sb.WriteString("policy_id = '")
		sb.WriteString(pct.PolicyID.String())
		sb.WriteString("' AND cluster_name  = '")
		sb.WriteString(pct.ClusterName)
		sb.WriteString("' OR ")
		return false
	})

	sb.WriteString("FALSE") // for the last OR

	return sb.String()
}

/* #nosec G404: Use of weak random number generator (math/rand instead of crypto/rand) */
func generatePolicyClusterSet(leafHubIndex, size int) set.Set {
	s := set.NewSet()

	for i := 0; i < size; i++ {
		policyID := policyUUIDs[rand.Intn(policiesNumber)]
		clusterIndex := leafHubIndex*clustersPerLeafHub + rand.Intn(clustersPerLeafHub)
		clusterName := fmt.Sprintf("cluster%d", clusterIndex)

		s.Add(policyClusterTuple{PolicyID: policyID, ClusterName: clusterName})
	}

	return s
}

func generateRowsFromTuples(policyClusterTuples set.Set, leafHubName string,
	compliant bool) []interface{} {
	rows := make([]interface{}, 0, policyClusterTuples.Cardinality()*columnSize)

	compliance := compliantString
	if !compliant {
		compliance = nonCompliantString
	}

	policyClusterTuples.Each(func(tuple interface{}) bool {
		pct, correctType := tuple.(policyClusterTuple)
		if !correctType {
			panic("policyClusterTuples contains a member of a wrong type")
		}

		errorValue, _, action := generateDerivedColumns()

		rows = append(rows, pct.PolicyID.String(), pct.ClusterName, leafHubName, errorValue,
			compliance, action)
		return false
	})

	return rows
}
