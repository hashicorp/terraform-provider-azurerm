---
subcategory: ""
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
  name                = "myPrivateLinkService"
  resource_group_name = "PrivateLinkServiceRG"
}

output "private_link_service_id" {
  value = "${data.azurerm_private_link_service.example.id}"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the private link service.

* `resource_group_name` - (Required) The name of the resource group in which the private link service resides.


## Attributes Reference

The following attributes are exported:

* `id` - The Azure resource ID of the Prviate Link Service.

* `location` - The supported Azure location where the resource exists.

* `alias` - The alias is a globally unique name for your private link service which Azure generates for you. Your can use this alias to request a connection to your private link service.

* `auto_approval_subscription_ids` - The list of subscription(s) globally unique identifiers that will be auto approved to use the private link service.

* `visibility_subscription_ids` - The list of subscription(s) globally unique identifiers(GUID) that will be able to see the private link service.

* `nat_ip_configuration` - The `nat_ip_configuration` block as defined below.

* `load_balancer_frontend_ip_configuration_ids` - The list of Standard Load Balancer(SLB) resource IDs. The Private Link service is tied to the frontend IP address of a SLB. All traffic destined for the private link service will reach the frontend of the SLB. You can configure SLB rules to direct this traffic to appropriate backend pools where your applications are running.

* `network_interfaces` - The list of network interface resource ids that are being used by the service.

* `tags` - A mapping of tags to assign to the resource.


---

The `nat_ip_configuration` block exports the following:

* `name` - The name of private link service NAT IP configuration.

* `private_ip_address` - The private IP address of the NAT IP configuration.

* `private_ip_address_version` - The ip address version of the `ip_configuration`.

* `subnet_id` - The resource ID of the subnet to be used by the service.

* `primary` - Value that indicates if the IP configuration is the primary configuration or not.

