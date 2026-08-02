---
subcategory: "Quota"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_quota_group_subscription_allocation"
description: |-
  Manages quota allocations from an Azure Quota Group to a specific subscription.
---

# azurerm_quota_group_subscription_allocation

Manages quota allocations from an Azure Quota Group to a specific subscription. Once a Quota Group holds a pooled limit (via `quota_request` in `azurerm_quota_group`), this resource carves out a share of that pool and assigns it to an individual subscription for a given location and resource provider.

~> **Note:** The subscription must already be associated with the Quota Group (via `associated_subscription_ids` on `azurerm_quota_group`) before quota can be allocated to it.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

data "azurerm_management_group" "example" {
  name = data.azurerm_client_config.current.tenant_id
}

resource "azurerm_quota_group" "example" {
  name                = "example-quota-group"
  management_group_id = data.azurerm_management_group.example.id

  associated_subscription_ids = [data.azurerm_client_config.current.subscription_id]

  quota_request {
    resource_name = "standardDDv4Family"
    location      = "eastus"
    limit         = 500
  }
}

resource "azurerm_quota_group_subscription_allocation" "example" {
  quota_group_id  = azurerm_quota_group.example.id
  subscription_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  location        = "eastus"

  allocation {
    resource_name = "standardDDv4Family"
    limit         = 200
  }
}
```

## Arguments Reference

The following arguments are supported:

* `quota_group_id` - (Required) The ID of the Quota Group. Changing this forces a new resource to be created.

* `subscription_id` - (Required) The ID of the subscription to allocate quota to (e.g. `/subscriptions/00000000-0000-0000-0000-000000000000`). Changing this forces a new resource to be created.

* `location` - (Required) The Azure region for which quota is being allocated (e.g. `eastus`). Changing this forces a new resource to be created.

* `allocation` - (Required) One or more `allocation` blocks as defined below.

---

* `resource_provider_name` - (Optional) The resource provider namespace. Defaults to `Microsoft.Compute`. Changing this forces a new resource to be created. ~> **Note:** At this time only `Microsoft.Compute` is supported by the Azure Quota Groups API.

---

An `allocation` block supports the following:

* `resource_name` - (Required) The name of the quota resource to allocate (e.g. `standardDDv4Family`).

* `limit` - (Required) The quota limit to allocate to the subscription from the group pool. ~> **Note:** Azure may silently cap the allocated amount to the group's `availableLimit` if the requested `limit` exceeds the group's remaining pool. The Terraform state will reflect what Azure actually allocated.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Quota Group Subscription Allocation.

The `allocation` block exports the following:

* `shareable_quota` - The portion of this subscription's allocated quota that can be returned to the group pool (computed as `limit` minus current resource usage). This is a read-only value set by the API.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Quota Group Subscription Allocation.
* `read` - (Defaults to 5 minutes) Used when retrieving the Quota Group Subscription Allocation.
* `update` - (Defaults to 30 minutes) Used when updating the Quota Group Subscription Allocation.
* `delete` - (Defaults to 30 minutes) Used when deleting the Quota Group Subscription Allocation (zeroes all limits to return quota to the group pool).

## Import

Quota Group Subscription Allocations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_quota_group_subscription_allocation.example \
  /providers/Microsoft.Management/managementGroups/example-mg/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Quota/groupQuotas/example-quota-group/resourceProviders/Microsoft.Compute/quotaAllocations/eastus
```
