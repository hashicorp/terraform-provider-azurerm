---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_service_plan"
description: |-
  Manages an App Service: Service Plan.
---

# azurerm_service_plan

Manages an App Service: Service Plan.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_service_plan" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  os_type             = "Linux"
  sku_name            = "P1v2"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Service Plan. Changing this forces a new Service Plan to be created.

* `location` - (Required) The Azure Region where the Service Plan should exist. Changing this forces a new Service Plan to be created.

* `os_type` - (Required) The O/S type for the App Services to be hosted in this plan. Possible values include `Windows`, `Linux`, and `WindowsContainer`. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Service Plan should exist. Changing this forces a new Service Plan to be created.

* `sku_name` - (Required) The SKU for the plan. Possible values include `B1`, `B2`, `B3`, `D1`, `F1`, `I1`, `I2`, `I3`, `I1v2`, `I1mv2`, `I2v2`, `I2mv2`, `I3v2`, `I3mv2`, `I4v2`, `I4mv2`, `I5v2`, `I5mv2`, `I6v2`, `P1v2`, `P2v2`, `P3v2`, `P0v3`, `P1v3`, `P2v3`, `P3v3`, `P1mv3`, `P2mv3`, `P3mv3`, `P4mv3`, `P5mv3`, `S1`, `S2`, `S3`, `SHARED`, `EP1`, `EP2`, `EP3`, `FC1`, `WS1`, `WS2`, `WS3`, and `Y1`.

~> **Note:** Isolated SKUs (`I1`, `I2`, `I3`, `I1v2`, `I1mv2`, `I2v2`, `I2mv2`, `I3v2`, `I3mv2`) can only be used with App Service Environments

~> **Note:** Elastic and Consumption SKUs (`Y1`, `FC1`, `EP1`, `EP2`, and `EP3`) are for use with Function Apps.

---

* `app_service_environment_id` - (Optional) The ID of the App Service Environment to create this Service Plan in.

~> **Note:** Requires an Isolated SKU for `azurerm_app_service_environment_v3`, supported values include `I1v2`, `I1mv2`, `I2v2`, `I2mv2`, `I3v2`, `I3mv2`, `I4v2`, `I4mv2`, `I5v2`, `I5mv2`, and `I6v2`.

* `premium_plan_auto_scale_enabled` - (Optional) Should automatic scaling be enabled for the Premium SKU Plan. Defaults to `false`. Cannot be set unless using a Premium SKU.

* `maximum_elastic_worker_count` - (Optional) The maximum number of workers to use in an Elastic SKU Plan or Premium Plan that have `premium_plan_auto_scale_enabled` set to `true`. Cannot be set unless using an Elastic or Premium SKU.

* `worker_count` - (Optional) The number of Workers (instances) to be allocated.

* `per_site_scaling_enabled` - (Optional) Should Per Site Scaling be enabled. Defaults to `false`.

* `zone_balancing_enabled` - (Optional) Should the Service Plan balance across Availability Zones in the region. Changing this forces a new resource to be created.

~> **Note:** If this setting is set to `true` and the `worker_count` value is specified, it should be set to a multiple of the number of availability zones in the region. Please see the Azure documentation for the number of Availability Zones in your region.

* `tags` - (Optional) A mapping of tags which should be assigned to the AppService.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Service Plan.

* `kind` - A string representing the Kind of Service Plan.

* `reserved` - Whether this is a reserved Service Plan Type. `true` if `os_type` is `Linux`, otherwise `false`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Service Plan.
* `read` - (Defaults to 5 minutes) Used when retrieving the Service Plan.
* `update` - (Defaults to 1 hour) Used when updating the Service Plan.
* `delete` - (Defaults to 1 hour) Used when deleting the Service Plan.

## Import

AppServices can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_service_plan.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/serverFarms/farm1
```
