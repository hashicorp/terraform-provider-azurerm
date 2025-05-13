---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_partner_configuration"
description: |-
  Manages an Event Grid Partner Configuration.
---

# azurerm_eventgrid_partner_configuration

Manages an Event Grid Partner Configuration.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_eventgrid_partner_configuration" "example" {
  resource_group_name                     = azurerm_resource_group.example.name
  default_maximum_expiration_time_in_days = 14

  partner_authorization {
    partner_registration_id              = "804a11ca-ce9b-4158-8e94-3c8dc7a072ec"
    partner_name                         = "Auth0"
    authorization_expiration_time_in_utc = "2025-02-05T00:00:00Z"
  }

  tags = {
    environment = "Production"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the Resource Group where the Event Grid Partner Configuration should exist. Changing this forces a new Event Grid Partner Configuration to be created.

---

* `default_maximum_expiration_time_in_days` - (Optional) Time used to validate the authorization expiration time for each authorized partner. Defaults to `7`.

* `partner_authorization` - (Optional) One or more `partner_authorization` blocks as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Event Grid Partner Configuration.

---

A `partner_authorization` block supports the following:

* `partner_name` - (Required) The partner name.

* `partner_registration_id` - (Required) The immutable id of the corresponding partner registration.

* `authorization_expiration_time_in_utc` - (Optional) Expiration time of the partner authorization. Value should be in RFC 3339 format in UTC time zone, for example: "2025-02-04T00:00:00Z".

-> **Note:** If the time from `authorization_expiration_time_in_utc` expires, any request from this partner to create, update or delete resources in the subscriber's context will fail. If not specified, the authorization will expire after `default_maximum_expiration_time_in_days`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Event Grid Partner Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Event Grid Partner Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Event Grid Partner Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the Event Grid Partner Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the Event Grid Partner Configuration.

## Import

Event Grid Partner Configurations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventgrid_partner_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1
```
