---
subcategory: "Machine Learning"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_machine_learning_workspace"
description: |-
  Manages a Azure Machine Learning Workspace.
---
# azurerm_machine_learning_workspace

Manages a Azure Machine Learning Workspace

~> **NOTE:** For examples on how to set up the Azure Machine Learning workspace, together with compute and integrated services, see [Terraform Quickstart](https://github.com/Azure/terraform/tree/master/quickstart)

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "workspace-example-ai"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_key_vault" "example" {
  name                = "workspaceexamplekeyvault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"
}

resource "azurerm_storage_account" "example" {
  name                     = "workspacestorageaccount"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_machine_learning_workspace" "example" {
  name                    = "example-workspace"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  application_insights_id = azurerm_application_insights.example.id
  key_vault_id            = azurerm_key_vault.example.id
  storage_account_id      = azurerm_storage_account.example.id

  identity {
    type = "SystemAssigned"
  }
}
```

## Example Usage with Data encryption

~> **NOTE:** The Key Vault must enable purge protection.

```hcl
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "workspace-example-ai"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_key_vault" "example" {
  name                     = "workspaceexamplekeyvault"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "premium"
  purge_protection_enabled = true
}
resource "azurerm_key_vault_access_policy" "example" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Get",
    "Delete",
    "Purge",
  ]
}

resource "azurerm_storage_account" "example" {
  name                     = "workspacestorageaccount"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_key_vault_key" "example" {
  name         = "workspaceexamplekeyvaultkey"
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

  depends_on = [azurerm_key_vault.example, azurerm_key_vault_access_policy.example]
}

resource "azurerm_machine_learning_workspace" "example" {
  name                    = "example-workspace"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  application_insights_id = azurerm_application_insights.example.id
  key_vault_id            = azurerm_key_vault.example.id
  storage_account_id      = azurerm_storage_account.example.id

  identity {
    type = "SystemAssigned"
  }

  encryption {
    key_vault_id = azurerm_key_vault.example.id
    key_id       = azurerm_key_vault_key.example.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Machine Learning Workspace. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group in which the Machine Learning Workspace should exist. Changing this forces a new resource to be created.

* `location` - (Optional) Specifies the supported Azure location where the Machine Learning Workspace should exist. Changing this forces a new resource to be created.

* `application_insights_id` - (Required) The ID of the Application Insights associated with this Machine Learning Workspace. Changing this forces a new resource to be created.

* `key_vault_id` - (Required) The ID of key vault associated with this Machine Learning Workspace. Changing this forces a new resource to be created.

* `storage_account_id` - (Required) The ID of the Storage Account associated with this Machine Learning Workspace. Changing this forces a new resource to be created.

-> **NOTE:** The `account_tier` cannot be `Premium` in order to associate the Storage Account to this Machine Learning Workspace.

* `identity` - (Required) An `identity` block as defined below.

* `container_registry_id` - (Optional) The ID of the container registry associated with this Machine Learning Workspace. Changing this forces a new resource to be created.

-> **NOTE:** The `admin_enabled` should be `true` in order to associate the Container Registry to this Machine Learning Workspace.

* `public_network_access_enabled` - (Optional) Enable public access when this Machine Learning Workspace is behind VNet.

* `image_build_compute_name` - (Optional) The compute name for image build of the Machine Learning Workspace.

* `description` - (Optional) The description of this Machine Learning Workspace.

* `friendly_name` - (Optional) Display name for this Machine Learning Workspace.

* `high_business_impact` - (Optional) Flag to signal High Business Impact (HBI) data in the workspace and reduce diagnostic data collected by the service

* `sku_name` - (Optional) SKU/edition of the Machine Learning Workspace, possible values are `Basic`. Defaults to `Basic`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `identity` block supports the following:

* `type` - (Required) The Type of Identity which should be used for this Azure Machine Learning workspace. At this time the only possible value is `SystemAssigned`.

---

An `encryption` block supports the following:

* `key_vault_id` - (Required) The ID of the keyVault where the customer owned encryption key is present.

* `key_id` - (Required) The Key Vault URI to access the encryption key.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Machine Learning Workspace.

* `discovery_url` - The url for the discovery service to identify regional endpoints for machine learning experimentation services.

---

An `identity` block exports the following:

* `principal_id` - The (Client) ID of the Service Principal.

* `tenant_id` - The ID of the Tenant the Service Principal is assigned in.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Machine Learning Workspace.
* `update` - (Defaults to 30 minutes) Used when updating the Machine Learning Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Machine Learning Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Machine Learning Workspace.

## Import

Machine Learning Workspace can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_machine_learning_workspace.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.MachineLearningServices/workspaces/workspace1
```
