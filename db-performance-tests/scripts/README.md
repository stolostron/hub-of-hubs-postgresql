# Running a command on multiple VMs

## Start multiple VMs by a name prefix

```
VM_NAME_PREFIX=veisenbe-postgresql-perf-client  SECURITY_GROUP=veisenbe-postgresql-sg ./start_instances.sh
```

## Run command on multiple VMs by a name prefix

```
VM_NAME_PREFIX=veisenbe-postgresql-perf-client COMMAND="export DATABASE_URL=\"$DATABASE_URL\"; export LEAF_HUBS_NUMBER=1; export START_LEAF_HUB_INDEX=0; export UPDATE=; hub-of-hubs-postgresql/db-performance-tests/bin/db-performance-tests" ./run_command.sh
```

## Stop multiple VMs by a name prefix

```
VM_NAME_PREFIX=veisenbe-postgresql-perf-client  SECURITY_GROUP=veisenbe-postgresql-sg ./stop_instances.sh
```

## Linting

**Prerequisite**: install the `shellcheck` tool (a Linter for shell):

```
brew install shellcheck
```

Run
```
shellcheck *.sh
```
