---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_resource_groups"
description: |-
  Get information about the available resource groups.
---

# Data Source: azurerm_resource_groups

Use this data source to access information about an existing Resource Group.

## Example Usage

```hcl
data "azurerm_resource_groups" "example" {

}

output "id" {
  value = data.azurerm_resource_groups.example.id
}
```

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Resource Group.

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

* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Group.