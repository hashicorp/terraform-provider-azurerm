---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_link_service_endpoint_connections"
sidebar_current: "docs-azurerm-datasource-private-link-service-endpoint-connections"
description: |-
  Use this data source to access endpoint connection information about an existing Private Link Service.
---

# Data Source: azurerm_private_link_service_endpoint_connections

Use this data source to access endpoint connection information about an existing Private Link Service.

-> **NOTE** Private Link is currently in Public Preview.

## Example Usage

```hcl
data "azurerm_private_link_service_endpoint_connections" "example" {
  service_id          = azurerm_private_link_service.example.id
  resource_group_name = azurerm_resource_group.example.name
}

output "private_endpoint_status" {
  value = data.azurerm_private_link_service_endpoint_connections.example.private_endpoint_connections.0.status
}
```


## Argument Reference

The following arguments are supported:

* `service_id` - (Required) The resource ID of the private link service.

* `resource_group_name` - (Required) The name of the resource group in which the private link service resides.


## Attributes Reference

* `service_name` - The name of the private link service.

The `private_endpoint_connections` block exports the following:

* `connection_id` - The resource id of the private link service connection between the private link service and the private link endpoint.

* `connection_name` - The name of the connection between the private link service and the private link endpoint.

* `private_endpoint_id` - The resource id of the private link endpoint.

* `private_endpoint_name` - The name of the private link endpoint.

* `action_required` - A message indicating if changes on the service provider require any updates or not.

* `description` -  The request for approval message or the reason for rejection message.

* `status` - Indicates the state of the connection between the private link service and the private link endpoint, possible values are `Pending`, `Approved` or `Rejected`.
