---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_link_service_endpoint_connection_status"
sidebar_current: "docs-azurerm-datasource-private-link-service-endpoint-connection-status"
description: |-
  Use this data source to access endpoint connection information about an existing Private Link Service.
---

# Data Source: azurerm_private_link_service_endpoint_connection_status

Use this data source to access endpoint connection information about an existing Private Link Service.


## Private Link Service Usage

```hcl
data "azurerm_private_link_service_endpoint_connection_status" "example" {
  name                = azurerm_private_link_service.example.name
  resource_group_name = azurerm_resource_group.example.name
}

output "connection_name" {
  value = data.azurerm_private_link_service_endpoint_connection_status.example.private_endpoint_connections.0.connection_name
}

output "private_endpoint_endpoint_name" {
  value = data.azurerm_private_link_service_endpoint_connection_status.example.private_endpoint_connections.0.private_endpoint_name
}

output "private_endpoint_description" {
  value = data.azurerm_private_link_service_endpoint_connection_status.example.private_endpoint_connections.0.description
}

output "private_endpoint_status" {
  value = data.azurerm_private_link_service_endpoint_connection_status.example.private_endpoint_connections.0.status
}

output "private_endpoint_action_required" {
  value = data.azurerm_private_link_service_endpoint_connection_status.example.private_endpoint_connections.0.action_required
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the private link service.

* `resource_group_name` - (Required) The name of the resource group in which the private link service resides.


## Attributes Reference

The `private_endpoint_connections` block exports the following:

* `connection_id` - The resource id of the private link service connection between the private link service and the private link endpoint.

* `connection_name` - The name of the connection between the private link service and the private link endpoint.

* `private_endpoint_id` - The resource id of the private link endpoint.

* `private_endpoint_name` - The name of the private link endpoint.

* `action_required` - A message indicating if changes on the service provider require any updates or not.

* `description` -  The request for approval message or the reason for rejection message.

* `status` - Indicates the state of the connection between the private link service and the private link endpoint, possible values are `Pending`, `Approved` or `Rejected`.
