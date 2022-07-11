---
subcategory: "Template"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_resource_group_template_deployment"
description: |-
  Gets information about an existing Resource Group Template Deployment.
---

# Data Source: azurerm_resource_group_template_deployment

Use this data source to access information about an existing Resource Group Template Deployment.

## Example Usage

```hcl
data "azurerm_resource_group_template_deployment" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_resource_group_template_deployment.example.id
}

output "example_output" {
  value = jsondecode(data.azurerm_resource_group_template_deployment.example.output_content).exampleOutput.value
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Resource Group Template Deployment.

* `resource_group_name` - (Required) The name of the Resource Group to which the Resource Group Template Deployment was applied.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Resource Group Template Deployment.

* `output_content` - The JSON Content of the Outputs of the ARM Template Deployment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Group Template Deployment.
