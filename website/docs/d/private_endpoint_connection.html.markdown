---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_endpoint_connection"
description: |-
  Gets the connection status information about an existing Private Endpoint
---

# Data Source: azurerm_private_endpoint_connection

Use this data source to access the connection status information about an existing Private Endpoint Connection.

## Example Usage

```hcl
data "azurerm_private_endpoint_connection" "example" {
  name                = "example-private-endpoint"
  resource_group_name = "example-rg"
}

output "private_endpoint_status" {
  value = data.azurerm_private_endpoint_connection.example.private_service_connection.0.status
}
```

## Argument Reference

The following arguments are supported:

* `name` - Specifies the Name of the private endpoint.
* `resource_group_name` - Specifies the Name of the Resource Group within which the private endpoint exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Private Endpoint.

* `location` - The supported Azure location where the resource exists.

A `private_service_connection` block exports the following:

* `name` - The name of the private endpoint.
* `status` - The current status of the private endpoint request, possible values will be `Pending`, `Approved`, `Rejected`, or `Disconnected`.
* `private_ip_address` - The private IP address associated with the private endpoint, note that you will have a private IP address assigned to the private endpoint even if the connection request was `Rejected`.
* `request_response` - Possible values are as follows:
  Value | Meaning
  -- | --
  `Auto-Approved` | The remote resource owner has added you to the `Auto-Approved` RBAC permission list for the remote resource, all private endpoint connection requests will be automatically `Approved`.
  `Deleted state` | The resource owner has `Rejected` the private endpoint connection request and has removed your private endpoint request from the remote resource.
  `request/response message` | If you submitted a manual private endpoint connection request, while in the `Pending` status the `request_response` will display the same text from your `request_message` in the `private_service_connection` block above. If the private endpoint connection request was `Rejected` by the owner of the remote resource, the text for the rejection will be displayed as the `request_response` text, if the private endpoint connection request was `Approved` by the owner of the remote resource, the text for the approval will be displayed as the `request_response` text

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Private Endpoint.
