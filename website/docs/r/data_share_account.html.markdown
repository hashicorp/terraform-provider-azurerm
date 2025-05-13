---
subcategory: "Data Share"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_share_account"
description: |-
  Manages a Data Share Account.
---

# azurerm_data_share_account

Manages a Data Share Account.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_share_account" "example" {
  name                = "example-dsa"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  identity {
    type = "SystemAssigned"
  }

  tags = {
    foo = "bar"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Data Share Account. Changing this forces a new Data Share Account to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Data Share Account should exist. Changing this forces a new Data Share Account to be created.

* `location` - (Required) The Azure Region where the Data Share Account should exist. Changing this forces a new Data Share Account to be created.

* `identity` - (Required) An `identity` block as defined below. Changing this forces a new resource to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Data Share Account.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Data Share Account. The only possible value is `SystemAssigned`. Changing this forces a new resource to be created.

~> **Note:** The assigned `principal_id` and `tenant_id` can be retrieved after the identity `type` has been set to `SystemAssigned` and the Data Share Account has been created. More details are available below.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Share Account.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this Data Share Account.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this Data Share Account.

-> **Note:** You can access the Principal ID via `${azurerm_data_share_account.example.identity[0].principal_id}` and the Tenant ID via `${azurerm_data_share_account.example.identity[0].tenant_id}`

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Share Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Share Account.
* `update` - (Defaults to 30 minutes) Used when updating the Data Share Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Share Account.

## Import

Data Share Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_share_account.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataShare/accounts/account1
```
