---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_signalr_service"
sidebar_current: "docs-azurerm-datasource-signalr-service"
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

* `name` - (Required) Specifies the name of the SignalR service.

* `resource_group_name` - (Required) Specifies the name of the resource group the SignalR service is located in.

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
