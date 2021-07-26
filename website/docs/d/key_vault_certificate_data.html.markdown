---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_certificate_data"
description: |-
  Gets data contained in an existing Key Vault Certificate.
---

# Data Source: azurerm_key_vault_certificate_data

Use this data source to access data stored in an existing Key Vault Certificate.

~> **Note:** All arguments including the secret value will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

~> **Note:** This data source uses the `GetSecret` function of the Azure API, to get the key of the certificate. Therefore you need secret/get permission

## Example Usage

```hcl
data "azurerm_key_vault" "example" {
  name                = "examplekv"
  resource_group_name = "some-resource-group"
}

data "azurerm_key_vault_certificate_data" "example" {
  name         = "secret-sauce"
  key_vault_id = data.azurerm_key_vault.example.id
}

output "example_pem" {
  value = data.azurerm_key_vault_certificate_data.example.pem
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Key Vault Secret.

* `key_vault_id` - (Required) Specifies the ID of the Key Vault instance where the Secret resides, available on the `azurerm_key_vault` Data Source / Resource.

* `version` - (Optional) Specifies the version of the certificate to look up.  (Defaults to latest) 

~> **NOTE:** The vault must be in the same subscription as the provider. If the vault is in another subscription, you must create an aliased provider for that subscription.

## Attributes Reference

The following attributes are exported:

* `hex` - The raw Key Vault Certificate data represented as a hexadecimal string. 

* `pem` - The Key Vault Certificate in PEM format. 

* `key` - The Key Vault Certificate Key. 

* `expires` - Expiry date of certificate in RFC3339 format. 

* `tags` - A mapping of tags to assign to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Certificate.
