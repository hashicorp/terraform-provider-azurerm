---
subcategory: "Managed Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_registration_definition"
description: |-
  Get information about an existing Registration Definition.
---

# Data Source: azurerm_registration_definition

Use this data source to access information about an existing Registration Definition.

## Example Usage

```hcl
data "azurerm_subscription" "primary" {
}

resource "azuread_application" "example" {
  name = "acctestspa-%d"
}

resource "azuread_service_principal" "example" {
  application_id = azuread_application.example.application_id
}

data "azurerm_role_definition" "builtin" {
  name = "Contributor"
}

resource "azurerm_registration_definition" "example" {
  name                  = "Sample registration"
  scope                 = data.azurerm_subscription.primary.id
  description           = "This is a registration definition created via Terraform"
  managed_by_tenant_id  = "00000000-0000-0000-0000-000000000000"

  authorization {
    principal_id        = azuread_service_principal.example.id
    role_definition_id  = data.azurerm_role_definition.builtin.name
  }
}

data "azurerm_registration_definition" "example" {
  registration_definition_id = azurerm_registration_definition.definition1.registration_definition_id
  scope                      = data.azurerm_subscription.primary.id # /subscriptions/00000000-0000-0000-0000-000000000000
}

output "registration_definition_id" {
  value = data.azurerm_registration_definition.example.id
}
```

## Argument Reference

* `registration_definition_id` - (Required) Specifies the ID of the Registration Definition as a UUID/GUID.
* `scope` - (Required) Specifies the Scope at which the Custom Role Definition exists.

## Attributes Reference

* `id` - the Registration Definition ID.
* `registration_definition_name` - the name of the Registration Definition.
* `description` - the description of the Registration Definition.
* `scope` - the scope at which the Registration Definition applies too, such as `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333`. Only subscription level scope is supported.
* `managed_by_tenant_id` - the ID of the managing tenant.
* `authorization` - the block as the following properties:
* `principal_id` - the principal Id of the security group/service principal/user that would be assigned permissions to the projected subscription or resource group.
* `role_definition_id` - the role definition identifier. This role will define all the permissions that the security group/service principal/user must have on the projected subscription or resource group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Registration Definition.
