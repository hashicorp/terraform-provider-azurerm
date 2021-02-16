---
subcategory: "Time Series Insights"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iot_time_series_insights_standard_environment"
description: |-
  Manages an Azure IoT Time Series Insights Reference Data Set.
---

# azurerm_iot_time_series_insights_reference_data_set

Manages an Azure IoT Time Series Insights Reference Data Set.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "northeurope"
}
resource "azurerm_iot_time_series_insights_standard_environment" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "S1_1"
  data_retention_time = "P30D"
}
resource "azurerm_iot_time_series_insights_reference_data_set" "example" {
  name                                = "example"
  time_series_insights_environment_id = azurerm_iot_time_series_insights_standard_environment.example.id
  location                            = azurerm_resource_group.example.location

  key_property {
    name = "keyProperty1"
    type = "String"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Azure IoT Time Series Insights Reference Data Set. Changing this forces a new resource to be created. Must be globally unique.

* `time_series_insights_environment_id` - (Required) The resource ID of the Azure IoT Time Series Insights Environment in which to create the Azure IoT Time Series Insights Reference Data Set. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `data_string_comparison_behavior` - (Optional) The comparison behavior that will be used to compare keys. Valid values include `Ordinal` and `OrdinalIgnoreCase`. Defaults to `Ordinal`. Changing this forces a new resource to be created.

* `key_property` - (Optional) A `key_property` block as defined below. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `key_property` block supports the following:

* `name`- (Required) The name of the key property. Changing this forces a new resource to be created.

* `type` - (Required) The data type of the key property. Valid values include `Bool`, `DateTime`, `Double`, `String`. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoT Time Series Insights Reference Data Set.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Time Series Insights Reference Data Set.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Time Series Insights Reference Data Set.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Time Series Insights Reference Data Set.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Time Series Insights Reference Data Set.

## Import

Azure IoT Time Series Insights Reference Data Set can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iot_time_series_insights_reference_data_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.TimeSeriesInsights/environments/example/referenceDataSets/example
```
