---
subcategory: "Application Insights"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_insights_workbook_template"
description: |-
  Manages an Application Insights Workbook Templates.
---

# azurerm_application_insights_workbook_template

Manages an Application Insights Workbook Templates.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_application_insights_workbook_template" "example" {
  name                = "example-aiwt"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
  author              = ""
  priority            = 0
  galleries {
    category      = ""
    name          = ""
    order         = 0
    resource_type = ""
    type          = ""
  }
  template_data = jsonencode({
    "key" = "value"
  })
  tags = {
    key = "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Application Insights Workbook Templates. Changing this forces a new Application Insights Workbook Templates to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Application Insights Workbook Templates should exist. Changing this forces a new Application Insights Workbook Templates to be created.

* `galleries` - (Required) A `galleries` block as defined below.

* `location` - (Required) Specifies the Azure Region where the Application Insights Workbook Templates should exist. Changing this forces a new Application Insights Workbook Templates to be created.

* `template_data` - (Required) Valid JSON object containing workbook template payload.

* `author` - (Optional) Information about the author of the workbook template.

* `priority` - (Optional) Priority of the template. Determines which template to open when a workbook gallery is opened in viewer mode. Defaults to 0.

* `tags` - (Optional) A mapping of tags which should be assigned to the Application Insights Workbook Templates.

---

A `galleries` block supports the following:

* `name` - (Required) Name of the workbook template in the gallery.

* `category` - (Required) Category for the gallery.

* `order` - (Optional) Order of the template within the gallery. Defaults to `0`.

* `resource_type` - (Optional) Azure resource type supported by the gallery. Defaults to `Azure Monitor`.

* `type` - (Optional) Type of workbook supported by the workbook template. Defaults to `workbook`.

~> **Note:** See [documentation](https://docs.microsoft.com/en-us/azure/azure-monitor/visualize/workbooks-automate#galleries) for more information of `resource_type` and `type`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Application Insights Workbook Templates.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Application Insights Workbook Templates.
* `read` - (Defaults to 5 minutes) Used when retrieving the Application Insights Workbook Templates.
* `update` - (Defaults to 30 minutes) Used when updating the Application Insights Workbook Templates.
* `delete` - (Defaults to 30 minutes) Used when deleting the Application Insights Workbook Templates.

## Import

Application Insights Workbook Templates can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_application_insights_workbook_template.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Insights/workbooktemplates/resource1
```
