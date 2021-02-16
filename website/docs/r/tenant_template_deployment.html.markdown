---
subcategory: "Template"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_tenant_template_deployment"
description: |-
  Manages a Template Deployment at the Tenant Scope.
---

# azurerm_tenant_template_deployment

Manages a Template Deployment at the Tenant Scope.

~> **Note:** Deleting a Deployment at the Tenant Scope will not delete any resources created by the deployment.

~> **Note:** Deployments to a Tenant are always Incrementally applied. Existing resources that are not part of the template will not be removed.

## Example Usage

```hcl
data "azurerm_template_spec_version" "example" {
  name                = "myTemplateForTenant"
  resource_group_name = "myResourceGroup"
  version             = "v0.1"
}

resource "azurerm_tenant_template_deployment" "example" {
  name                     = "example"
  location                 = "West Europe"
  template_spec_version_id = data.azurerm_template_spec_version.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Template should exist. Changing this forces a new Template to be created.

* `name` - (Required) The name which should be used for this Template. Changing this forces a new Template to be created.

---

* `debug_level` - (Optional) The Debug Level which should be used for this Resource Group Template Deployment. Possible values are `none`, `requestContent`, `responseContent` and `requestContent, responseContent`.

* `parameters_content` - (Optional) The contents of the ARM Template parameters file - containing a JSON list of parameters.

* `template_content` - (Optional) The contents of the ARM Template which should be deployed into this Resource Group. Cannot be specified with `template_spec_version_id`.

* `template_spec_version_id` - (Optional) The ID of the Template Spec Version to deploy. Cannot be specified with `template_content`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Template.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Tenant Template Deployment.

* `output_content` - The JSON Content of the Outputs of the ARM Template Deployment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Template.
* `read` - (Defaults to 5 minutes) Used when retrieving the Template.
* `update` - (Defaults to 3 hours) Used when updating the Template.
* `delete` - (Defaults to 3 hours) Used when deleting the Template.

## Import

Templates can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_tenant_template_deployment.example /providers/Microsoft.Resources/deployments/deploy1
```
