---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logic_app_integration_account_session"
description: |-
  Manages a Logic App Integration Account Session.
---

# azurerm_logic_app_integration_account_session

Manages a Logic App Integration Account Session.

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

resource "azurerm_logic_app_integration_account_session" "example" {
  name                     = "example-ias"
  resource_group_name      = azurerm_resource_group.example.name
  integration_account_name = azurerm_logic_app_integration_account.example.name

  content = <<CONTENT
	{
       "controlNumber": "1234"
    }
  CONTENT
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Logic App Integration Account Session. Changing this forces a new Logic App Integration Account Session to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Logic App Integration Account Session should exist. Changing this forces a new Logic App Integration Account Session to be created.

* `integration_account_name` - (Required) The name of the Logic App Integration Account. Changing this forces a new Logic App Integration Account Session to be created.

* `content` - (Optional) The content of the Logic App Integration Account Session.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Logic App Integration Account Session.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Logic App Integration Account Session.
* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App Integration Account Session.
* `update` - (Defaults to 30 minutes) Used when updating the Logic App Integration Account Session.
* `delete` - (Defaults to 30 minutes) Used when deleting the Logic App Integration Account Session.

## Import

Logic App Integration Account Sessions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logic_app_integration_account_session.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Logic/integrationAccounts/account1/sessions/session1
```
