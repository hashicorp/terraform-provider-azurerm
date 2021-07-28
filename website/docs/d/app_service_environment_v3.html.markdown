---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_app_service_environment_v3"
description: |-
  Gets information about an existing 3rd Generation (v3) App Service Environment.
---

# Data Source: azurerm_app_service_environment_v3

Use this data source to access information about an existing 3rd Generation (v3) App Service Environment.

## Example Usage

```hcl
data "azurerm_app_service_environment_v3" "example" {
  name                = "example-ASE"
  resource_group_name = "example-resource-group"
}

output "id" {
  value = data.azurerm_app_service_environment_v3.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this v3 App Service Environment.

* `resource_group_name` - (Required) The name of the Resource Group where the v3 App Service Environment exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the v3 App Service Environment.

* `allow_new_private_endpoint_connections` - Are new Private Endpoint Connections allowed. 

* `cluster_setting` - A `cluster_setting` block as defined below.

* `dedicated_host_count` - The number of Dedicated Hosts used by this ASEv3. 

* `dns_suffix` - the DNS suffix for this App Service Environment V3.

* `inbound_network_dependencies` - An Inbound Network Dependencies block as defined below.

* `internal_load_balancing_mode` - The Internal Load Balancing Mode of this ASEv3.  

* `ip_ssl_address_count` - The number of IP SSL addresses reserved for the App Service Environment V3.

* `linux_outbound_ip_addresses` - The list of Outbound IP Addresses of Linux based Apps in this App Service Environment V3.

* `location` - The location where the App Service Environment exists.

* `pricing_tier` - Pricing tier for the front end instances.

* `subnet_id` - The ID of the v3 App Service Environment Subnet.

* `windows_outbound_ip_addresses` - Outbound addresses of Windows based Apps in this App Service Environment V3.

* `tags` - A mapping of tags assigned to the v3 App Service Environment.

---

A `cluster_setting` block exports the following:

* `name` - The name of the Cluster Setting.

* `value` - The value for the Cluster Setting.

--- 

An `inbound_network_dependencies` block exports the following:

* `description` - A short description of the purpose of the network traffic.

* `ip_addresses` - A list of IP addresses that network traffic will originate from in CIDR notation.

* `ports` - The ports that network traffic will arrive to the App Service Environment V3 on.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the 3rd Generation (v3) App Service Environment.
