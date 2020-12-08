---
subcategory: "Time Series Insights"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iot_time_series_insights_gen2_environment"
description: |-
  Manages an Azure IoT Time Series Insights Gen2 Environment.
---

# azurerm_time_series_insights_gen2_environment

Manages an Azure IoT Time Series Insights Gen2 Environment.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "northeurope"
}
resource "azurerm_storage_account" "storage" {
  name                     = "example"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
resource "azurerm_iot_time_series_insights_gen2_environment" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "L1"
  data_retention_time = "P30D"
  property {
    ids = ["id"]
  }
  storage {
    name = azurerm_storage_account.storage.name
    key  = azurerm_storage_account.storage.primary_access_key
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Azure IoT Time Series Insights Gen2 Environment. Changing this forces a new resource to be created. Must be globally unique.

* `resource_group_name` - (Required) The name of the resource group in which to create the Azure IoT Time Series Insights Gen2 Environment.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku_name` - (Required) Specifies the SKU Name for this IoT Time Series Insights Gen2 Environment. Currently it supports only `L1`. For gen2, capacity cannot be specified.

* `data_retention_time` - (Required) Specifies the ISO8601 timespan specifying the minimum number of days the environment's events will be available for query. Changing this forces a new resource to be created.

* `property` - (Required) A `property` block as defined below.

* `storage` - (Required) A `storage` block as defined below.

* `id_properties` - (Required) A list of property ids for the Azure IoT Time Series Insights Gen2 Environment

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `storage` block supports the following:

* `name` - (Required) Name of storage account for Azure IoT Time Series Insights Gen2 Environment

* `key` - (Required) Access key of storage account for Azure IoT Time Series Insights Gen2 Environment


## Attributes Reference

* `id` - The ID of the IoT Time Series Insights Gen2 Environment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Time Series Insights Gen2 Environment.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Time Series Insights Gen2 Environment.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Time Series Insights Gen2 Environment.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Time Series Insights Gen2 Environment.

## Import

Azure IoT Time Series Insights Gen2 Environment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iot_time_series_insights_gen2_environment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.TimeSeriesInsights/environments/example
```
