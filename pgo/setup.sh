#!/bin/bash

# ensure the pgo operator is deleted first to start its deployment from scratch
kubectl delete -k ./high-availability --ignore-not-found=true 2>/dev/null

# ensure the pgo operator crd and other stuff is deleted first to start its deployment from scratch
kubectl delete -k ./install --ignore-not-found=true 2>/dev/null

pg_namespace="hoh-postgres"

postgres_user_secrets=("hoh-pguser-hoh-process-user" "hoh-pguser-postgres" "hoh-pguser-transport-bridge-user")

# ensure all the user secrets are deleted
for secret in ${postgres_user_secrets[@]}
do
    matched=$(kubectl get secret $secret -n $pg_namespace --ignore-not-found=true)
    while [ ! -z "$matched" ]; do
	echo "Waiting for secret $secret to be deleted from namespace $pg_namespace"
	matched=$(kubectl get secret $secret -n $pg_namespace --ignore-not-found=true)
	sleep 10
    done
done

# install the pgo operator to postgres-operator
kubectl apply -k ./install

# set up hoh database hoh
kubectl apply -k ./high-availability

# ensure all the user secrets are created
for secret in ${postgres_user_secrets[@]}
do
    matched=$(kubectl get secret $secret -n $pg_namespace --ignore-not-found=true)
    while [ -z "$matched" ]; do
	echo "Waiting for secret $secret to be created in namespace $pg_namespace"
	matched=$(kubectl get secret $secret -n $pg_namespace --ignore-not-found=true)
	sleep 10
    done
done

kubectl delete -f ./postgres-job.yaml --ignore-not-found=true

envsubst < ./postgres-job.yaml | kubectl apply -f -

kubectl wait --for=condition=complete job/postgres-init -n $pg_namespace --timeout=300s

kubectl logs $(kubectl get pods --field-selector status.phase=Succeeded  --selector=job-name=postgres-init -n $pg_namespace  --output=jsonpath='{.items[*].metadata.name}') -n $pg_namespace
