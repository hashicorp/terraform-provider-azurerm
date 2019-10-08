---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_link_service"
sidebar_current: "docs-azurerm-datasource-private-link-service"
description: |-
  Use this data source to access information about an existing Private Link Service.
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

* `resource_group_name` - (Required) The name of the resource group in which the private link service resides. Changing this forces a new resource to be created.


## Attributes Reference

The following attributes are exported:

* `location` - Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `alias` - The alias of the private link service.

* `auto_approval_subscription_ids` - A list of subscription globally unique identifiers that will be auto approved to use this private link service.

* `nat_ip_configurations` - One or more `ip_configuration` blocks as defined below.

* `load_balancer_frontend_ip_configuration_ids` - A list of `Standard` Load Balancer resource IDs to direct the service network traffic toward.

* `network_interfaces` - A list of network interface resource ids that are being used by the service.

* `private_endpoint_connection` - One or more `private_endpoint_connection` blocks as defined below.

* `visibility_subscription_ids` - A list of subscription globally unique identifiers(GUID) that will be able to see this service. If left undefined all Azure subscriptions will be able to see this service.

* `tags` - A mapping of tags to assign to the resource. Changing this forces a new resource to be created


---

The `nat_ip_configuration` block contains the following:

* `name` - The name of private link service ip configuration.

* `private_ip_address` - The private IP address of the IP configuration.

* `private_ip_allocation_method` - The private IP address allocation method.

* `subnet_id` - The resource ID of the subnet to be used by the service.

* `private_ip_address_version` - The ip address version of the `ip_configuration`.

---

The `private_endpoint_connection` block contains the following:

* `name` - (Required) The name of the resource that is unique within a resource group. This name can be used to access the resource.

* `private_endpoint_location` - (Required) The resource location of the `private_endpoint`.

* `private_endpoint_id` - (Required) The Private Endpoint ID.

* `state_action_required` - (Computed) A message indicating if changes on the service provider require any updates by the consumer.

* `state_description` - (Computed) The reason for approval/rejection of the connection.

* `state_status` - (Computed) Indicates whether the connection has been Approved/Rejected/Removed by the owner of the service.

