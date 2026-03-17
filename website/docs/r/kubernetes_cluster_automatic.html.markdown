---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_cluster_automatic"
description: |-
  Manages an AKS Automatic Cluster (Azure Kubernetes Service Automatic)
---

# azurerm_kubernetes_cluster_automatic

Manages an AKS Automatic Cluster (Azure Kubernetes Service Automatic).

AKS Automatic is a managed Kubernetes experience that automates cluster setup, node management, scaling, security, and other operations. It provides a fast and frictionless Kubernetes experience while preserving flexibility and consistency. For more information, see [AKS Automatic documentation](https://learn.microsoft.com/azure/aks/intro-aks-automatic).

-> **Note:** AKS Automatic clusters are preconfigured with production-ready settings including Azure CNI Overlay powered by Cilium, Managed NGINX ingress, Azure RBAC for Kubernetes authorization, Workload Identity, and more.

~> **Note:** AKS Automatic clusters require Azure regions that support at least three availability zones and API Server VNet Integration.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_kubernetes_cluster_automatic" "example" {
  name                = "example-aks-automatic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Environment = "Production"
  }
}

output "kube_config" {
  value     = azurerm_kubernetes_cluster_automatic.example.kube_config_raw
  sensitive = true
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the AKS Automatic Cluster to create. Changing this forces a new resource to be created.

* `location` - (Required) The location where the AKS Automatic Cluster should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Resource Group where the AKS Automatic Cluster should exist. Changing this forces a new resource to be created.

---

* `identity` - (Required) An `identity` block as defined below.

* `kubernetes_version` - (Optional) Version of Kubernetes specified when creating the AKS Automatic cluster. If not specified, the latest recommended version will be used at provisioning time.

* `node_resource_group` - (Optional) The name of the Resource Group where the Kubernetes Nodes should exist. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this AKS Automatic Cluster. Possible values are `SystemAssigned` or `UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this AKS Automatic Cluster.

~> **Note:** This is required when `type` is set to `UserAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kubernetes Managed Cluster.

* `fqdn` - The FQDN of the AKS Automatic Cluster.

* `kube_config_raw` - Raw Kubernetes config to be used by [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) and other compatible tools. This is only available when Role Based Access Control with Microsoft Entra ID is enabled and local accounts are enabled.

* `portal_fqdn` - The FQDN for the Azure Portal to access the Managed Cluster. This is only visible from the Azure Portal.

---

The `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the AKS Automatic Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the AKS Automatic Cluster.
* `update` - (Defaults to 90 minutes) Used when updating the AKS Automatic Cluster.
* `delete` - (Defaults to 90 minutes) Used when deleting the AKS Automatic Cluster.

## Import

AKS Automatic Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_cluster_automatic.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ContainerService` - 2025-10-01
