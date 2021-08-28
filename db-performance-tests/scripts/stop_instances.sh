#!/bin/bash

# Copyright (c) 2021 Red Hat, Inc.
# Copyright Contributors to the Open Cluster Management project

set -o errexit
set -o nounset

postgresql_port=5432

instance_ids=$(aws ec2 describe-instances --filters "Name=tag:Name,Values=$VM_NAME_PREFIX*"  --output json | jq -r '.Reservations[].Instances[].InstanceId')
echo handling "$instance_ids"

# shellcheck disable=SC2086
public_ips=$(aws ec2 describe-instances  --instance-ids $instance_ids --output json | jq -r '.Reservations[].Instances[].PublicIpAddress')

for ip in $public_ips
do
# shellcheck disable=SC2086
	  aws ec2 revoke-security-group-ingress --group-name "$SECURITY_GROUP" --protocol tcp --port $postgresql_port --cidr $ip/32
done

# shellcheck disable=SC2086
aws ec2 stop-instances --instance-ids $instance_ids > /dev/null

# shellcheck disable=SC2086
aws ec2 wait instance-stopped --instance-ids $instance_ids
echo all the instances stopped running
