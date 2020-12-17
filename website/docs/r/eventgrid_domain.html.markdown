---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_domain"
description: |-
  Manages an EventGrid Domain

---

# azurerm_eventgrid_domain

Manages an EventGrid Domain

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "resourceGroup1"
  location = "West US 2"
}

resource "azurerm_eventgrid_domain" "example" {
  name                = "my-eventgrid-domain"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the EventGrid Domain resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the EventGrid Domain exists. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `input_schema` - (Optional) Specifies the schema in which incoming events will be published to this domain. Allowed values are `CloudEventSchemaV1_0`, `CustomEventSchema`, or `EventGridSchema`. Defaults to `eventgridschema`. Changing this forces a new resource to be created.

* `input_mapping_fields` - (Optional) A `input_mapping_fields` block as defined below.

* `input_mapping_default_values` - (Optional) A `input_mapping_default_values` block as defined below.

* `public_network_access_enabled` - (Optional) Whether or not public network access is allowed for this server. Defaults to `true`.

* `inbound_ip_rule` - (Optional) One or more `inbound_ip_rule` blocks as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `input_mapping_fields` supports the following:

* `id` - (Optional) Specifies the id of the EventGrid Event to associate with the domain. Changing this forces a new resource to be created.

* `topic` - (Optional) Specifies the topic of the EventGrid Event to associate with the domain. Changing this forces a new resource to be created.

* `event_type` - (Optional) Specifies the event type of the EventGrid Event to associate with the domain. Changing this forces a new resource to be created.

* `event_time` - (Optional) Specifies the event time of the EventGrid Event to associate with the domain. Changing this forces a new resource to be created.

* `data_version` - (Optional) Specifies the data version of the EventGrid Event to associate with the domain. Changing this forces a new resource to be created.

* `subject` - (Optional) Specifies the subject of the EventGrid Event to associate with the domain. Changing this forces a new resource to be created.

---

A `input_mapping_default_values` supports the following:

* `event_type` - (Optional) Specifies the default event type of the EventGrid Event to associate with the domain. Changing this forces a new resource to be created.

* `data_version` - (Optional) Specifies the default data version of the EventGrid Event to associate with the domain. Changing this forces a new resource to be created.

* `subject` - (Optional) Specifies the default subject of the EventGrid Event to associate with the domain. Changing this forces a new resource to be created.

---

A `inbound_ip_rule` block supports the following:

* `ip_mask` - (Required) The ip mask (CIDR) to match on.

* `action` - (Optional) The action to take when the rule is matched. Possible values are `Allow`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the EventGrid Domain.

* `endpoint` - The Endpoint associated with the EventGrid Domain.

* `primary_access_key` - The Primary Shared Access Key associated with the EventGrid Domain.

* `secondary_access_key` - The Secondary Shared Access Key associated with the EventGrid Domain.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the EventGrid Domain.
* `update` - (Defaults to 30 minutes) Used when updating the EventGrid Domain.
* `read` - (Defaults to 5 minutes) Used when retrieving the EventGrid Domain.
* `delete` - (Defaults to 30 minutes) Used when deleting the EventGrid Domain.

## Import

EventGrid Domains can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventgrid_domain.domain1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventGrid/domains/domain1
```
