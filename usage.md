# Commands to work with Hub-of-Hubs database

## Prerequisites

* [install psql on your system](https://blog.timescale.com/tutorials/how-to-install-psql-on-mac-ubuntu-debian-windows/)

* Configure `root.crt` on the client machine. For example, for
[Let's encrypt](https://letsencrypt.org/) certificates, run `curl -s https://letsencrypt.org/certs/trustid-x3-root.pem --output ~/.postgresql/root.crt`.

* Connect to the database by `psql`:

```
PGSSLMODE=verify-full psql -h hohdb.dev10.red-chesterfield.com -U hoh_process_user -W -d hoh
```

# Queries
* Insert a Policy from a YAML file into `policies` table:

```
hoh=> \set policy `curl -s https://raw.githubusercontent.com/open-cluster-management/policy-collection/main/community/AU-Audit-and-Accountability/policy-openshift-audit-logs-sample.yaml | yq eval 'select(.kind == "Policy")' - -j`
insert into spec.policies (payload) values(:'policy');
```

Insert policies, placementrules and placementbindings using instructions in [test/insertpolicies.psql](test/insertpolicies.psql):
```
PGSSLMODE=verify-full psql -h hohdb.dev10.red-chesterfield.com -U hoh_process_user -W -d hoh -f test/insertpolicies.psql
```

* Select policy name from the `policies` table:

```sql
select payload -> 'metadata' -> 'name' as name from spec.policies;
             name              
-------------------------------
 "policy-openshift-audit-logs"
 "policy-disallowed-roles"
(2 rows)

```

* Select policy name and remediation action from the `policies` table:

```sql
hoh=> select payload -> 'metadata' -> 'name' as name, payload -> 'spec' -> 'remediationAction' as action from spec.policies where payload -> 'metadata' ->> 'name' = 'policy-disallowed-roles';
             name              |  action  
-------------------------------+----------
 "policy-disallowed-roles"     | "inform"
 (1 row)

```

* select GVK
```sql
select payload -> 'metadata' ->> 'name' as name, (payload ->> 'apiVersion')||'/'||(payload ->> 'kind') as gvk,payload -> 'spec' ->> 'remediationAction' as action from spec.policies where payload -> 'metadata' ->> 'name' = 'policy-disallowed-roles';
          name           |                     gvk                     | action  
-------------------------+---------------------------------------------+---------
 policy-disallowed-roles | policy.open-cluster-management.io/v1/Policy | enforce
```

* select group and kind
```sql
select payload -> 'metadata' ->> 'name' as name, split_part(payload ->> 'apiVersion', '/', 1) as group, payload ->> 'kind' as kind,payload -> 'spec' ->> 'remediationAction' as action from spec.policies where payload -> 'metadata' ->> 'name' = 'policy-disallowed-roles';
```

* select id, name, namespace, remediation policy, created at,updated at, deleted:

```
select id, created_at, updated_at, payload -> 'metadata' -> 'name' as name, payload -> 'metadata' -> 'namespace' as namespace,payload -> 'spec' -> 'remediationAction' as action, deleted from spec.policies;
                  id                  |         created_at         |         updated_at         |            name            |  namespace  |  action  | deleted 
--------------------------------------+----------------------------+----------------------------+----------------------------+-------------+----------+---------
 6261f452-4df5-4ce3-959a-513f0a558866 | 2021-05-23 20:20:47.489224 | 2021-05-23 20:20:47.489224 | "policy-podsecuritypolicy" | "myproject" | "inform" | f
```

* Update remediation action of a policy to be `enforce`:

```sql
update spec.policies set payload = jsonb_set(payload, '{spec,remediationAction}', '"enforce"', true) where payload -> 'metadata' ->> 'name' = 'policy-disallowed-roles';
```

* Find matching placement rules and placement bindings:

```sql
select pr.payload -> 'metadata' -> 'name' as policyrulename, pb.payload -> 'metadata' -> 'name' as placementbindingname from spec.placementrules pr INNER JOIN  spec.placementbindings pb ON pr.payload -> 'metadata' ->> 'name' = pb.payload -> 'placementRef' ->> 'name' AND pr.payload ->> 'kind' = pb.payload -> 'placementRef' ->> 'kind' AND split_part(pr.payload ->> 'apiVersion', '/', 1) = pb.payload -> 'placementRef' ->> 'apiGroup';
             policyrulename              |         placementbindingname          
-----------------------------------------+---------------------------------------
 "placement-policy-disallowed-roles"     | "binding-policy-disallowed-roles"
 "placement-policy-openshift-audit-logs" | "binding-policy-openshift-audit-logs"
(2 rows)

```

* select name, kind, group as json

```sql
select json_build_object( 'name',p.payload -> 'metadata' ->> 'name', 'kind', p.payload ->> 'kind', 'apiGroup', split_part(p.payload ->> 'apiVersion', '/',1)) as name_kind_group from spec.policies p;
                                                name_kind_group                                                
---------------------------------------------------------------------------------------------------------------
 {"name" : "policy-openshift-audit-logs", "kind" : "Policy", "apiGroup" : "policy.open-cluster-management.io"}
 {"name" : "policy-disallowed-roles", "kind" : "Policy", "apiGroup" : "policy.open-cluster-management.io"}
(2 rows)
```

* select matching policy, placement rule and placement binding
```sql
select p.payload -> 'metadata' ->> 'name' as policy, pb.payload -> 'metadata' ->> 'name' as binding, pr.payload -> 'metadata' ->> 'name' as placementrule from spec.policies p INNER JOIN spec.placementbindings pb ON pb.payload -> 'subjects' @> json_build_array(json_build_object( 'name',p.payload -> 'metadata' ->> 'name', 'kind', p.payload ->> 'kind', 'apiGroup', split_part(p.payload ->> 'apiVersion', '/',1)))::jsonb INNER JOIN spec.placementrules pr ON pr.payload -> 'metadata' ->> 'name' = pb.payload -> 'placementRef' ->> 'name' AND pr.payload ->> 'kind' = pb.payload -> 'placementRef' ->> 'kind' AND split_part(pr.payload ->> 'apiVersion', '/', 1) = pb.payload -> 'placementRef' ->> 'apiGroup';
           policy            |               binding               |             placementrule             
-----------------------------+-------------------------------------+---------------------------------------
 policy-disallowed-roles     | binding-policy-disallowed-roles     | placement-policy-disallowed-roles
 policy-openshift-audit-logs | binding-policy-openshift-audit-logs | placement-policy-openshift-audit-logs

```

* select policies updated within the last hour
```sql
select created_at, updated_at, payload -> 'metadata' -> 'name' as name, payload -> 'spec' -> 'remediationAction' as action from spec.policies where updated_at > now() - interval '1 hour';
         created_at         |         updated_at         |           name            |  action  
----------------------------+----------------------------+---------------------------+----------
 2021-05-08 14:28:42.948361 | 2021-05-09 07:04:48.173653 | "policy-disallowed-roles" | "inform"
(1 row)

```

* print field from type jsonb as pretty json 
```sql
select jsonb_pretty(payload) as payload from status.managed_clusters where cluster_name='cluster6';
                              payload                              
-------------------------------------------------------------------
 {                                                                +
     "kind": "ManagedCluster",                                    +
     "spec": {                                                    +
         "hubAcceptsClient": true,                                +
         "leaseDurationSeconds": 60                               +
     },                                                           +
     "status": {                                                  +
         "version": {                                             +
             "kubernetes": "v1.19.1"                              +
         },                                                       +
         "capacity": {                                            +
             "cpu": "8",                                          +
....

```

* select policy name, namespace, cluster and leaf hub names from the status.compliance table

```sql
select p.payload -> 'metadata' ->> 'name' as policyname, p.payload -> 'metadata' ->> 'namespace' as policynamespace, c.cluster_name, c.leaf_hub_name, c.compliance from spec.policies p INNER JOIN status.compliance c ON p.id = c.policy_id;
        policyname        | policynamespace | cluster_name | leaf_hub_name |  compliance   
--------------------------+-----------------+--------------+---------------+---------------
 policy-podsecuritypolicy | myproject       | cluster0     | hub1          | compliant
 policy-podsecuritypolicy | myproject       | cluster3     | hub1          | non_compliant
```

* select non compliant clusters

```sql
select p.payload -> 'metadata' ->> 'name' as policyname, p.payload -> 'metadata' ->> 'namespace' as policynamespace, c.cluster_name, c.leaf_hub_name, c.compliance from spec.policies p INNER JOIN status.compliance c ON p.id = c.policy_id where c.compliance = 'non_compliant';
   policyname        | policynamespace | cluster_name | leaf_hub_name |  compliance   
--------------------------+-----------------+--------------+---------------+---------------
 policy-podsecuritypolicy | myproject       | cluster3     | hub1          | non_compliant
```

* count the number of non-compliant clusters (each non-compliant cluster is counted once)

```sql
select count(distinct cluster_name) from status.compliance where compliance = 'non_compliant';
```
* select policies with counts of non_compliant clusters

```sql
select p.payload -> 'metadata' ->> 'name' as policyname, p.payload -> 'metadata' ->> 'namespace' as policynamespace, count(c.cluster_name) as non_compliant_count from spec.policies p INNER JOIN status.compliance c ON p.id = c.policy_id where c.compliance = 'non_compliant' GROUP BY policyname, policynamespace;
        policyname        | policynamespace | non_compliant_count 
--------------------------+-----------------+---------------------
 policy-podsecuritypolicy | myproject       |                   2
```

* select policies with counts of compliant clusters:

```sql
select p.payload -> 'metadata' ->> 'name' as policyname, p.payload -> 'metadata' ->> 'namespace' as policynamespace, count(c.cluster_name) as compliant_count from spec.policies p INNER JOIN status.compliance c ON p.id = c.policy_id where c.compliance = 'compliant' GROUP BY policyname, policynamespace;
        policyname        | policynamespace | compliant_count 
--------------------------+-----------------+-----------------
 policy-podsecuritypolicy | myproject       |               1

```

* select distinct clusters that are non compliant
```sql
 select distinct cluster_name from status.compliance where compliance = 'non_compliant';
 cluster_name 
--------------
 cluster3
 cluster7
(2 rows)
```

* select distinct leaf hubs that are non compliant
```sql
select distinct leaf_hub_name from status.compliance where compliance = 'non_compliant';
 leaf_hub_name 
---------------
 hub1
 hub2
(2 rows)
```
