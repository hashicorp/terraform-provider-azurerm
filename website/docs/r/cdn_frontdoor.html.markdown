---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor"
description: |-
  Manages an Azure Front Door (Standard/Premium) instance. (currently in public preview)
---

# azurerm_cdn_frontdoor_profile

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "afdpremv2"
  location            = "global"
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Premium_AzureFrontDoor"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Front Door service. Changing this forces a new resource to be created.

* `location` - (Required) Must be set to `global`.

* `resource_group_name` - (Required) Specifies the name of the Resource Group in which the Front Door service should exist. Changing this forces a new resource to be created.

* `sku` - (Required) Can be either `Premium_AzureFrontDoor` or `Standard_AzureFrontDoor`. Changing this forces a new resource to be created.