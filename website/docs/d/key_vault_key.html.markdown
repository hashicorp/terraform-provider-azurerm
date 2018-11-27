---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_key"
sidebar_current: "docs-azurerm-data-source-key-vault-key"
description: |-
  Gets information about an existing Key Vault Key.

---

# Data Source: azurerm_key_vault_key

Use this data source to access information about an existing Key Vault Key.

~> **Note:** All arguments including the secret value will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
data "azurerm_key_vault_key" "test" {
  name      = "secret-sauce"
  vault_uri = "https://rickslab.vault.azure.net/"
}

output "key_type" {
  value = "${data.azurerm_key_vault_secret.test.key_type}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Key Vault Key.

* `vault_uri` - (Required) Specifies the URI used to access the Key Vault instance, available on the `azurerm_key_vault` Data Source / Resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Key Vault Key.

* `e` - The RSA public exponent of this Key Vault Key.

* `key_type` - Specifies the Key Type of this Key Vault Key

* `key_size` - Specifies the Size of this Key Vault Key.

* `key_opts` - A list of JSON web key operations assigned to this Key Vault Key

* `n` - The RSA modulus of this Key Vault Key.

* `tags` - A mapping of tags assigned to this Key Vault Key.

* `version` - The current version of the Key Vault Key.

