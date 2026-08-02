---
subcategory: "Quota"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_quota_group"
description: |-
  Gets information about an existing Azure Quota Group.
---

# Data Source: azurerm_quota_group

Gets information about an existing Azure Quota Group.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

data "azurerm_management_group" "example" {
  name = data.azurerm_client_config.current.tenant_id
}

data "azurerm_quota_group" "example" {
  name                = "example-quota-group"
  management_group_id = data.azurerm_management_group.example.id
}

output "quota_group_display_name" {
  value = data.azurerm_quota_group.example.display_name
}
```

## Example Usage with Quota Requests

Providing `location` returns all quota limit items for that `(resource_provider_name, location)` scope. This is useful with `for_each` to inspect quota limits across multiple regions.

```hcl
data "azurerm_client_config" "current" {}

data "azurerm_management_group" "example" {
  name = data.azurerm_client_config.current.tenant_id
}

data "azurerm_quota_group" "by_region" {
  for_each = toset(["eastus", "westus2"])

  name                   = "example-quota-group"
  management_group_id    = data.azurerm_management_group.example.id
  location               = each.value
  resource_provider_name = "Microsoft.Compute"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Quota Group. Must match `^[a-z][a-z0-9]*$` and be between 3 and 63 characters.

* `management_group_id` - (Required) The ID of the Management Group this Quota Group is scoped to.

* `location` - (Optional) The Azure region to retrieve quota limit items for (e.g. `eastus`). When provided, `quota_request` is populated with all quota limits for the given `(resource_provider_name, location)` scope. When omitted, `quota_request` will be empty.

* `resource_provider_name` - (Optional) The resource provider namespace to query quota limits for. Defaults to `Microsoft.Compute`. Only meaningful when `location` is also set. ~> **Note:** At this time only `Microsoft.Compute` is supported by the Azure Quota Groups API.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Quota Group.

* `display_name` - The human-readable display name of the Quota Group.

* `associated_subscription_ids` - A list of subscription IDs enrolled in this Quota Group.

* `quota_request` - A list of quota limit items for the given `location` and `resource_provider_name` scope. Empty when `location` is not specified.

The `quota_request` block exports the following:

* `resource_name` - The name of the quota resource (e.g. `standardDDv4Family`).

* `location` - The Azure region this quota limit applies to.

* `resource_provider_name` - The resource provider namespace.

* `limit` - The approved quota limit for the resource in the given location.

* `comment` - A comment associated with the quota request.

* `available_limit` - The portion of the group's quota pool that has not yet been allocated to subscriptions.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Quota Group.
