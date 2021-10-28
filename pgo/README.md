# Why?
If we can run the hoh postgres DB inside cluster as an operator to test/try our the Hub-of-Hubs setup.

# What does it do
1. installs [PostgreSQL Operator](https://access.crunchydata.com/documentation/postgres-operator/v5/)
2. set up hoh database and expose the db as a `LoadBalancer` service
3. the database with 2 users created `./high-availability/ha-postgres.yaml:10,13`, these users will be used by the anisble-playbook as well.

# How to do
1. make sure your `KUBECONFIG` is pointing the HoH cluster
2. run `./setup.sh`

If the command above does not produce any errors, you should able to connect to the Hoh DB sits inside your hoh cluster.

You can use the following commands to get HoH DB credentials and to connect to the database.
```
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

# connect to db
PGPASSWORD=$PGPASSWORD PGSSLMODE=require psql -h $PGHOST -U $DB_NAME -d $DB_NAME

# load test file from tpl to db hoh
PGPASSWORD=$PGPASSWORD PGSSLMODE=require psql -h $PGHOST -U $DB_NAME -d $DB_NAME  -f ../test/insertpolicies.psql
```

# What's next
Use ansible-playbook to initialize the database, e.g, creates tables, set up permissions, etc.

0. install the `ansible-playbook` dependency `psycopg2`

`pip3 install psycopg2`

To use the ansible-playbook, you need to
1. create an inventory file, `production`, such as:

```
[local]
localhost   ansible_connection=local
```

2. create varialbes for the `local` host defined in the above inventory file
The local variables file should be located at `../group_vars/local/vars.yaml`, these varaiables are expected by be set via environment variables.

```
db_login_host: "'{{ lookup('env', 'DB_LOGIN_HOST') }}'"
db_login_user: "'{{ lookup('env', 'DB_LOGIN_USER') }}'"
db_login_password: "'{{ lookup('env', 'DB_LOGIN_PASSWORD') }}'"
db_ssl_mode: require
```

3. run the anislbe-playbook with the following command

```
sudo ANSIBLE_PYTHON_INTERPRETER=/usr/local/bin/python3 ansible-playbook ../pgo.yaml -i ../production --ask-vault-pass -l local
```