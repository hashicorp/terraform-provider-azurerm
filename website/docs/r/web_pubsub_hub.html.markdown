---
subcategory: "Web Pubsub"
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

resource "azurerm_user_assigned_identity" "test" {
  name                = "tfex-uai"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_web_pubsub" "example" {
  name                = "tfex-webpubsub"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    name     = "Standard_S1"
    capacity = 1
  }
}

resource "azurerm_web_pubsub_hub" "test" {
  name                = "tfex-wpsh"
  web_pubsub_name     = azurerm_web_pubsub.exmaple.name
  resource_group_name = azurerm_resource_group.example.name
  event_handler {
    url_template       = "https://test.com/api/{hub}/{event}"
    user_event_pattern = "event1"
    system_events      = ["connect", "connected"]
    auth {
      type                      = "ManagedIdentity"
      managed_identity_resource = azurerm_user_assigned_identity.example.id
    }
  }
}
```

##Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Web Pubsub hub service. Changing this forces a new resource to be created.

* `web_pubsub_name` - (Required) The default action to control the network access when no other rule matches. Possible values are `Allow` and `Deny`.

* `resource_group_name` - (Required) The name of the resource group in which to create the Web Pubsub Hub resource. Changing this forces a new resource to be created.

* `event_handler` - (Required) An `event_handler` block as defined below.

* `anonymous_connect_policy` - (Optional) Is anonymous connections are allowed for this hub? Defaults to `Deny`. Possible value are `Allow`, `Deny`

---

An `event_handler` block supports the following:

* `url_template` - (Required) The Event Handler URL Template. You can use a predefined parameter {hub} and {event} inside the template, the value of the EventHandler URL is dynamically calculated when the client request comes in. For example, UrlTemplate can be `http://example.com/api/{hub}/{event}`. 

-> **NOTE:** The host part can't contain parameters.

* `user_event_pattern` - (Optional) Specify the matching event names. There are 3 kind of patterns supported:
    - `*` matches any event name 
    - `,` Combine multiple events with `,` for example `event1,event2`, it matches event `event1` and `event2`
    - The single event name, for example `event1`, it matches `event1`.

* `system_events` - (Optional) Specify the list of system events. Supported values are `connect`, `connected` and `disconnected`.

* `auth` - (Optional) An `auth` block as defined below.

---

An `auth` block supports the following:

* `type` - (Optional) Specify the auth type. Possible values are `None`, `ManagedIdentity`. Defaults to `None`

* `managed_identity_resource` - (Optional) Specify the App ID URI of the target resource. It also appears in the aud (audience) claim of the issued token.

-> **NOTE:** `managed_identity_resource` is required if the auth type is set to ManagedIdentity

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
terraform import azurerm_web_pubsub_hub.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/terraform-webpubsub/providers/Microsoft.SignalRService/webPubsub/tfex-webpubsub/hubs/tfex-wpsh
```
