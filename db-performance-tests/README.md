# Performance tests for Hub-of-hubs database

## Build

```
go build
```

## Run

Set the following environment variable:

* DATABASE_URL

Set the `DATABASE_URL` according to the PostgreSQL URL format, plus the maximal pool connection size : `postgres://YourUserName:YourURLEscapedPassword@YourHostname:5432/YourDatabaseName?sslmode=verify-full&pool_max_conns=50`.

:exclamation: Remember to URL-escape the password, you can do it in bash:

```
python -c "import sys, urllib as ul; print ul.quote_plus(sys.argv[1])" 'YourPassword'
```

Run the tests.

Insert by inserting multiple values:

```
ROWS_NUMBER=1000 INSERT_MULTIPLE_VALUES= ./bin/db-performance-tests
```

Insert by the COPY protocol:

```
ROWS_NUMBER=1000 INSERT_COPY= ./bin/db-performance-tests
```

# Build image

Define the `REGISTRY` environment variable.

```
make build-images
```

## Deploy to a cluster

1.  Create a secret with your database url:

    ```
    kubectl create secret generic hub-of-hubs-database-secret --kubeconfig $TOP_HUB_CONFIG --from-literal=url=$DATABASE_URL -n myproject
    ```

1.  Deploy the container:

    ```
    IMAGE_TAG=latest COMPONENT=$(basename $(pwd)) envsubst < deploy/operator.yaml.template | kubectl apply -n myproject --kubeconfig $TOP_HUB_CONFIG -f -
    ```

1. Connect to the pod:

```
kubectl exec -it -n myproject --kubeconfig $TOP_HUB_CONFIG $(kubectl get pods -l name=db-performance-tests -n myproject --kubeconfig $TOP_HUB_CONFIG  -o jsonpath='{.items..metadata.name}') -- bash
```
