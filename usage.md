# Commands to work with Hub-of-Hubs database

* Connect to the database by `psql`:

```
PGSSLMODE=verify-full psql -h hohdb.dev10.red-chesterfield.com -U hoh_process_user -W -d hoh
```

* Insert a Policy from a YAML file into `policies` table:

```
hoh=> \set policy `yq r -d'*'  ~/dev/ACMSamples/policies/pod_security_policy/policy-psp.yaml -j | jq -c '.[] | select (.kind == "Policy")'`
insert into spec.policies (payload) values(:'policy');
```

* Select policy name from the `policies` table:

```
select payload -> 'metadata' -> 'name' as name from spec.policies;
            name            
----------------------------
 "policy-podsecuritypolicy"

```

* Select policy name and remediation action from the `policies` table:

```
hoh=> select payload -> 'metadata' -> 'name' as name, payload -> 'spec' -> 'remediationAction' as action from spec.policies;
            name            |  action  
----------------------------+----------
 "policy-podsecuritypolicy" | "inform"
(1 row)

```

* Update remediation action of a policy to be `enforce`:

```
update spec.policies set payload = jsonb_set(payload, '{spec,remediationAction}', '"enforce"', true) where payload -> 'metadata' ->> 'name' = 'policy-podsecuritypolicy';
```
