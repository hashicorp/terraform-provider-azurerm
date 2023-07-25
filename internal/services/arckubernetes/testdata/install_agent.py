# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

import argparse
import json
import logging as logger
import os
import platform
import shutil
import stat
import subprocess
import time
import urllib
from subprocess import PIPE, Popen
from urllib import request

HELM_VERSION = 'v3.6.3'
HELM_STORAGE_URL = "https://k8connecthelm.azureedge.net"
Pre_Onboarding_Helm_Charts_Folder_Name = 'PreOnboardingChecksCharts'


def get_helm_registry(config_dp_endpoint):
    # Setting uri
    get_chart_location_url = "{}/{}/GetLatestHelmPackagePath?api-version=2019-11-01-preview".format(
        config_dp_endpoint, 'azure-arc-k8sagents')

    try:
        response = urllib.request.urlopen(
            request.Request(get_chart_location_url, method="POST"))
    except Exception as e:
        raise Exception("Failed to get helm registry." + str(e))

    try:
        return json.load(response).get('repositoryPath')
    except Exception as e:
        raise Exception(
            "Error while fetching helm chart registry path from JSON response: " + str(e))


def pull_helm_chart(registry_path, kube_config, kube_context, helm_client_location, chart_name='azure-arc-k8sagents', retry_count=5, retry_delay=3):
    cmd_helm_chart_pull = [helm_client_location,
                           "chart", "pull", registry_path]
    if kube_config:
        cmd_helm_chart_pull.extend(["--kubeconfig", kube_config])
    if kube_context:
        cmd_helm_chart_pull.extend(["--kube-context", kube_context])
    for i in range(retry_count):
        response_helm_chart_pull = subprocess.Popen(
            cmd_helm_chart_pull, stdout=PIPE, stderr=PIPE)
        _, error_helm_chart_pull = response_helm_chart_pull.communicate()
        if response_helm_chart_pull.returncode != 0:
            if i == retry_count - 1:
                raise Exception("Unable to pull {} helm chart from the registry '{}': ".format(
                    chart_name, registry_path) + error_helm_chart_pull.decode("ascii"))
            time.sleep(retry_delay)
        else:
            break


def export_helm_chart(registry_path, chart_export_path, kube_config, kube_context, helm_client_location, chart_name='azure-arc-k8sagents'):
    cmd_helm_chart_export = [helm_client_location, "chart",
                             "export", registry_path, "--destination", chart_export_path]
    if kube_config:
        cmd_helm_chart_export.extend(["--kubeconfig", kube_config])
    if kube_context:
        cmd_helm_chart_export.extend(["--kube-context", kube_context])
    response_helm_chart_export = subprocess.Popen(
        cmd_helm_chart_export, stdout=PIPE, stderr=PIPE)
    _, error_helm_chart_export = response_helm_chart_export.communicate()
    if response_helm_chart_export.returncode != 0:
        raise Exception("Unable to export {} helm chart from the registry '{}': ".format(
            chart_name, registry_path) + error_helm_chart_export.decode("ascii"))


def get_chart_path(registry_path, kube_config, kube_context, helm_client_location, chart_folder_name='AzureArcCharts', chart_name='azure-arc-k8sagents'):
    # Pulling helm chart from registry
    os.environ['HELM_EXPERIMENTAL_OCI'] = '1'
    pull_helm_chart(registry_path, kube_config, kube_context,
                    helm_client_location, chart_name)

    # Exporting helm chart after cleanup
    chart_export_path = os.path.join(
        os.path.expanduser('~'), '.azure', chart_folder_name)
    try:
        if os.path.isdir(chart_export_path):
            shutil.rmtree(chart_export_path)
    except:
        logger.warning("Unable to cleanup the {} already present on the machine. In case of failure, please cleanup the directory '{}' and try again.".format(
            chart_folder_name, chart_export_path))

    export_helm_chart(registry_path, chart_export_path, kube_config,
                      kube_context, helm_client_location, chart_name)

    # Returning helm chart path
    helm_chart_path = os.path.join(chart_export_path, chart_name)
    if chart_folder_name == Pre_Onboarding_Helm_Charts_Folder_Name:
        chart_path = helm_chart_path
    else:
        chart_path = os.getenv('HELMCHART') if os.getenv(
            'HELMCHART') else helm_chart_path

    return chart_path


