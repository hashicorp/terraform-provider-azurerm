---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_account"
description: |-
  Manages a Automation Account.
---

# azurerm_automation_account

Manages a Automation Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_automation_account" "example" {
  name                = "automationAccount1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku_name = "Basic"

  tags = {
    environment = "development"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Automation Account. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Automation Account is created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Optional **Deprecated**)) A `sku` block as described below.

* `sku_name` - (Optional) The SKU name of the account - only `Basic` is supported at this time.

* `tags` - (Optional) A mapping of tags to assign to the resource.

----

A `sku` block supports the following:

* `name` - (Required) The SKU name of the account - only `Basic` is supported at this time.

----

## Attributes Reference

The following attributes are exported:

* `id` - The Automation Account ID.

* `dsc_server_endpoint` - The DSC Server Endpoint associated with this Automation Account.

* `dsc_primary_access_key` - The Primary Access Key for the DSC Endpoint associated with this Automation Account.

* `dsc_secondary_access_key` - The Secondary Access Key for the DSC Endpoint associated with this Automation Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation Account.
* `update` - (Defaults to 30 minutes) Used when updating the Automation Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Automation Account.

## Import

Automation Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_account.account1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1
```
