---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_link_endpoint_connection"
sidebar_current: "docs-azurerm-datasource-private-endpoint-connection"
description: |-
  Gets the connecton status information about an existing Private Link Endpoint
---
-> **NOTE:** The 'azurerm_private_link_endpoint_connection' resource is being deprecated in favour of the renamed version 'azurerm_private_endpoint_connection'.
Information on migrating to the renamed resource can be found here: https://terraform.io/docs/providers/azurerm/guides/migrating-between-renamed-resources.html
As such the existing 'azurerm_private_link_endpoint_connection' resource is deprecated and will be removed in the next major version of the AzureRM Provider (2.0).

# Data Source: azurerm_private_link_endpoint_connection

Use this data source to access the connection status information about an existing Private Link Endpoint.

-> **NOTE** Private Link is currently in Public Preview.

## Example Usage

```hcl
data "azurerm_private_link_endpoint_connection" "example" {
  name                = "example-private-endpoint"
  resource_group_name = "example-rg"
}

output "private_link_endpoint_status" {
  value = data.azurerm_private_link_endpoint_connection.example.private_service_connection.0.status
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the Name of the private link endpoint.
* `resource_group_name` - (Required) Specifies the Name of the Resource Group within which the private link endpoint exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Azure resource ID of the Prviate Link Endpoint.
* `location` - The supported Azure location where the resource exists.

A `private_service_connection` block exports the following:

* `name` - The name of the private linke endpoint.
* `status` - The current status of the private link endpoint request, possible values will be `Pending`, `Approved`, `Rejected`, or `Disconnected`.
* `private_ip_address` - The private IP address associated with the private link endpoint, note that you will have a private IP address assigned to the private link endpoint even if the connection request was `Rejected`.
* `request_response` - Possible values are as follows:
  Value | Meaning
  -- | --
  `Auto-Approved` | The remote resource owner has added you to the `Auto-Approved` RBAC permission list for the remote resource, all private link endpoint connection requests will be automatically `Approved`.
  `Deleted state` | The resource owner has `Rejected` the private link endpoint connection request and has removed your private link endpoint request from the remote resource.
  `request/response message` | If you submitted a manual private link endpoint connection request, while in the `Pending` status the `request_response` will display the same text from your `request_message` in the `private_service_connection` block above. If the private link endpoint connection request was `Rejected` by the owner of the remote resource, the text for the rejection will be displayed as the `request_response` text, if the private link endpoint connection request was `Approved` by the owner of the remote resource, the text for the approval will be displayed as the `request_response` text
