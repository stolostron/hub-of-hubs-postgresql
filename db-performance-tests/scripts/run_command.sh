#!/bin/bash

# Copyright (c) 2021 Red Hat, Inc.
# Copyright Contributors to the Open Cluster Management project

set -o errexit
set -o nounset

instance_ids=$(aws ec2 describe-instances --filters "Name=tag:Name,Values=$VM_NAME"  --output json | jq -r '.Reservations[].Instances[].InstanceId')

instance_ids_array=($instance_ids)
instance_ids_number=${#instance_ids_array[@]}

echo handling ${instance_ids_number} instances:
echo "$instance_ids"

# shellcheck disable=SC2086
public_dns_names=$(aws ec2 describe-instances  --instance-ids $instance_ids --output json | jq -r '.Reservations[].Instances[].PublicDnsName')

start_time=$(date +"%s")

total_leaf_hubs_number=${TOTAL_LEAF_HUBS:-1000}
max_connections_number=48
max_connections_number_per_instance=$(( $max_connections_number / $instance_ids_number ))
leaf_hubs_number_per_instance=$(( $total_leaf_hubs_number  / $instance_ids_number ))
start_leaf_hub_index=0

for dns_name in $public_dns_names
do
# shellcheck disable=SC2086
    ssh -i ${AWS_SSH_KEY} -o StrictHostKeyChecking=no ec2-user@$dns_name \
	 "export DATABASE_URL=\"$DATABASE_URL&pool_max_conns=$max_connections_number_per_instance\";" \
	 "export LEAF_HUBS_NUMBER=$leaf_hubs_number_per_instance;" \
	 "export START_LEAF_HUB_INDEX=$start_leaf_hub_index;" $COMMAND &
    start_leaf_hub_index=$(( $start_leaf_hub_index + $leaf_hubs_number_per_instance ))
done

wait
end_time=$(date +"%s")

echo elapsed $(($end_time - $start_time)) seconds
