---
subcategory: "Trusted Signing"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_trusted_signing_account"
description: |-
  Manages a Trusted Signing Account.
---

# azurerm_trusted_signing_account

Manages a Trusted Signing Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_trusted_signing_account" "example" {
  name                = "example-account"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
  sku_name            = "Basic"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Trusted Signing Account. Changing this forces a new Trusted Signing Account to be created.

* `location` - (Required) The Azure Region where the Trusted Signing Account should exist. Changing this forces a new Trusted Signing Account to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Trusted Signing Account should exist. Changing this forces a new Trusted Signing Account to be created.

* `sku_name` - (Required) The sku name of this Trusted Signing Account. Possible values are `Basic` and `Premium`.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Trusted Signing Account.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Trusted Signing Account.

* `account_uri` - The URI of the Trusted Signing Account which is used during signing files.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Trusted Signing Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Trusted Signing Account.
* `update` - (Defaults to 10 minutes) Used when updating the Trusted Signing Account.
* `delete` - (Defaults to 10 minutes) Used when deleting the Trusted Signing Account.

## Import

Trusted Signing Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_trusted_signing_account.example /subscriptions/0000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.CodeSigning/codeSigningAccounts/example-account
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.CodeSigning`: 2024-09-30-preview
