---
subcategory: "Confluent"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_confluent_environment"
description: |-
  Manages a Confluent Environment.
---

# azurerm_confluent_environment

Manages a Confluent Environment within a Confluent Organization on Azure.

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
  }
}

resource "azurerm_confluent_environment" "example" {
  environment_id      = "env-12345"
  organization_id     = azurerm_confluent_organization.example.name
  resource_group_name = azurerm_resource_group.example.name
  display_name        = "production-environment"

  stream_governance {
    package = "ESSENTIALS"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `environment_id` - (Required) The Confluent Environment ID. Changing this forces a new resource to be created.

* `organization_id` - (Required) The name of the parent Confluent Organization. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Confluent Environment should exist. Changing this forces a new resource to be created.

* `display_name` - (Optional) The display name for the environment. Changing this forces a new resource to be created.

* `stream_governance` - (Optional) A `stream_governance` block as defined below. Changing this forces a new resource to be created.

---

A `stream_governance` block supports the following:

* `package` - (Required) The Stream Governance package type. Possible values are `ESSENTIALS` and `ADVANCED`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Confluent Environment.

* `kind` - The kind of the resource.

* `metadata` - A `metadata` block as defined below.

---

A `metadata` block exports the following:

* `self` - The self-referencing link for this environment.

* `resource_name` - The Confluent resource name for this environment.

* `created_timestamp` - The timestamp when the environment was created.

* `updated_timestamp` - The timestamp when the environment was last updated.

* `deleted_timestamp` - The timestamp when the environment was deleted (if applicable).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Confluent Environment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Confluent Environment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Confluent Environment.

## Import

Confluent Environments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_confluent_environment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Confluent/organizations/org1/environments/env-12345
```
