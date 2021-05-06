# Commands to work with Hub-of-Hubs database

* Configure `root.crt` on the client machine. For example, for
[Let's encrypt](https://letsencrypt.org/) certificates, run `curl -s https://letsencrypt.org/certs/trustid-x3-root.pem --output ~/.postgresql/root.crt`.

* Connect to the database by `psql`:

```
PGSSLMODE=verify-full psql -h hohdb.dev10.red-chesterfield.com -U hoh_process_user -W -d hoh
```

* Insert a Policy from a YAML file into `policies` table:

```
hoh=> \set policy `curl -s https://raw.githubusercontent.com/open-cluster-management/policy-collection/main/community/AC-Access-Control/policy-roles-no-wildcards.yaml | yq r -d'*' - -j | jq -c '.[] | select (.kind == "Policy")'`
insert into spec.policies (payload) values(:'policy');
```

Another policy:
```
hoh=> \set policy `curl -s https://raw.githubusercontent.com/open-cluster-management/policy-collection/main/community/AU-Audit-and-Accountability/policy-openshift-audit-logs-sample.yaml | yq r -d'*' - -j | jq -c '.[] | select (.kind == "Policy")'`
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
