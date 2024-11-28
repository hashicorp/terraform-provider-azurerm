---
subcategory: "AI Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_ai_services_project"
description: |-
  Manages an AI Services Project.
---

# azurerm_ai_services_project

Manages an AI Services Project.

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

resource "azurerm_ai_services" "example" {
  name                = "exampleaiservices"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "S0"
}

resource "azurerm_ai_services_hub" "example" {
  name                = "exampleaihub"
  location            = azurerm_ai_services.example.location
  resource_group_name = azurerm_resource_group.example.name
  storage_account_id  = azurerm_storage_account.example.id
  key_vault_id        = azurerm_key_vault.example.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_ai_services_project" "example" {
  name               = "example"
  location           = azurerm_ai_services_hub.example.location
  ai_services_hub_id = azurerm_ai_services_hub.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this AI Services Project. Changing this forces a new AI Services Project to be created.

* `location` - (Required) The Azure Region where the AI Services Project should exist. Changing this forces a new AI Services Project to be created.

* `ai_services_hub_id` - (Required) The AI Services Hub ID under which this Project should be created. Changing this forces a new AI Services Project to be created.

---

* `description` - (Optional) The description of this AI Services Project.

* `friendly_name` - (Optional) The display name of this AI Services Project.

* `high_business_impact_enabled` - (Optional) Whether High Business Impact (HBI) should be enabled or not. Enabling this setting will reduce diagnostic data collected by the service. Changing this forces a new AI Services Project to be created. Defaults to `false`.

* `identity` - (Optional) A `identity` block as defined below.

* `image_build_compute_name` - (Optional) The compute name for image build of the AI Services Project.

* `tags` - (Optional) A mapping of tags which should be assigned to the AI Services Project.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this AI Services Project. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this AI Services Project.

~> **NOTE:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the AI Services Project.

* `project_id` - The immutable project ID associated with this AI Services Project.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 hour) Used when creating the AI Services Project.
* `read` - (Defaults to 5 minutes) Used when retrieving the AI Services Project.
* `update` - (Defaults to 30 minutes) Used when updating the AI Services Project.
* `delete` - (Defaults to 30 minutes) Used when deleting the AI Services Project.

## Import

AI Services Projects can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_ai_services_project.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.MachineLearningServices/workspaces/project1
```
