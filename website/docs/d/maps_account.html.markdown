---
subcategory: "Maps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_maps_account"
sidebar_current: "docs-azurerm-datasource-maps-account"
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
  value = "${data.azurerm_maps_account.example.id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Maps Account.

* `resource_group_name` - (Required) Specifies the name of the Resource Group in which the Maps Account is located.

## Attributes Reference

* `id` - The ID of the Maps Account.

* `sku_name` - The sku of the Azure Maps Account.

* `primary_access_key` - The primary key used to authenticate and authorize access to the Maps REST APIs.

* `secondary_access_key` - The primary key used to authenticate and authorize access to the Maps REST APIs. The second key is given to provide seamless key regeneration.

* `x_ms_client_id` - A unique identifier for the Maps Account.

