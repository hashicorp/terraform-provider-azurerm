---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_link_service"
sidebar_current: "docs-azurerm-datasource-private-link-service"
description: |-
  Gets information about an existing Private Link Service
---

# Data Source: azurerm_private_link_service

Use this data source to access information about an existing Private Link Service.


## Private Link Service Usage

```hcl
data "azurerm_private_link_service" "example" {
  resource_group_name = "acctestRG"
  name                = "acctestpls"
}

output "private_link_service_id" {
  value = "${data.azurerm_private_link_service.example.id}"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the private link service.

* `resource_group_name` - (Required) The name of the resource group.


## Attributes Reference

The following attributes are exported:

* `location` - Resource location.

* `alias` - The alias of the private link service.

* `auto_approval` - One `auto_approval` block defined below.

* `fqdns` - The list of Fqdn.

* `ip_configurations` - One or more `ip_configuration` block defined below.

* `load_balancer_frontend_ip_configurations` - One or more `load_balancer_frontend_ip_configuration` block defined below.

* `network_interfaces` - One or more `network_interface` block defined below.

* `private_endpoint_connections` - One or more `private_endpoint_connection` block defined below.

* `type` - Resource type.

* `visibility` - One `visibility` block defined below.

* `tags` - Resource tags.


---

The `auto_approval` block contains the following:

* `subscriptions` - The list of subscriptions.

---

The `ip_configuration` block contains the following:

* `private_ip_address` - The private IP address of the IP configuration.

* `private_ip_allocation_method` - The private IP address allocation method.

* `subnet_id` - Resource ID.

* `private_ip_address_version` - Available from Api-Version 2016-03-30 onwards, it represents whether the specific ipconfiguration is IPv4 or IPv6. Default is taken as IPv4.

* `name` - The name of private link service ip configuration.

---

The `load_balancer_frontend_ip_configuration` block contains the following:

* `id` - Resource ID.

---

The `network_interface` block contains the following:

* `id` - Resource ID.

---

The `private_endpoint_connection` block contains the following:

* `id` - Resource ID.

* `private_endpoint` - One `private_endpoint` block defined below.

* `private_link_service_connection_state` - One `private_link_service_connection_state` block defined below.

* `name` - The name of the resource that is unique within a resource group. This name can be used to access the resource.


---

The `private_endpoint` block contains the following:

* `id` - Resource ID.

* `location` - Resource location.

* `tags` - Resource tags.

---

The `private_link_service_connection_state` block contains the following:

* `status` - Indicates whether the connection has been Approved/Rejected/Removed by the owner of the service.

* `description` - The reason for approval/rejection of the connection.

* `action_required` - A message indicating if changes on the service provider require any updates on the consumer.

---

The `visibility` block contains the following:

* `subscriptions` - The list of subscriptions.
