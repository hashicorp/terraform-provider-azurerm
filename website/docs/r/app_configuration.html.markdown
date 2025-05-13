---
subcategory: "App Configuration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_configuration"
description: |-
  Manages an Azure App Configuration.

---

# azurerm_app_configuration

Manages an Azure App Configuration.

## Disclaimers

-> **Note:** Version 3.27.0 and later of the Azure Provider include a Feature Toggle which will purge an App Configuration resource on destroy, rather than the default soft-delete. The Provider will automatically recover a soft-deleted App Configuration during creation if one is found. See [the Features block documentation](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block) for more information on Feature Toggles within Terraform.

-> **Note:** Reading and purging soft-deleted App Configurations requires the `Microsoft.AppConfiguration/locations/deletedConfigurationStores/read` and `Microsoft.AppConfiguration/locations/deletedConfigurationStores/purge/action` permission on Subscription scope. Recovering a soft-deleted App Configuration requires the `Microsoft.AppConfiguration/configurationStores/write` permission on Subscription or Resource Group scope. [More information can be found in the Azure Documentation for App Configuration](https://learn.microsoft.com/en-us/azure/azure-app-configuration/concept-soft-delete#permissions-to-recover-a-deleted-store). See the following links for more information on assigning [Azure custom roles](https://learn.microsoft.com/en-us/azure/role-based-access-control/custom-roles) or using the [`azurerm_role_assignment`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/role_assignment) resource to assign a custom role.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_app_configuration" "appconf" {
  name                = "appConf1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Example Usage (encryption)

```hcl
provider "azurerm" {
  features {
    app_configuration {
      purge_soft_delete_on_destroy = true
      recover_soft_deleted         = true
    }
  }
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example-identity"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "example" {
  name                       = "exampleKVt123"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  purge_protection_enabled   = true
}

resource "azurerm_key_vault_access_policy" "server" {
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
  name         = "exampleKVkey"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey"
  ]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.server,
  ]
}

resource "azurerm_app_configuration" "example" {
  name                       = "appConf2"
  resource_group_name        = azurerm_resource_group.example.name
  location                   = azurerm_resource_group.example.location
  sku                        = "standard"
  local_auth_enabled         = true
  public_network_access      = "Enabled"
  purge_protection_enabled   = false
  soft_delete_retention_days = 1

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.example.id,
    ]
  }

  encryption {
    key_vault_key_identifier = azurerm_key_vault_key.example.id
    identity_client_id       = azurerm_user_assigned_identity.example.client_id
  }

  replica {
    name     = "replica1"
    location = "West US"
  }

  tags = {
    environment = "development"
  }

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.server,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the App Configuration. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the App Configuration. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `data_plane_proxy_authentication_mode` - The data plane proxy authentication mode. Possible values are `Local` and `Pass-through`. Defaults to `Local`.

* `data_plane_proxy_private_link_delegation_enabled` - Whether data plane proxy private link delegation is enabled. Defaults to `false`.

~> **Note:** `data_plane_proxy_private_link_delegation_enabled` cannot be set to `true` when `data_plane_proxy_authentication_mode` is set to `Local`.

* `encryption` - (Optional) An `encryption` block as defined below.

* `identity` - (Optional) An `identity` block as defined below.

* `local_auth_enabled` - (Optional) Whether local authentication methods is enabled. Defaults to `true`.

* `public_network_access` - (Optional) The Public Network Access setting of the App Configuration. Possible values are `Enabled` and `Disabled`.

~> **Note:** If `public_network_access` is not specified, the App Configuration will be created as  `Automatic`. However, once a different value is defined, can not be set again as automatic.

* `purge_protection_enabled` - (Optional) Whether Purge Protection is enabled. This field only works for `standard` sku. Defaults to `false`.

!> **Note:** Once Purge Protection has been enabled it's not possible to disable it. Deleting the App Configuration with Purge Protection enabled will schedule the App Configuration to be deleted (which will happen by Azure in the configured number of days).

* `replica` - (Optional) One or more `replica` blocks as defined below.

* `sku` - (Optional) The SKU name of the App Configuration. Possible values are `free`, `developer`, `standard` and `premium`. Defaults to `free`.

~> **Note:** Azure does not support downgrading `sku` to a lower tier, except from `premium` to `standard`. Downgrading will force a new resource to be created.

* `soft_delete_retention_days` - (Optional) The number of days that items should be retained for once soft-deleted. This field only works for `standard` sku. This value can be between `1` and `7` days. Defaults to `7`. Changing this forces a new resource to be created.

~> **Note:** If Purge Protection is enabled, this field can only be configured one time and cannot be updated.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `encryption` block supports the following:

* `identity_client_id` - (Optional) Specifies the client ID of the identity which will be used to access key vault.

* `key_vault_key_identifier` - (Optional) Specifies the URI of the key vault key used to encrypt data.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this App Configuration. Possible values are `SystemAssigned`, `UserAssigned`, and `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) A list of User Assigned Managed Identity IDs to be assigned to this App Configuration.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `replica` block supports the following:

* `name` - (Required) Specifies the name of the replica.

* `location` - (Required) Specifies the supported Azure location where the replica exists.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The App Configuration ID.

* `endpoint` - The URL of the App Configuration.

* `identity` - An `identity` block as defined below.

* `primary_read_key` - A `primary_read_key` block as defined below containing the primary read access key.

* `primary_write_key` - A `primary_write_key` block as defined below containing the primary write access key.

* `replica` - A `replica` block as defined below.

* `secondary_read_key` - A `secondary_read_key` block as defined below containing the secondary read access key.

* `secondary_write_key` - A `secondary_write_key` block as defined below containing the secondary write access key.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

---

A `replica` block exports the following:

* `id` - The ID of the App Configuration Replica.

* `endpoint` - The URL of the App Configuration Replica.

---

A `primary_read_key` block exports the following:

* `id` - The ID of the Access Key.

* `connection_string` - The Connection String for this Access Key - consisting of the Endpoint, ID, and Secret.

* `secret` - The Secret of the Access Key.

---

A `primary_write_key` block exports the following:

* `id` - The ID of the Access Key.

* `connection_string` - The Connection String for this Access Key - consisting of the Endpoint, ID, and Secret.

* `secret` - The Secret of the Access Key.

---

A `secondary_read_key` block exports the following:

* `id` - The ID of the Access Key.

* `connection_string` - The Connection String for this Access Key - consisting of the Endpoint, ID, and Secret.

* `secret` - The Secret of the Access Key.

---

A `secondary_write_key` block exports the following:

* `id` - The ID of the Access Key.

* `connection_string` - The Connection String for this Access Key - consisting of the Endpoint, ID, and Secret.

* `secret` - The Secret of the Access Key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the App Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Configuration.
* `update` - (Defaults to 1 hour) Used when updating the App Configuration.
* `delete` - (Defaults to 1 hour) Used when deleting the App Configuration.

## Import

App Configurations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_configuration.appconf /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.AppConfiguration/configurationStores/appConf1
```
