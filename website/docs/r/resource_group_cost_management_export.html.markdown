---
subcategory: "Cost Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_group_cost_management_export"
description: |-
  Manages an Azure Cost Management Export for a Resource Group.
---

# azurerm_resource_group_cost_management_export

Manages an Azure Cost Management Export for a Resource Group.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                = "example-storage-account"
  resource_group_name = azurerm_resource_group.example.name

  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_resource_group_cost_management_export" "example" {
  name = "example"
  resource_group_id       = azurerm_resource_group.example.id
  recurrence_type         = "Monthly"
  recurrence_period_start = "2020-08-18T00:00:00Z"
  recurrence_period_end   = "2020-09-18T00:00:00Z"

  export_data_storage_location {
    storage_account_id = azurerm_storage_account.example.id
    container_name     = "examplecontainer"
    root_folder_path   = "/root/updated"    
  }

  export_data_definition {
    type       = "Usage"
    time_frame = "WeekToDate"   
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cost Management Export. Changing this forces a new resource to be created.

* `resource_group_id` - (Required) The id of the resource group on which to create an export.

* `recurrence_type` - (Required) How often the requested information will be exported. Valid values include `Annually`, `Daily`, `Monthly`, `Weekly`.

* `recurrence_period_start` - (Required) The date the export will start capturing information.

* `recurrence_period_end` - (Required) The date the export will stop capturing information.

* `delivery_info` - (Required) A `delivery_info` block as defined below.

* `query` - (Required) A `query` block as defined below.

* `active` - (Optional) Is the cost management export active? Default is `true`.

---

A `export_data_storage_location` block supports the following:

* `storage_account_id` - (Required) The storage account id where exports will be delivered.

* `container_name` - (Required) The name of the container where exports will be uploaded.

* `root_folder_path` - (Required) The path of the directory where exports will be uploaded.

---

A `export_data_definition` block supports the following:

* `type` - (Required) The type of the query.

* `time_frame` - (Required) The time frame for pulling data for the query. If custom, then a specific time period must be provided. Possible values include: `WeekToDate`, `MonthToDate`, `YearToDate`, `TheLastWeek`, `TheLastMonth`, `TheLastYear`, `Custom`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Cost Management.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cost Management.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cost Management.
* `update` - (Defaults to 30 minutes) Used when updating the Cost Management.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cost Management.

## Import

Cost Management Export for a Resource Group can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_resource_group_cost_management_export.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.CostManagement/exports/export1
```