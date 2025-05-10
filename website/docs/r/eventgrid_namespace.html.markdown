---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_namespace"
description: |-
  Manages an EventGrid Namespace

---

# azurerm_eventgrid_namespace

Manages an EventGrid Namespace

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_eventgrid_namespace" "example" {
  name                = "my-eventgrid-namespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Event Grid Namespace resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Event Grid Namespace should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource should exist. Changing this forces a new resource to be created.

* `capacity` - (Optional) Specifies the Capacity / Throughput Units for an Eventgrid Namespace. Valid values can be between `1` and `40`.

* `identity` - (Optional) An `identity` block as defined below.

* `inbound_ip_rule` - (Optional) One or more `inbound_ip_rule` blocks as defined below.

* `public_network_access` - (Optional) Whether or not public network access is allowed for this server. Defaults to `Enabled`.

* `sku` - (Optional) Defines which tier to use. The only possible value is `Standard`. Defaults to `Standard`.

* `topic_spaces_configuration` - (Optional) A `topic_spaces_configuration` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Event Grid Namespace. Possible values are `SystemAssigned`, `UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Event Grid Namespace.

~> **Note:** This is required when `type` is set to `UserAssigned`

---

An `inbound_ip_rule` block supports the following:

* `ip_mask` - (Required) The IP mask (CIDR) to match on.

* `action` - (Optional) The action to take when the rule is matched. Possible values are `Allow`. Defaults to `Allow`.

---

A `topic_spaces_configuration` block supports the following:

* `alternative_authentication_name_source` - (Optional) Specifies a list of alternative sources for the client authentication name from the client certificate. Possible values are `ClientCertificateDns`, `ClientCertificateEmail`, `ClientCertificateIp`, `ClientCertificateSubject` and `ClientCertificateUri`.

* `maximum_client_sessions_per_authentication_name` - (Optional) Specifies the maximum number of client sessions per authentication name. Valid values can be between `1` and `100`.

* `maximum_session_expiry_in_hours` - (Optional) Specifies the maximum session expiry interval allowed for all MQTT clients connecting to the Event Grid namespace. Valid values can be between `1` and `8`.

* `route_topic_id` - (Optional) Specifies the Event Grid topic resource ID to route messages to.

* `dynamic_routing_enrichment` - One or more `dynamic_routing_enrichment` blocks as defined below.

* `static_routing_enrichment` - One or more `static_routing_enrichment` blocks as defined below.

---

A `dynamic_routing_enrichment` block supports the following:

* `key` - (Required) The enrichment key.

* `value` - (Required) The enrichment value.

---

A `static_routing_enrichment` block supports the following:

* `key` - (Required) The enrichment key.

* `value` - (Required) The enrichment value.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The EventGrid Namespace ID.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the EventGrid Namespace.
* `read` - (Defaults to 5 minutes) Used when retrieving the EventGrid Namespace.
* `update` - (Defaults to 30 minutes) Used when updating the EventGrid Namespace.
* `delete` - (Defaults to 30 minutes) Used when deleting the EventGrid Namespace.

## Import

EventGrid Namespace's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventgrid_namespace.namespace1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventGrid/namespaces/namespace1
```
