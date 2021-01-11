---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_logic_app_integration_account"
description: |-
  Gets information about an existing Logic App Integration Account.
---

# Data Source: azurerm_logic_app_integration_account

Use this data source to access information about an existing Logic App Integration Account.

## Example Usage

```hcl
data "azurerm_logic_app_integration_account" "example" {
  name                = "example-account"
  resource_group_name = "example-resource-group"
}

output "id" {
  value = data.azurerm_logic_app_integration_account.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Logic App Integration Account.

* `resource_group_name` - (Required) The name of the Resource Group where the Logic App Integration Account exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Logic App Integration Account.

* `location` - The Azure Region where the Logic App Integration Account exists.

* `sku_name` - The sku name of the Logic App Integration Account.

* `tags` - A mapping of tags assigned to the Logic App Integration Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App Integration Account.
