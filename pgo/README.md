# Why?
If we can the hoh postgres DB runs inside cluster as an operator, it will make it easier for people to test/try our Hoh.

# What does it do
1. install [PostgreSQL Operator](https://access.crunchydata.com/documentation/postgres-operator/v5/)
2. set up hoh database and expose the db as a `LoadBalancer` service
3. use the exposed db service to load `tpl`

# How to do
1. make sure your `KUBECONFIG` is pointing to hoh cluster
2. run `./setup.sh`

If above run without issue, you should able to connect to the hoh DB which sits inside your hoh cluster. Followings are command for you
get hoh DB creds and how to connect to it.
```
PG_NAMESPACE="hoh-postgres"
DB_NAME="hoh"
HOH_PGBOUNCER="pgbouncer"
PG_CLUSTER_USER_SECRET_NAME="$DB_NAME-pguser-$DB_NAME"
PGPASSWORD="$(kubectl get secrets -n "${PG_NAMESPACE}" "${PG_CLUSTER_USER_SECRET_NAME}" -o go-template='{{.data.password | base64decode}}')"
PGUSER="$(kubectl get secrets -n "${PG_NAMESPACE}" "${PG_CLUSTER_USER_SECRET_NAME}" -o go-template='{{.data.user | base64decode}}')"
PG_URI="$(kubectl get secrets -n "${PG_NAMESPACE}" "${PG_CLUSTER_USER_SECRET_NAME}" -o go-template='{{.data.uri | base64decode}}')"

PGHOST="$(kubectl get svc -n "${PG_NAMESPACE}" "${HOH_PGBOUNCER}" -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')"
echo "postgres secret name: " $PG_CLUSTER_USER_SECRET_NAME
echo "postgres user: " $PGUSER
echo "postgres password: " $PGPASSWORD
echo "postgres host: "$PGHOST

# connect to db
PGPASSWORD=$PGPASSWORD PGSSLMODE=require psql -h $PGHOST -U $DB_NAME -d $DB_NAME

# load test file from tpl to db hoh
PGPASSWORD=$PGPASSWORD PGSSLMODE=require psql -h $PGHOST -U $DB_NAME -d $DB_NAME  -f ../test/insertpolicies.psql
```

# What's next
1. how to unify `ansible` and `pgo` flow

I guess, we can leverage `tpl` files, then `ansible` and `pgo` will just load from tpl files
```
    - hosts: db_hosts
      become_user: postgres
      tasks:
      - copy: src=file.dmp dest=/home/postgres/file.dmp
      - shell: psql dbname < file.dmp

```