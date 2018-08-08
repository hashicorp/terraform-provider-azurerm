---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_secret"
sidebar_current: "docs-azurerm-resource-key-vault-secret"
description: |-
  Manages a Key Vault Secret.

---

# azurerm_key_vault_secret

Manages a Key Vault Secret.

~> **Note:** All arguments including the secret value will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_key_vault" "example" {
  # ...
}

resource "azurerm_key_vault_secret" "example" {
  name      = "secret-sauce"
  value     = "szechuan"
  vault_uri = "${azurerm_key_vault.example.vault_uri}"

  tags {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Key Vault Secret. Changing this forces a new resource to be created.

* `value` - (Required) Specifies the value of the Key Vault Secret.

* `vault_uri` - (Required) Specifies the URI used to access the Key Vault instance, available on the `azurerm_key_vault` resource.

* `content_type` - (Optional) Specifies the content type for the Key Vault Secret.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Key Vault Secret ID.
* `version` - The current version of the Key Vault Secret.

## Import

Key Vault Secrets which are Enabled can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_secret.test https://example-keyvault.vault.azure.net/secrets/example/fdf067c93bbb4b22bff4d8b7a9a56217
```
