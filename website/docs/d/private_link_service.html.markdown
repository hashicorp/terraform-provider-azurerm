---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_link_service"
description: |-
  Use this data source to access information about an existing Private Link Service.
---

# Data Source: azurerm_private_link_service

Use this data source to access information about an existing Private Link Service.

## Example Usage

```hcl
data "azurerm_private_link_service" "example" {
  name                = "myPrivateLinkService"
  resource_group_name = "PrivateLinkServiceRG"
}

output "private_link_service_id" {
  value = data.azurerm_private_link_service.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the private link service.

* `resource_group_name` - The name of the resource group in which the private link service resides.

## Attributes Reference

The following attributes are exported:

* `id` - The Azure resource ID of the Private Link Service.

* `alias` - The alias is a globally unique name for your private link service which Azure generates for you. Your can use this alias to request a connection to your private link service.

* `auto_approval_subscription_ids` - The list of subscription(s) globally unique identifiers that will be auto approved to use the private link service.

* `enable_proxy_protocol` - Does the Private Link Service support the Proxy Protocol?

* `load_balancer_frontend_ip_configuration_ids` - The list of Standard Load Balancer(SLB) resource IDs. The Private Link service is tied to the frontend IP address of a SLB. All traffic destined for the private link service will reach the frontend of the SLB. You can configure SLB rules to direct this traffic to appropriate backend pools where your applications are running.

* `location` - The supported Azure location where the resource exists.

* `nat_ip_configuration` - The `nat_ip_configuration` block as defined below.

* `tags` - A mapping of tags to assign to the resource.

* `visibility_subscription_ids` - The list of subscription(s) globally unique identifiers(GUID) that will be able to see the private link service.

---

The `nat_ip_configuration` block exports the following:

* `name` - The name of private link service NAT IP configuration.

* `private_ip_address` - The private IP address of the NAT IP configuration.

* `private_ip_address_version` - The version of the IP Protocol.

* `subnet_id` - The ID of the subnet to be used by the service.

* `primary` - Value that indicates if the IP configuration is the primary configuration or not.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Private Link Service.
