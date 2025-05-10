---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_topic"
description: |-
  Manages an EventGrid Topic

---

# azurerm_eventgrid_topic

Manages an EventGrid Topic

~> **Note:** at this time EventGrid Topic's are only available in a limited number of regions.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_eventgrid_topic" "example" {
  name                = "my-eventgrid-topic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the EventGrid Topic resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the EventGrid Topic exists. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `input_schema` - (Optional) Specifies the schema in which incoming events will be published to this domain. Allowed values are `CloudEventSchemaV1_0`, `CustomEventSchema`, or `EventGridSchema`. Defaults to `EventGridSchema`. Changing this forces a new resource to be created.

* `input_mapping_fields` - (Optional) A `input_mapping_fields` block as defined below. Changing this forces a new resource to be created.

* `input_mapping_default_values` - (Optional) A `input_mapping_default_values` block as defined below. Changing this forces a new resource to be created.

* `public_network_access_enabled` - (Optional) Whether or not public network access is allowed for this server. Defaults to `true`.

* `local_auth_enabled` - (Optional) Whether local authentication methods is enabled for the EventGrid Topic. Defaults to `true`.

* `inbound_ip_rule` - (Optional) One or more `inbound_ip_rule` blocks as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Event Grid Topic. Possible values are `SystemAssigned`, `UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Event Grid Topic.

~> **Note:** This is required when `type` is set to `UserAssigned`

~> **Note:** When `type` is set to `SystemAssigned`, The assigned `principal_id` and `tenant_id` can be retrieved after the Event Grid Topic has been created. More details are available below.

---

A `input_mapping_fields` block supports the following:

* `id` - (Optional) Specifies the id of the EventGrid Event to associate with the domain. Changing this forces a new resource to be created.

* `topic` - (Optional) Specifies the topic of the EventGrid Event to associate with the domain. Changing this forces a new resource to be created.

* `event_type` - (Optional) Specifies the event type of the EventGrid Event to associate with the domain. Changing this forces a new resource to be created.

* `event_time` - (Optional) Specifies the event time of the EventGrid Event to associate with the domain. Changing this forces a new resource to be created.

* `data_version` - (Optional) Specifies the data version of the EventGrid Event to associate with the domain. Changing this forces a new resource to be created.

* `subject` - (Optional) Specifies the subject of the EventGrid Event to associate with the domain. Changing this forces a new resource to be created.

---

A `input_mapping_default_values` block supports the following:

* `event_type` - (Optional) Specifies the default event type of the EventGrid Event to associate with the domain. Changing this forces a new resource to be created.

* `data_version` - (Optional) Specifies the default data version of the EventGrid Event to associate with the domain. Changing this forces a new resource to be created.

* `subject` - (Optional) Specifies the default subject of the EventGrid Event to associate with the domain. Changing this forces a new resource to be created.

---

A `inbound_ip_rule` block supports the following:

* `ip_mask` - (Required) The IP mask (CIDR) to match on.

* `action` - (Optional) The action to take when the rule is matched. Possible values are `Allow`. Defaults to `Allow`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The EventGrid Topic ID.

* `endpoint` - The Endpoint associated with the EventGrid Topic.

* `primary_access_key` - The Primary Shared Access Key associated with the EventGrid Topic.

* `secondary_access_key` - The Secondary Shared Access Key associated with the EventGrid Topic.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the EventGrid Topic.
* `read` - (Defaults to 5 minutes) Used when retrieving the EventGrid Topic.
* `update` - (Defaults to 30 minutes) Used when updating the EventGrid Topic.
* `delete` - (Defaults to 30 minutes) Used when deleting the EventGrid Topic.

## Import

EventGrid Topic's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventgrid_topic.topic1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventGrid/topics/topic1
```
