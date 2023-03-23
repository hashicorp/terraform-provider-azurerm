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

## Argument Reference

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

* `primary_access_key` - The primary access key of the SignalR service.

* `primary_connection_string` - The primary connection string of the SignalR service.

* `secondary_access_key` - The secondary access key of the SignalR service.

* `secondary_connection_string` - The secondary connection string of the SignalR service.

* `public_network_access_enabled` - Is public network access enabled for this SignalR service?

* `local_auth_enabled` - Is local auth enable for this SignalR serviced?

* `aad_auth_enabled` - Is aad auth enabled for this SignalR service?

* `tls_client_cert_enabled` - Is tls client cert enabled for this SignalR service?

* `serverless_connection_timeout_in_seconds` - The serverless connection timeout of this SignalR service.

* `identity` - An `identity` block as documented below.

---

The `identity` block exports the following:

* `type` - The type of identity used for the signalR service.

* `user_assigned_identity_id` - The ID of the User Assigned Identity. This value will be empty when using system assigned identity.

* `principal_id` - The principal id of the system assigned identity.

* `tenant_id` - The tenant id of the system assigned identity.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the SignalR service.
