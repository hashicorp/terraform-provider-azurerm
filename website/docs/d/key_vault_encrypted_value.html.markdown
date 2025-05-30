---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_key_vault_encrypted_value"
description: |-
    Encrypts or Decrypts a value using a Key Vault Key.
---

# Data Source: azurerm_key_vault_encrypted_value

Encrypts or Decrypts a value using a Key Vault Key.

## Example Usage

```hcl
data "azurerm_key_vault" "example" {
  name                = "mykeyvault"
  resource_group_name = "some-resource-group"
}

data "azurerm_key_vault_key" "example" {
  name         = "some-key"
  key_vault_id = data.azurerm_key_vault.example.id
}

data "azurerm_key_vault_encrypted_value" "encrypted" {
  key_vault_key_id = azurerm_key_vault_key.test.id
  algorithm        = "RSA1_5"
  plain_text_value = base64encode("some-encrypted-value")
}

data "azurerm_key_vault_encrypted_value" "decrypted" {
  key_vault_key_id = azurerm_key_vault_key.test.id
  algorithm        = "RSA1_5"
  encrypted_data   = data.azurerm_key_vault_encrypted_value.encrypted.encrypted_data
}

output "id" {
  value = data.azurerm_key_vault_encrypted_value.example.encrypted_data
}

output "decrypted_text" {
  value = nonsensitive(data.azurerm_key_vault_encrypted_value.decrypted.decoded_plain_text_value)
}
```

## Arguments Reference

The following arguments are supported:

* `algorithm` - (Required) The Algorithm which should be used to Decrypt/Encrypt this Value. Possible values are `RSA1_5`, `RSA-OAEP` and `RSA-OAEP-256`.

* `key_vault_key_id` - (Required) The ID of the Key Vault Key which should be used to Decrypt/Encrypt this Value.

---

* `encrypted_data` - (Optional) The Base64 URL Encoded Encrypted Data which should be decrypted into `plain_text_value`.

* `plain_text_value` - (Optional) The plain-text value which should be Encrypted into `encrypted_data`.

-> **Note:** One of either `encrypted_data` or `plain_text_value` must be specified and is used to populate the encrypted/decrypted value for the other field.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of this Encrypted Value

* `decoded_plain_text_value` - The Base64URL decoded string of `plain_text_value`. Because the API would remove padding characters of `plain_text_value` when encrypting, this attribute is useful to get the original value.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Encrypted Value
