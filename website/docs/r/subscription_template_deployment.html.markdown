---
subcategory: "Template"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subscription_template_deployment"
description: |-
  Manages a Subscription Template Deployment.
---

# azurerm_subscription_template_deployment

Manages a Subscription Template Deployment.

## Example Usage

```hcl
resource "azurerm_subscription_template_deployment" "example" {
  name             = "example-deployment"
  location         = "West Europe"
  template_content = <<TEMPLATE
 {
   "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
   "contentVersion": "1.0.0.0",
   "parameters": {},
   "variables": {},
   "resources": [
     {
       "type": "Microsoft.Resources/resourceGroups",
       "apiVersion": "2018-05-01",
       "location": "West Europe",
       "name": "some-resource-group",
       "properties": {}
     }
   ]
 }
 TEMPLATE

  // NOTE: whilst we show an inline template here, we recommend
  // sourcing this from a file for readability/editor support
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Subscription Template Deployment should exist. Changing this forces a new Subscription Template Deployment to be created.

* `name` - (Required) The name which should be used for this Subscription Template Deployment. Changing this forces a new Subscription Template Deployment to be created.

---

* `debug_level` - (Optional) The Debug Level which should be used for this Subscription Template Deployment. Possible values are `none`, `requestContent`, `responseContent` and `requestContent, responseContent`.

* `expression_evaluation_option` - (Optional) An `expression_evaluation_option` block as defined below.

* `parameters_content` - (Optional) The contents of the ARM Template parameters file - containing a JSON list of parameters.

* `parameters_link` - (Optional)  A `parameters_link` block as defined below.

* `template_content` - (Optional) The contents of the ARM Template which should be deployed into this Subscription.

* `template_link` - (Optional)  A `template_link` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Subscription Template Deployment.

---

An `expression_evaluation_option` block exports the following:

* `scope` - (Optional) The scope to be used for evaluation of parameters, variables and functions in the nested ARM Template. Possible values are `Inner`, `Outer` and `NotSpecified`. Defaults to `Outer`.

---

An `parameters_link` block exports the following:

* `uri` - (Required) The URI of the ARM Template parameters file.

* `content_version` - (Optional) The content version of the ARM Template parameters file.

---

An `template_link` block exports the following:

* `uri` - (Required) The URI of the ARM Template file.

* `content_version` - (Optional) The content version of the ARM Template file.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Subscription Template Deployment.

* `output_content` - The JSON Content of the Outputs of the ARM Template Deployment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Subscription Template Deployment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Subscription Template Deployment.
* `update` - (Defaults to 3 hours) Used when updating the Subscription Template Deployment.
* `delete` - (Defaults to 3 hours) Used when deleting the Subscription Template Deployment.

## Import

Subscription Template Deployments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_subscription_template_deployment.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Resources/deployments/template1
```
