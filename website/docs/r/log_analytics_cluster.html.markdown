---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_cluster"
description: |-
  Manages a Log Analytics Cluster.
---

# azurerm_log_analytics_cluster

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
  name = "example-cluster"
  resource_group_name = azurerm_resource_group.example.name
  location = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Log Analytics Cluster. Changing this forces a new Log Analytics Cluster to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Log Analytics Cluster should exist. Changing this forces a new Log Analytics Cluster to be created.

* `location` - (Required) The Azure Region where the Log Analytics Cluster should exist. Changing this forces a new Log Analytics Cluster to be created.

---

* `identity` - (Optional)  A `identity` block as defined below. Changing this forces a new Log Analytics Cluster to be created.

* `next_link` - (Optional) The link used to get the next page of recommendations. Changing this forces a new Log Analytics Cluster to be created.

* `key_vault_property` - (Optional)  A `key_vault_property` block as defined below.

* `sku` - (Optional)  A `sku` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Log Analytics Cluster.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the Log Analytics Cluster. At this time the only allowed value is `SystemAssigned`.

~> **NOTE:** The assigned `principal_id` and `tenant_id` can be retrieved after the identity `type` has been set to `SystemAssigned` and the Log Analytics Cluster has been created. More details are available below.

---

An `key_vault_property` block exports the following:

* `key_name` - (Optional) The name of the key associated with the Log Analytics cluster.

* `key_vault_uri` - (Optional) The Key Vault uri which holds they key associated with the Log Analytics cluster.

* `key_version` - (Optional) The version of the key associated with the Log Analytics cluster.

---

An `sku` block exports the following:

* `name` - (Optional) The name which should be used for this sku. Possible value is "CapacityReservation" is allowed.

* `capacity` - (Optional) The capacity value.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Log Analytics Cluster.

* `identity` - A `identity` block as defined below.

* `cluster_id` - The ID of the cluster.

* `type` - The type of the resource. Ex- Microsoft.Compute/virtualMachines or Microsoft.Storage/storageAccounts.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this Log Analytics Cluster.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this Log Analytics Cluster.

-> You can access the Principal ID via `azurerm_log_analytics_cluster.example.identity.0.principal_id` and the Tenant ID via `azurerm_log_analytics_cluster.example.identity.0.tenant_id`

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Log Analytics Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the Log Analytics Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Cluster.

## Import

Log Analytics Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/group1/providers/Microsoft.OperationalInsights/clusters/cluster1
```