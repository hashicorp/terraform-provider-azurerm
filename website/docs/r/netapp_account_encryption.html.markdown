---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_account_encryption"
description: |-
  Manages a NetApp Account Encryption Resource.
---

# azurerm_netapp_account_encryption

Manages a NetApp Account Encryption Resource.

For more information about Azure NetApp Files Customer-Managed Keys feature, please refer to [Configure customer-managed keys for Azure NetApp Files volume encryption](https://learn.microsoft.com/en-us/azure/azure-netapp-files/configure-customer-managed-keys)

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_client_config" "current" {
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "anf-user-assigned-identity"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_key_vault" "example" {
  name                            = "anfcmkakv"
  location                        = azurerm_resource_group.example.location
  resource_group_name             = azurerm_resource_group.example.name
  enabled_for_disk_encryption     = true
  enabled_for_deployment          = true
  enabled_for_template_deployment = true
  purge_protection_enabled        = true
  tenant_id                       = "00000000-0000-0000-0000-000000000000"

  sku_name = "standard"

  access_policy {
    tenant_id = "00000000-0000-0000-0000-000000000000"
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Get",
      "Create",
      "Delete",
      "WrapKey",
      "UnwrapKey",
      "GetRotationPolicy",
      "SetRotationPolicy",
    ]
  }

  access_policy {
    tenant_id = "00000000-0000-0000-0000-000000000000"
    object_id = azurerm_user_assigned_identity.example.principal_id

    key_permissions = [
      "Get",
      "Encrypt",
      "Decrypt"
    ]
  }
}

resource "azurerm_key_vault_key" "example" {
  name         = "anfencryptionkey"
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
}

resource "azurerm_netapp_account" "example" {
  name                = "netappaccount"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.example.id
    ]
  }
}

resource "azurerm_netapp_account_encryption" "example" {
  netapp_account_id = azurerm_netapp_account.example.id

  user_assigned_identity_id = azurerm_user_assigned_identity.example.id

  encryption_key = azurerm_key_vault_key.example.versionless_id
}
```

## Arguments Reference

The following arguments are supported:

* `encryption_key` - (Required) Specify the versionless ID of the encryption key.

* `netapp_account_id` - (Required) The ID of the NetApp account where volume under it will have customer managed keys-based encryption enabled.

---

* `system_assigned_identity_principal_id` - (Optional) The ID of the System Assigned Manged Identity. Conflicts with `user_assigned_identity_id`.

* `user_assigned_identity_id` - (Optional) The ID of the User Assigned Managed Identity. Conflicts with `system_assigned_identity_principal_id`.

---



A full example of the `azurerm_netapp_account_encryption` resource and NetApp Volume with customer-managed keys encryption enabled can be found in [the `./examples/netapp/nfsv3_volume_cmk_userassigned` directory within the GitHub Repository](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/netapp/nfsv3_volume_cmk_userassigned)

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Account Encryption Resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Account Encryption Resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the Account Encryption Resource.
* `update` - (Defaults to 2 hours) Used when updating the Account Encryption Resource.
* `delete` - (Defaults to 2 hours) Used when deleting the Account Encryption Resource.

## Import

Account Encryption Resources can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_netapp_account_encryption.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1
```
