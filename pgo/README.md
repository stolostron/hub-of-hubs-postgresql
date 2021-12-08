# Why?
If we can run the hoh postgres DB inside cluster as an operator to test/try our the Hub-of-Hubs setup.

# What does it do
1. installs [PostgreSQL Operator](https://access.crunchydata.com/documentation/postgres-operator/v5/)
2. set up hoh database
3. the database with 2 users created `./high-availability/ha-postgres.yaml:11,14`, these users will be used by the anisble-playbook as well.
4. all hoh related schema and tables will be created via `job/postgres-init`

# How to do
1. make sure your `KUBECONFIG` is pointing the HoH cluster. Ask your cluster's admin to give you appropriate permissions.
2. set the `USER_NAME` environment variable to hold the username part of your docker registry:
    ```
    $ export USER_NAME=...
    ```
3. set the `IMAGE_TAG` environment variable to hold the tag of your image:
    ```
    $ export IMAGE_TAG=latest
    ```
4. set the `REGISTRY` environment variable to hold your docker registry:
    ```
    $ export REGISTRY=quay.io/$USER_NAME/postgre-ansible:$IMAGE_TAG
    ```
5. run `docker build -f Dockerfile -t $REGISTRY .` and then `docker push $REGISTRY` from the project root folder
6. run `./setup.sh`

If the command above does not produce any errors, you should able to connect to the Hoh DB sits inside your hoh cluster.

You can use the following commands to get HoH DB credentials and to connect to the database.
```
PG_NAMESPACE="hoh-postgres"
DB_NAME="hoh"
HOH_PGBOUNCER="pgbouncer"
PGBOUNCER_SERVICE="hoh-pgbouncer"
PG_CLUSTER_USER_SECRET_NAME="$DB_NAME-pguser-hoh-process-user"
PGPASSWORD="$(kubectl get secrets -n "${PG_NAMESPACE}" "${PG_CLUSTER_USER_SECRET_NAME}" -o go-template='{{.data.password | base64decode}}')"
PGUSER="$(kubectl get secrets -n "${PG_NAMESPACE}" "${PG_CLUSTER_USER_SECRET_NAME}" -o go-template='{{.data.user | base64decode}}')"
PG_URI="$(kubectl get secrets -n "${PG_NAMESPACE}" "${PG_CLUSTER_USER_SECRET_NAME}" -o go-template='{{.data.uri | base64decode}}')"

kubectl expose service $PGBOUNCER_SERVICE -n $PG_NAMESPACE --port=5432 --target-port=5432 --name=$HOH_PGBOUNCER --type=LoadBalancer

PGHOST="$(kubectl get svc -n "${PG_NAMESPACE}" "${HOH_PGBOUNCER}" -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')"
echo "postgres secret name: " $PG_CLUSTER_USER_SECRET_NAME
echo "postgres user: " $PGUSER
echo "postgres password: " $PGPASSWORD
echo "postgres host: "$PGHOST

# connect to db
PGPASSWORD=$PGPASSWORD PGSSLMODE=require psql -h $PGHOST -U $DB_NAME -d $DB_NAME

```

# What's next
Run the Hoh deployments and play around.
