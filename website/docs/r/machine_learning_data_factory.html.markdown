---
subcategory: "Machine Learning"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_machine_learning_data_factory"
description: |-
  Manages a Machine Learning Data Factory.
---

# azurerm_machine_learning_data_factory

Manages a Machine Learning Data Factory.

## Example Usage

```hcl

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "west europe"
  tags = {
    "stage" = "example"
  }
}

resource "azurerm_application_insights" "example" {
  name                = "example-ai"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_key_vault" "example" {
  name                = "example-kv"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  purge_protection_enabled = true
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_machine_learning_workspace" "example" {
  name                    = "example-mlw"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  application_insights_id = azurerm_application_insights.example.id
  key_vault_id            = azurerm_key_vault.example.id
  storage_account_id      = azurerm_storage_account.example.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_machine_learning_data_factory" "example" {
  name                          = "example"
  location                      = azurerm_resource_group.example.location
  machine_learning_workspace_id = azurerm_machine_learning_workspace.example.id
  data_factory_id               = azurerm_data_factory.example.id
  identity {
    type = "SystemAssigned"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Machine Learning Data Factory. Changing this forces a new Machine Learning Data Factory to be created.
  
* `machine_learning_workspace_id` - (Required) The ID of the Machine Learning Workspace. Changing this forces a new Machine Learning Data Factory to be created.
  
* `location` - (Required) The Azure Region where the Machine Learning Data Factory should exist. Changing this forces a new Machine Learning Data Factory to be created.

* `data_factory_id` - (Required) The ID of the Data Factory. Changing this forces a new Machine Learning Data Factory to be created.

---

* `description` - (Optional) The description of the Machine Learning Data Factory. Changing this forces a new Machine Learning Data Factory to be created.

* `identity` - (Optional) A `identity` block as defined below. Changing this forces a new Machine Learning Data Factory to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Machine Learning Data Factory. Changing this forces a new Machine Learning Data Factory to be created.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on the Machine Learning Data Factory. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned,UserAssigned` (to enable both).

* `identity_ids` - (Optional) A list of IDs for User Assigned Managed Identity resources to be assigned.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Machine Learning Data Factory.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this Machine Learning Data Factory.

---

The `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this Machine Learning Data Factory.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this Machine Learning Data Factory.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Machine Learning Data Factory.
* `read` - (Defaults to 5 minutes) Used when retrieving the Machine Learning Data Factory.
* `delete` - (Defaults to 30 minutes) Used when deleting the Machine Learning Data Factory.

## Import

Machine Learning Data Factory can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_machine_learning_data_factory.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/computes/compute1
```
