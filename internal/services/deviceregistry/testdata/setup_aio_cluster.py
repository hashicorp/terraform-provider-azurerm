# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# This python script is used to arc enable and install Azure IoT Operations extension.
# This should be identical to the script used in the AIO onboarding flow, plus setting custom location enabled.
# See this link for the exact commands: https://review.learn.microsoft.com/en-us/azure/iot-operations/get-started-end-to-end-sample/quickstart-deploy?branch=main

import argparse
import json
import logging as logger
import os
import platform
import shutil
import signal
import stat
import subprocess
import time
import urllib
from subprocess import PIPE, Popen
from urllib import request

def register_az_providers():
    print("Registering required Azure providers...")
    required_providers = [
        "Microsoft.ExtendedLocation",
        "Microsoft.Kubernetes",
        "Microsoft.KubernetesConfiguration",
        "Microsoft.IoTOperations",
        "Microsoft.DeviceRegistry",
        "Microsoft.SecretSyncController"
    ]
    for provider in required_providers:
        # Equivalent of `az provider register -n <provider>`
        print("Registering provider " + provider + "...")
        response = Popen(["az", "provider", "register", "-n", provider], stdout=PIPE, stderr=PIPE)
        _, error = response.communicate()
        if response.returncode != 0:
            raise Exception("Failed to register provider " + provider + ": " + error.decode("ascii"))
        print("Successfully registered provider " + provider)
    
    print("Successfully registered required Azure providers")
    return

def onboard_k8s_cluster(resource_group_name, cluster_name, location):
    print("Connecting k8s cluster " + cluster_name + " to Azure...")

    # Equivalent of `az connectedk8s connect --name $CLUSTER_NAME --location $LOCATION --resource-group $RESOURCE_GROUP`
    response = Popen(["az", "connectedk8s", "connect", "--name", cluster_name, "--location", location, "--resource-group", resource_group_name], stdout=PIPE, stderr=PIPE)
    _, error = response.communicate()
    if response.returncode != 0:
        raise Exception("Failed to connect k8s cluster " + cluster_name + ": " + error.decode("ascii"))
    
    print("Successfully connected k8s cluster " + cluster_name + " to Azure. Enabling features on the cluster...")

    # Equivalent of `$OBJECT_ID=$(az ad sp show --id bc313c14-388c-4e7d-a58e-70017303ee3b --query id -o tsv)`
    response = Popen(["az", "ad", "sp", "show", "--id", "bc313c14-388c-4e7d-a58e-70017303ee3b", "--query", "id", "-o", "tsv"], stdout=PIPE, stderr=PIPE)
    object_id, error = response.communicate()
    if response.returncode != 0:
        raise Exception("Failed to get object id of Microsoft Entra ID Application for Azure Arc service: " + error.decode("ascii"))

    print("Enabling features on the cluster...")
    # Equivalent of `az connectedk8s enable-features -n $CLUSTER_NAME -g $RESOURCE_GROUP --custom-locations-oid $OBJECT_ID --features cluster-connect custom-locations`
    response = Popen(["az", "connectedk8s", "enable-features", "-n", cluster_name, "-g", resource_group_name, "--custom-locations-oid", object_id.strip(), "--features", "cluster-connect", "custom-locations"], stdout=PIPE, stderr=PIPE)
    _, error = response.communicate()
    if response.returncode != 0:
        raise Exception("Failed to enable features on k8s cluster " + cluster_name + ": " + error.decode("ascii"))
    
    print("Successfully enabled features on k8s cluster " + cluster_name)
    return

def setup_schema_registry(storage_account, schema_registry, schema_registry_namespace, resource_group_name, location):
    # print("Setting up storage account " + storage_account + " for schema registry " + schema_registry + "...")
    # Equivalent of `az storage account create --name $STORAGE_ACCOUNT --location $LOCATION --resource-group $RESOURCE_GROUP --enable-hierarchical-namespace`
    response = Popen(["az", "storage", "account", "create", "--name", storage_account, "--location", location, "--resource-group", resource_group_name, "--enable-hierarchical-namespace"], stdout=PIPE, stderr=PIPE)
    _, error = response.communicate()
    if response.returncode != 0:
        raise Exception("Failed to create storage account " + storage_account + ": " + error.decode("ascii"))

    print("Successfully created storage account " + storage_account + ". Setting up schema registry " + schema_registry)

    # Equivalent of `az storage account show --name $STORAGE_ACCOUNT -o tsv --query id`
    response = Popen(["az", "storage", "account", "show", "--name", storage_account, "-o", "tsv", "--query", "id"], stdout=PIPE, stderr=PIPE)
    storage_account_id, error = response.communicate()
    if response.returncode != 0:
        raise Exception("Failed to get storage account id of " + storage_account + ": " + error.decode("ascii"))

    # Equivalent of `az iot ops schema registry create --name $SCHEMA_REGISTRY --resource-group $RESOURCE_GROUP --registry-namespace $SCHEMA_REGISTRY_NAMESPACE --sa-resource-id $STORAGE_ACCOUNT_ID`
    response = Popen(["az", "iot", "ops", "schema", "registry", "create", "--name", schema_registry, "--resource-group", resource_group_name, "--registry-namespace", schema_registry_namespace, "--sa-resource-id", storage_account_id.strip()], stdout=PIPE, stderr=PIPE)
    _, error = response.communicate()
    if response.returncode != 0:
        raise Exception("Failed to create schema registry " + schema_registry + ": " + error.decode("ascii"))
    
    print("Successfully created schema registry " + schema_registry)
    return

