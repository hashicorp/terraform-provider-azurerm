---
subcategory: "Cost Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cost_management_export_resource_group"
description: |-
  Manages an Azure Cost Management Export for a Resource Group.
---

# azurerm_cost_management_export_resource_group

Manages an Azure Cost Management Export for a Resource Group.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "northeurope"
}

resource "azurerm_storage_account" "example" {
  name                = "example-storage-account"
  resource_group_name = azurerm_resource_group.example.name

  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_cost_management_export_resource_group" "example" {
  name                    = "example"
  resource_group_id       = azurerm_resource_group.example.id
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

* `name` - (Required) Specifies the name of the Cost Management Export. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/data-factory/naming-rules) for all restrictions.

* `resource_group_id` - (Required) The id of the resource group in which to export information.

* `recurrence_type` - (Required) How often the requested information will be exported. Valid values include `Annually`, `Daily`, `Monthly`, `Weekly`.

* `recurrence_period_start` - (Required) The date the export will start capturing information.

* `recurrence_period_end` - (Required) The date the export will stop capturing information. 

* `delivery_info` - (Required) A `delivery_info` block as defined below.

---

A `delivery_info` block supports the following:

* `account_name` - (Required) Specifies the GitHub account name.

* `branch_name` - (Required) Specifies the branch of the repository to get code from.

* `git_url` - (Required) Specifies the GitHub Enterprise host name. For example: https://github.mydomain.com. Use https://github.com for open source repositories.

* `repository_name` - (Required) Specifies the name of the git repository.

* `root_folder` - (Required) Specifies the root folder within the repository. Set to `/` for the top level.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Data Factory.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the resource.
* `update` - (Defaults to 30 minutes) Used when updating the resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the resource.

## Import

Cost Management Export for a Resource Group can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cost_management_export_resource_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example
```
