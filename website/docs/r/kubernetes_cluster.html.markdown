---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_cluster"
sidebar_current: "docs-azurerm-resource-container-kubernetes-cluster"
description: |-
  Manages a managed Kubernetes Cluster (AKS)
---

# azurerm_kubernetes_cluster

Manages a managed Kubernetes Cluster (AKS)

~> **Note:** All arguments including the client secret will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).


## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acctestRG1"
  location = "East US"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  dns_prefix          = "acctestagent1"

  linux_profile {
    admin_username = "acctestuser1"

    ssh_key {
      key_data = "ssh-rsa ..."
    }
  }

  agent_pool_profile {
    name            = "default"
    count           = 1
    vm_size         = "Standard_D1_v2"
    os_type         = "Linux"
    os_disk_size_gb = 30
  }

  service_principal {
    client_id     = "00000000-0000-0000-0000-000000000000"
    client_secret = "00000000000000000000000000000000"
  }

  tags {
    Environment = "Production"
  }
}

output "id" {
    value = "${azurerm_kubernetes_cluster.test.id}"
}

output "kube_config" {
  value = "${azurerm_kubernetes_cluster.test.kube_config_raw}"
}

output "client_key" {
  value = "${azurerm_kubernetes_cluster.test.kube_config.0.client_key}"
}

output "client_certificate" {
  value = "${azurerm_kubernetes_cluster.test.kube_config.0.client_certificate}"
}

output "cluster_ca_certificate" {
  value = "${azurerm_kubernetes_cluster.test.kube_config.0.cluster_ca_certificate}"
}

output "host" {
  value = "${azurerm_kubernetes_cluster.test.kube_config.0.host}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the AKS Managed Cluster instance to create. Changing this forces a new resource to be created.

* `location` - (Required) The location where the AKS Managed Cluster instance should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the resource group where the resource exists. Changing this forces a new resource to be created.

* `dns_prefix` - (Required) DNS prefix specified when creating the managed cluster.

* `kubernetes_version` - (Optional) Version of Kubernetes specified when creating the AKS managed cluster. If not specified, the latest recommended version will be used at provisioning time (but won't auto-upgrade).

* `linux_profile` - (Required) A Linux Profile block as documented below.

* `agent_pool_profile` - (Required) One or more Agent Pool Profile's block as documented below.

* `service_principal` - (Required) A Service Principal block as documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

`linux_profile` supports the following:

* `admin_username` - (Required) The Admin Username for the Cluster. Changing this forces a new resource to be created.
* `ssh_key` - (Required) An SSH Key block as documented below.

`ssh_key` supports the following:

* `key_data` - (Required) The Public SSH Key used to access the cluster. Changing this forces a new resource to be created.

`agent_pool_profile` supports the following:

* `name` - (Required) Unique name of the Agent Pool Profile in the context of the Subscription and Resource Group. Changing this forces a new resource to be created.
* `count` - (Required) Number of Agents (VMs) in the Pool. Possible values must be in the range of 1 to 50 (inclusive). Defaults to `1`.
* `vm_size` - (Required) The size of each VM in the Agent Pool (e.g. `Standard_F1`). Changing this forces a new resource to be created.
* `os_disk_size_gb` - (Optional) The Agent Operating System disk size in GB. Changing this forces a new resource to be created.
* `os_type` - (Optional) The Operating System used for the Agents. Possible values are `Linux` and `Windows`.  Changing this forces a new resource to be created. Defaults to `Linux`.
* `vnet_subnet_id` - (Optional) The ID of the Subnet where the Agents in the Pool should be provisioned. Changing this forces a new resource to be created.

~> **NOTE:** There's a known issue where Agents connected to an Internal Network (e.g. on a Subnet) have their network routing configured incorrectly; such that Pods cannot communicate across nodes. This is a bug in the Azure API - and will be fixed there in the future.

`service_principal` supports the following:

* `client_id` - (Required) The Client ID for the Service Principal.
* `client_secret` - (Required) The Client Secret for the Service Principal.

## Attributes Reference

The following attributes are exported:

* `id` - The Kubernetes Managed Cluster ID.

* `fqdn` - The FQDN of the Azure Kubernetes Managed Cluster.

* `kube_config_raw` - Base64 encoded Kubernetes configuration

* `kube_config` - Kubernetes configuration, sub-attributes defined below:

  * `client_key` - Base64 encoded private key used by clients to authenticate to the Kubernetes cluster.

  * `client_certificate` - Base64 encoded public certificate used by clients to authenticate to the Kubernetes cluster.

  * `cluster_ca_certificate` - Base64 encoded public CA certificate used as the root of trust for the Kubernetes cluster.

  * `host` - The Kubernetes cluster server host.

  * `username` - A username used to authenticate to the Kubernetes cluster.

  * `password` - A password or token used to authenticate to the Kubernetes cluster.

-> **NOTE:** It's possible to use these credentials with [the Kubernetes Provider](/docs/providers/kubernetes/index.html) like so:

```
provider "kubernetes" {
  host                   = "${azurerm_kubernetes_cluster.main.kube_config.0.host}"
  username               = "${azurerm_kubernetes_cluster.main.kube_config.0.username}"
  password               = "${azurerm_kubernetes_cluster.main.kube_config.0.password}"
  client_certificate     = "${base64decode(azurerm_kubernetes_cluster.main.kube_config.0.client_certificate)}"
  client_key             = "${base64decode(azurerm_kubernetes_cluster.main.kube_config.0.client_key)}"
  cluster_ca_certificate = "${base64decode(azurerm_kubernetes_cluster.main.kube_config.0.cluster_ca_certificate)}"
}
```

## Import

Kubernetes Managed Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_cluster.cluster1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1
```
