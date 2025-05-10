---
subcategory: "Lighthouse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lighthouse_definition"
description: |-
  Manages a Lighthouse Definition.

---

# azurerm_lighthouse_definition

Manages a [Lighthouse](https://docs.microsoft.com/azure/lighthouse) Definition.

## Example Usage

```hcl
data "azurerm_role_definition" "contributor" {
  role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c"
}

resource "azurerm_lighthouse_definition" "example" {
  name               = "Sample definition"
  description        = "This is a lighthouse definition created via Terraform"
  managing_tenant_id = "00000000-0000-0000-0000-000000000000"
  scope              = "/subscriptions/00000000-0000-0000-0000-000000000000"

  authorization {
    principal_id           = "00000000-0000-0000-0000-000000000000"
    role_definition_id     = data.azurerm_role_definition.contributor.role_definition_id
    principal_display_name = "Tier 1 Support"
  }
}
```

## Argument Reference

The following arguments are supported:

* `lighthouse_definition_id` - (Optional) A unique UUID/GUID which identifies this lighthouse definition - one will be generated if not specified. Changing this forces a new resource to be created.

* `name` - (Required) The name of the Lighthouse Definition. Changing this forces a new resource to be created.

* `managing_tenant_id` - (Required) The ID of the managing tenant. Changing this forces a new resource to be created.

* `scope` - (Required) The ID of the managed subscription. Changing this forces a new resource to be created.

* `authorization` - (Required) An `authorization` block as defined below.
  
* `description` - (Optional) A description of the Lighthouse Definition.

* `eligible_authorization` - (Optional) An `eligible_authorization` block as defined below.

* `plan` - (Optional) A `plan` block as defined below.

---

An `authorization` block supports the following:

* `principal_id` - (Required) Principal ID of the security group/service principal/user that would be assigned permissions to the projected subscription.

* `role_definition_id` - (Required) The role definition identifier. This role will define the permissions that are granted to the principal. This cannot be an `Owner` role.

* `delegated_role_definition_ids` - (Optional) The set of role definition ids which define all the permissions that the principal id can assign.
  
* `principal_display_name` - (Optional) The display name of the security group/service principal/user that would be assigned permissions to the projected subscription.

---

An `eligible_authorization` block supports the following:

* `principal_id` - (Required) Principal ID of the security group/service principal/user that would be assigned permissions to the projected subscription.

* `role_definition_id` - (Required) The Principal ID of the Azure built-in role that defines the permissions that the Azure Active Directory will have on the projected scope.

* `just_in_time_access_policy` - (Optional) A `just_in_time_access_policy` block as defined below.

* `principal_display_name` - (Optional) The display name of the Azure Active Directory Principal.

---

A `just_in_time_access_policy` block supports the following:

* `multi_factor_auth_provider` - (Optional) The multi-factor authorization provider to be used for just-in-time access requests. Possible value is `Azure`.

~> **Note:** When this property isn't set, it would be set to `None`.

* `maximum_activation_duration` - (Optional) The maximum access duration in ISO 8601 format for just-in-time access requests. Defaults to `PT8H`.

* `approver` - (Optional) An `approver` block as defined below.

---

An `approver` block supports the following:

* `principal_id` - (Required) The Principal ID of the Azure Active Directory principal for the approver.

* `principal_display_name` - (Optional) The display name of the Azure Active Directory Principal for the approver.

---

A `plan` block supports the following:

* `name` - (Required) The plan name of the marketplace offer.

* `publisher` - (Required) The publisher ID of the plan.

* `product` - (Required) The product code of the plan.

* `version` - (Required) The version of the plan.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - the fully qualified ID of the Lighthouse Definition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Lighthouse Definition.
* `read` - (Defaults to 5 minutes) Used when retrieving the Lighthouse Definition.
* `update` - (Defaults to 30 minutes) Used when updating the Lighthouse Definition.
* `delete` - (Defaults to 30 minutes) Used when deleting the Lighthouse Definition.

## Import

Lighthouse Definitions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_lighthouse_definition.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.ManagedServices/registrationDefinitions/00000000-0000-0000-0000-000000000000
```
