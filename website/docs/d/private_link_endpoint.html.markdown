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

* `manual_private_link_service_connection` - A grouping of information about the connection to the remote resource. One or more `manual_private_link_service_connection` block as defined below.

* `private_link_service_connection` - A grouping of information about the connection to the remote resource. One or more `private_link_service_connection` block as defined below.

* `network_interface_ids` - A list of network interfaces created for this private endpoint.

* `subnet_id` - The ID of the subnet from which the private IP will be allocated.

* `tags` - A mapping of tags assigned to the resource.

---

The `manual_private_link_service_connection` block exports the following:

* `name` - Specifies the name of the `manual_private_link_service_connection`. This name can be used to access the resource.

* `private_link_service_id` - The resource id of the Private Link Service.

* `group_ids` - The ID(s) of the group(s) obtained from the remote resource that this private endpoint should connect to.

* `request_message` - A message passed to the owner of the remote resource with this connection request. Restricted to `140` chars.

* `state_action_required` - A message indicating if changes on the `Private Link Service` provider require any updates on the `Private Link Endpoint`.

* `state_description` - The reason for `approval`/`rejection` of the connection.

* `state_status` - Indicates whether the `Private Link Service` connection has been `Approved`, `Rejected` or `Removed` by the owner of the `Private Link Service`.

---

The `private_link_service_connection` block exports the following:

* `name` - Specifies the name of the `private_link_service_connection`. This name can be used to access the resource.

* `private_link_service_id` - The resource id of Private Link Service.

* `group_ids` - The ID(s) of the group(s) obtained from the remote resource that this private endpoint should connect to.

* `request_message` - A message passed to the owner of the remote resource with this connection request. Restricted to `140` chars.

* `state_action_required` - A message indicating if changes on the `Private Link Service` provider require any updates on the `Private Link Endpoint`.

* `state_description` - The reason for `approval`/`rejection` of the connection.

* `state_status` - Indicates whether the `Private Link Service` connection has been `Approved`, `Rejected` or `Removed` by the owner of the `Private Link Service`.

