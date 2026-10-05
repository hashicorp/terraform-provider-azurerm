---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_web_pubsub"
description: |-
  Gets information about an existing Azure Web Pubsub service.
---

# Data Source: azurerm_web_pubsub

Use this data source to access information about an existing Azure Web Pubsub service.

## Example Usage

```hcl
data "azurerm_web_pubsub" "example" {
  name                = "test-webpubsub"
  resource_group_name = "wps-resource-group"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - Specifies the name of the Web Pubsub service.

* `resource_group_name` - Specifies the name of the resource group the Web Pubsub service is located in.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Web Pubsub service.

* `hostname` - The FQDN of the Web Pubsub service.

* `ip_address` - The publicly accessible IP of the Web Pubsub service.

* `identity` - An `identity` block as documented below.

* `location` - The Azure location where the Web Pubsub service exists.

* `live_trace` - A `live_trace` block as documented below.

* `public_port` - The publicly accessible port of the Web Pubsub service which is designed for browser/client use.

* `server_port` - The publicly accessible port of the Web Pubsub service which is designed for customer server side use.

* `primary_access_key` - The primary access key of the Web Pubsub service.

* `primary_connection_string` - The primary connection string of the Web Pubsub service.

* `secondary_access_key` - The secondary access key of the Web Pubsub service.

* `secondary_connection_string` - The secondary connection string of the Web Pubsub service.

---

A `live_trace` block exports the following:

* `enabled` - Whether the live trace is enabled.

* `messaging_logs_enabled` - Whether the log category `MessagingLogs` is enabled.

* `connectivity_logs_enabled` - Whether the log category `ConnectivityLogs` is enabled.

* `http_request_logs_enabled` - Whether the log category `HttpRequestLogs` is enabled.

---

The `identity` block exports the following:

* `type` - The type of identity used for the web pubsub.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to the Web Pubsub service.

* `principal_id` - The principal id of the system assigned identity.

* `tenant_id` - The tenant id of the system assigned identity.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Web Pubsub service.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.SignalRService` - 2024-03-01
