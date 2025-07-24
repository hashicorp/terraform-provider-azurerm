---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_partner_registration"
description: |-
  Gets information about an existing EventGrid Partner Registration

---

# Data Source: azurerm_eventgrid_partner_registration

Use this data source to access information about an existing EventGrid Partner Registration

## Example Usage

```hcl
data "azurerm_eventgrid_partner_registration" "example" {
  name                = "my-eventgrid-partner-registration"
  resource_group_name = "example-resources"
}

output "eventgrid_partner_registration_id" {
  value = data.azurerm_eventgrid_partner_registration.example.partner_registration_id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the EventGrid Partner Registration resource.

* `resource_group_name` - (Required) The name of the resource group in which the EventGrid Partner Registration exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the EventGrid Partner Registration.

* `partner_registration_id` - The immutable id of the corresponding partner registration.

* `tags` - A mapping of tags which are assigned to the EventGrid Partner Registration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the EventGrid Partner Registration.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.EventGrid`: 2022-06-15
