---
subcategory: "Data Lake"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_lake_store_virtual_network_rule"
description: |-
  Manages a Azure Data Lake Store Virtual Network Rule.
---

# azurerm_data_lake_store_virtual_network_rule

Allows you to add, update, or remove an Azure Data Lake Store to a subnet of a virtual network.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-data-lake-store-vnet-rule"
  location = "northeurope"
}

resource "azurerm_virtual_network" "vnet" {
  name                = "example-vnet"
  address_space       = ["10.7.29.0/29"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "subnet" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = ["10.7.29.0/29"]
  service_endpoints    = ["Microsoft.AzureActiveDirectory"]
}

resource "azurerm_data_lake_store" "example" {
  name                = "exampledatalake"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_data_lake_store_virtual_network_rule" "adlsvnetrule" {
  name                = "adls-vnet-rule"
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_data_lake_store.example.name
  subnet_id           = azurerm_subnet.subnet.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Data Lake Store virtual network rule. Changing this forces a new resource to be created. Cannot be empty and must only contain alphanumeric characters, underscores, periods and hyphens. Cannot start with a period, underscore or hyphen, and cannot end with a period and a hyphen.

~> **NOTE:** `name` must be between 2-64 characters long and must satisfy all of the requirements below:

1. Contains only alphanumeric characters, periods, underscores or hyphens
2. Cannot start with a period, underscore or hyphen
3. Cannot end with a period and a hyphen

* `resource_group_name` - (Required) The name of the resource group where the Data Lake Store resides. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the Data Lake Store to which this Data Lake Store virtual network rule will be applied to. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the subnet that the Data Lake Store will be connected to.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Data Lake Store virtual network rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Lake Store Virtual Network Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Lake Store Virtual Network Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Data Lake Store Virtual Network Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Lake Store Virtual Network Rule.

## Import

Data Lake Store Virtual Network Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_lake_store_virtual_network_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.DataLakeStore/accounts/myaccount/virtualNetworkRules/vnetrulename
```
