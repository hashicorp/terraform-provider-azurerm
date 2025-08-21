---
subcategory: "Machine Learning"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_ai_hub"
description: |-
  Manages an AI Hub.
---

# azurerm_ai_hub

Manages an AI Hub.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "westeurope"
}

resource "azurerm_key_vault" "example" {
  name                = "examplekv"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Get",
    "Delete",
    "Purge",
    "GetRotationPolicy",
  ]
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_ai_foundry" "example" {
  name                = "exampleaiservices"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "S0"
}

resource "azurerm_ai_hub" "example" {
  name                = "exampleaihub"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  storage_account_id  = azurerm_storage_account.example.id
  key_vault_id        = azurerm_key_vault.example.id

  identity {
    type = "SystemAssigned"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this AI Hub. Changing this forces a new AI Hub to be created.

* `location` - (Required) The Azure Region where the AI Hub should exist. Changing this forces a new AI Hub to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the AI Hub should exist. Changing this forces a new AI Hub to be created.

* `identity` - (Required) A `identity` block as defined below.

* `key_vault_id` - (Required) The Key Vault ID that should be used by this AI Hub. Changing this forces a new AI Hub to be created.

* `storage_account_id` - (Required) The Storage Account ID that should be used by this AI Hub. Changing this forces a new AI Hub to be created.

---

* `application_insights_id` - (Optional) The Application Insights ID that should be used by this AI Hub.

* `container_registry_id` - (Optional) The Container Registry ID that should be used by this AI Hub.

* `description` - (Optional) The description of this AI Hub.

* `encryption` - (Optional) An `encryption` block as defined below. Changing this forces a new AI Hub to be created.

* `friendly_name` - (Optional) The display name of this AI Hub.

* `high_business_impact_enabled` - (Optional) Whether High Business Impact (HBI) should be enabled or not. Enabling this setting will reduce diagnostic data collected by the service. Changing this forces a new AI Hub to be created. Defaults to `false`.

-> **Note:** `high_business_impact_enabled` will be enabled by default when creating an AI Hub with `encryption` enabled.

* `managed_network` - (Optional) A `managed_network` block as defined below.

* `primary_user_assigned_identity` - (Optional) The user assigned identity ID that represents the AI Hub identity. This must be set when enabling encryption with a user assigned identity.

* `public_network_access` - (Optional) Whether public network access for this AI Service Hub should be enabled. Possible values include `Enabled` and `Disabled`. Defaults to `Enabled`.

* `tags` - (Optional) A mapping of tags which should be assigned to the AI Hub.

---

A `encryption` block supports the following:

* `key_id` - (Required) The Key Vault URI to access the encryption key.

* `key_vault_id` - (Required) The Key Vault ID where the customer owned encryption key exists.

* `user_assigned_identity_id` - (Optional) The user assigned identity ID that has access to the encryption key.

~> **Note:** `user_assigned_identity_id` must be set when`identity.type` is `UserAssigned` in order for the service to find the assigned permissions.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this AI Hub. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this AI Hub.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `managed_network` block supports the following:

* `isolation_mode` - (Optional) The isolation mode of the AI Hub. Possible values are `Disabled`, `AllowOnlyApprovedOutbound`, and `AllowInternetOutbound`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the AI Hub.

* `discovery_url` - The URL for the discovery service to identify regional endpoints for AI Hub services.

* `workspace_id` - The immutable ID associated with this AI Hub.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the AI Hub.
* `read` - (Defaults to 5 minutes) Used when retrieving the AI Hub.
* `update` - (Defaults to 30 minutes) Used when updating the AI Hub.
* `delete` - (Defaults to 30 minutes) Used when deleting the AI Hub.

## Import

AI Hubs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_ai_hub.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.MachineLearningServices/workspaces/hub1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.MachineLearningServices` - 2025-06-01
