---
subcategory: "Artifacts Signing"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_artifacts_signing_account"
description: |-
  Manages a Artifacts Signing Account.
---

# azurerm_artifacts_signing_account

Manages a Artifacts Signing Account.

~> **Note:** The `azurerm_artifacts_signing_account` resource has been deprecated in favour of `azurerm_artifacts_signing_account` and will be removed in v5.0 of the AzureRM Provider.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_artifacts_signing_account" "example" {
  name                = "example-account"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
  sku_name            = "Basic"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Artifacts Signing Account. Changing this forces a new Artifacts Signing Account to be created.

* `location` - (Required) The Azure Region where the Artifacts Signing Account should exist. Changing this forces a new Artifacts Signing Account to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Artifacts Signing Account should exist. Changing this forces a new Artifacts Signing Account to be created.

* `sku_name` - (Required) The sku name of this Artifacts Signing Account. Possible values are `Basic` and `Premium`.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Artifacts Signing Account.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Artifacts Signing Account.

* `account_uri` - The URI of the Artifacts Signing Account which is used during signing files.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Artifacts Signing Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Artifacts Signing Account.
* `update` - (Defaults to 10 minutes) Used when updating the Artifacts Signing Account.
* `delete` - (Defaults to 10 minutes) Used when deleting the Artifacts Signing Account.

## Import

Artifacts Signing Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_Artifacts_signing_account.example /subscriptions/0000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.CodeSigning/codeSigningAccounts/example-account
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.CodeSigning` - 2025-10-13
