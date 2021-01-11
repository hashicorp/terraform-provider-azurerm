---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_linked_service_synapse"
description: |-
  Manages a Linked Service (connection) between Synapse and Azure Data Factory.
---

# azurerm_data_factory_linked_service_synapse

Manages a Linked Service (connection) between Synapse and Azure Data Factory.

~> **Note:** All arguments including the client secret will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "northeurope"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_data_factory_linked_service_synapse" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  data_factory_name   = azurerm_data_factory.example.name

  connection_string = "Integrated Security=False;Data Source=test;Initial Catalog=test;User ID=test;Password=test"
}
```

## Example Usage with Password in Key Vault

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "northeurope"
}

resource "azurerm_key_vault" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_data_factory_linked_service_key_vault" "example" {
  name                = "kvlink"
  resource_group_name = azurerm_resource_group.example.name
  data_factory_name   = azurerm_data_factory.example.name
  key_vault_id        = azurerm_key_vault.example.id
}

resource "azurerm_data_factory_linked_service_synapse" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  data_factory_name   = azurerm_data_factory.example.name

  connection_string = "Integrated Security=False;Data Source=test;Initial Catalog=test;User ID=test;"
  key_vault_password {
    linked_service_name = azurerm_data_factory_linked_service_key_vault.example.name
    secret_name         = "secret"
  }
}
```

## Argument Reference

The following supported arguments are common across all Azure Data Factory Linked Services:

* `name` - (Required) Specifies the name of the Data Factory Linked Service Synapse. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/data-factory/naming-rules) for all restrictions.

* `resource_group_name` - (Required) The name of the resource group in which to create the Data Factory Linked Service Synapse. Changing this forces a new resource

* `data_factory_name` - (Required) The Data Factory name in which to associate the Linked Service with. Changing this forces a new resource.

* `description` - (Optional) The description for the Data Factory Linked Service Synapse.

* `integration_runtime_name` - (Optional) The integration runtime reference to associate with the Data Factory Linked Service Synapse.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Linked Service Synapse.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Linked Service Synapse.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Linked Service Synapse.

The following supported arguments are specific to Data Factory Synapse Linked Service:

* `connection_string` - (Required) The connection string in which to authenticate with the Synapse.

* `key_vault_password` - (Optional) A `key_vault_password` block as defined below. Use this argument to store Synapse password in an existing Key Vault. It needs an existing Key Vault Data Factory Linked Service.

---

A `key_vault_password` block supports the following:

* `linked_service_name` - (Required) Specifies the name of an existing Key Vault Data Factory Linked Service.

* `secret_name` - (Required) Specifies the secret name in Azure Key Vault that stores Synapse password.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Data Factory Synapse Linked Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Synapse Linked Service.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Synapse Linked Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Synapse Linked Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Synapse Linked Service.

## Import

Data Factory Synapse Linked Service's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_linked_service_synapse.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/linkedservices/example
```
