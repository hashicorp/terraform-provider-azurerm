---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_certificate_issuer"
description: |-
  Manages a Key Vault Certificate Issuer.
---

# azurerm_key_vault_certificate_issuer

Manages a Key Vault Certificate Issuer.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_key_vault" "example" {
  name                = "examplekeyvault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "standard"
  tenant_id           = data.azurerm_client_config.current.tenant_id
}

resource "azurerm_key_vault_certificate_issuer" "example" {
  name          = "example-issuer"
  org_id        = "ExampleOrgName"
  key_vault_id  = azurerm_key_vault.example.id
  provider_name = "DigiCert"
  account_id    = "0000"
  password      = "example-password"
}
```

## Arguments Reference

The following arguments are supported:

* `key_vault_id` - (Required) The ID of the Key Vault in which to create the Certificate Issuer. Changing this forces a new resource to be created.

* `name` - (Required) The name which should be used for this Key Vault Certificate Issuer. Changing this forces a new Key Vault Certificate Issuer to be created.

* `provider_name` - (Required) The name of the third-party Certificate Issuer. Possible values are: `DigiCert`, `GlobalSign`, `OneCertV2-PrivateCA`, `OneCertV2-PublicCA` and `SslAdminV2`.

* `org_id` - (Optional) The ID of the organization as provided to the issuer.

* `account_id` - (Optional) The account number with the third-party Certificate Issuer.

* `admin` - (Optional) One or more `admin` blocks as defined below.

* `password` - (Optional) The password associated with the account and organization ID at the third-party Certificate Issuer. If not specified, will not overwrite any previous value.

---

An `admin` block supports the following:

* `email_address` - (Required) E-mail address of the admin.

* `first_name` - (Optional) First name of the admin.

* `last_name` - (Optional) Last name of the admin.

* `phone` - (Optional) Phone number of the admin.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Key Vault Certificate Issuer.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Key Vault Certificate Issuer.
* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Certificate Issuer.
* `update` - (Defaults to 30 minutes) Used when updating the Key Vault Certificate Issuer.
* `delete` - (Defaults to 30 minutes) Used when deleting the Key Vault Certificate Issuer.

## Import

Key Vault Certificate Issuers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_certificate_issuer.example "https://key-vault-name.vault.azure.net/certificates/issuers/example"
```
