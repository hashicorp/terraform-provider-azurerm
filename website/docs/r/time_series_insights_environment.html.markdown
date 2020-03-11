---
subcategory: "Time Series Insights
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_time_series_insights_environment"
description: |-
  Manages an Azure Time Series Insights Environment.
---

# azurerm_time_series_insights_environment

Manages an Azure Time Series Insights Environment.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "northeurope"
}
resource "azurerm_time_series_insights_environment" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Azure Time Series Insights Environment. Changing this forces a new resource to be created. Must be globally unique.

* `resource_group_name` - (Required) The name of the resource group in which to create the Azure Time Series Insights Environment.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Azure Time Series Insights Environment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spatial Anchors Account.
* `update` - (Defaults to 30 minutes) Used when updating the Spatial Anchors Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spatial Anchors Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spatial Anchors Account.

## Import

Azure Time Series Insights Environment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_time_series_environment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.MixedReality/spatialAnchorsAccounts/example
```