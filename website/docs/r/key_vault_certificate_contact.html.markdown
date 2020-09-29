---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_certificate_contact"
description: |-
  Manages Key Vault Certificate Contacts.
---

# azurerm_key_vault_certificate_contact

Manages a Key Vault Certificate Contacts.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West US"
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "example" {
  name                = "example-key-vault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "standard"
  tenant_id           = data.azurerm_client_config.current.tenant_id

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "delete",
      "import",
      "get",
      "ManageContacts",
    ]
  }
}

resource "azurerm_key_vault_certificate_contact" "example" {
  key_vault_id = data.azurerm_key_vault.example.id
  contact {
    email = "example@example.com"
    name  = "example"
    phone = "0123456789"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `key_vault_id` - (Required) The ID of the Key Vault in which to create the Certificate Contact. Changing this forces a new Key Vault Certificate Contact to be created.

* `contact` - (Required) One or more `contact` block as defined below.

---

An `contact` block supports the following:

* `email` - (Required) E-mail address of the contact.

* `name` - (Optional) Name of the contact.

* `phone` - (Optional) Phone number of the contact.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Key Vault Certificate Contact.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Key Vault Certificate Contacts.
* `update` - (Defaults to 30 minutes) Used when updating the Key Vault Certificate Contacts.
* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Certificate Contacts.
* `delete` - (Defaults to 30 minutes) Used when deleting the Key Vault Certificate Contacts.

## Import

Key Vault Certificate Contact can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_certificate_contact.example "https://key-vault-name.vault.azure.net/certificates/contacts"
```
