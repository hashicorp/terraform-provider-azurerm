---
subcategory: "Planetary Computer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_geo_catalog"
description: |-
    Lists Microsoft Planetary Computer Pro GeoCatalog resources.
---

# List resource: azurerm_geo_catalog

Lists Microsoft Planetary Computer Pro GeoCatalog resources.

## Example Usage

### List all Microsoft Planetary Computer Pro GeoCatalogs in the subscription

```hcl
list "azurerm_geo_catalog" "example" {
  provider = azurerm
  config {
  }
}
```

### List all Microsoft Planetary Computer Pro GeoCatalogs in a Resource Group

```hcl
list "azurerm_geo_catalog" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `subscription_id` - (Optional) The ID of the Subscription to query. Defaults to the value specified in the Provider Configuration.

* `resource_group_name` - (Optional) The name of the Resource Group to query.
