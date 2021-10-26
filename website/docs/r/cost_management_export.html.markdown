---
subcategory: "Cost Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cost_management_export"
description: |-
  Manages an Azure Cost Management Export.
---

# azurerm_cost_management_export

Manages an Azure Cost Management Export.

## Example Usage

```hcl
data "azurerm_subscription" "current" {
}

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

resource "azurerm_cost_management_export" "example" {
  name                    = "example"
  resource_group_id       = data.azurerm_subscription.current.id
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

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cost Management Export. Changing this forces a new resource to be created.

* `scope` - (Required) The ID of the scope in which to export information from, such as for `Subscription`, `Resource group`, `Billing Account`, `Department` ... Changing this forces a new resource to be created.

~> **NOTE:** The ID format differs depending on the type of `scope`:
- for scope `Subscription`, the id format is `/subscriptions/{subscriptionId}/`
- for scope `Resource Group`, the id format is `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`
- for scope `Billing Accounts`, the id format is `/providers/Microsoft.Billing/billingAccounts/{billingAccountId}`
- for scope `Department`, the id format is `/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/departments/{departmentId}`
- for scope `Enrollment Account`, the id format is `/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/billingProfiles/{billingProfileId}`
- for scope `Management Group`, the id format is `/providers/Microsoft.Management/managementGroups/{managementGroupId}`
- for scope `Billing Profile`, the id format is `/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/billingProfiles/{billingProfileId}`
- for scope `Invoice Section`, the id format is `/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/billingProfiles/{billingProfileId}/invoiceSections/{invoiceSectionId}`
- for scope `Partner Customer`, the id format is `/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/customers/{customerId}`
...

* `recurrence_type` - (Required) How often the requested information will be exported. Accepted values are `Annually`, `Daily`, `Monthly`, or `Weekly`.

* `recurrence_period_start` - (Required) The date the export will start capturing information.

* `recurrence_period_end` - (Required) The date the export will stop capturing information. 

* `export_data_storage_location` - (Required) A `export_data_storage_location` block as defined below.

* `export_data_definition` - (Required) A `export_data_definition` block as defined below.

* `active` - (Optional) Is the cost management export active? Default is `true`.

---

A `export_data_storage_location` block supports the following:

* `storage_account_id` - (Required) The storage account id where exports will be delivered.

* `container_name` - (Required) The name of the container where exports will be uploaded.

* `root_folder_path` - (Required) The path of the directory where exports will be uploaded.

---

A `export_data_definition` block supports the following:

* `type` - (Required) The type of the query. Accepted values are `Usage`, `ActualCost`, or `AmortizedCost`.

* `time_frame` - (Required) The time frame for pulling data for the query. If custom, then a specific time period must be provided. Accepted values are `WeekToDate`, `MonthToDate`, `TheLastBillingMonth`, `BillingMonthToDate`, `TheLastMonth`, `MonthToDate` or `Custom`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Cost Management Export.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the resource.
* `update` - (Defaults to 30 minutes) Used when updating the resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the resource.

## Import

Cost Management Export for a Resource Group can be imported using the `resource id`, e.g.

### Management Group
```shell
terraform import azurerm_cost_management_export_resource_group.example /providers/Microsoft.Management/managementGroups/example/providers/Microsoft.CostManagement/exports/example
```

### Subscription
```shell
terraform import azurerm_cost_management_export_resource_group.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.CostManagement/exports/example
```

### Resource Group
```shell
terraform import azurerm_cost_management_export_resource_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.CostManagement/exports/example
```
