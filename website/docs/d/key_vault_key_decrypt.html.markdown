---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_key_decrypt"
description: |-
  Decrypt data encrypted with a Key Vault Key.

---

# Data Source: azurerm_key_vault_key_decrypt

Use this data source to decrypt data encrypted with a Key Vault Key.

## Example Usage

```hcl
data "azurerm_key_vault_key_decrypt" "example" {
  key_vault_key_id         = azurerm_key_vault_key.example.id
  encrypted_base64url_data = var.encrypted_base64url_data
  algorithm                = "RSA1_5"
}

output "decrypted_data" {
  value = data.azurerm_key_vault_key_decrypt.example.plaintext
}
```

## Argument Reference

The following arguments are supported:

* `key_vault_key_id` - (Required) Specifies the ID of the Key Vault key which is used to decrypt. 

* `algorithm` - (Required) Specifies the Algorithm which is used to decrypt. Possible values are `RSA1_5`, `RSA-OAEP` and `RSA-OAEP-256`.

* `encrypted_base64url_data` - (Required) Specifies the data to be decrypted.

**NOTE:** The format of `encrypted_base64url_data` is `base64url`, please refer to RFC: https://tools.ietf.org/html/rfc4648#section-5

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the data decrypted.

* `plaintext` - The decrypted data.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Key decrypted data.
