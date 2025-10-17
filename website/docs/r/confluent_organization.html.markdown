---
subcategory: "Confluent"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_confluent_organization"
description: |-
  Manages a Confluent Organization.
---

# azurerm_confluent_organization

Manages a Confluent Organization on Azure.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_confluent_organization" "example" {
  name                = "example-confluent-org"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  offer_detail {
    id           = "confluent-cloud-azure-prod"
    plan_id      = "confluent-cloud-azure-payg-prod"
    plan_name    = "Confluent Cloud - Pay as you Go"
    publisher_id = "confluentinc"
    term_unit    = "P1M"
  }

  user_detail {
    email_address = "user@example.com"
    first_name    = "Example"
    last_name     = "User"
  }

  tags = {
    environment = "Production"
  }
}
```

### With Link to Existing Confluent Organization

```hcl
resource "azurerm_confluent_organization" "example" {
  name                = "example-confluent-org"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  offer_detail {
    id           = "confluent-cloud-azure-prod"
    plan_id      = "confluent-cloud-azure-payg-prod"
    plan_name    = "Confluent Cloud - Pay as you Go"
    publisher_id = "confluentinc"
    term_unit    = "P1M"
  }

  user_detail {
    email_address = "user@example.com"
    first_name    = "Example"
    last_name     = "User"
  }

  link_organization {
    token = "confluent-link-token"
  }

  tags = {
    environment = "Production"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Confluent Organization. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Confluent Organization should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Confluent Organization should exist. Changing this forces a new resource to be created.

* `offer_detail` - (Required) An `offer_detail` block as defined below.

* `user_detail` - (Required) A `user_detail` block as defined below.

* `link_organization` - (Optional) A `link_organization` block as defined below. Used to link to an existing Confluent Cloud organization.

* `tags` - (Optional) A mapping of tags which should be assigned to the Confluent Organization.

---

An `offer_detail` block supports the following:

* `id` - (Required) The ID of the Confluent offer.

* `plan_id` - (Required) The plan ID for the Confluent organization.

* `plan_name` - (Required) The plan name for the Confluent organization.

* `publisher_id` - (Required) The publisher ID for the Confluent offer.

* `term_unit` - (Required) The term unit for the Confluent offer (e.g., `P1M` for monthly).

* `private_offer_id` - (Optional) The private offer ID if using a private marketplace offer.

* `private_offer_ids` - (Optional) A list of private offer IDs.

* `term_id` - (Optional) The term ID for the offer.

---

A `user_detail` block supports the following:

* `email_address` - (Required) The email address of the user.

* `first_name` - (Optional) The first name of the user.

* `last_name` - (Optional) The last name of the user.

---

A `link_organization` block supports the following:

* `token` - (Required) The linking token from an existing Confluent Cloud organization. This is a sensitive value.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Confluent Organization.

* `organization_id` - The Confluent Organization ID.

* `sso_url` - The Single Sign-On URL for the organization.

* `created_time` - The time when the organization was created.

* `provisioning_state` - The provisioning state of the organization.

---

An `offer_detail` block exports the following:

* `status` - The status of the SaaS offer subscription.

---

A `user_detail` block exports the following:

* `aad_email` - The Azure Active Directory email of the user.

* `user_principal_name` - The user principal name.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Confluent Organization.
* `read` - (Defaults to 5 minutes) Used when retrieving the Confluent Organization.
* `update` - (Defaults to 30 minutes) Used when updating the Confluent Organization.
* `delete` - (Defaults to 30 minutes) Used when deleting the Confluent Organization.

## Import

Confluent Organizations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_confluent_organization.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Confluent/organizations/org1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Confluent` - 2024-07-01
