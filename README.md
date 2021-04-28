# Install and configure Hub-of-Hubs database

PostgreSQL serves as the database of Hub-of-Hubs
The host group in the commands is `acm_aws`.

## To install

Run:

```
ansible-playbook install.yaml -i production --ask-vault-pass -l acm_aws
```

## To uninstall

Run:

```
ansible-playbook uninstall_spokes.yaml -i production --ask-vault-pass -l <host group>
```
