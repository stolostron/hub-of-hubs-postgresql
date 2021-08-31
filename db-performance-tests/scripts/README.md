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

Set `DATABASE_URL` before running. 

1.  Insert 100 million rows:

    ```
    VM_NAME='veisenbe-postgresql-perf-client*' COMMAND="export BATCH_SIZE=10000; export INSERT_COPY=; hub-of-hubs-postgresql/db-performance-tests/bin/db-performance-tests" ./run_command.sh
    ```
    
1.  Report compliance (select/update/upsert 100 thousand rows):

    ```
    VM_NAME='veisenbe-postgresql-perf-client*' COMMAND="export UPDATE=; hub-of-hubs-postgresql/db-performance-tests/bin/db-performance-tests" ./run_command.sh
    ```

1.  Update all to be compliant (100 million rows):

    ```
    VM_NAME='veisenbe-postgresql-perf-client*' COMMAND="export UPDATE_ALL=; hub-of-hubs-postgresql/db-performance-tests/bin/db-performance-tests" ./run_command.sh
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