def install_azure_iot_ops_extension(cluster_name, resource_group_name, schema_registry, custom_location, aio_create_timeout=15*60):
    print("Initializing cluster for installing Azure IoT Operations extension...")

    # Run `az iot ops init` to initialize + prepare cluster for aio ext install
    # Equivalent of `az iot ops init --cluster $CLUSTER_NAME --resource-group $RESOURCE_GROUP --no-progress`
    response = Popen(["az", "iot", "ops", "init", "--cluster", cluster_name, "--resource-group", resource_group_name, "--no-progress"], stdout=PIPE, stderr=PIPE)
    _, error = response.communicate()
    if response.returncode != 0:
        raise Exception("Failed to initialize Azure IoT Operations extension on cluster " + cluster_name + ": " + error.decode("ascii"))
    
    print("Successfully initialized cluster. Installing Azure IoT Operations extension...")

    # Get the schema registry resource ID
    # Equivalent of `az iot ops schema registry show --name $SCHEMA_REGISTRY --resource-group $RESOURCE_GROUP -o tsv --query id`
    response = Popen(["az", "iot", "ops", "schema", "registry", "show", "--name", schema_registry, "--resource-group", resource_group_name, "-o", "tsv", "--query", "id"], stdout=PIPE, stderr=PIPE)
    schema_registry_id, error = response.communicate()
    if response.returncode != 0:
        raise Exception("Failed to get schema registry id of " + schema_registry + ": " + error.decode("ascii"))

    # Install the Azure IoT Operations extension onto cluster.
        # Note: there is a known issue when trying to install the extension on a Kind cluster.
        # The schema registry portion of the extension will timeout after 30+ minutes and fail to install.
        # But this will not block other parts of AIO installation and the acceptance tests for device registry,
        # so set a timeout cancellation for this step to avoid blocking the rest of the acceptance tests.
    try:
        # Equivalent of `az iot ops create --cluster $CLUSTER_NAME --resource-group $RESOURCE_GROUP --name ${CLUSTER_NAME}-instance  --sr-resource-id $SCHEMA_REGISTRY_ID --broker-frontend-replicas 1 --broker-frontend-workers 1  --broker-backend-part 1  --broker-backend-workers 1 --broker-backend-rf 2 --broker-mem-profile Low --custom-location $CUSTOM_LOCATION --no-progress --yes`, 
        # along with a 15 minute timeout using subprocess.run()
        response = subprocess.run(["az", "iot", "ops", "create", "--cluster", cluster_name, "--resource-group", resource_group_name, "--name", cluster_name + "-instance", "--sr-resource-id", schema_registry_id.strip(), "--broker-frontend-replicas", "1", "--broker-frontend-workers", "1", "--broker-backend-part", "1", "--broker-backend-workers", "1", "--broker-backend-rf", "2", "--broker-mem-profile", "Low", "--custom-location", custom_location, "--no-progress", "--yes"], timeout=aio_create_timeout, stdout=PIPE, stderr=PIPE)
        # _, error = response.communicate()
        print(response)
        print("Successfully installed Azure IoT Operations extension on cluster " + cluster_name)
        return
    except subprocess.TimeoutExpired:
        # Swallow timeout exception and continue with the rest of the acceptance tests
        print("Warning: installing Azure IoT Operations extension on cluster " + cluster_name + " timed out after " + str(aio_create_timeout) + " seconds but should be working for tests. Continuing with the rest of the acceptance tests")
        return
    except Exception as e:
        raise Exception("Failed to install Azure IoT Operations extension on cluster " + cluster_name + ": " + str(e))

def setup_aio_arc_enabled_cluster():
    parser = argparse.ArgumentParser(
        description='Install AIO extension onto cluster')
    parser.add_argument('--subscriptionId', type=str, required=True)
    parser.add_argument('--resourceGroupName', type=str, required=True)
    parser.add_argument('--clusterName', type=str, required=True)
    parser.add_argument('--location', type=str, required=True)
    parser.add_argument('--customLocation', type=str, required=True)
    parser.add_argument('--storageAccount', type=str, required=True)
    parser.add_argument('--schemaRegistry', type=str, required=True)
    parser.add_argument('--schemaRegistryNamespace', type=str, required=True)
    parser.add_argument('--tenantId', type=str, required=True)
    # parser.add_argument('--privatePem', type=str, required=True)

    try:
        args = parser.parse_args()
    except Exception as e:
        raise Exception("Failed to parse arguments." + str(e))

    # try:
    #     with open(args.privatePem, "r") as f:
    #         privateKey = f.read()
    # except Exception as e:
    #     raise Exception("Failed to get private key." + str(e))

    register_az_providers()

    onboard_k8s_cluster(args.resourceGroupName, args.clusterName, args.location)

    setup_schema_registry(args.storageAccount, args.schemaRegistry, args.schemaRegistryNamespace, args.resourceGroupName, args.location)

    install_azure_iot_ops_extension(args.clusterName, args.resourceGroupName, args.schemaRegistry, args.customLocation)

if __name__ == "__main__":
    setup_aio_arc_enabled_cluster()
