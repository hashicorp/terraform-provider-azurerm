---
subcategory: "Managed Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lighthouse_definition"
description: |-
  Get information about an existing Lighthouse Definition.
---

# Data Source: azurerm_lighthouse_definition

Use this data source to access information about an existing Lighthouse Definition.

## Example Usage

```hcl
resource "azurerm_lighthouse_definition" "example" {
  registration_definition_name = "Sample definition"
  description                  = "This is a lighthouse definition created via Terraform"
  managed_by_tenant_id         = "00000000-0000-0000-0000-000000000000"

  authorization {
    principal_id       = "00000000-0000-0000-0000-000000000000"
    role_definition_id = data.azurerm_role_definition.builtin.name
  }
}

data "azurerm_lighthouse_definition" "example" {
  registration_definition_id = azurerm_lighthouse_definition.example.registration_definition_id
}

output "registration_definition_id" {
  value = data.azurerm_lighthouse_definition.example.id
}
```

## Argument Reference

* `registration_definition_id` - (Required) Specifies the ID of the Lighthouse Definition as a UUID/GUID.

## Attributes Reference

* `id` - the Lighthouse Definition ID.
* `registration_definition_name` - the name of the Lighthouse Definition.
* `description` - the description of the Lighthouse Definition.
* `scope` - the scope at which the Lighthouse Definition applies too, such as `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333`. Only subscription level scope is supported.
* `managed_by_tenant_id` - the ID of the managing tenant.
* `authorization` - the block as the following properties:
* `principal_id` - the principal Id of the security group/service principal/user that would be assigned permissions to the projected subscription or resource group.
* `role_definition_id` - the role definition identifier. This role will define all the permissions that the security group/service principal/user must have on the projected subscription or resource group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Lighthouse Definition.
