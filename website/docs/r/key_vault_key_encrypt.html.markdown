---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_key_encrypt"
description: |-
  Encrypt data with a Key Vault Key.

---

# azurerm_key_vault_key_encrypt

Use this resource to encrypt data with a Key Vault Key.

## Example Usage

```hcl
resource "azurerm_key_vault_key_encrypt" "example" {
  name             = "example_encrypt"
  key_vault_key_id = azurerm_key_vault_key.example.id
  plaintext        = "testData"
  algorithm        = "RSA1_5"
}

output "encrypted_data" {
  value = azurerm_key_vault_key_encrypt.example.encrypted_data_in_base64url
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this resource. Changing this forces a new resource to be created.

* `key_vault_key_id` - (Required) Specifies the ID of the Key Vault key which is used to encrypt. Changing this forces a new resource to be created.

* `plaintext` - (Required) Specifies the Data to be encrypted. Changing this forces a new resource to be created.

* `algorithm` - (Required) Specifies the Algorithm which is used to encrypt. Possible values are `RSA1_5`, `RSA-OAEP` and `RSA-OAEP-256`. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the data encrypted.

* `encrypted_data_in_base64url` - The encrypted data.

**NOTE:** The format of `encrypted_data_in_base64url` is `base64url`, please refer to RFC: https://tools.ietf.org/html/rfc4648#section-5

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when encrypting data using the Key Vault Key.

## Import

this resource does not support `import` operation.
