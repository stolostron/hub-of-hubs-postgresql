# Running a command on multiple VMs

In the following commands, `VM_NAME` specifies the instance names pattern. Examples:

* `VM_NAME=veisenbe-postgresql-perf-client1`
* `VM_NAME='veisenbe-postgresql-perf-client1,VM_NAME_PREFIX=veisenbe-postgresql-perf-clien2'`
* `VM_NAME='veisenbe-postgresql-perf-client*'`

## Start multiple VMs

```
VM_NAME='veisenbe-postgresql-perf-client*'  SECURITY_GROUP=veisenbe-postgresql-sg ./start_instances.sh
```

## Run command on multiple VMs

```
VM_NAME='veisenbe-postgresql-perf-client*' COMMAND="export UPDATE=; hub-of-hubs-postgresql/db-performance-tests/bin/db-performance-tests" ./run_command.sh
```

### git pull and build

```
VM_NAME='veisenbe-postgresql-perf-client*' COMMAND="cd hub-of-hubs-postgresql/db-performance-tests; git checkout perf; git pull; make build" ./run_command.sh
```

## Stop multiple VMs

```
VM_NAME='veisenbe-postgresql-perf-client*'  SECURITY_GROUP=veisenbe-postgresql-sg ./stop_instances.sh
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
