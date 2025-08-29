---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_partner_registration"
description: |-
  Manages an EventGrid Partner Registration.
---

# azurerm_eventgrid_partner_registration

Manages an EventGrid Partner Registration.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_eventgrid_partner_registration" "example" {
  name                = "example-partner-registration"
  resource_group_name = azurerm_resource_group.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this EventGrid Partner Registration. Changing this forces a new EventGrid Partner Registration to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the EventGrid Partner Registration should exist. Changing this forces a new EventGrid Partner Registration to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the EventGrid Partner Registration.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the EventGrid Partner Registration.

* `partner_registration_id` - The immutable id of the corresponding partner registration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the EventGrid Partner Registration.
* `read` - (Defaults to 5 minutes) Used when retrieving the EventGrid Partner Registration.
* `update` - (Defaults to 30 minutes) Used when updating the EventGrid Partner Registration.
* `delete` - (Defaults to 30 minutes) Used when deleting the EventGrid Partner Registration.

## Import

EventGrid Partner Registrations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventgrid_partner_registration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.EventGrid/partnerRegistrations/example
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.EventGrid` - 2025-02-15
