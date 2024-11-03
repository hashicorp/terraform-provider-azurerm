---
subcategory: "astro"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_astro_organization"
description: |-
  Manages an Astro Organization.
---

# azurerm_astro_organization

Manages an Astro Organization.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_astro_organization" "example" {
  name                = "example-ao"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = []
  }
  marketplace {
    subscription_id     = "00000000-0000-0000-0000-000000000000"
    subscription_status = "Subscribed"
    offer {
      offer_id     = "example-offer-id"
      plan_id      = "example-plan-id"
      plan_name    = "example-plan-name"
      publisher_id = "example-publisher-id"
      term_id      = "example-term-id"
      term_unit    = "example-term-unit"
    }
  }
  partner_organization {
    organization_id   = "example-organization-id"
    organization_name = "example-organization-name"
    workspace_id      = "example-workspace-id"
    workspace_name    = "example-workspace-name"
    single_sign_on {
      enterprise_app_id    = "00000000-0000-0000-0000-000000000000"
      single_sign_on_state = "Enable"
      single_sign_on_url   = "https://example.com/sso"
      aad_domains          = ["example.com"]
    }
  }
  user {
    email_address  = "user@example.com"
    first_name     = "John"
    last_name      = "Doe"
    phone_number   = "+1234567890"
    principal_name = "john.doe@example.com"
  }
  tags = {
    environment = "production"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Astro Organization. Changing this forces a new Astro Organization to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Astro Organization should exist. Changing this forces a new Astro Organization to be created.

* `location` - (Required) Specifies the Azure Region where the Astro Organization should exist. Changing this forces a new Astro Organization to be created.

* `marketplace` - (Required) A `marketplace` block as defined below.

* `user` - (Required) A `user` block as defined below.

* `identity` - (Optional) An `identity` block as defined below.

* `partner_organization` - (Optional) A `partner_organization` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Astro Organization.

---

A `marketplace` block supports the following:

* `offer` - (Required) An `offer` block as defined below.

* `subscription_id` - (Optional) Azure subscription id for the the marketplace offer is purchased from. Changing this forces a new Astro Organization to be created.

* `subscription_status` - (Optional) Marketplace subscription status. Changing this forces a new Astro Organization to be created.

---

An `offer` block supports the following:

* `offer_id` - (Required) Offer Id for the marketplace offer. Changing this forces a new Astro Organization to be created.

* `plan_id` - (Required) Plan Id for the marketplace offer. Changing this forces a new Astro Organization to be created.

* `plan_name` - (Optional) Plan Name for the marketplace offer. Changing this forces a new Astro Organization to be created.

* `publisher_id` - (Required) Publisher Id for the marketplace offer. Changing this forces a new Astro Organization to be created.

* `term_id` - (Optional) Plan Display Name for the marketplace offer. Changing this forces a new Astro Organization to be created.

* `term_unit` - (Optional) Plan Display Name for the marketplace offer. Changing this forces a new Astro Organization to be created.

---

A `user` block supports the following:

* `email_address` - (Required) Email address of the user.

* `first_name` - (Required) First name of the user.

* `last_name` - (Required) Last name of the user.

* `phone_number` - (Optional) User's phone number.

* `principal_name` - (Optional) User's principal name.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) A list of IDs for User Assigned Managed Identity resources to be assigned.

---

A `partner_organization` block supports the following:

* `organization_id` - (Optional) Organization Id in partner's system.

* `organization_name` - (Required) Organization name in partner's system.

* `single_sign_on` - (Optional) A `single_sign_on` block as defined below.

* `workspace_id` - (Optional) Workspace Id in partner's system.

* `workspace_name` - (Optional) Workspace name in partner's system.

---

A `single_sign_on` block supports the following:

* `aad_domains` - (Optional) List of AAD domains fetched from Microsoft Graph for user.

* `enterprise_app_id` - (Optional) AAD enterprise application Id used to setup SSO.

* `single_sign_on_state` - (Optional) State of the Single Sign On for the organization.

* `single_sign_on_url` - (Optional) URL for SSO to be used by the partner to redirect the user to their system.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Astro Organization.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Astro Organization.
* `read` - (Defaults to 5 minutes) Used when retrieving the Astro Organization.
* `update` - (Defaults to 30 minutes) Used when updating the Astro Organization.
* `delete` - (Defaults to 30 minutes) Used when deleting the Astro Organization.

## Import

Astro Organization can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_astro_organization.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Astronomer.astro/organizations/organization1
```
