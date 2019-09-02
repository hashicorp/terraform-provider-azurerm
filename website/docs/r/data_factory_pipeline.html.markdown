---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_pipeline"
sidebar_current: "docs-azurerm-resource-data-factory-pipeline"
description: |-
  Manages a Pipeline inside a Azure Data Factory.
---

# azurerm_data_factory_pipeline

Manages a Pipeline inside a Azure Data Factory.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "northeurope"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_data_factory_pipeline" "example" {
  name                = "example"
  resource_group_name = "${azurerm_resource_group.example.name}"
  data_factory_name   = "${azurerm_data_factory.example.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory Pipeline. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/data-factory/naming-rules) for all restrictions.

* `resource_group_name` - (Required) The name of the resource group in which to create the Data Factory Pipeline. Changing this forces a new resource

* `data_factory_name` - (Required) The Data Factory name in which to associate the Pipeline with. Changing this forces a new resource.

* `description` - (Optional) The description for the Data Factory Pipeline.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Pipeline.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Pipeline.

* `variables` - (Optional) A map of variables to associate with the Data Factory Pipeline.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Data Factory Pipeline.

## Import

Data Factory Pipeline can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_pipeline.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/pipelines/example
```
