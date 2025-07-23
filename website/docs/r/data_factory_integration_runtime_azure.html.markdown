---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_integration_runtime_azure"
description: |-
  Manages a Data Factory Azure Integration Runtime.
---

# azurerm_data_factory_integration_runtime_azure

Manages a Data Factory Azure Integration Runtime.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_data_factory_integration_runtime_azure" "example" {
  name            = "example"
  data_factory_id = azurerm_data_factory.example.id
  location        = azurerm_resource_group.example.location
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Managed Integration Runtime. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/azure/data-factory/naming-rules) for all restrictions.

* `data_factory_id` - (Required) The Data Factory ID in which to associate the Linked Service with. Changing this forces a new resource.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Use `AutoResolve` to create an auto-resolve integration runtime. Changing this forces a new resource to be created.

* `description` - (Optional) Integration runtime description.

* `cleanup_enabled` - (Optional) Cluster will not be recycled and it will be used in next data flow activity run until TTL (time to live) is reached if this is set as `false`. Defaults to `true`.

* `compute_type` - (Optional) Compute type of the cluster which will execute data flow job. Valid values are `General`, `ComputeOptimized` and `MemoryOptimized`. Defaults to `General`.

* `core_count` - (Optional) Core count of the cluster which will execute data flow job. Valid values are `8`, `16`, `32`, `48`, `80`, `144` and `272`. Defaults to `8`.

* `time_to_live_min` - (Optional) Time to live (in minutes) setting of the cluster which will execute data flow job. Defaults to `0`.

* `virtual_network_enabled` - (Optional) Is Integration Runtime compute provisioned within Managed Virtual Network? Changing this forces a new resource to be created.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Integration Runtime Azure.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Integration Runtime Azure.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Integration Runtime Azure.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Integration Runtime Azure.

## Import

Data Factory Azure Integration Runtimes can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_integration_runtime_azure.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/integrationruntimes/example
```
