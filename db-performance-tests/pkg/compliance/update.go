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

	err = updateCompliance(ctx, dbConnectionPool, tuplesToBecomeCompliant, leafHubName, true)
	if err != nil {
		return fmt.Errorf("failed to make previous non-compliant tuples compliant: %w", err)
	}

	err = updateComplianceOrInsert(ctx, dbConnectionPool, newTuplesToBecomeNonCompliant, leafHubName, false)
	if err != nil {
		return fmt.Errorf("failed to make new tuples non-compliant: %w", err)
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

func updateCompliance(ctx context.Context, dbConnectionPool *pgxpool.Pool, policyClusterTuples set.Set,
	leafHubName string, compliant bool) error {
	if policyClusterTuples.Cardinality() < 1 {
		return nil
	}

	compliance := compliantString
	if !compliant {
		compliance = nonCompliantString
	}

	_, err := dbConnectionPool.Exec(ctx, fmt.Sprintf(`UPDATE status.compliance SET compliance = '%s'
		WHERE leaf_hub_name = '%s' AND (%s)`, compliance, leafHubName, sqlConditionFromTuples(policyClusterTuples)))
	if err != nil {
		return fmt.Errorf("error in updating compliance: %w", err)
	}

	return nil
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
	compliant bool) [][]interface{} {
	rows := make([][]interface{}, 0, policyClusterTuples.Cardinality())

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

		rows = append(rows, []interface{}{
			pct.PolicyID.String(), pct.ClusterName, leafHubName, errorValue,
			compliance, action,
		})
		return false
	})

	return rows
}

func updateComplianceOrInsert(ctx context.Context, dbConnectionPool *pgxpool.Pool, policyClusterTuples set.Set,
	leafHubName string, compliant bool) error {
	existingRows, err := dbConnectionPool.Query(ctx, fmt.Sprintf(`SELECT policy_id, cluster_name FROM status.compliance
		WHERE leaf_hub_name = '%s' AND (%s)`, leafHubName, sqlConditionFromTuples(policyClusterTuples)))
	if err != nil {
		return fmt.Errorf("error in getting non_compliant clusters: %w", err)
	}

	existingPolicyClusterTuples, err := getPolicyClusterTupleSet(existingRows)
	if err != nil {
		return fmt.Errorf("error in scanning existing rows: %w", err)
	}

	tuplesToInsert := policyClusterTuples.Difference(existingPolicyClusterTuples)
	tuplesToUpdate := existingPolicyClusterTuples.Intersect(policyClusterTuples)

	err = updateCompliance(ctx, dbConnectionPool, tuplesToUpdate, leafHubName, compliant)
	if err != nil {
		return fmt.Errorf("failed to update compliance of existing tuples: %w", err)
	}

	err = insert(ctx, dbConnectionPool, tuplesToInsert, leafHubName, compliant)
	if err != nil {
		return fmt.Errorf("failed to insert new tuples: %w", err)
	}

	return nil
}

func insert(ctx context.Context, dbConnectionPool *pgxpool.Pool, policyClusterTuples set.Set,
	leafHubName string, compliant bool) error {
	if policyClusterTuples.Cardinality() < 1 {
		return nil
	}

	rows := generateRowsFromTuples(policyClusterTuples, leafHubName, compliant)

	err := insertRowsByInsertWithMultipleValues(ctx, dbConnectionPool, rows)
	if err != nil {
		return fmt.Errorf("error in inserting non_compliant clusters: %w", err)
	}

	return nil
}
