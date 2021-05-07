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

  delivery_info {
    storage_account_id = azurerm_storage_account.example.id
    container_name     = "examplecontainer"
    root_folder_path   = "/root/updated"
  }

  query {
    type       = "Usage"
    time_frame = "WeekToDate"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cost Management Export. Changing this forces a new resource to be created.

* `scope` - (Required) The ID of the scope in which to export information from. This includes `'/subscriptions/{subscriptionId}/'` for subscription scope, '`/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`' for resourceGroup scope, `'/providers/Microsoft.Billing/billingAccounts/{billingAccountId}'` for Billing Account scope and `'/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/departments/{departmentId}'` for Department scope, `'/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/enrollmentAccounts/{enrollmentAccountId}'` for EnrollmentAccount scope, `'/providers/Microsoft.Management/managementGroups/{managementGroupId}` for Management Group scope, `'/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/billingProfiles/{billingProfileId}'` for billingProfile scope, `'/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/billingProfiles/{billingProfileId}/invoiceSections/{invoiceSectionId}'` for invoiceSection scope, and `'/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/customers/{customerId}'` specific for partners.

* `recurrence_type` - (Required) How often the requested information will be exported. Valid values include `Annually`, `Daily`, `Monthly`, `Weekly`.

* `recurrence_period_start` - (Required) The date the export will start capturing information.

* `recurrence_period_end` - (Required) The date the export will stop capturing information. 

* `delivery_info` - (Required) A `delivery_info` block as defined below.

* `query` - (Required) A `query` block as defined below.

* `active` - (Optional) Is the cost management export active? Default is `true`.

---

A `delivery_info` block supports the following:

* `storage_account_id` - (Required) The storage account id where exports will be delivered.

* `container_name` - (Required) The name of the container where exports will be uploaded.

* `root_folder_path` - (Required) The path of the directory where exports will be uploaded.

---

A `query` block supports the following:

* `type` - (Required) The type of the query.

* `time_frame` - (Required) The time frame for pulling data for the query. If custom, then a specific time period must be provided. Possible values include: `WeekToDate`, `MonthToDate`, `YearToDate`, `TheLastWeek`, `TheLastMonth`, `TheLastYear`, `Custom`.

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
