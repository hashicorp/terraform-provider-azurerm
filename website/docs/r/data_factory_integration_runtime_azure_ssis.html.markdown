---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_integration_runtime_azure_ssis"
description: |-
  Manages a Data Factory Azure-SSIS Integration Runtime.
---

# azurerm_data_factory_integration_runtime_azure_ssis

Manages a Data Factory Azure-SSIS Integration Runtime.

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

resource "azurerm_data_factory_integration_runtime_azure_ssis" "example" {
  name                = "example"
  data_factory_name   = azurerm_data_factory.example.name
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  node_size = "Standard_D8_v3"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Azure-SSIS Integration Runtime. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/data-factory/naming-rules) for all restrictions.

* `data_factory_name` - (Required) Specifies the name of the Data Factory the Azure-SSIS Integration Runtime belongs to. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Azure-SSIS Integration Runtime. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `node_size` - (Required) The size of the nodes on which the Azure-SSIS Integration Runtime runs. Valid values are: `Standard_D2_v3`, `Standard_D4_v3`, `Standard_D8_v3`, `Standard_D16_v3`, `Standard_D32_v3`, `Standard_D64_v3`, `Standard_E2_v3`, `Standard_E4_v3`, `Standard_E8_v3`, `Standard_E16_v3`, `Standard_E32_v3`, `Standard_E64_v3`, `Standard_D1_v2`, `Standard_D2_v2`, `Standard_D3_v2`, `Standard_D4_v2`, `Standard_A4_v2` and `Standard_A8_v2`

* `number_of_nodes` - (Optional) Number of nodes for the Azure-SSIS Integration Runtime. Max is `10`. Defaults to `1`.

* `max_parallel_executions_per_node` - (Optional) Defines the maximum parallel executions per node. Defaults to `1`. Max is `16`.

* `edition` - (Optional) The Azure-SSIS Integration Runtime edition. Valid values are `Standard` and `Enterprise`. Defaults to `Standard`.

* `license_type` - (Optional) The type of the license that is used. Valid values are `LicenseIncluded` and `BasePrice`. Defaults to `LicenseIncluded`.

* `catalog_info` - (Optional) A `catalog_info` block as defined below.

* `custom_setup_script` - (Optional) A `custom_setup_script` block as defined below.

* `express_custom_setup` - (Optional) An `express_custom_setup` block as defined below.

* `package_store` - (Optional) One or more `package_store` block as defined below.
  
* `proxy` - (Optional) A `proxy` block as defined below.

* `vnet_integration` - (Optional) A `vnet_integration` block as defined below.

* `description` - (Optional) Integration runtime description.

---

A `catalog_info` block supports the following:

* `server_endpoint` - (Required) The endpoint of an Azure SQL Server that will be used to host the SSIS catalog.

* `administrator_login` - (Optional) Administrator login name for the SQL Server.

* `administrator_password` - (Optional) Administrator login password for the SQL Server.

* `pricing_tier` - (Optional) Pricing tier for the database that will be created for the SSIS catalog. Valid values are: `Basic`, `Standard`, `Premium` and `PremiumRS`.

* `dual_standby_pair_name` - (Optional) The dual standby Azure-SSIS Integration Runtime pair with SSISDB failover.

---

A `custom_setup_script` block supports the following:

* `blob_container_uri` - (Required) The blob endpoint for the container which contains a custom setup script that will be run on every node on startup. See [https://docs.microsoft.com/en-us/azure/data-factory/how-to-configure-azure-ssis-ir-custom-setup](https://docs.microsoft.com/en-us/azure/data-factory/how-to-configure-azure-ssis-ir-custom-setup) for more information.

* `sas_token` - (Required) A container SAS token that gives access to the files. See [https://docs.microsoft.com/en-us/azure/data-factory/how-to-configure-azure-ssis-ir-custom-setup](https://docs.microsoft.com/en-us/azure/data-factory/how-to-configure-azure-ssis-ir-custom-setup) for more information.

---

An `express_custom_setup` block supports the following:

* `command_key` - (Optional) One or more `command_key` blocks as defined below.

* `component` - (Optional) One or more `component` blocks as defined below.

* `environment` - (Optional) The Environment Variables for the Azure-SSIS Integration Runtime.

* `powershell_version` - (Optional) The version of Azure Powershell installed for the Azure-SSIS Integration Runtime.

~> **NOTE** At least one of `env`, `powershell_version`, `component` and `command_key` should be specified.

---

A `command_key` block supports the following:

* `target_name` - (Required) The target computer or domain name.

* `user_name` - (Required) The username for the target device.

* `password` - (Required) The password for the target device.

---

A `component` block supports the following:

* `name` - (Required) The Component Name installed for the Azure-SSIS Integration Runtime.

* `license` - (Optional) The license used for the Component.

---

A `package_store` block supports the following:

* `name` - (Required) Name of the package store.

* `linked_service_name` - (Required) Name of the Linked Service to associate with the packages.

---

A `proxy` block supports the following:

* `self_hosted_integration_runtime_name` - (Required) Name of Self Hosted Integration Runtime as a proxy.

* `staging_storage_linked_service_name` - (Required)  Name of Azure Blob Storage linked service to reference the staging data store to be used when moving data between self-hosted and Azure-SSIS integration runtimes.

* `path` - (Optional) The path in the data store to be used when moving data between Self-Hosted and Azure-SSIS Integration Runtimes.

---

A `vnet_integration` block supports the following:

* `vnet_id` - (Required) ID of the virtual network to which the nodes of the Azure-SSIS Integration Runtime will be added.

* `subnet_name` - (Required) Name of the subnet to which the nodes of the Azure-SSIS Integration Runtime will be added.

* `public_ips` - (Optional) Static public IP addresses for the Azure-SSIS Integration Runtime. The size must be 2.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Data Factory Azure-SSIS Integration Runtime.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Azure-SSIS Integration Runtime.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Azure-SSIS Integration Runtime.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Azure-SSIS Integration Runtime.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Azure-SSIS Integration Runtime.

## Import

Data Factory Azure-SSIS Integration Runtimes can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_integration_runtime_azure_ssis.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/integrationruntimes/example
```
