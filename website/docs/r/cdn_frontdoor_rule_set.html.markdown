---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_ruleset"
description: |-
  Manages an Azure Front Door (Standard/Premium) instance. (currently in public preview)
---

# azurerm_cdn_frontdoor_rule_set

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

resource "azurerm_cdn_frontdoor_rule_set" "example" {
  name       = "sampleruleset"
  profile_id = azurerm_cdn_frontdoor_profile.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Front Door Rule Set service. Changing this forces a new resource to be created.

* `profile_id` - (Required) Azure Front Door Profile ID.