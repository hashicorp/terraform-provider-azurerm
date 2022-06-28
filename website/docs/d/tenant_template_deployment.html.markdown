---
subcategory: "Template"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_tenant_template_deployment"
description: |-
  Gets information about an existing Tenant Template Deployment.
---

# Data Source: azurerm_tenant_template_deployment

Use this data source to access information about an existing Tenant Template Deployment.

## Example Usage

```hcl
data "azurerm_tenant_template_deployment" "example" {
  name = "existing"
}

output "id" {
  value = data.azurerm_tenant_template_deployment.example.id
}

output "example_output" {
  value = jsondecode(data.azurerm_tenant_template_deployment.example.output_content).exampleOutput.value
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Tenant Template Deployment.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Tenant Template Deployment.

* `output_content` - The JSON Content of the Outputs of the ARM Template Deployment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Tenant Template Deployment.
