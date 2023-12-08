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

provider "azurerm" {
  features {}
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

resource "azurerm_role_assignment" "client" {
  scope                = azurerm_resource_group.example.id
  role_definition_name = "Key Vault Crypto Officer"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_role_assignment" "userAssignedIdentity" {
  scope                = azurerm_resource_group.example.id
  role_definition_name = "Key Vault Crypto Service Encryption User"
  principal_id         = azurerm_user_assigned_identity.example.principal_id
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

  depends_on = [azurerm_role_assignment.userAssignedIdentity, azurerm_role_assignment.client]
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

* `encryption_type` - (Optional) Type of encryption. Possible values are `EncryptionAtRestWithCustomerManagedKey` and `EncryptionAtRestWithPlatformKey`. Defaults to `EncryptionAtRestWithPlatformKey`.

* `encryption` - (Optional) An `encryption` block as defined below.

**NOTE:** The `encryption` block can only be set when `encryption_type` is set to `EncryptionAtRestWithCustomerManagedKey`.

* `identity` - (Optional) An `identity` block as defined below. Specifies the Managed Identity which should be assigned to this Elastic SAN Volume Group.

* `network_rule` - (Optional) One or more `network_rule` blocks as defined below.

* `protocol_type` - (Optional) The protocol type of the Elastic SAN Volume Group. The only possible value is `Iscsi`. Defaults to `Iscsi`.

---

An `encryption` block supports the following arguments:

* `identity` - (Optional) An `identity` block as defined below. 
* `key_vault_key_id` - (Optional) 

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

---

An `encryption` block exports the following arguments:

* `current_versioned_key_expiration_timestamp` - Current Versioned Key Expiration Timestamp.

* `current_versioned_key_id` - Current Versioned Key ID.

* `last_key_rotation_timestamp` - Last Key Rotation Timestamp.

---

An `identity` block exports the following arguments:

* `principal_id` - The Principal ID for the System-Assigned Managed Identity assigned to this Elastic SAN Volume Group.

* `tenant_id` - The Tenant ID for the System-Assigned Managed Identity assigned to this Elastic SAN Volume Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating this Elastic SAN Volume Group.
* `delete` - (Defaults to 30 minutes) Used when deleting this Elastic SAN Volume Group.
* `read` - (Defaults to 5 minutes) Used when retrieving this Elastic SAN Volume Group.
* `update` - (Defaults to 30 minutes) Used when updating this Elastic SAN Volume Group.

## Import

An existing Elastic SAN Volume Group can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_elastic_san_volume_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ElasticSan/elasticSans/esan1/volumeGroups/vg1
```
