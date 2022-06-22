---
subcategory: "Cost Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subscription_cost_management_export"
description: |-
  Manages an Azure Cost Management Export for a Subscription.
---

# azurerm_subscription_cost_management_export

Manages a Cost Management Export for a Subscription.

## Example Usage

```hcl
data "azurerm_subscription" "example" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name

  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                 = "examplecontainer"
  storage_account_name = azurerm_storage_account.example.name
}

resource "azurerm_subscription_cost_management_export" "example" {
  name                         = "example"
  subscription_id              = data.azurerm_subscription.example.id
  recurrence_type              = "Monthly"
  recurrence_period_start_date = "2020-08-18T00:00:00Z"
  recurrence_period_end_date   = "2020-09-18T00:00:00Z"

  export_data_storage_location {
    container_id     = azurerm_storage_container.example.resource_manager_id
    root_folder_path = "/root/updated"
  }

  export_data_options {
    type       = "Usage"
    time_frame = "WeekToDate"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cost Management Export. Changing this forces a new resource to be created.

* `subscription_id` - (Required) The id of the subscription on which to create an export.

* `recurrence_type` - (Required) How often the requested information will be exported. Valid values include `Annually`, `Daily`, `Monthly`, `Weekly`.

* `recurrence_period_start_date` - (Required) The date the export will start capturing information.

* `recurrence_period_end` - (Required) The date the export will stop capturing information.

* `export_data_storage_location` - (Required) A `export_data_storage_location` block as defined below.

* `export_data_options` - (Required) A `export_data_options` block as defined below.

* `active` - (Optional) Is the cost management export active? Default is `true`.

---

A `export_data_storage_location` block supports the following:

* `container_id` - (Required) The Resource Manager ID of the container where exports will be uploaded.

* `root_folder_path` - (Required) The path of the directory where exports will be uploaded.

**Note:** The Resource Manager ID of a Storage Container is exposed via the `resource_manager_id` attribute of the `azurerm_storage_container` resource.

---

A `export_data_options` block supports the following:

* `type` - (Required) The type of the query.

* `time_frame` - (Required) The time frame for pulling data for the query. If custom, then a specific time period must be provided. Possible values include: `WeekToDate`, `MonthToDate`, `BillingMonthToDate`, `TheLastWeek`, `TheLastMonth`, `TheLastBillingMonth`, `Custom`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Cost Management Export for this Subscription.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Subscription Cost Management Export.
* `read` - (Defaults to 5 minutes) Used when retrieving the Subscription Cost Management Export.
* `update` - (Defaults to 30 minutes) Used when updating the Subscription Cost Management Export.
* `delete` - (Defaults to 30 minutes) Used when deleting the Subscription Cost Management Export.

## Import

Subscription Cost Management Exports can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_subscription_cost_management_export.example /subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.CostManagement/exports/export1
```
