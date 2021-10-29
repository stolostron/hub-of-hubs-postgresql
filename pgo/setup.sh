#!/bin/bash

# install the pgo operator to postgres-operator
kubectl apply -k ./install

# set up hoh database hoh
kubectl apply -k ./high-availability

pg_namespace="hoh-postgres"
db_name="hoh"

postgres_secert_name="$db_name-pguser-postgres"

while [ -z "$matched" ]; do
    echo "Waiting for ($pg_namespace/$postgres_secert_name) to be created"
    matched=$(kubectl get secret $postgres_secert_name -n $pg_namespace --ignore-not-found=true)
    sleep 10
done

kubectl apply -f ./postgres-job.yaml

kubectl wait --for=condition=complete job/postgres-init -n $pg_namespace --timeout=120s

kubectl logs job/postgres-init -n $pg_namespace
