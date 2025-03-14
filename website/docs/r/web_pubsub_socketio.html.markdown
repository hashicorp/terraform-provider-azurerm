---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_web_pubsub_socketio"
description: |-
  Manages a Web PubSub Service for Socket.IO.
---

# azurerm_web_pubsub_socketio

Manages a Web PubSub Service for Socket.IO.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_web_pubsub_socketio" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Free_F1"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Web PubSub Service. Changing this forces a new Web PubSub Service to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Web PubSub Service should exist. Changing this forces a new Web PubSub Service to be created.

* `location` - (Required) The Azure Region where the Web PubSub Service should exist. Changing this forces a new Web PubSub Service to be created.

* `sku` - (Required) The SKU to use for this Web PubSub Service. Possible values are `Free_F1`, `Standard_S1`, `Premium_P1`, and `Premium_P2`.

---

* `aad_auth_enabled` - (Optional) Whether Azure Active Directory authentication is enabled. Defaults to `true`.

* `capacity` - (Optional) The number of units associated with this Web PubSub Service. Defaults to `1`. Possible values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `20`, `30`, `40`, `50`, `60`, `70`, `80`, `90`, `100`, `200`, `300`, `400`, `500`, `600`, `700`, `800`, `900` and `1000`.

~> **Note:** The valid range depends on which `sku` is used. For `Free_F1` only `1` is supported, for `Standard_S1` and `Premium_P1` `1` through `100` are supported, and for `Premium_P2` the minimum capacity is `100`.

* `identity` - (Optional) An `identity` block as defined below.

* `live_trace_enabled` - (Optional) Whether the live trace tool is enabled. Defaults to `true`.

* `live_trace_connectivity_logs_enabled` - (Optional) Whether the connectivity log category for live trace is enabled. Defaults to `true`.

* `live_trace_http_request_logs_enabled` - (Optional) Whether the HTTP request log category for live trace is enabled. Defaults to `true`.

* `live_trace_messaging_logs_enabled` - (Optional) Whether the messaging log category for live trace is enabled. Defaults to `true`.

* `local_auth_enabled` - (Optional) Whether local authentication using an access key is enabled. Defaults to `true`.

* `public_network_access` - (Optional) Whether public network access is enabled. Defaults to `Enabled`. Possible values are `Enabled` and `Disabled`.

~> **Note:** `public_network_access` cannot be set to `Disabled` when `sku` is `Free_F1`.

* `service_mode` - (Optional) The service mode of this Web PubSub Service. Defaults to `Default`. Possible values are `Default` and `Serverless`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Web PubSub Service.

* `tls_client_cert_enabled` - (Optional) Whether the service should request a client certificate during a TLS handshake. Defaults to `false`.

~> **Note:** `tls_client_cert_enabled` cannot be set to `true` when `sku` is `Free_F1`.

---

A `identity` block supports the following:

* `type` - (Required) The type of Managed Identity for this Web PubSub Service. Possible Values are `SystemAssigned` and `UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs for this Web PubSub Service.

~> **Note:** `identity_ids` is required when `type` is `UserAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Web PubSub Service.

* `external_ip` - The publicly accessible IP address of the Web PubSub Service.

* `hostname` - The FQDN of the Web PubSub Service.

* `primary_access_key` - The primary access key for the Web PubSub Service.

* `primary_connection_string` - The primary connection string for the Web PubSub Service.

* `public_port` - The publicly accessible port for client-side usage of the Web PubSub Service.

* `secondary_access_key` - The secondary access key for the Web PubSub Service.

* `secondary_connection_string` - The secondary connection string for the Web PubSub Service.

* `server_port` - The publicly accessible port for server-side usage of the Web PubSub Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Web PubSub Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Web PubSub Service.
* `update` - (Defaults to 30 minutes) Used when updating the Web PubSub Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Web PubSub Service.

## Import

Web PubSub Service for Socket.IOs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_web_pubsub_socketio.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.SignalRService/webPubSub/pubsub1
```
