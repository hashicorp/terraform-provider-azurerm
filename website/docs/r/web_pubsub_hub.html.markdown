---
subcategory: "Web PubSub"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_web_pubsub_hub"
description: |-
  Manages the hub settings for a Web Pubsub service.
---

# azurerm_web_pubsub_hub

Manages the hub settings for a Web Pubsub.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "terraform-webpubsub"
  location = "east us"
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "tfex-uai"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_web_pubsub" "example" {
  name                = "tfex-webpubsub"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku      = "Standard_S1"
  capacity = 1
}

resource "azurerm_web_pubsub_hub" "example" {
  name          = "tfex_wpsh"
  web_pubsub_id = azurerm_web_pubsub.example.id
  event_handler {
    url_template       = "https://test.com/api/{hub}/{event}"
    user_event_pattern = "*"
    system_events      = ["connect", "connected"]
  }

  event_handler {
    url_template       = "https://test.com/api/{hub}/{event}"
    user_event_pattern = "event1, event2"
    system_events      = ["connected"]
    auth {
      managed_identity_id = azurerm_user_assigned_identity.example.id
    }
  }
  anonymous_connections_enabled = true

  depends_on = [
    azurerm_web_pubsub.example
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Web Pubsub hub service. Changing this forces a new resource to be created.

* `web_pubsub_id` - (Required) Specify the id of the Web Pubsub. Changing this forces a new resource to be created.

* `anonymous_connections_enabled` - (Optional) Is anonymous connections are allowed for this hub? Defaults to `false`.
  Possible values are `true`, `false`.

* `event_handler` - (Optional) An `event_handler` block as defined below.

---

An `event_handler` block supports the following:

* `url_template` - (Required) The Event Handler URL Template. Two predefined parameters `{hub}` and `{event}` are
  available to use in the template. The value of the EventHandler URL is dynamically calculated when the client request
  comes in. Example: `http://example.com/api/{hub}/{event}`.

* `user_event_pattern` - (Optional) Specify the matching event names. There are 3 kind of patterns supported:
    - `*` matches any event name
    - `,` Combine multiple events with `,` for example `event1,event2`, it matches event `event1` and `event2`
    - The single event name, for example `event1`, it matches `event1`.

* `system_events` - (Optional) Specify the list of system events. Supported values are `connect`, `connected`
  and `disconnected`.

* `auth` - (Optional) An `auth` block as defined below.

---

An `auth` block supports the following:

* `managed_identity_id` - (Required) Specify the identity ID of the target resource.

-> **NOTE:** `managed_identity_id` is required if the auth block is defined

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Web Pubsub Hub resource.

* `name` - The name of the Web Pubsub Hub resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Web Pubsub Resource.
* `update` - (Defaults to 30 minutes) Used when updating the Web Pubsub Resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the Web Pubsub Resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the Web Pubsub Resource.

## Import

Web Pubsub Hub can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_web_pubsub_hub.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.SignalRService/webPubsub/webpubsub1/hubs/webpubsubhub1
```
