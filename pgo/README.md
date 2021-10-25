# Why?
If we can the hoh postgres DB runs inside cluster as an operator, it will make it easier for people to test/try our Hoh.

# What does it do
1. install [PostgreSQL Operator](https://access.crunchydata.com/documentation/postgres-operator/v5/)
2. set up hoh database and expose the db as a `LoadBalancer` service
3. the database with 2 users created `./high-availability/ha-postgres.yaml:10,13`, these user will be used at the anisble-playbook as well.

# How to do
1. make sure your `KUBECONFIG` is pointing to hoh cluster
2. run `./setup.sh`

If above run without issue, you should able to connect to the hoh DB which sits inside your hoh cluster. Followings are command for you
get hoh DB creds and how to connect to it.
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

To use the ansible-playbook, you need to
1. create an inventory file, `production`, such as:

```
[local]
localhost   ansible_connection=local
```

2. create varialbes for the `local` host defined in the above inventory file
The local variables file should be located at `../group_vars/local/vars.yaml`

```
db_login_host: a9637aa2200a642e7889d424cec29df3-1992409284.us-east-1.elb.amazonaws.com
db_login_user: hoh
db_login_password: "Oy:E@^yiKpH.6?CsK((?I<eo"
db_ssl_mode: require
```

3. run the anislbe-playbook with the following command

```
sudo ANSIBLE_PYTHON_INTERPRETER=/usr/local/bin/python3 ansible-playbook ../install.yaml -i production --ask-vault-pass -l local
```