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

* `enable_rbac` - Whether Role Based Access Control is currently enabled.

* `kube_config_raw` - Base64 encoded Kubernetes configuration.

* `node_resource_group` - Auto-generated Resource Group containing AKS Cluster resources.

* `kube_config` - A `kube_config` block as defined below.

* `location` - The Azure Region in which the managed Kubernetes Cluster exists.

* `dns_prefix` - The DNS Prefix of the managed Kubernetes cluster.

* `kubernetes_version` - The version of Kubernetes used on the managed Kubernetes Cluster.

* `linux_profile` - A `linux_profile` block as documented below.

* `agent_pool_profile` - One or more `agent_profile_pool` blocks as documented below.

* `addon_profile` - A `addon_profile` block as documented below.

* `service_principal` - A `service_principal` block as documented below.

* `network_profile` - A `network_profile` block as documented below.

* `aad_profile` - If AzureAD integration with RBAC is in use, a `aad_profile` block as documented below.

* `tags` - A mapping of tags assigned to this resource.

---

A `addon_profile` block exports the following:

* `http_application_routing` - A `http_application_routing` block.
* `oms_agent` - A `oms_agent` block.

---

A `agent_pool_profile` block exports the following:

* `name` - The name assigned to this pool of agents
* `count` - The number of Agents (VM's) in the Pool.
* `vm_size` - The size of each VM in the Agent Pool (e.g. `Standard_F1`).
* `os_disk_size_gb` - The size of the Agent VM's Operating System Disk in GB.
* `os_type` - The Operating System used for the Agents.
* `vnet_subnet_id` - The ID of the Subnet where the Agents in the Pool are provisioned.
* `max_pods` - The maximum number of pods that can run on each agent.

---

A `http_application_routing` block exports the following:

* `enabled` - Is HTTP Application Routing Enabled?

* `http_application_routing_zone_name` - The Zone Name of the HTTP Application Routing.

---

A `kube_config` block exports the following:

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

---

A `linux_profile` block exports the following:

* `admin_username` - The username associated with the administrator account of the managed Kubernetes Cluster.
* `ssh_key` - One or more `ssh_key` blocks as defined below.

---

A `network_profile` block exports the following:

* `network_plugin` - Network plugin used such as `azure` or `kubenet`.
* `service_cidr` - Network range used by the Kubernetes service.
* `dns_service_ip` - IP address within the Kubernetes service address range used by cluster service discovery (kube-dns).
* `docker_bridge_cidr` - IP address (in CIDR notation) used as the Docker bridge IP address on nodes.
* `pod_cidr` - The CIDR used for pod IP addresses.

A `oms_agent` block exports the following:

* `enabled` - Is the OMS Agent Enabled?

* `log_analytics_workspace_id` - The ID of the Log Analytics Workspace which the OMS Agent should send data to.

---

A `aad_profile` block exports the following:

* `server_app_id` - AzureAD Server Application ID.

* `client_id` - AzureAD Client Application ID.

* `tenant_id` - AzureAD Tenant ID.

---

A `service_principal` block supports the following:

* `client_id` - The Client ID of the Service Principal used by this Managed Kubernetes Cluster.

---

A `ssh_key` block exports the following:

* `key_data` - The Public SSH Key used to access the cluster.
