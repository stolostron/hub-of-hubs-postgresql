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
ROWS_NUMBER=100000 BATCH_SIZE=10000 INSERT_MULTIPLE_VALUES= ./bin/db-performance-tests
```

Insert by the COPY protocol:

```
ROWS_NUMBER=100000 BATCH_SIZE=10000 INSERT_COPY= ./bin/db-performance-tests
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

1.  Connect to the pod:

    ```
    kubectl exec -it -n myproject --kubeconfig $TOP_HUB_CONFIG $(kubectl get pods -l name=db-performance-tests -n myproject --kubeconfig $TOP_HUB_CONFIG  -o jsonpath='{.items..metadata.name}') -- bash
    ```

1.  Run the tests (the executable is `db-performance-tests`, available via the `PATH` environment variable)

## Setting a client on RHEL

    ```
    sudo dnf module install postgresql:13
    ```

    ```
    sudo dnf install git wget emacs make -y
    git clone https://github.com/open-cluster-management/hub-of-hubs-postgresql
    cd hub-of-hubs-postgresql/
    git checkout perf
    cd db-performance-tests/
    wget https://golang.org/dl/go1.16.6.linux-amd64.tar.gz
    sha256sum go1.16.6.linux-amd64.tar.gz
    ```

    Verify the checksum at [Go downloads](https://golang.org/dl/).

    ```
    sudo tar -zxvf go1.16.6.linux-amd64.tar.gz -C /usr/local/
    rm -rf go1.16.6.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    source ~/.bashrc
    make build
    ```
