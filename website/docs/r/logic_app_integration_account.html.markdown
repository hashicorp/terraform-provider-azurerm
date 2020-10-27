---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logic_app_integration_account"
description: |-
  Manages a Logic App Integration Account.
---

# azurerm_logic_app_integration_account

Manages a Logic App Integration Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_logic_app_integration_account" "example" {
  name                = "example-ia"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "Standard"
  tags = {
    foo = "bar"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Logic App Integration Account. Changing this forces a new Logic App Integration Account to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Logic App Integration Account should exist. Changing this forces a new Logic App Integration Account to be created.

* `location` - (Required) The Azure Region where the Logic App Integration Account should exist. Changing this forces a new Logic App Integration Account to be created.

* `sku_name` - (Required) The sku name of the Logic App Integration Account. Possible Values are `Basic`, `Free` and `Standard`.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Logic App Integration Account.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Logic App Integration Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Logic App Integration Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App Integration Account.
* `update` - (Defaults to 30 minutes) Used when updating the Logic App Integration Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Logic App Integration Account.

## Import

Logic App Integration Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logic_app_integration_account.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Logic/integrationAccounts/account1
```
