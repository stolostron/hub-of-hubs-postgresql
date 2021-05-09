The database contains a single `database` construct - `hoh`. The `hoh` database contains two schemas `spec` and `status`. 
The Hub-of-Hubs logic will update the tables in the `spec` schema, the transport logic will update tables in the `status` schema. 
Both components will read tables from the both schemas.

The database will contain tables corresponding to the Hub-of-Hubs CRDs, for example `policies`, `placementrules`, `placementbindings`. 
Each table will contain CRs in the `paylod` field of type [jsonb](https://www.postgresql.org/docs/9.4/datatype-json.html). In addition to the `payload` field,
each table will contain a [surrogate key](https://en.wikipedia.org/wiki/Surrogate_key) `id` and two timestamps `created_at` and `updated_at`
that will be updated automatically by the database. See [example queries](https://github.com/open-cluster-management/hub-of-hubs-postgresql/blob/main/usage.md).

