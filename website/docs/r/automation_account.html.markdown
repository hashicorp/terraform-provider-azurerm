---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_account"
sidebar_current: "docs-azurerm-resource-automation-account"
description: |-
  Creates a new Automation Account.
---

# azurerm\_automation\_account

Creates a new Automation Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
 name = "resourceGroup1"
 location = "West Europe"
}

resource "azurerm_automation_account" "example" {
  name                = "automationAccount1"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  sku {
    name = "Basic"
  }

  tags {
    environment = "development"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Automation Account. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Automation Account is created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

`sku` supports the following:

* `name` - (Optional) The SKU name of the account - only `Basic` is supported at this time. Defaults to `Basic`.

## Attributes Reference

The following attributes are exported:

* `id` - The Automation Account ID.

## Import

Automation Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_account.account1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1
```
