---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_partner_registration"
description: |-
  Manages an Event Grid Partner Registration.
---

# azurerm_eventgrid_partner_registration

Manages an Event Grid Partner Registration.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_eventgrid_partner_registration" "example" {
  name                = "example-partner-registration"
  resource_group_name = azurerm_resource_group.example.name

  tags = {
    environment = "Production"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The Name which should be used for this Partner Registration. Changing this forces a new Partner Registration to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Partner Registration exists.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Partner Registration.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Partner Registration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Resource Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Group.
* `update` - (Defaults to 30 minutes) Used when updating the Resource Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Resource Group.

## Import

Resource Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventgrid_partner_registration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.EventGrid/partnerRegistrations/example
```
