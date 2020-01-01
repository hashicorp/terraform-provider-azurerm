---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_integration_runtime_managed"
sidebar_current: "docs-azurerm_data_factory_integration_runtime_managed"
description: |-
  Manages an Azure Data Factory Managed Integration Runtime.
---

# azurerm_data_factory_integration_runtime_managed

Manages an Azure Data Factory Managed Integration Runtime.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "northeurope"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_data_factory_integration_runtime_managed" "example" {
  name                = "example"
  data_factory_name   = "${azurerm_data_factory.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"

  node_size = "Standard_D8_v3"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Managed Integration Runtime. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/data-factory/naming-rules) for all restrictions.

* `data_factory_name` - (Required) Specifies the name of the Data Factory the Managed Integration Runtime belongs to. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Managed Integration Runtime. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `node_size` - (Required) The size of the nodes on which the Managed Integration Runtime runs. Valid values are: `Standard_D2_v3`, `Standard_D4_v3`, `Standard_D8_v3`, `Standard_D16_v3`, `Standard_D32_v3`, `Standard_D64_v3`, `Standard_E2_v3`, `Standard_E4_v3`, `Standard_E8_v3`, `Standard_E16_v3`, `Standard_E32_v3`, `Standard_E64_v3`, `Standard_D1_v2`, `Standard_D2_v2`, `Standard_D3_v2`, `Standard_D4_v2`, `Standard_A4_v2` and `Standard_A8_v2`

* `number_of_nodes` - (Optional) Number of nodes for the Managed Integration Runtime. Max is `10`. Defaults to `1`.

* `max_parallel_executions_per_node` - (Optional) Defines the maximum parallel executions per node. Defaults to `1`. Max is `16`.

* `edition` - (Optional) The Managed Integration Runtime edition. Valid values are `Standard` and `Enterprise`. Defaults to `Standard`.

* `license_type` - (Optional) The type of the license that is used. Valid values are `LicenseIncluded` and `BasePrize`. Defaults to `LicenseIncluded`.

* `catalog_info` - (Optional) A `catalog_info` block as defined below.

* `custom_setup_script` - (Optional) A `custom_setup_script` block as defined below.

* `vnet_integration` - (Optional) A `vnet_integration` block as defined below.

---

A `catalog_info` block supports the following:

* `server_endpoint` - (Required) The endpoint of an Azure SQL Server that will be used to host the SSIS catalog.

* `administrator_login` - (Required) Administrator login name for the SQL Server.

* `administrator_password` - (Required) Administrator login password for the SQL Server.

* `pricing_tier` - (Required) Pricing tier for the database that will be created for the SSIS catalog. Valid values are: `Basic`, `Standard`, `Premium` and `PremiumRS`.

---

A `custom_setup_script` block supports the following:

* `blob_container_uri` - (Required) The blob endpoint for the container which contains a custom setup script that will be run on every node on startup. See [https://docs.microsoft.com/en-us/azure/data-factory/how-to-configure-azure-ssis-ir-custom-setup](https://docs.microsoft.com/en-us/azure/data-factory/how-to-configure-azure-ssis-ir-custom-setup) for more information.

* `sas_token` - (Required) A container SAS token that gives access to the files. See [https://docs.microsoft.com/en-us/azure/data-factory/how-to-configure-azure-ssis-ir-custom-setup](https://docs.microsoft.com/en-us/azure/data-factory/how-to-configure-azure-ssis-ir-custom-setup) for more information.

---

A `vnet_integration` block supports the following:

* `vnet_id` - (Required) ID of the virtual network to which the nodes of the Managed Integration Runtime will be added.

* `subnet_name` - (Required) Name of the subnet to which the nodes of the Managed Integration Runtime will be added.


## Attributes Reference

The following attributes are exported:

* `id` - The Managed Integration Runtime ID.


## Import

Managed Integration Runtime can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/integrationruntimes/example
```
