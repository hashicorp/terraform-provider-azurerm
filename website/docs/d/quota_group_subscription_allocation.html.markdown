---
subcategory: "Quota"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_quota_group_subscription_allocation"
description: |-
  Gets information about existing quota allocations for a subscription within an Azure Quota Group.
---

# Data Source: azurerm_quota_group_subscription_allocation

Gets information about existing quota allocations for a subscription within an Azure Quota Group.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

data "azurerm_management_group" "example" {
  name = data.azurerm_client_config.current.tenant_id
}

data "azurerm_quota_group_subscription_allocation" "example" {
  quota_group_id  = "/providers/Microsoft.Management/managementGroups/${data.azurerm_management_group.example.name}/providers/Microsoft.Quota/groupQuotas/example-quota-group"
  subscription_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  location        = "eastus"
}

output "allocations" {
  value = data.azurerm_quota_group_subscription_allocation.example.allocation
}
```

## Arguments Reference

The following arguments are supported:

* `quota_group_id` - (Required) The ID of the Quota Group.

* `subscription_id` - (Required) The ID of the subscription whose allocations are to be read.

* `location` - (Required) The Azure region to read quota allocations for (e.g. `eastus`).

* `resource_provider_name` - (Optional) The resource provider namespace. Defaults to `Microsoft.Compute`. ~> **Note:** At this time only `Microsoft.Compute` is supported by the Azure Quota Groups API.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Quota Group Subscription Allocation.

* `allocation` - A list of quota allocations for the subscription in the given location.

The `allocation` block exports the following:

* `resource_name` - The name of the quota resource (e.g. `standardDDv4Family`).

* `limit` - The quota limit allocated to this subscription from the group pool.

* `shareable_quota` - The portion of this subscription's allocated quota that can be returned to the group pool (computed as `limit` minus current resource usage). This is a read-only value set by the API.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Quota Group Subscription Allocation.