def install_helm_client():
    # Fetch system related info
    operating_system = platform.system().lower()
    platform.machine()

    # Set helm binary download & install locations
    if (operating_system == 'windows'):
        download_location_string = f'.azure\\helm\\{HELM_VERSION}\\helm-{HELM_VERSION}-{operating_system}-amd64.zip'
        install_location_string = f'.azure\\helm\\{HELM_VERSION}\\{operating_system}-amd64\\helm.exe'
        requestUri = f'{HELM_STORAGE_URL}/helm/helm-{HELM_VERSION}-{operating_system}-amd64.zip'
    elif (operating_system == 'linux' or operating_system == 'darwin'):
        download_location_string = f'.azure/helm/{HELM_VERSION}/helm-{HELM_VERSION}-{operating_system}-amd64.tar.gz'
        install_location_string = f'.azure/helm/{HELM_VERSION}/{operating_system}-amd64/helm'
        requestUri = f'{HELM_STORAGE_URL}/helm/helm-{HELM_VERSION}-{operating_system}-amd64.tar.gz'
    else:
        raise Exception(
            f'The {operating_system} platform is not currently supported for installing helm client.')

    download_location = os.path.expanduser(
        os.path.join('~', download_location_string))
    download_dir = os.path.dirname(download_location)
    install_location = os.path.expanduser(
        os.path.join('~', install_location_string))

    # Download compressed halm binary if not already present
    if not os.path.isfile(download_location):
        # Creating the helm folder if it doesnt exist
        if not os.path.exists(download_dir):
            try:
                os.makedirs(download_dir)
            except Exception as e:
                raise Exception("Failed to create helm directory." + str(e))

        # Downloading compressed helm client executable
        logger.warning(
            "Downloading helm client for first time. This can take few minutes...")
        try:
            response = urllib.request.urlopen(requestUri)
        except Exception as e:
            raise Exception("Failed to download helm client." + str(e))

        responseContent = response.read()
        response.close()

        # Creating the compressed helm binaries
        try:
            with open(download_location, 'wb') as f:
                f.write(responseContent)
        except Exception as e:
            raise Exception("Failed to create helm executable." + str(e))

    # Extract compressed helm binary
    if not os.path.isfile(install_location):
        try:
            shutil.unpack_archive(download_location, download_dir)
            os.chmod(install_location, os.stat(
                install_location).st_mode | stat.S_IXUSR)
        except Exception as e:
            raise Exception("Failed to extract helm executable." + str(e))

    return install_location


def helm_install_release(chart_path, subscription_id, kubernetes_distro, kubernetes_infra, resource_group_name, cluster_name,
                         location, onboarding_tenant_id, private_key_pem,
                         no_wait, cloud_name, helm_client_location, onboarding_timeout="600"):
    cmd_helm_install = [helm_client_location, "upgrade", "--install", "azure-arc", chart_path,
                        "--set", "global.subscriptionId={}".format(
                            subscription_id),
                        "--set", "global.kubernetesDistro={}".format(
                            kubernetes_distro),
                        "--set", "global.kubernetesInfra={}".format(
                            kubernetes_infra),
                        "--set", "global.resourceGroupName={}".format(
                            resource_group_name),
                        "--set", "global.resourceName={}".format(cluster_name),
                        "--set", "global.location={}".format(location),
                        "--set", "global.tenantId={}".format(
                            onboarding_tenant_id),
                        "--set", "global.onboardingPrivateKey={}".format(
                            private_key_pem),
                        "--set", "systemDefaultValues.spnOnboarding=false",
                        "--set", "global.azureEnvironment={}".format(
                            cloud_name),
                        "--set", "systemDefaultValues.clusterconnect-agent.enabled=true",
                        "--namespace", "{}".format("azure-arc-release"),
                        "--create-namespace",
                        "--output", "json"]

    if not no_wait:
        # Change --timeout format for helm client to understand
        onboarding_timeout = onboarding_timeout + "s"
        cmd_helm_install.extend(
            ["--wait", "--timeout", "{}".format(onboarding_timeout)])
    response_helm_install = Popen(cmd_helm_install, stdout=PIPE, stderr=PIPE)
    _, error_helm_install = response_helm_install.communicate()
    if response_helm_install.returncode != 0:
        raise Exception("Unable to install helm release" + error_helm_install.decode("ascii"))


def install_agent():
    parser = argparse.ArgumentParser(
        description='Install Connected Cluster Agent')
    parser.add_argument('--subscriptionId', type=str, required=True)
    parser.add_argument('--resourceGroupName', type=str, required=True)
    parser.add_argument('--clusterName', type=str, required=True)
    parser.add_argument('--location', type=str, required=True)
    parser.add_argument('--tenantId', type=str, required=True)
    parser.add_argument('--privatePem', type=str, required=True)

    try:
        args = parser.parse_args()
    except Exception as e:
        raise Exception("Failed to parse arguments." + str(e))

    try:
        with open(args.privatePem, "r") as f:
            privateKey = f.read()
    except Exception as e:
        raise Exception("Failed to get private key." + str(e))

    # Install helm client
    helm_client_location = install_helm_client()

    # Retrieving Helm chart OCI Artifact location
    registry_path = get_helm_registry("https://westeurope.dp.kubernetesconfiguration.azure.com")
    
    # Get helm chart path
    chart_path = get_chart_path(
        registry_path, None, None, helm_client_location)

    helm_install_release(chart_path,
                         args.subscriptionId,
                         "generic",
                         "generic",
                         args.resourceGroupName,
                         args.clusterName,
                         args.location,
                         args.tenantId,
                         privateKey,
                         False,
                         "AZUREPUBLICCLOUD",
                         helm_client_location)


if __name__ == "__main__":
    install_agent()
