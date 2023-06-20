---
subcategory: "Application Insights"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_insights_workbook"
description: |-
  Manages an Azure Workbook.
---

# azurerm_application_insights_workbook

Manages an Azure Workbook.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_application_insights_workbook" "example" {
  name                = "85b3e8bb-fc93-40be-83f2-98f6bec18ba0"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  display_name        = "workbook1"
  data_json = jsonencode({
    "version" = "Notebook/1.0",
    "items" = [
      {
        "type" = 1,
        "content" = {
          "json" = "Test2022"
        },
        "name" = "text - 0"
      }
    ],
    "isLocked" = false,
    "fallbackResourceIds" = [
      "Azure Monitor"
    ]
  })

  tags = {
    ENV = "Test"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Workbook as a UUID/GUID. It should not contain any uppercase letters. Changing this forces a new Workbook to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Workbook should exist. Changing this forces a new Workbook to be created.

* `location` - (Required) Specifies the Azure Region where the Workbook should exist. Changing this forces a new Workbook to be created.

* `display_name` - (Required) Specifies the user-defined name (display name) of the workbook.

* `data_json` - (Required) Configuration of this particular workbook. Configuration data is a string containing valid JSON.

* `source_id` - (Optional) Resource ID for a source resource. It should not contain any uppercase letters. Defaults to `azure monitor`.

* `category` - (Optional) Workbook category, as defined by the user at creation time. There may be additional category types beyond the following: `workbook`, `sentinel`. Defaults to `workbook`.

* `description` - (Optional) Specifies the description of the workbook.

* `identity` - (Optional) An `identity` block as defined below. Changing this forces a new Workbook to be created.

* `storage_container_id` - (Optional) Specifies the Resource Manager ID of the Storage Container when bring your own storage is used. Changing this forces a new Workbook to be created.

-> **Note:** This is the Resource Manager ID of the Storage Container, rather than the regular ID - and can be accessed on the `azurerm_storage_container` Data Source/Resource as `resource_manager_id`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Workbook.

---

An `identity` block exports the following:

* `type` - (Required) The type of Managed Service Identity that is configured on this Workbook. Possible values are `UserAssigned`, `SystemAssigned` and `SystemAssigned, UserAssigned`. Changing this forces a new resource to be created.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Workbook.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Workbook.

* `identity_ids` - (Optional) The list of User Assigned Managed Identity IDs assigned to this Workbook. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Workbook.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Workbook.
* `read` - (Defaults to 5 minutes) Used when retrieving the Workbook.
* `update` - (Defaults to 30 minutes) Used when updating the Workbook.
* `delete` - (Defaults to 30 minutes) Used when deleting the Workbook.

## Import

Workbooks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_application_insights_workbook.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Insights/workbooks/resource1
```
