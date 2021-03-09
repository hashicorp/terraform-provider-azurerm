---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_linked_service_key_vault"
description: |-
  Manages a Linked Service (connection) between Key Vault and Azure Data Factory.
---

# azurerm_data_factory_linked_service_key_vault

Manages a Linked Service (connection) between Key Vault and Azure Data Factory.

## Example Usage

```hcl
data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
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
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  data_factory_name   = azurerm_data_factory.example.name
  key_vault_id        = azurerm_key_vault.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory Linked Service Key Vault. Changing this forces a new resource to be created. Must be unique within a data
  factory. See the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/data-factory/naming-rules) for all restrictions.

* `resource_group_name` - (Required) The name of the resource group in which to create the Data Factory Linked Service Key Vault. Changing this forces a new resource

* `data_factory_name` - (Required) The Data Factory name in which to associate the Linked Service with. Changing this forces a new resource.

* `key_vault_id` - (Required) The ID the Azure Key Vault resource.

* `description` - (Optional) The description for the Data Factory Linked Service Key Vault.

* `integration_runtime_name` - (Optional) The integration runtime reference to associate with the Data Factory Linked Service Key Vault.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Linked Service Key Vault.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Linked Service Key Vault.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Linked Service Key Vault.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Data Factory Key Vault Linked Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Key Vault Linked Service.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Key Vault Linked Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Key Vault Linked Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Key Vault Linked Service.

## Import

Data Factory Key Vault Linked Service's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_linked_service_key_vault.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/linkedservices/example
```
