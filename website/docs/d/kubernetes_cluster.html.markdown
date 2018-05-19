---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_cluster"
sidebar_current: "docs-azurerm-data-source-kubernetes-cluster"
description: |-
  Gets information about a managed Kubernetes Cluster (AKS)
---

# Data Source: azurerm_kubernetes_cluster

Gets information about a managed Kubernetes Cluster (AKS)

~> **Note:** All arguments including the client secret will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).


## Example Usage

```hcl
data "azurerm_kubernetes_cluster" "test" {
  name                = "myakscluster"
  resource_group_name = "my-example-resource-group"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the managed Kubernetes Cluster.

* `resource_group_name` - (Required) The name of the Resource Group in which the managed Kubernetes Cluster exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Kubernetes Managed Cluster ID.

* `fqdn` - The FQDN of the Azure Kubernetes Managed Cluster.

* `kube_config_raw` - Base64 encoded Kubernetes configuration.

* `kube_config` - A `kube_config` block as defined below.

* `location` - The Azure Region in which the managed Kubernetes Cluster exists.

* `dns_prefix` - The DNS Prefix of the managed Kubernetes cluster.

* `kubernetes_version` - The version of Kubernetes used on the managed Kubernetes Cluster.

* `linux_profile` - A `linux_profile` block as documented below.

* `agent_pool_profile` - One or more `agent_profile_pool` blocks as documented below.

* `service_principal` - A `service_principal` block as documented below.

* `tags` - A mapping of tags assigned to this resource.

---

`kube_config` exports the following:

* `client_key` - Base64 encoded private key used by clients to authenticate to the Kubernetes cluster.

* `client_certificate` - Base64 encoded public certificate used by clients to authenticate to the Kubernetes cluster.

* `cluster_ca_certificate` - Base64 encoded public CA certificate used as the root of trust for the Kubernetes cluster.

* `host` - The Kubernetes cluster server host.

* `username` - A username used to authenticate to the Kubernetes cluster.

* `password` - A password or token used to authenticate to the Kubernetes cluster.

-> **NOTE:** It's possible to use these credentials with [the Kubernetes Provider](/docs/providers/kubernetes/index.html) like so:

```
provider "kubernetes" {
  host                   = "${data.azurerm_kubernetes_cluster.main.kube_config.0.host}"
  username               = "${data.azurerm_kubernetes_cluster.main.kube_config.0.username}"
  password               = "${data.azurerm_kubernetes_cluster.main.kube_config.0.password}"
  client_certificate     = "${base64decode(data.azurerm_kubernetes_cluster.main.kube_config.0.client_certificate)}"
  client_key             = "${base64decode(data.azurerm_kubernetes_cluster.main.kube_config.0.client_key)}"
  cluster_ca_certificate = "${base64decode(data.azurerm_kubernetes_cluster.main.kube_config.0.cluster_ca_certificate)}"
}
```

`linux_profile` exports the following:

* `admin_username` - The username associated with the administrator account of the managed Kubernetes Cluster.
* `ssh_key` - One or more `ssh_key` blocks as defined below.

`ssh_key` exports the following:

* `key_data` - The Public SSH Key used to access the cluster.

`agent_pool_profile` exports the following:

* `name` - The name assigned to this pool of agents
* `count` - The number of Agents (VM's) in the Pool.
* `vm_size` - The size of each VM in the Agent Pool (e.g. `Standard_F1`).
* `os_disk_size_gb` - The size of the Agent VM's Operating System Disk in GB.
* `os_type` - The Operating System used for the Agents.
* `vnet_subnet_id` - The ID of the Subnet where the Agents in the Pool are provisioned.

`service_principal` supports the following:

* `client_id` - The Client ID of the Service Principal used by this Managed Kubernetes Cluster.
