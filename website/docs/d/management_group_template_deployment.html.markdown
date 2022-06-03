---
subcategory: "Template"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_management_group_template_deployment"
description: |-
  Gets information about an existing Management Group Template Deployment.
---

# Data Source: azurerm_management_group_template_deployment

Use this data source to access information about an existing Management Group Template Deployment.

## Example Usage

```hcl
data "azurerm_management_group_template_deployment" "example" {
  name                = "existing"
  management_group_id = "00000000-0000-0000-000000000000"
}

output "id" {
  value = data.azurerm_management_group_template_deployment.example.id
}

output "example_output" {
  value = jsondecode(data.azurerm_management_group_template_deployment.example.output_content).exampleOutput.value
}
```

## Arguments Reference

The following arguments are supported:

* `management_group_id` - (Required) The ID of the Management Group to which this template was applied.

* `name` - (Required) The name of this Management Group Template Deployment.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Management Group Template Deployment.

* `output_content` - The JSON Content of the Outputs of the ARM Template Deployment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Management Group Template Deployment.
