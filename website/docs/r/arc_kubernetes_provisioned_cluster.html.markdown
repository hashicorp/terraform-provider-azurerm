---
subcategory: "ArcKubernetes"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_arc_kubernetes_provisioned_cluster"
description: |-
  Manages an Arc Kubernetes Provisioned Cluster.
---

# azurerm_arc_kubernetes_provisioned_cluster

Manages an Arc Kubernetes Provisioned Cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_client_config" "current" {}

resource "azuread_group" "example" {
  display_name     = "example-adg"
  owners           = [data.azurerm_client_config.current.object_id]
  security_enabled = true
}

resource "azurerm_arc_kubernetes_provisioned_cluster" "example" {
  name                = "example-akpc"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  azure_active_directory {
    azure_rbac_enabled     = true
    admin_group_object_ids = [azuread_group.example.id]
    tenant_id              = data.azurerm_client_config.current.tenant_id
  }

  identity {
    type = "SystemAssigned"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Arc Kubernetes Provisioned Cluster. Changing this forces a new Arc Kubernetes Provisioned Cluster to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Arc Kubernetes Provisioned Cluster should exist. Changing this forces a new Arc Kubernetes Provisioned Cluster to be created.

* `identity` - (Required) An `identity` block as defined below. Changing this forces a new Arc Kubernetes Provisioned Cluster to be created.

* `location` - (Required) The Azure Region where the Arc Kubernetes Provisioned Cluster should exist. Changing this forces a new Arc Kubernetes Provisioned Cluster to be created.

---

* `arc_agent_auto_upgrade_enabled` - (Optional) Whether the Arc agents will be upgraded automatically to the latest version. Defaults to `true`.

* `arc_agent_desired_version` - (Optional) The version of the Arc agents to be installed on the cluster.

* `azure_active_directory` - (Optional) An `azure_active_directory` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Arc Kubernetes Provisioned Cluster.

---

An `azure_active_directory` block supports the following:

* `admin_group_object_ids` - (Optional) A list of IDs of Microsoft Entra ID Groups. All members of the specified Microsoft Entra ID Groups have the cluster administrator access to the Kubernetes cluster.

* `azure_rbac_enabled` - (Optional) Whether to enable Azure RBAC for Kubernetes authorization. Defaults to `false`.

* `tenant_id` - (Optional) The Tenant ID to use for authentication. If not specified, the Tenant of the Arc Kubernetes Cluster will be used.

---

An `identity` block supports the following:

* `type` - (Required) The type of the Managed Identity. The only possible value is `SystemAssigned`. Changing this forces a new Arc Kubernetes Provisioned Cluster to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Arc Kubernetes Provisioned Cluster.

* `agent_version` - The version of the agent running on the cluster resource.

* `distribution` - The distribution running on this Arc Kubernetes Provisioned Cluster.

* `identity` - An `identity` block as defined below.

* `infrastructure` - The infrastructure on which the Arc Kubernetes Provisioned Cluster is running on.

* `kubernetes_version` - The Kubernetes version of the cluster resource.

* `offering` - The cluster offering.

* `total_core_count` - The number of CPU cores present in the cluster resource.

* `total_node_count` - The number of nodes present in the cluster resource.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Arc Kubernetes Provisioned Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Arc Kubernetes Provisioned Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the Arc Kubernetes Provisioned Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Arc Kubernetes Provisioned Cluster.

## Import

Arc Kubernetes Provisioned Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_arc_kubernetes_provisioned_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Kubernetes/connectedClusters/cluster1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Kubernetes`: 2024-01-01
