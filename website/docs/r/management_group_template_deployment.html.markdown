---
subcategory: "Template"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_management_group_template_deployment"
description: |-
  Manages a Template Deployment at a Management Group Scope.
---

# azurerm_management_group_template_deployment

Manages a Template Deployment at a Management Group Scope.

~> **Note:** Deleting a Deployment at the Management Group Scope will not delete any resources created by the deployment. 

~> **Note:** Deployments to a Management Group are always Incrementally applied. Existing resources that are not part of the template will not be removed.

## Example Usage

```hcl
data "azurerm_management_group" "example" {
  name = "00000000-0000-0000-0000-000000000000"
}

resource "azurerm_management_group_template_deployment" "example" {
  name                = "example"
  location            = "West Europe"
  management_group_id = data.azurerm_management_group.example.id
  template_content    = <<TEMPLATE
{
  "$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    "policyAssignmentName": {
      "type": "string",
      "defaultValue": "[guid(parameters('policyDefinitionID'), resourceGroup().name)]",
      "metadata": {
        "description": "Specifies the name of the policy assignment, can be used defined or an idempotent name as the defaultValue provides."
      }
    },
    "policyDefinitionID": {
      "type": "string",
      "metadata": {
        "description": "Specifies the ID of the policy definition or policy set definition being assigned."
      }
    }
  },
  "resources": [
    {
      "type": "Microsoft.Authorization/policyAssignments",
      "name": "[parameters('policyAssignmentName')]",
      "apiVersion": "2019-09-01",
      "properties": {
        "scope": "[subscriptionResourceId('Microsoft.Resources/resourceGroups', resourceGroup().name)]",
        "policyDefinitionId": "[parameters('policyDefinitionID')]"
      }
    }
  ]
}
TEMPLATE

  parameters_content = <<PARAMS
{
  "$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentParameters.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    "policyDefinitionID": {
      "value": "/providers/Microsoft.Authorization/policyDefinitions/0a914e76-4921-4c19-b460-a2d36003525a"
    }
  }
}
PARAMS
}
```

```hcl

data "azurerm_management_group" "example" {
  name = "00000000-0000-0000-0000-000000000000"
}

resource "azurerm_management_group_template_deployment" "example" {
  name                = "example"
  location            = "West Europe"
  management_group_id = data.azurerm_management_group.example.id
  template_content    = file("templates/example-deploy-template.json")
  parameters_content  = file("templates/example-deploy-params.json")
}
```

```hcl

data "azurerm_management_group" "example" {
  name = "00000000-0000-0000-0000-000000000000"
}

data "azurerm_template_spec_version" "example" {
  name                = "exampleTemplateForManagementGroup"
  resource_group_name = "exampleResourceGroup"
  version             = "v1.0.9"
}

resource "azurerm_management_group_template_deployment" "example" {
  name                     = "example"
  location                 = "West Europe"
  management_group_id      = data.azurerm_management_group.example.id
  template_spec_version_id = data.azurerm_template_spec_version.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Template should exist. Changing this forces a new Template to be created.

* `management_group_name` - (Required) The Name of the Management Group to apply the Deployment Template to.

* `name` - (Required) The name which should be used for this Template Deployment. Changing this forces a new Template Deployment to be created.

---

* `debug_level` - (Optional) The Debug Level which should be used for this Resource Group Template Deployment. Possible values are `none`, `requestContent`, `responseContent` and `requestContent, responseContent`. 

* `parameters_content` - (Optional) The contents of the ARM Template parameters file - containing a JSON list of parameters.

* `template_content` - (Optional) The contents of the ARM Template which should be deployed into this Resource Group. Cannot be specified with `template_spec_version_id`. 

* `template_spec_version_id` - (Optional) The ID of the Template Spec Version to deploy. Cannot be specified with `template_content`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Template.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Management Group Template Deployment.

* `output_content` - The JSON Content of the Outputs of the ARM Template Deployment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Management Group Template Deployment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Management Group Template Deployment.
* `update` - (Defaults to 3 hours) Used when updating the Management Group Template Deployment.
* `delete` - (Defaults to 3 hours) Used when deleting the Management Group Template Deployment.

## Import

Management Group Template Deployments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_management_group_template_deployment.example /providers/Microsoft.Management/managementGroups/my-management-group-id/providers/Microsoft.Resources/deployments/deploy1
```
