# Commands to work with Hub-of-Hubs database

* Configure `root.crt` on the client machine. For example, for
[Let's encrypt](https://letsencrypt.org/) certificates, run `curl -s https://letsencrypt.org/certs/trustid-x3-root.pem --output ~/.postgresql/root.crt`.

* Connect to the database by `psql`:

```
PGSSLMODE=verify-full psql -h hohdb.dev10.red-chesterfield.com -U hoh_process_user -W -d hoh
```

* Insert a Policy from a YAML file into `policies` table:

```
hoh=> \set policy `curl -s https://raw.githubusercontent.com/open-cluster-management/policy-collection/main/community/AU-Audit-and-Accountability/policy-openshift-audit-logs-sample.yaml | yq eval 'select(.kind == "Policy")' - -j`
insert into spec.policies (payload) values(:'policy');
```

Another policy:
```
hoh=> \set policy `curl -s https://raw.githubusercontent.com/open-cluster-management/policy-collection/main/community/AC-Access-Control/policy-roles-no-wildcards.yaml | yq eval 'select(documentIndex == 0)' - -j`
insert into spec.policies (payload) values(:'policy');
```


* Select policy name from the `policies` table:

```
select payload -> 'metadata' -> 'name' as name from spec.policies;
             name              
-------------------------------
 "policy-openshift-audit-logs"
 "policy-disallowed-roles"
(2 rows)

```

* Select policy name and remediation action from the `policies` table:

```
hoh=> select payload -> 'metadata' -> 'name' as name, payload -> 'spec' -> 'remediationAction' as action from spec.policies where payload -> 'metadata' ->> 'name' = 'policy-disallowed-roles';
             name              |  action  
-------------------------------+----------
 "policy-disallowed-roles"     | "inform"
 (1 row)

```

* select GVK
```
select payload -> 'metadata' ->> 'name' as name, (payload ->> 'apiVersion')||'/'||(payload ->> 'kind') as gvk,payload -> 'spec' ->> 'remediationAction' as action from spec.policies where payload -> 'metadata' ->> 'name' = 'policy-disallowed-roles';
          name           |                     gvk                     | action  
-------------------------+---------------------------------------------+---------
 policy-disallowed-roles | policy.open-cluster-management.io/v1/Policy | enforce
```

* select group and kind
```
select payload -> 'metadata' ->> 'name' as name, split_part(payload ->> 'apiVersion', '/', 1) as group, payload ->> 'kind' as kind,payload -> 'spec' ->> 'remediationAction' as action from spec.policies where payload -> 'metadata' ->> 'name' = 'policy-disallowed-roles';
```

* Update remediation action of a policy to be `enforce`:

```
update spec.policies set payload = jsonb_set(payload, '{spec,remediationAction}', '"enforce"', true) where payload -> 'metadata' ->> 'name' = 'policy-disallowed-roles';
```

* Find matching placement rules and placement bindings:

```
select pr.payload -> 'metadata' -> 'name' as policyrulename, pb.payload -> 'metadata' -> 'name' as placementbindingname from spec.placementrules pr INNER JOIN  spec.placementbindings pb ON pr.payload -> 'metadata' ->> 'name' = pb.payload -> 'placementRef' ->> 'name' AND pr.payload ->> 'kind' = pb.payload -> 'placementRef' ->> 'kind' AND split_part(pr.payload ->> 'apiVersion', '/', 1) = pb.payload -> 'placementRef' ->> 'apiGroup';
             policyrulename              |         placementbindingname          
-----------------------------------------+---------------------------------------
 "placement-policy-disallowed-roles"     | "binding-policy-disallowed-roles"
 "placement-policy-openshift-audit-logs" | "binding-policy-openshift-audit-logs"
(2 rows)

```

* select name, kind, group as json

```
select json_build_object( 'name',p.payload -> 'metadata' ->> 'name', 'kind', p.payload ->> 'kind', 'apiGroup', split_part(p.payload ->> 'apiVersion', '/',1)) as name_kind_group from spec.policies p;
                                                name_kind_group                                                
---------------------------------------------------------------------------------------------------------------
 {"name" : "policy-openshift-audit-logs", "kind" : "Policy", "apiGroup" : "policy.open-cluster-management.io"}
 {"name" : "policy-disallowed-roles", "kind" : "Policy", "apiGroup" : "policy.open-cluster-management.io"}
(2 rows)
```

* select matching policy, placement rule and placement binding
```
select p.payload -> 'metadata' ->> 'name' as policy, pb.payload -> 'metadata' ->> 'name' as binding, pr.payload -> 'metadata' ->> 'name' as placementrule from spec.policies p INNER JOIN spec.placementbindings pb ON pb.payload -> 'subjects' @> json_build_array(json_build_object( 'name',p.payload -> 'metadata' ->> 'name', 'kind', p.payload ->> 'kind', 'apiGroup', split_part(p.payload ->> 'apiVersion', '/',1)))::jsonb INNER JOIN spec.placementrules pr ON pr.payload -> 'metadata' ->> 'name' = pb.payload -> 'placementRef' ->> 'name' AND pr.payload ->> 'kind' = pb.payload -> 'placementRef' ->> 'kind' AND split_part(pr.payload ->> 'apiVersion', '/', 1) = pb.payload -> 'placementRef' ->> 'apiGroup';

```
