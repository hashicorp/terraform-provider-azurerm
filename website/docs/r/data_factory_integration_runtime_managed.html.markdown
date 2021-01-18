---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_integration_runtime_managed"
description: |-
  Manages an Azure Data Factory Managed Integration Runtime.
---

# azurerm_data_factory_integration_runtime_managed_ssis

Manages an Azure Data Factory Managed Integration Runtime.

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

resource "azurerm_data_factory_integration_runtime_managed" "example" {
  name                = "example"
  data_factory_name   = azurerm_data_factory.example.name
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Managed Integration Runtime. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/data-factory/naming-rules) for all restrictions.

* `data_factory_name` - (Required) Specifies the name of the Data Factory the Managed Integration Runtime belongs to. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Managed Integration Runtime. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `compute_type` - (Optional) TODO. Valid values are `General`, `ComputeOptimized` and `MemoryOptimized`. Defaults to `General`.

* `core_count` - (Optional) TODO. Valid values are `8`, `16`, `32`, `48`, `80`, `144` and `272`. Defaults to `8`.

* `time_to_live` - (Optional) TODO. Defaults to `0`.

* `description` - (Optional) Integration runtime description.

---

A `custom_setup_script` block supports the following:

* `blob_container_uri` - (Required) The blob endpoint for the container which contains a custom setup script that will be run on every node on startup. See [https://docs.microsoft.com/en-us/azure/data-factory/how-to-configure-azure-ssis-ir-custom-setup](https://docs.microsoft.com/en-us/azure/data-factory/how-to-configure-azure-ssis-ir-custom-setup) for more information.

* `sas_token` - (Required) A container SAS token that gives access to the files. See [https://docs.microsoft.com/en-us/azure/data-factory/how-to-configure-azure-ssis-ir-custom-setup](https://docs.microsoft.com/en-us/azure/data-factory/how-to-configure-azure-ssis-ir-custom-setup) for more information.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Integration Managed Runtime.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Integration Managed Runtime.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Integration Managed Runtime.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Integration Managed Runtime.

## Import

Data Factory Integration Managed Runtimes can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_integration_runtime_managed.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/integrationruntimes/example
```
