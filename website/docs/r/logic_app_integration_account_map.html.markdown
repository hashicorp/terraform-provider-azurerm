---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logic_app_integration_account_map"
description: |-
  Manages a Logic App Integration Account Map.
---

# azurerm_logic_app_integration_account_map

Manages a Logic App Integration Account Map.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_logic_app_integration_account" "example" {
  name                = "example-ia"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Standard"
}

resource "azurerm_logic_app_integration_account_map" "example" {
  name                     = "example-iamap"
  resource_group_name      = azurerm_resource_group.example.name
  integration_account_name = azurerm_logic_app_integration_account.example.name
  map_type                 = "Xslt"
  content                  = file("testdata/integration_account_map_content.xsd")
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Logic App Integration Account Map. Changing this forces a new Logic App Integration Account Map to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Logic App Integration Account Map should exist. Changing this forces a new Logic App Integration Account Map to be created.

* `integration_account_name` - (Required) The name of the Logic App Integration Account. Changing this forces a new Logic App Integration Account Map to be created.

* `content` - (Required) The content of the Logic App Integration Account Map.

* `map_type` - (Required) The type of the Logic App Integration Account Map. Possible values are `Liquid`, `NotSpecified`, `Xslt`, `Xslt30` and `Xslt20`.

* `metadata` - (Optional) The metadata of the Logic App Integration Account Map.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Logic App Integration Account Map.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Logic App Integration Account Map.
* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App Integration Account Map.
* `update` - (Defaults to 30 minutes) Used when updating the Logic App Integration Account Map.
* `delete` - (Defaults to 30 minutes) Used when deleting the Logic App Integration Account Map.

## Import

Logic App Integration Account Maps can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logic_app_integration_account_map.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Logic/integrationAccounts/account1/maps/map1                                                                   
```
