---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_link_endpoint"
sidebar_current: "docs-azurerm-datasource-private-endpoint"
description: |-
  Gets information about an existing Private Link Endpoint
---

# Data Source: azurerm_private_link_endpoint

Use this data source to access information about an existing Private Link Endpoint.

## Example Usage

```hcl
data "azurerm_private_link_endpoint" "example" {
  resource_group_name = "example-rg"
  name                = "example-private-endpoint"
}

output "subnet_id" {
  value = "${data.azurerm_private_link_endpoint.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the Name of the `Private Link Endpoint`.
* `resource_group_name` - (Required) Specifies the Name of the Resource Group within which the `Private Link Endpoint` exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Prviate Link Endpoint.
* `location` - The supported Azure location where the resource exists.
* `network_interface_ids` - A list of network interfaces IDs.
* `subnet_id` - The subnet ID.
* `tags` - A mapping of tags assigned to the resource.

A `private_service_connection` block exports the following:

* `request_response` - Possible values are as follows:
   * `Auto-Approved` - The remote resource owner has added you to the `Auto-Approved` RBAC permission list for the remote resource, all private link endpoint connection requests will be automatically `Approved`.
   * `Deleted state` - The resource owner has `Rejected` the private link endpoint connection request and has removed your private link endpoint request from the remote resource.
   * `request/response message` - If you submitted a manual private link endpoint connection request, while in the `Pending` status the `request_response` will display the same text from your `request_message` in the `private_service_connection` block. If the private link endpoint connection request was `Rejected` by the owner of the remote resource, the text for the rejection will be displayed as the `request_response` text, if If the private link endpoint connection request was `Approved` by the owner of the remote resource, the text for the approval will be displayed as the `request_response` text
* `status` - The current status of the private link endpoint request, possible values will be `Pending`, `Approved`, `Rejected`, or `Disconnected`.
* `private_ip_address` - The private IP address associated with the private link endpoint, note that you will have a private IP address assigned to the private link endpoint even if the connection request was `Rejected`.
