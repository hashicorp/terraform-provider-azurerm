---
subcategory: "Machine Learning"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_machine_learning_workspace"
description: |-
  Manages a Azure Machine Learning Workspace.
---
# azurerm_machine_learning_workspace

Manages a Azure Machine Learning Workspace

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

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Machine Learning Workspace. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group in which the Machine Learning Workspace should exist. Changing this forces a new resource to be created.

* `location` - (Optional) Specifies the supported Azure location where the Machine Learning Workspace should exist. Changing this forces a new resource to be created.

* `application_insights_id` - (Required) The ID of the Application Insights associated with this Machine Learning Workspace. Changing this forces a new resource to be created.

* `key_vault_id` - (Required) The ID of key vault associated with this Machine Learning Workspace. Changing this forces a new resource to be created.

* `storage_account_id` - (Required) The ID of the Storage Account associated with this Machine Learning Workspace. Changing this forces a new resource to be created.

-> **NOTE:** The `account_tier` cannot be `Premium` in order to associate the Storage Account to this Machine Learning Workspace.

* `identity` - (Required) An `identity` block defined below.

* `container_registry_id` - (Optional) The ID of the container registry associated with this Machine Learning Workspace. Changing this forces a new resource to be created.

-> **NOTE:** The `admin_enabled` should be `true` in order to associate the Container Registry to this Machine Learning Workspace.

* `description` - (Optional) The description of this Machine Learning Workspace.

* `discovery_url` - (Optional) The URL for the discovery service to identify regional endpoints for machine learning experimentation services.

* `friendly_name` - (Optional) Friendly name for this Machine Learning Workspace.

* `high_business_impact` - (Optional) Flag to signal High Business Impact (HBI) data in the workspace and reduce diagnostic data collected by the service

* `sku_name` - (Optional) SKU/edition of the Machine Learning Workspace, possible values are `Basic` for a basic workspace or `Enterprise` for a feature rich workspace. Defaults to `Basic`.

* `tags` - (Optional) A mapping of tags to assign to the resource. Changing this forces a new resource to be created.

---

An `identity` block supports the following:

* `type` - (Required) The Type of Identity which should be used for this Disk Encryption Set. At this time the only possible value is `SystemAssigned`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Machine Learning Workspace.

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
