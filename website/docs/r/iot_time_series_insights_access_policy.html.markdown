---
subcategory: "Time Series Insights"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iot_time_series_insights_access_policy"
description: |-
  Manages an Azure IoT Time Series Insights Access Policy.
---

# azurerm_iot_time_series_insights_access_policy

Manages an Azure IoT Time Series Insights Access Policy.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_iot_time_series_insights_standard_environment" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "S1_1"
  data_retention_time = "P30D"
}
resource "azurerm_iot_time_series_insights_access_policy" "example" {
  name                                = "example"
  time_series_insights_environment_id = azurerm_iot_time_series_insights_standard_environment.example.name

  principal_object_id = "aGUID"
  roles               = ["Reader"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Azure IoT Time Series Insights Access Policy. Changing this forces a new resource to be created. Must be globally unique.

* `time_series_insights_environment_id` - (Required) The resource ID of the Azure IoT Time Series Insights Environment in which to create the Azure IoT Time Series Insights Reference Data Set. Changing this forces a new resource to be created.

* `principal_object_id` - (Required) The id of the principal in Azure Active Directory. Changing this forces a new resource to be created.

* `roles` - (Required) A list of roles to apply to the Access Policy. Valid values include `Contributor` and `Reader`.

* `description` - (Optional) The description of the Azure IoT Time Series Insights Access Policy.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the IoT Time Series Insights Access Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Time Series Insights Access Policy.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Time Series Insights Access Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Time Series Insights Access Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Time Series Insights Access Policy.

## Import

Azure IoT Time Series Insights Access Policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iot_time_series_insights_access_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.TimeSeriesInsights/environments/environment1/accessPolicies/example
```
