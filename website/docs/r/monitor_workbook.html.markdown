---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_workbook"
description: |-
  Manages a Monitor Workbook.
---

# azurerm_monitor_workbook

Manages a Monitor Workbook.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_monitor_workbook" "example" {
  name                = "85b3e8bb-fc93-40be-83f2-98f6bec18ba0"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  display_name        = "workbook1"
  source_id           = azurerm_resource_group.example.id
  serialized_data = jsonencode({
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

* `name` - (Required) Specifies the name which should be used for this Monitor Workbook. Changing this forces a new Monitor Workbook to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Monitor Workbook should exist. Changing this forces a new Monitor Workbook to be created.

* `location` - (Required) Specifies the Azure Region where the Monitor Workbook should exist. Changing this forces a new Monitor Workbook to be created.

* `display_name` - (Required) Specifies the user-defined name (display name) of the workbook.

* `serialized_data` - (Required) Configuration of this particular workbook. Configuration data is a string containing valid JSON.

* `source_id` - (Required) ResourceId for a source resource. Changing this forces a new Monitor Workbook to be created.

* `category` - (Optional) Workbook category, as defined by the user at creation time. Defaults to `workbook`.

* `description` - (Optional) Specifies the description of the workbook.

* `identity` - (Optional) An `identity` block as defined below. Changing this forces a new Monitor Workbook to be created.

* `storage_uri` - (Optional) Specifies the resourceId to the storage account when bring your own storage is used. Changing this forces a new Monitor Workbook to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Monitor Workbook.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Monitor Workbook.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Monitor Workbook.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Monitor Workbook.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Monitor Workbook.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Monitor Workbook.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Monitor Workbook.
* `read` - (Defaults to 5 minutes) Used when retrieving the Monitor Workbook.
* `update` - (Defaults to 30 minutes) Used when updating the Monitor Workbook.
* `delete` - (Defaults to 30 minutes) Used when deleting the Monitor Workbook.

## Import

Monitor Workbooks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_workbook.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Insights/workbooks/resource1
```
