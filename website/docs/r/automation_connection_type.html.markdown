---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_connection_type"
description: |-
  Manages an Automation Connection Type.
---

# azurerm_automation_connection_type

Manages anAutomation Connection Type.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "resourceGroup-example"
  location = "West Europe"
}

data "azurerm_client_config" "example" {}

resource "azurerm_automation_account" "example" {
  name                = "account-example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Basic"
}

resource "azurerm_automation_connection_type" "example" {
  name                    = "example"
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name

  field {
    name = "example"
    type = "string"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Automation Connection Type. Changing this forces a new Automation to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Automation should exist. Changing this forces a new Automation to be created.

* `automation_account_name` - (Required) The name of the automation account in which the Connection is created. Changing this forces a new resource to be created.

* `field` - (Required) One or more `field` blocks as defined below. Changing this forces a new Automation to be created.

---

* `is_global` - (Optional) Whether the connection type is global. Changing this forces a new Automation to be created.

---

A `field` block supports the following:

* `name` - (Required) The name which should be used for this connection field definition.

* `type` - (Required) The type of the connection field definition.

* `is_encrypted` - (Optional) Whether to set the isEncrypted flag of the connection field definition.

* `is_optional` - (Optional) Whether to set the isOptional flag of the connection field definition.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The the Automation Connection Type ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation.
* `delete` - (Defaults to 10 minutes) Used when deleting the Automation.

## Import

Automations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_connection_type.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/connectionTypes/type1
```
