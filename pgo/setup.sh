#!/bin/bash

# install the pgo operator to postgres-operator
kubectl apply -f ./install

# set up hoh database hoh
kubectl apply -f ./high-availability

PG_NAMESPACE="hoh-postgres"
DB_NAME="hoh"
HOH_PGBOUNCER="pgbouncer"
PGBOUNCER_SERVICE="hoh-pgbouncer"
PG_CLUSTER_USER_SECRET_NAME="$DB_NAME-pguser-$DB_NAME"
PGPASSWORD="$(kubectl get secrets -n "${PG_NAMESPACE}" "${PG_CLUSTER_USER_SECRET_NAME}" -o go-template='{{.data.password | base64decode}}')"
PGUSER="$(kubectl get secrets -n "${PG_NAMESPACE}" "${PG_CLUSTER_USER_SECRET_NAME}" -o go-template='{{.data.user | base64decode}}')"
PG_URI="$(kubectl get secrets -n "${PG_NAMESPACE}" "${PG_CLUSTER_USER_SECRET_NAME}" -o go-template='{{.data.uri | base64decode}}')"

kubectl expose service $PGBOUNCER_SERVICE -n $PG_NAMESPACE --port=5432 --target-port=5432 --name=$HOH_PGBOUNCER --type=LoadBalancer

PGHOST="$(kubectl get svc -n "${PG_NAMESPACE}" "${HOH_PGBOUNCER}" -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')"
echo "postgres secret name: " $PG_CLUSTER_USER_SECRET_NAME
echo "postgres user: " $PGUSER
echo "postgres password: " $PGPASSWORD
echo "postgres host: "$PGHOST

# load file from tpl to db hoh
echo "PGPASSWORD='$PGPASSWORD' PGSSLMODE=require psql -h $PGHOST -U $DB_NAME -d $DB_NAME"
PGPASSWORD='$PGPASSWORD' PGSSLMODE=require psql -h $PGHOST -U $DB_NAME -d $DB_NAME