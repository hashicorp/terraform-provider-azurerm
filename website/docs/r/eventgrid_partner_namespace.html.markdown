---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_partner_namespace"
description: |-
  Manages an EventGrid Partner Namespace.
---

# azurerm_eventgrid_partner_namespace

Manages an EventGrid Partner Namespace.

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

resource "azurerm_eventgrid_partner_namespace" "example" {
  name                         = "example-partner-namespace"
  location                     = azurerm_resource_group.example.location
  resource_group_name          = azurerm_resource_group.example.name
  partner_registration_id      = azurerm_eventgrid_partner_registration.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this EventGrid Partner Namespace. Changing this forces a new EventGrid Partner Namespace to be created.

* `location` - (Required) Specifies the Azure Region where the EventGrid Partner Namespace exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the EventGrid Partner Namespace should exist. Changing this forces a new EventGrid Partner Namespace to be created.

* `partner_registration_id` - (Required) The fully qualified ARM ID of the EventGrid Partner Registration that this namespace is associated with. Changing this forces a new EventGrid Partner Namespace to be created.

---

* `inbound_ip_rule` - (Optional) One or more `inbound_ip_rule` blocks as defined below.

* `local_auth_enabled` - (Optional) Whether local authentication methods is enabled for the EventGrid Partner Namespace. Defaults to `true`.

* `partner_topic_routing_mode` - (Optional) The partner topic routing mode. Possible values are `ChannelNameHeader` and `SourceEventAttribute`. Defaults to `ChannelNameHeader`. Changing this forces a new EventGrid Partner Namespace to be created.

* `public_network_access_enabled` - (Optional) Whether or not public network access is allowed for this server. Defaults to `true`.

* `tags` - (Optional) A mapping of tags which should be assigned to the EventGrid Partner Namespace.

---

An `inbound_ip_rule` block supports the following:

* `ip_mask` - (Required) The IP mask (CIDR) to match on.

* `action` - (Optional) The action to take when the rule is matched. Possible values are `Allow`. Defaults to `Allow`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the EventGrid Partner Namespace.

* `endpoint` - The endpoint for the EventGrid Partner Namespace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the EventGrid Partner Namespace.
* `read` - (Defaults to 5 minutes) Used when retrieving the EventGrid Partner Namespace.
* `update` - (Defaults to 30 minutes) Used when updating the EventGrid Partner Namespace.
* `delete` - (Defaults to 30 minutes) Used when deleting the EventGrid Partner Namespace.

## Import

EventGrid Partner Namespaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventgrid_partner_namespace.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.EventGrid/partnerNamespaces/example
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.EventGrid`: 2022-06-15