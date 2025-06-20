---
subcategory: "Maps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_maps_account"
description: |-
  Gets information about an existing Azure Maps Account.
---

# Data Source: azurerm_maps_account

Use this data source to access information about an existing Azure Maps Account.

## Example Usage

```hcl
data "azurerm_maps_account" "example" {
  name                = "production"
  resource_group_name = "maps"
}

output "maps_account_id" {
  value = data.azurerm_maps_account.example.id
}
```

## Argument Reference

* `name` - Specifies the name of the Maps Account.

* `resource_group_name` - Specifies the name of the Resource Group in which the Maps Account is located.

## Attributes Reference

* `id` - The ID of the Maps Account.

* `sku_name` - The SKU of the Azure Maps Account.

* `primary_access_key` - The primary key used to authenticate and authorize access to the Maps REST APIs.

* `secondary_access_key` - The primary key used to authenticate and authorize access to the Maps REST APIs. The second key is given to provide seamless key regeneration.

* `x_ms_client_id` - A unique identifier for the Maps Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Maps Account.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Maps`: 2023-06-01
