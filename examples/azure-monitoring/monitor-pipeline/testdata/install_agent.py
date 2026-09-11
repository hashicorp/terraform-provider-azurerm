# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

import argparse
import json
import logging
import os
import pathlib
import shutil
import stat
import subprocess
import time
import urllib.request


HELM_VERSION = "v3.6.3"
HELM_STORAGE_URL = "https://k8connecthelm.azureedge.net"
CHART_NAME = "azure-arc-k8sagents"


def install_helm_client(retry_count=3, retry_delay=3):
    archive = pathlib.Path.home() / ".azure" / "helm" / HELM_VERSION / f"helm-{HELM_VERSION}-linux-amd64.tar.gz"
    executable = archive.parent / "linux-amd64" / "helm"

    if executable.is_file():
        return str(executable)

    archive.parent.mkdir(parents=True, exist_ok=True)
    for attempt in range(retry_count):
        try:
            if not archive.is_file():
                logging.warning("Downloading Helm client. This can take a few minutes...")
                temporary_archive = pathlib.Path(f"{archive}.tmp")
                try:
                    with urllib.request.urlopen(f"{HELM_STORAGE_URL}/helm/{archive.name}") as response, temporary_archive.open("wb") as output:
                        shutil.copyfileobj(response, output)
                    temporary_archive.replace(archive)
                except Exception:
                    temporary_archive.unlink(missing_ok=True)
                    raise

            shutil.unpack_archive(archive, archive.parent)
            executable.chmod(executable.stat().st_mode | stat.S_IXUSR)
            return str(executable)
        except Exception:
            archive.unlink(missing_ok=True)
            if attempt == retry_count - 1:
                raise
            logging.warning("Failed to install Helm client; retrying download...")
            time.sleep(retry_delay)


def get_helm_registry():
    endpoint = "https://westeurope.dp.kubernetesconfiguration.azure.com/azure-arc-k8sagents/GetLatestHelmPackagePath?api-version=2019-11-01-preview"
    with urllib.request.urlopen(urllib.request.Request(endpoint, method="POST")) as response:
        return json.load(response)["repositoryPath"]


def get_chart_path(registry_path, helm_client):
    os.environ["HELM_EXPERIMENTAL_OCI"] = "1"
    for attempt in range(5):
        result = subprocess.run([helm_client, "chart", "pull", registry_path], capture_output=True, text=True)
        if result.returncode == 0:
            break
        if attempt == 4:
            raise RuntimeError(f"Unable to pull {CHART_NAME} Helm chart: {result.stderr}")
        time.sleep(3)

    export_path = pathlib.Path.home() / ".azure" / "AzureArcCharts"
    shutil.rmtree(export_path, ignore_errors=True)
    subprocess.run(
        [helm_client, "chart", "export", registry_path, "--destination", str(export_path)],
        check=True,
    )
    return str(export_path / CHART_NAME)


def install_agent():
    parser = argparse.ArgumentParser(description="Install Connected Cluster Agent")
    parser.add_argument("--subscriptionId", required=True)
    parser.add_argument("--resourceGroupName", required=True)
    parser.add_argument("--clusterName", required=True)
    parser.add_argument("--location", required=True)
    parser.add_argument("--tenantId", required=True)
    parser.add_argument("--privatePemEnvironmentVariable", required=True)
    parser.add_argument("--kubeConfig", required=True)
    parser.add_argument("--customLocationsOid", required=True)
    args = parser.parse_args()

    helm_client = install_helm_client()
    chart_path = get_chart_path(get_helm_registry(), helm_client)
    subprocess.run(
        [
            helm_client,
            "upgrade",
            "--install",
            "azure-arc",
            chart_path,
            "--kubeconfig",
            args.kubeConfig,
            "--set",
            f"global.subscriptionId={args.subscriptionId}",
            "--set",
            "global.kubernetesDistro=aks",
            "--set",
            "global.kubernetesInfra=azure",
            "--set",
            f"global.resourceGroupName={args.resourceGroupName}",
            "--set",
            f"global.resourceName={args.clusterName}",
            "--set",
            f"global.location={args.location}",
            "--set",
            f"global.tenantId={args.tenantId}",
            "--set",
            f"global.onboardingPrivateKey={os.environ[args.privatePemEnvironmentVariable]}",
            "--set",
            "systemDefaultValues.spnOnboarding=false",
            "--set",
            "global.azureEnvironment=AZUREPUBLICCLOUD",
            "--set",
            "systemDefaultValues.clusterconnect-agent.enabled=true",
            "--set",
            "systemDefaultValues.customLocations.enabled=true",
            "--set",
            f"systemDefaultValues.customLocations.oid={args.customLocationsOid}",
            "--namespace",
            "azure-arc-release",
            "--create-namespace",
            "--output",
            "json",
            "--wait",
            "--timeout",
            "1800s",
        ],
        check=True,
    )


if __name__ == "__main__":
    install_agent()