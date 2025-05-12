---
subcategory: "Machine Learning"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint"
description: |-
  Manages an Azure Machine Learning Workspace Network Outbound Rule Private Endpoint.
---

# azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint

Manages an Azure Machine Learning Workspace Network Outbound Rule Private Endpoint.

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

  managed_network {
    isolation_mode = "AllowOnlyApprovedOutbound"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_storage_account" "example2" {
  name                     = "example-sa"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint" "example" {
  name                = "example-outboundrule"
  workspace_id        = azurerm_machine_learning_workspace.example.id
  service_resource_id = azurerm_storage_account.example2.id
  sub_resource_target = "blob"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Machine Learning Workspace Network Outbound Rule Private Endpoint. Changing this forces a new resource to be created.

* `workspace_id` - (Required) Specifies the ID of the Machine Learning Workspace. Changing this forces a new resource to be created.

* `service_resource_id` - (Required) Specifies the Service Resource ID to connect. Changing this forces a new resource to be created.

~> **Note:** Supported service resources: **Key Vault**, **Storage Account**, **Machine Learning Workspace**, **Redis**.

* `sub_resource_target` - (Required) Specifies the Sub Resource of the service resource to connect to. Possible values are `vault`,`amlworkspace`,`blob`,`table`,`queue`,`file`,`web`,`dfs`, `redisCache`. Changing this forces a new resource to be created.
  
  | Service                    | Sub Resource Type                         |
  |----------------------------|-------------------------------------------|
  | Machine Learning Workspace | `amlworkspace`                            |
  | Redis                      | `redisCache`                              |
  | Storage Account            | `blob`,`table`,`queue`,`file`,`web`,`dfs` |
  | Key Vault                  | `vault`                                   |

* `spark_enabled` - (Optional) Whether to enable an additional private endpoint to be used by jobs running on Spark. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Machine Learning Workspace Network Outbound Rule Private Endpoint.

### Timeouts

The `timeouts` block allows you to
specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Machine Learning Workspace Network Outbound Rule Private Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Machine Learning Workspace Network Outbound Rule Private Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Machine Learning Workspace Network Outbound Rule Private Endpoint.

## Import

Machine Learning Workspace Network Outbound Rule Private Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/outboundRules/rule1
```
