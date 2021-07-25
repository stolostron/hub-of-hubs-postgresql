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
ROWS_NUMBER=1000 INSERT_MULTIPLE_VALUES= ./bin/performance_tests
```

Insert by the COPY protocol:

```
ROWS_NUMBER=1000 INSERT_COPY= ./bin/performance_tests
```
