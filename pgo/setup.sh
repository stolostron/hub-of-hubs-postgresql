#!/bin/bash

# ensure the pgo operator is deleted first to start its deployment from scratch
kubectl delete -k ./high-availability --ignore-not-found=true 2>/dev/null

# ensure the pgo operator crd and other stuff is deleted first to start its deployment from scratch
kubectl delete -k ./install --ignore-not-found=true 2>/dev/null

# install the pgo operator to postgres-operator
kubectl apply -k ./install

# set up hoh database hoh
kubectl apply -k ./high-availability

pg_namespace="hoh-postgres"

postgres_secret_name="hoh-pguser-postgres"

while [ -z "$matched" ]; do
    echo "Waiting for ($pg_namespace/$postgres_secret_name) to be created"
    matched=$(kubectl get secret $postgres_secret_name -n $pg_namespace --ignore-not-found=true)
    sleep 10
done

kubectl delete -f ./postgres-job.yaml --ignore-not-found=true

envsubst < ./postgres-job.yaml | kubectl apply -f -

kubectl wait --for=condition=complete job/postgres-init -n $pg_namespace --timeout=300s

kubectl logs $(kubectl get pods --field-selector status.phase=Succeeded  --selector=job-name=postgres-init -n $pg_namespace  --output=jsonpath='{.items[*].metadata.name}') -n $pg_namespace
