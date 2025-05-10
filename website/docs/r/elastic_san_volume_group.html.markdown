---
subcategory: "Elastic SAN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_elastic_san_volume_group"
description: |-
  Manages an Elastic SAN Volume Group resource.
---

# azurerm_elastic_san_volume_group

Manages an Elastic SAN Volume Group resource.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_elastic_san" "example" {
  name                = "examplees-es"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  base_size_in_tib    = 1
  sku {
    name = "Premium_LRS"
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example-uai"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
  service_endpoints    = ["Microsoft.Storage.Global"]

}

resource "azurerm_key_vault" "example" {
  name                        = "examplekv"
  location                    = azurerm_resource_group.example.location
  resource_group_name         = azurerm_resource_group.example.name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = true
  sku_name                    = "standard"
}

resource "azurerm_key_vault_access_policy" "userAssignedIdentity" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_user_assigned_identity.example.principal_id

  key_permissions    = ["Get", "UnwrapKey", "WrapKey"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions    = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_key" "example" {
  name         = "example-kvk"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = [azurerm_key_vault_access_policy.userAssignedIdentity, azurerm_key_vault_access_policy.client]
}

resource "azurerm_elastic_san_volume_group" "example" {
  name            = "example-esvg"
  elastic_san_id  = azurerm_elastic_san.example.id
  encryption_type = "EncryptionAtRestWithCustomerManagedKey"

  encryption {
    key_vault_key_id          = azurerm_key_vault_key.example.versionless_id
    user_assigned_identity_id = azurerm_user_assigned_identity.example.id
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }

  network_rule {
    subnet_id = azurerm_subnet.example.id
    action    = "Allow"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Elastic SAN Volume Group. Changing this forces a new resource to be created.

* `elastic_san_id` - (Required) Specifies the Elastic SAN ID within which this Elastic SAN Volume Group should exist. Changing this forces a new resource to be created.

* `encryption_type` - (Optional) Specifies the type of the key used to encrypt the data of the disk. Possible values are `EncryptionAtRestWithCustomerManagedKey` and `EncryptionAtRestWithPlatformKey`. Defaults to `EncryptionAtRestWithPlatformKey`.

* `encryption` - (Optional) An `encryption` block as defined below.

-> **Note:** The `encryption` block can only be set when `encryption_type` is set to `EncryptionAtRestWithCustomerManagedKey`.

* `identity` - (Optional) An `identity` block as defined below. Specifies the Managed Identity which should be assigned to this Elastic SAN Volume Group.

* `network_rule` - (Optional) One or more `network_rule` blocks as defined below.

* `protocol_type` - (Optional) Specifies the type of the storage target. The only possible value is `Iscsi`. Defaults to `Iscsi`.

---

An `encryption` block supports the following arguments:

* `key_vault_key_id` - (Required) The Key Vault key URI for Customer Managed Key encryption, which can be either a full URI or a versionless URI.

* `user_assigned_identity_id` - (Optional) The ID of the User Assigned Identity used by this Elastic SAN Volume Group.

---

An `identity` block supports the following arguments:

* `type` - (Required) Specifies the type of Managed Identity that should be assigned to this Elastic SAN Volume Group. Possible values are `SystemAssigned` and `UserAssigned`.

* `identity_ids` - (Optional) A list of the User Assigned Identity IDs that should be assigned to this Elastic SAN Volume Group.

---

A `network_rule` block supports the following arguments:

* `subnet_id` - (Required) The ID of the Subnet which should be allowed to access this Elastic SAN Volume Group.

* `action` - (Optional) The action to take when the Subnet attempts to access this Elastic SAN Volume Group. The only possible value is `Allow`. Defaults to `Allow`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Elastic SAN Volume Group.

* `encryption` - An `encryption` block as defined below.

* `identity` - An `identity` block as defined below.

---

An `encryption` block exports the following arguments:

* `current_versioned_key_expiration_timestamp` - The timestamp of the expiration time for the current version of the customer managed key.

* `current_versioned_key_id` - The ID of the current versioned Key Vault Key in use.

* `last_key_rotation_timestamp` - The timestamp of the last rotation of the Key Vault Key.

---

An `identity` block exports the following arguments:

* `principal_id` - The Principal ID associated with the Managed Service Identity assigned to this Elastic SAN Volume Group.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity assigned to this Elastic SAN Volume Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Elastic SAN Volume Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Elastic SAN Volume Group.
* `update` - (Defaults to 30 minutes) Used when updating the Elastic SAN Volume Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Elastic SAN Volume Group.

## Import

An existing Elastic SAN Volume Group can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_elastic_san_volume_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ElasticSan/elasticSans/esan1/volumeGroups/vg1
```
