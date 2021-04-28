# Install and configure Hub-of-Hubs database

PostgreSQL serves as the database of Hub-of-Hubs.

## Initial setup

1.  Create `production` file in the main directory with the host name of your machine to install the database, under the `acm` host group `acm`:
    ```
    [acm]
    <your host>
    ```
1.  Create `vault` file with following variables:
    - `ansible_user`: contains the user of the machine where you install the database
    - `vault_ansible_ssh_private_key_file`: the path to the SSSH private key file to connect to the machine

## To install

Run:

```
ansible-playbook install.yaml -i production --ask-vault-pass -l acm
```

## To uninstall

Run:

```
ansible-playbook uninstall.yaml -i production --ask-vault-pass -l acm
```
