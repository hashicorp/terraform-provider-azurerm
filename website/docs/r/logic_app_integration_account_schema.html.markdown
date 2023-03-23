---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logic_app_integration_account_schema"
description: |-
  Manages a Logic App Integration Account Schema.
---

# azurerm_logic_app_integration_account_schema

Manages a Logic App Integration Account Schema.

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
  sku_name            = "Basic"
}

resource "azurerm_logic_app_integration_account_schema" "example" {
  name                     = "example-ias"
  resource_group_name      = azurerm_resource_group.example.name
  integration_account_name = azurerm_logic_app_integration_account.example.name
  content                  = file("testdata/integration_account_schema_content.xsd")
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Logic App Integration Account Schema. Changing this forces a new Logic App Integration Account Schema to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Logic App Integration Account Schema should exist. Changing this forces a new Logic App Integration Account Schema to be created.

* `integration_account_name` - (Required) The name of the Logic App Integration Account. Changing this forces a new Logic App Integration Account Schema to be created.

* `content` - (Required) The content of the Logic App Integration Account Schema.

* `file_name` - (Optional) The file name of the Logic App Integration Account Schema.

* `metadata` - (Optional) The metadata of the Logic App Integration Account Schema.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Logic App Integration Account Schema.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Logic App Integration Account Schema.
* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App Integration Account Schema.
* `update` - (Defaults to 30 minutes) Used when updating the Logic App Integration Account Schema.
* `delete` - (Defaults to 30 minutes) Used when deleting the Logic App Integration Account Schema.

## Import

Logic App Integration Account Schemas can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logic_app_integration_account_schema.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Logic/integrationAccounts/account1/schemas/schema1
```
