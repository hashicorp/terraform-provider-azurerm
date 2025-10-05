---
subcategory: "Template"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_subscription_template_deployment"
description: |-
  Gets information about an existing Subscription Template Deployment.
---

# Data Source: azurerm_subscription_template_deployment

Use this data source to access information about an existing Subscription Template Deployment.

## Example Usage

```hcl
data "azurerm_subscription_template_deployment" "example" {
  name = "existing"
}

output "id" {
  value = data.azurerm_subscription_template_deployment.example.id
}

output "example_output" {
  value = jsondecode(data.azurerm_subscription_template_deployment.example.output_content).exampleOutput.value
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of this Subscription Template Deployment.

## Attribute Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Subscription Template Deployment.

* `output_content` - The JSON Content of the Outputs of the ARM Template Deployment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Subscription Template Deployment.
