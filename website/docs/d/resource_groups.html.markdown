---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_resource_groups"
description: |-
  Gets information about existing Resource Groups.
---

# Data Source: azurerm_resource_groups

Use this data source to access information about existing Resource Groups.

## Example Usage

```hcl
data "azurerm_client_config" "current" {
}

data "azurerm_resource_groups" "test" {
  subscription_ids = [data.azurerm_client_config.current.subscription_id]
}

output "id" {
  value = data.azurerm_resource_groups.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `subscription_ids` - (Required) Specifies a list of Subscription IDs.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Resource Groups.

* `resource_groups` - A `resource_groups` block as defined below.

---

A `resource_groups` block exports the following:

* `id` - The id of the resource group.

* `location` - The Azure Region where the Resource Group exists.

* `name` - The name of this resource group.

* `subscription_id` - The ID of the subscription the resource group resides in.

* `tags` - A mapping of tags assigned to the Resource Group.

* `tenant_id` - The ID of the tenant.

* `type` - The Microsoft resource type of the resource group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Groups.
