---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_signalr_service"
description: |-
  Manages an Azure SignalR service.
---

# azurerm_signalr_service

Manages an Azure SignalR service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "terraform-signalr"
  location = "West Europe"
}

resource "azurerm_signalr_service" "example" {
  name                = "tfex-signalr"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    name     = "Free_F1"
    capacity = 1
  }

  cors {
    allowed_origins = ["http://www.example.com"]
  }

  public_network_access_enabled = false

  connectivity_logs_enabled = true
  messaging_logs_enabled    = true
  service_mode              = "Default"

  upstream_endpoint {
    category_pattern = ["connections", "messages"]
    event_pattern    = ["*"]
    hub_pattern      = ["hub1"]
    url_template     = "http://foo.com"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SignalR service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the SignalR service. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the SignalR service exists. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as documented below.

* `cors` - (Optional) A `cors` block as documented below.

* `connectivity_logs_enabled` - (Optional) Specifies if Connectivity Logs are enabled or not. Defaults to `false`.

* `messaging_logs_enabled` - (Optional) Specifies if Messaging Logs are enabled or not. Defaults to `false`.

* `http_request_logs_enabled` - (Optional) Specifies if Http Request Logs are enabled or not. Defaults to `false`.

* `identity` - (Optional) An `identity` block as defined below.

* `public_network_access_enabled` - (Optional) Whether to enable public network access? Defaults to `true`.

~> **Note:** `public_network_access_enabled` cannot be set to `false` in `Free` sku tier.

* `local_auth_enabled` - (Optional) Whether to enable local auth? Defaults to `true`.

* `aad_auth_enabled` - (Optional) Whether to enable AAD auth? Defaults to `true`.

* `tls_client_cert_enabled` - (Optional) Whether to request client certificate during TLS handshake? Defaults to `false`.

~> **Note:** `tls_client_cert_enabled` cannot be set to `true` in `Free` sku tier.

* `serverless_connection_timeout_in_seconds` - (Optional) Specifies the client connection timeout. Defaults to `30`.

* `service_mode` - (Optional) Specifies the service mode. Possible values are `Classic`, `Default` and `Serverless`. Defaults to `Default`.

* `upstream_endpoint` - (Optional) An `upstream_endpoint` block as documented below. Using this block requires the SignalR service to be Serverless. When creating multiple blocks they will be processed in the order they are defined in.

* `live_trace` - (Optional) A `live_trace` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `cors` block supports the following:

* `allowed_origins` - (Required) A list of origins which should be able to make cross-origin calls. `*` can be used to allow all calls.

---

An `upstream_endpoint` block supports the following:

* `url_template` - (Required) The upstream URL Template. This can be a url or a template such as `http://host.com/{hub}/api/{category}/{event}`.

* `category_pattern` - (Required) The categories to match on, or `*` for all.

* `event_pattern` - (Required) The events to match on, or `*` for all.

* `hub_pattern` - (Required) The hubs to match on, or `*` for all.

* `user_assigned_identity_id` - (Optional) Specifies the Managed Identity IDs to be assigned to this signalR upstream setting by using resource uuid as both system assigned and user assigned identity is supported. 

---

A `live_trace` block supports the following:

* `enabled` - (Optional) Whether the live trace is enabled? Defaults to `true`.

* `messaging_logs_enabled` - (Optional) Whether the log category `MessagingLogs` is enabled? Defaults to `true`

* `connectivity_logs_enabled` - (Optional) Whether the log category `ConnectivityLogs` is enabled? Defaults to `true`

* `http_request_logs_enabled` - (Optional) Whether the log category `HttpRequestLogs` is enabled? Defaults to `true`

---

A `sku` block supports the following:

* `name` - (Required) Specifies which tier to use. Valid values are `Free_F1`, `Standard_S1`, `Premium_P1` and `Premium_P2`.

* `capacity` - (Required) Specifies the number of units associated with this SignalR service. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `20`, `30`, `40`, `50`, `60`, `70`, `80`, `90`, `100`, `200`, `300`, `400`, `500`, `600`, `700`, `800`, `900` and `1000`.

~> **Note:** The valid capacity range for sku `Free_F1` is `1`, for sku `Premium_P2` is from `100` to `1000`, and from `1` to `100` for sku `Standard_S1` and `Premium_P1`.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this signalR. Possible values are `SystemAssigned`, `UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this signalR.

~> **Note:** This is required when `type` is set to `UserAssigned`

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the SignalR service.

* `hostname` - The FQDN of the SignalR service.

* `ip_address` - The publicly accessible IP of the SignalR service.

* `public_port` - The publicly accessible port of the SignalR service which is designed for browser/client use.

* `server_port` - The publicly accessible port of the SignalR service which is designed for customer server side use.

* `primary_access_key` - The primary access key for the SignalR service.

* `primary_connection_string` - The primary connection string for the SignalR service.

* `secondary_access_key` - The secondary access key for the SignalR service.

* `secondary_connection_string` - The secondary connection string for the SignalR service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the SignalR Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the SignalR Service.
* `update` - (Defaults to 30 minutes) Used when updating the SignalR Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the SignalR Service.

## Import

SignalR services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_signalr_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/terraform-signalr/providers/Microsoft.SignalRService/signalR/tfex-signalr
```
