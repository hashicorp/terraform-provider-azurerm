---
subcategory: "Template"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subscription_deployment_stack"
description: |-
  Manages a Subscription Deployment Stack.
---

# azurerm_subscription_deployment_stack

Manages a Subscription Deployment Stack.

## Example Usage

```hcl
resource "azurerm_subscription_deployment_stack" "example" {
  name     = "example-stack"
  location = "West Europe"

  template_content = jsonencode({
    "$schema"      = "https://schema.management.azure.com/schemas/2018-05-01/subscriptionDeploymentTemplate.json#"
    contentVersion = "1.0.0.0"
    resources      = []
  })

  action_on_unmanage {
    resources = "delete"
  }

  deny_settings {
    mode = "none"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `action_on_unmanage` - (Required) An `action_on_unmanage` block as defined below.

* `deny_settings` - (Required) A `deny_settings` block as defined below.

* `location` - (Required) The Azure Region where the Subscription Deployment Stack should exist. Changing this forces a new Subscription Deployment Stack to be created.

* `name` - (Required) The name which should be used for this Subscription Deployment Stack. Changing this forces a new Subscription Deployment Stack to be created.

---

* `deployment_resource_group_name` - (Optional) The name of the Resource Group to deploy the resources to. Changing this forces a new Subscription Deployment Stack to be created.

* `description` - (Optional) The description of the Deployment Stack.

* `parameters_content` - (Optional) The JSON content of the ARM Template parameters file.

* `tags` - (Optional) A mapping of tags which should be assigned to the Subscription Deployment Stack.

* `template_content` - (Optional) The JSON content of the ARM Template. Exactly one of `template_content` or `template_spec_version_id` must be specified.

* `template_spec_version_id` - (Optional) The ID of the Template Spec Version. Exactly one of `template_content` or `template_spec_version_id` must be specified.

---

An `action_on_unmanage` block supports the following:

* `resources` - (Required) Specifies the action to take on resources that are no longer managed by the deployment stack. Possible values are `delete` and `detach`.

* `management_groups` - (Optional) Specifies the action to take on management groups that are no longer managed by the deployment stack. Possible values are `delete` and `detach`. Defaults to `detach`.

* `resource_groups` - (Optional) Specifies the action to take on resource groups that are no longer managed by the deployment stack. Possible values are `delete` and `detach`. Defaults to `detach`.

---

A `deny_settings` block supports the following:

* `mode` - (Required) Specifies the deny settings mode. Possible values are `none`, `denyDelete`, and `denyWriteAndDelete`.

* `apply_to_child_scopes` - (Optional) Specifies whether to apply the deny settings to child scopes. Defaults to `false`.

* `excluded_actions` - (Optional) Specifies a list of role-based access control (RBAC) management operations that are excluded from the deny settings. Each entry must be a valid Azure RBAC action string.

* `excluded_principals` - (Optional) Specifies a list of Azure Active Directory principal IDs that are excluded from the deny settings.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Subscription Deployment Stack.

* `deployment_id` - The ID of the underlying Deployment resource.

* `duration` - The duration of the last deployment operation.

* `output_content` - The JSON content of the deployment outputs.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Subscription Deployment Stack.
* `read` - (Defaults to 5 minutes) Used when retrieving the Subscription Deployment Stack.
* `update` - (Defaults to 3 hours) Used when updating the Subscription Deployment Stack.
* `delete` - (Defaults to 3 hours) Used when deleting the Subscription Deployment Stack.

## Import

Subscription Deployment Stacks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_subscription_deployment_stack.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Resources/deploymentStacks/stack1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Resources` - 2024-03-01
