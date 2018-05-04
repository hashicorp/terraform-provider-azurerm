---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_secret"
sidebar_current: "docs-azurerm-data-source-key-vault-secret"
description: |-
  Returns information about the specified Key Vault Secret.

---

# azurerm_key_vault_secret

Returns information about the specified Key Vault Secret.

## Example Usage

```hcl
data "azurerm_key_vault_secret" "test" {
  name      = "secret-sauce"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"
}

output "secret_value" {
  value = "${data.azurerm_key_vault_secret.test.value}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Key Vault Secret.

* `vault_uri` - (Required) Specifies the URI used to access the Key Vault instance, available on the `azurerm_key_vault` Data Source / Resource.


## Attributes Reference

The following attributes are exported:

* `id` - The Key Vault Secret ID.
* `value` - The value of the Key Vault Secret
* `version` - The current version of the Key Vault Secret.
* `content_type` - The content type for the Key Vault Secret.
* `tags` - Any tags assigned to this resource.
