---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_link_endpoint"
sidebar_current: "docs-azurerm-datasource-private-endpoint"
description: |-
  Gets information about an existing Private Endpoint
---

# Data Source: azurerm_private_link_endpoint

Use this data source to access information about an existing Private Endpoint.

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

* `name` - (Required) The name of the private endpoint.

* `resource_group_name` - (Required) The Name of the Resource Group where the App Service exists.

## Attributes Reference

The following attributes are exported:

* `location` - Resource location.

* `manual_private_link_service_connections` - A grouping of information about the connection to the remote resource. Used when the network admin does not have access to approve connections to the remote resource. One or more `manual_private_link_service_connection` block defined below.

* `network_interfaces` - Gets an array of references to the network interfaces created for this private endpoint. One or more `network_interface` block defined below.

* `private_link_service_connections` - A grouping of information about the connection to the remote resource. One or more `private_link_service_connection` block defined below.

* `subnet_id` - The ID of the subnet from which the private IP will be allocated.

* `tags` - Resource tags.

---

The `manual_private_link_service_connection` block contains the following:

* `private_link_service_id` - The resource id of private link service.

* `group_ids` - The ID(s) of the group(s) obtained from the remote resource that this private endpoint should connect to.

* `request_message` - A message passed to the owner of the remote resource with this connection request. Restricted to 140 chars.

* `name` - The name of the resource that is unique within a resource group. This name can be used to access the resource.

* `status` - Indicates whether the connection has been Approved/Rejected/Removed by the owner of the service.

---

The `network_interface` block contains the following:

* `id` - Resource ID.

---

The `private_link_service_connection` block contains the following:

* `private_link_service_id` - The resource id of private link service.

* `group_ids` - The ID(s) of the group(s) obtained from the remote resource that this private endpoint should connect to.

* `request_message` - A message passed to the owner of the remote resource with this connection request. Restricted to 140 chars.

* `name` - The name of the resource that is unique within a resource group. This name can be used to access the resource.

* `status` - Indicates whether the connection has been Approved/Rejected/Removed by the owner of the service.
