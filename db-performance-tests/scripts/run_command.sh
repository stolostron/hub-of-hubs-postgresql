#!/bin/bash

# Copyright (c) 2021 Red Hat, Inc.
# Copyright Contributors to the Open Cluster Management project

set -o errexit
set -o nounset

instance_ids=$(aws ec2 describe-instances --filters "Name=tag:Name,Values=$VM_NAME_PREFIX*"  --output json | jq -r '.Reservations[].Instances[].InstanceId')

echo handling "$instance_ids"

# shellcheck disable=SC2086
public_dns_names=$(aws ec2 describe-instances  --instance-ids $instance_ids --output json | jq -r '.Reservations[].Instances[].PublicDnsName')

for dns_name in $public_dns_names
do
# shellcheck disable=SC2086
    ssh -i ${AWS_SSH_KEY} -o StrictHostKeyChecking=no ec2-user@$dns_name $COMMAND &
done
