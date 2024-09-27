---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_certificate_contacts"
description: |-
  Manages Key Vault Certificate Contacts.
---

# azurerm_key_vault_certificate_contacts

Manages Key Vault Certificate Contacts.

## Disclaimers

<!-- TODO: Remove Note in 4.0 -->
~> **Note:** It's possible to define Key Vault Certificate Contacts both within [the `azurerm_key_vault` resource](key_vault.html) via the `contact` block and by using [the `azurerm_key_vault_certificate_contacts` resource](key_vault_certificate_contacts.html). However it's not possible to use both methods to manage Certificate Contacts within a KeyVault, since there'll be conflicts.

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
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"
}

resource "azurerm_key_vault_access_policy" "example" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  certificate_permissions = [
    "ManageContacts",
  ]

  key_permissions = [
    "Create",
  ]

  secret_permissions = [
    "Set",
  ]
}

resource "azurerm_key_vault_certificate_contacts" "example" {
  key_vault_id = azurerm_key_vault.example.id

  contact {
    email = "example@example.com"
    name  = "example"
    phone = "01234567890"
  }

  contact {
    email = "example2@example.com"
  }

  depends_on = [
    azurerm_key_vault_access_policy.example
  ]
}

```

## Arguments Reference

The following arguments are supported:

* `key_vault_id` - (Required) The ID of the Key Vault. Changing this forces a new resource to be created.

* `contact` - (Required) One or more `contact` blocks as defined below.
<!-- TODO: Update in 4.0
* `contact` - (Optional) One or more `contact` blocks as defined below.
-->

---

A `contact` block supports the following:

* `email` - (Required) E-mail address of the contact.

* `name` - (Optional) Name of the contact.

* `phone` - (Optional) Phone number of the contact.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Key Vault Certificate Contacts.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Key Vault Certificate Contacts.
* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Certificate Contacts.
* `update` - (Defaults to 30 minutes) Used when updating the Key Vault Certificate Contacts.
* `delete` - (Defaults to 30 minutes) Used when deleting the Key Vault Certificate Contacts.

## Import

Key Vault Certificate Contacts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_certificate_contacts.example https://example-keyvault.vault.azure.net/certificates/contacts
```
