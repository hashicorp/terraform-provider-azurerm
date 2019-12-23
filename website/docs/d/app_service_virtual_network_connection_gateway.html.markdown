---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_virtual_network_connection_gateway"
sidebar_current: "docs-azurerm-resource-app-service-virtual-network-connection-gateway"
description: |-
  Get Gateway Required Virtual Network Connection (Integrate with an App Service).

---

# Data Source: azurerm_app_service_virtual_network_connection_gateway

Use this data source to access information about an existing App Service integrated Gateway Required Virtual Network Connection.

## Example Usage

```hcl
data "azurerm_app_service_virtual_network_connection_gateway" "example" {
  resource_group_name   = "example-resource-group"
  app_service_name      = "example-appservice"
  virtual_network_name  = "example-virtual-network"
}

output "app_service_virtual_network_connection_gateway_id" {
  value = "${data.azurerm_app_service_virtual_network_connection_gateway.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group which the app service belongs to.

* `app_service_name` - (Required) Specifies the name of the App Service.

* `virtual_network_name` - (Required) Specifies the virtual network name that is connected with an app service.

## Attributes Reference

The following attributes are exported:

* `id` - The id of the App Service Virtual Network Connection.

* `name` - The name of the App Service Virtual Network Connection.

* `certificate_blob` - A certificate file (.cer) blob containing the public key of the private key used to authenticate a Point-To-Site VPN connection.

* `certificate_thumbprint` - The client certificate thumbprint.

* `dns_servers` - DNS servers to be used by this Virtual Network. It is a list of IP addresses.

* `resync_required` - is resync required

* `virtual_network_id` - The Virtual Network's resource ID.

* `routes` - (Array) One or more `route` block defined below.
  
---

A `route` block supports the following:
* `name` - Resource Name.
* `route_type` - The type of route this is: DEFAULT - By default, every app has routes to the local address ranges specified by RFC1918 INHERITED - Routes inherited from the real Virtual Network routes STATIC - Static route set on the app only. Valid values are `DEFAULT`, `INHERITED`, `STATIC`
* `start_address` The starting address for this route. This may also include a CIDR notation, in which case the end address must not be specified.
* `end_address` - The ending address for this route. If the start address is specified in CIDR notation, this must be omitted.
---
