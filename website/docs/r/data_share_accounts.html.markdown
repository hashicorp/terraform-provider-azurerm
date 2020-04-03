---
subcategory: "DataShare"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_share_account"
description: |-
  Manages a DataShare Account
---

# azurerm_data_share_account

Manages a DataShare Account

~> **NOTE:** Azure allows only one active directory can be joined to a single subscription at a time for DataShare Account.

## DataShare Account Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_share_account" "example" {
  name = "example-dsa"
  location = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tags = {
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DataShare Account. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the DataShare Account should exist. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the DataShare Account should exist. Changing this forces a new resource to be created.

* `created_at` - (Optional) Time at which the account was created.

* `user_email` - (Optional) Email of the user who created the resource.

* `user_name` - (Optional) Name of the user who created the resource.

* `tags` - (Required) Tags on the azure resource. Changing this forces a new resource to be created.
---

## Attributes Reference

The following attributes are exported:

* `id` - The Resource ID of the DataShare Account.

## Timeouts

~> **Note:** Custom Timeouts are available [as an opt-in Beta in version 1.43 & 1.44 of the Azure Provider](/docs/providers/azurerm/guides/2.0-beta.html) and will be enabled by default in version 2.0 of the Azure Provider.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DataShare Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the DataShare Account.
* `update` - (Defaults to 30 minutes) Used when updating the DataShare Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the DataShare Account.

## Import

DataShare Account can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_data_share_account.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataShare/accounts/account1
```