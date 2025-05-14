---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_cluster"
description: |-
  Manages a Log Analytics Cluster.
---

# azurerm_log_analytics_cluster

~> **Note:** Log Analytics Clusters are subject to 14-day soft delete policy. Clusters created with the same resource group & name as a previously deleted cluster will be recovered rather than creating anew.

Manages a Log Analytics Cluster.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_cluster" "example" {
  name                = "example-cluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  identity {
    type = "SystemAssigned"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Log Analytics Cluster. Changing this forces a new Log Analytics Cluster to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Log Analytics Cluster should exist. Changing this forces a new Log Analytics Cluster to be created.

* `location` - (Required) The Azure Region where the Log Analytics Cluster should exist. Changing this forces a new Log Analytics Cluster to be created.

* `identity` - (Required) An `identity` block as defined below. Changing this forces a new Log Analytics Cluster to be created.

* `size_gb` - (Optional) The capacity of the Log Analytics Cluster is specified in GB/day. Possible values include `100`, `200`, `300`, `400`, `500`, `1000`, `2000`, `5000`, `10000`, `25000`, or `50000`. Defaults to `100`.

~> **Note:** The cluster capacity must start at 100 GB and can be set to 500, 1000, 2000 or 5000 GB/day. For more information on cluster costs, see [Dedicated clusters](https://docs.microsoft.com/en-us/azure/azure-monitor/logs/cost-logs#dedicated-clusters).

* `tags` - (Optional) A mapping of tags which should be assigned to the Log Analytics Cluster.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Log Analytics Cluster. Possible values are `SystemAssigned` and  `UserAssigned`. Changing this forces a new resource to be created.

~> **Note:** The assigned `principal_id` and `tenant_id` can be retrieved after the identity `type` has been set to `SystemAssigned` and the Log Analytics Cluster has been created. More details are available below.

* `identity_ids` - (Optional) A list of User Assigned Managed Identity IDs to be assigned to this Windows Web App Slot.

~> **Note:** This is required when `type` is set to `UserAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Log Analytics Cluster.

* `identity` - A `identity` block as defined below.

* `cluster_id` - The GUID of the cluster.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

* `type` - The identity type of this Managed Service Identity.

-> **Note:** You can access the Principal ID via `azurerm_log_analytics_cluster.example.identity[0].principal_id` and the Tenant ID via `azurerm_log_analytics_cluster.example.identity[0].tenant_id`

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 6 hours) Used when creating the Log Analytics Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Cluster.
* `update` - (Defaults to 6 hours) Used when updating the Log Analytics Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Cluster.

## Import

Log Analytics Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/clusters/cluster1
```
