---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_signalr_service"
description: |-
  Gets information about an existing Azure SignalR service.
---

# Data Source: azurerm_signalr_service

Use this data source to access information about an existing Azure SignalR service.

## Example Usage

```hcl
data "azurerm_signalr_service" "example" {
  name                = "test-signalr"
  resource_group_name = "signalr-resource-group"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - Specifies the name of the SignalR service.

* `resource_group_name` - Specifies the name of the resource group the SignalR service is located in.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the SignalR service.

* `hostname` - The FQDN of the SignalR service.

* `ip_address` - The publicly accessible IP of the SignalR service.

* `location` - Specifies the supported Azure location where the SignalR service exists.

* `public_port` - The publicly accessible port of the SignalR service which is designed for browser/client use.

* `server_port` - The publicly accessible port of the SignalR service which is designed for customer server side use.

* `cors` - A `cors` block as documented below.

* `connectivity_logs_enabled` - Specifies if Connectivity Logs are enabled or not.

* `messaging_logs_enabled` - Specifies if Messaging Logs are enabled or not.

* `http_request_logs_enabled` - Specifies if Http Request Logs are enabled or not.

* `primary_access_key` - The primary access key of the SignalR service.

* `primary_connection_string` - The primary connection string of the SignalR service.

* `secondary_access_key` - The secondary access key of the SignalR service.

* `secondary_connection_string` - The secondary connection string of the SignalR service.

* `sku` - A `sku` block as documented below.

* `public_network_access_enabled` - Is public network access enabled for this SignalR service?

* `local_auth_enabled` - Is local auth enable for this SignalR serviced?

* `aad_auth_enabled` - Is aad auth enabled for this SignalR service?

* `tls_client_cert_enabled` - Is tls client cert enabled for this SignalR service?

* `serverless_connection_timeout_in_seconds` - The serverless connection timeout of this SignalR service.

* `service_mode` - Specifies the service mode.

* `upstream_endpoint` - One or more `upstream_endpoint` blocks as documented below.

* `live_trace` - A `live_trace` block as defined below.

* `identity` - An `identity` block as documented below.

---

The `sku` block exports the following:

* `name` - The name of the SKU.

* `capacity` - The capacity of the SKU.

---

A `cors` block exports the following:

* `allowed_origins` - A list of origins which should be able to make cross-origin calls.

---

An `upstream_endpoint` block exports the following:

* `url_template` - The upstream URL Template.

* `category_pattern` - The categories to match on, or `*` for all.

* `event_pattern` - The events to match on, or `*` for all.

* `hub_pattern` - The hubs to match on, or `*` for all.

* `user_assigned_identity_id` - The Managed Identity ID assigned to this SignalR upstream setting.

---

A `live_trace` block exports the following:

* `enabled` - Whether live trace is enabled.

* `messaging_logs_enabled` - Whether the log category `MessagingLogs` is enabled.

* `connectivity_logs_enabled` - Whether the log category `ConnectivityLogs` is enabled.

* `http_request_logs_enabled` - Whether the log category `HttpRequestLogs` is enabled.

---

The `identity` block exports the following:

* `type` - The type of identity used for the signalR service.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to the SignalR service.

* `principal_id` - The principal id of the system assigned identity.

* `tenant_id` - The tenant id of the system assigned identity.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the SignalR service.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.SignalRService` - 2024-03-01
