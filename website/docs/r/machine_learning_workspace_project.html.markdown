---
subcategory: "Machine Learning"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_machine_learning_workspace_hub"
description: |-
  Manages a Azure Machine Learning Workspace Hub.
---
# azurerm_machine_learning_workspace_project

Manages an Azure Machine Learning Workspace Project 

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

resource "azurerm_machine_learning_workspace_hub" "example" {
  name                    = "example-workspace-hub"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  application_insights_id = azurerm_application_insights.example.id
  key_vault_id            = azurerm_key_vault.example.id
  storage_account_id      = azurerm_storage_account.example.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_machine_learning_workspace_project" "example" {
  name                = "example-workspace-project"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  workspace_hub_id    = azurerm_machine_learning_workspace_hub.example.id
  identity {
    type = "SystemAssigned"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Machine Learning Workspace Project. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group in which the Machine Learning Workspace Project should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the Machine Learning Workspace Project should exist. Changing this forces a new resource to be created.

* `workspace_hub_id` - (Required) The ID of the machine learning workspace Hub where the Machine Learning Workspace Project should exist. Changing this forces a new resource to be created.

* `identity` - (Required) An `identity` block as defined below. Changing this forces a new resource to be created.

* `friendly_name` - (Optional) Display name for this Machine Learning Workspace Project.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Machine Learning Workspace Project. Possible values are `SystemAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Machine Learning Workspace Project.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Machine Learning Workspace Project.

* `workspace_id` - The immutable id associated with this Hub.
 
* `public_network_access` -  Whether or not public endpoint access is allowed for this Machine Learning Workspace Project.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Machine Learning Workspace Project.
* `update` - (Defaults to 30 minutes) Used when updating the Machine Learning Workspace Project.
* `read` - (Defaults to 5 minutes) Used when retrieving the Machine Learning Workspace Project.
* `delete` - (Defaults to 30 minutes) Used when deleting the Machine Learning Workspace Project.

## Import

Machine Learning Workspace Project can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_machine_learning_workspace_project.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.MachineLearningServices/workspaces/workspace1
```
