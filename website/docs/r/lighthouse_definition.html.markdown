---
subcategory: "Lighthouse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lighthouse_definition"
description: |-
  Manages a Lighthouse Definition.

---

# azurerm_lighthouse_definition

Manages a [Lighthouse](https://docs.microsoft.com/en-us/azure/lighthouse) Definition.

## Example Usage

```hcl
data "azurerm_role_definition" "contributor" {
  role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c"
}

resource "azurerm_lighthouse_definition" "example" {
  name               = "Sample definition"
  description        = "This is a lighthouse definition created via Terraform"
  managing_tenant_id = "00000000-0000-0000-0000-000000000000"

  authorization {
    principal_id       = "00000000-0000-0000-0000-000000000000"
    role_definition_id = data.azurerm_role_definition.contributor.role_definition_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `lighthouse_definition_id` - (Optional) A unique UUID/GUID which identifies this lighthouse definition - one will be generated if not specified. Changing this forces a new resource to be created.

* `name` - (Required) The name of the Lighthouse Definition.

* `managing_tenant_id` - (Required) The ID of the managing tenant.

* `description` - (Optional) A description of the Lighthouse Definition.

* `authorization` - (Required) An authorization block as defined below.  

--- 

An `authorization` block supports the following:

* `principal_id` - (Required) Principal ID of the security group/service principal/user that would be assigned permissions to the projected subscription.

* `role_definition_id` - (Required) The role definition identifier. This role will define the permissions that are granted to the principal. This cannot be an `Owner` role.

## Attributes Reference

The following attributes are exported:

* `id` - the fully qualified ID of the Lighthouse Definition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Lighthouse Definition.
* `update` - (Defaults to 30 minutes) Used when updating the Lighthouse Definition.
* `read` - (Defaults to 5 minutes) Used when retrieving the Lighthouse Definition.
* `delete` - (Defaults to 30 minutes) Used when deleting the Lighthouse Definition.

## Import

Lighthouse Definitions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_lighthouse_definition.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.ManagedServices/registrationDefinitions/00000000-0000-0000-0000-000000000000
```
