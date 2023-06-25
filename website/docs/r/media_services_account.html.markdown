---
subcategory: "Media"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_media_services_account"
description: |-
  Manages a Media Services Account.
---

# azurerm_media_services_account

Manages a Media Services Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "media-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestoracc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services_account" "example" {
  name                = "examplemediaacc"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  storage_account {
    id         = azurerm_storage_account.example.id
    is_primary = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Media Services Account. Only lowercase Alphanumeric characters allowed. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Media Services Account. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `storage_account` - (Required) One or more `storage_account` blocks as defined below.

* `encryption` - (Optional) An `encryption` block as defined below.

* `identity` - (Optional) An `identity` block as defined below.
 
* `public_network_access_enabled` - (Optional) Whether public network access is allowed for this server. Defaults to `true`.

* `storage_authentication_type` - (Optional) Specifies the storage authentication type. Possible value is `ManagedIdentity` or `System`.

* `key_delivery_access_control` - (Optional) A `key_delivery_access_control` block as defined below.

* `tags` - (Optional) A mapping of tags assigned to the resource.

---

A `storage_account` block supports the following:

* `id` - (Required) Specifies the ID of the Storage Account that will be associated with the Media Services instance.

* `is_primary` - (Optional) Specifies whether the storage account should be the primary account or not. Defaults to `false`.

~> **NOTE:** Whilst multiple `storage_account` blocks can be specified - one of them must be set to the primary

* `managed_identity` - (Optional) A `managed_identity` block as defined below.

---

A `encryption` block supports the following:

* `type` - (Optional) Specifies the type of key used to encrypt the account data. Possible values are `SystemKey` and `CustomerKey`.

* `key_vault_key_identifier` - (Optional) Specifies the URI of the Key Vault Key used to encrypt data. The key may either be versioned (for example https://vault/keys/mykey/version1) or reference a key without a version (for example https://vault/keys/mykey).

* `managed_identity` - (Optional) A `managed_identity` block as defined below.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Media Services Account. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Media Services Account.

---

A `key_delivery_access_control` block supports the following:

* `default_action` - (Optional) The Default Action to use when no rules match from `ip_allow_list`. Possible values are `Allow` and `Deny`.

* `ip_allow_list` - (Optional) One or more IP Addresses, or CIDR Blocks which should be able to access the Key Delivery.

---

A `managed_identity` block supports the following:

* `user_assigned_identity_id` - (Optional) The ID of the User Assigned Identity. This value can only be set when `use_system_assigned_identity` is `false`

* `use_system_assigned_identity` - (Optional) Whether to use System Assigned Identity. Possible Values are `true` and `false`.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Media Services Account.

* `identity` - An `identity` block as defined below.

---

An `encryption` block exports the following:

* `current_key_identifier` - The current key used to encrypt the Media Services Account, including the key version.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Media Services Account.
* `update` - (Defaults to 30 minutes) Used when updating the Media Services Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Media Services Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Media Services Account.

## Import

Media Services Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_media_services_account.account /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Media/mediaServices/account1
```
