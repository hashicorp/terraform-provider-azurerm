---
subcategory: "Template"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_group_template_deployment"
description: |-
  Manages a Resource Group Template Deployment.
---

# azurerm_resource_group_template_deployment

Manages a Resource Group Template Deployment.

~> **Note:** This resource will automatically attempt to delete resources deployed by the ARM Template when it is deleted. You can opt-out of this by setting the `delete_nested_items_during_deletion` field within the `template_deployment` block of the `features` block to `false`.

## Example Usage

```hcl
resource "azurerm_resource_group_template_deployment" "example" {
  name                = "example-deploy"
  resource_group_name = "example-group"
  deployment_mode     = "Complete"
  template_content    = <<TEMPLATE
{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {},
  "variables": {},
  "resources": [
    {
      "type": "Microsoft.Network/virtualNetworks",
      "apiVersion": "2020-05-01",
      "name": "acctest-network",
      "location": "[resourceGroup().location]",
      "properties": {
        "addressSpace": {
          "addressPrefixes": [
            "10.0.0.0/16"
          ]
        }
      }
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

* `deployment_mode` - (Required) The Deployment Mode for this Resource Group Template Deployment. Possible values are `Complete` (where resources in the Resource Group not specified in the ARM Template will be destroyed) and `Incremental` (where resources are additive only).

* `name` - (Required) The name which should be used for this Resource Group Template Deployment. Changing this forces a new Resource Group Template Deployment to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Resource Group Template Deployment should exist. Changing this forces a new Resource Group Template Deployment to be created.

---

* `debug_level` - (Optional) The Debug Level which should be used for this Resource Group Template Deployment. Possible values are `none`, `requestContent`, `responseContent` and `requestContent, responseContent`.

* `expression_evaluation_option` - (Optional) An `expression_evaluation_option` block as defined below.

* `on_error_deployment` - (Optional) An `on_error_deployment` block as defined below.

* `parameters_content` - (Optional) The contents of the ARM Template parameters file - containing a JSON list of parameters.

* `parameters_link` - (Optional)  A `parameters_link` block as defined below.

* `template_content` - (Optional) The contents of the ARM Template which should be deployed into this Resource Group.

~> **Note:** If `deployment_mode` is set to `Complete` then resources within this Resource Group which are not defined in the ARM Template will be deleted.

* `template_link` - (Optional)  A `template_link` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Resource Group Template Deployment.

---

An `expression_evaluation_option` block exports the following:

* `scope` - (Optional) The scope to be used for evaluation of parameters, variables and functions in the nested ARM Template. Possible values are `Inner`, `Outer` and `NotSpecified`. Defaults to `Outer`.

---

An `on_error_deployment` block exports the following:

* `type` - (Required) The type of the previous ARM Template Deployment to be used when provisioning the current ARM Template Deployment is failed. Possible values are `LastSuccessful` and `SpecificDeployment`.

* `deployment_name` - (Optional) The name of the previous ARM Template Deployment to be used when provisioning the current ARM Template Deployment is failed.

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

* `id` - The ID of the Resource Group Template Deployment.

* `output_content` - The JSON Content of the Outputs of the ARM Template Deployment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Resource Group Template Deployment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Group Template Deployment.
* `update` - (Defaults to 3 hours) Used when updating the Resource Group Template Deployment.
* `delete` - (Defaults to 3 hours) Used when deleting the Resource Group Template Deployment.

## Import

Resource Group Template Deployments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_resource_group_template_deployment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Resources/deployments/template1
```
