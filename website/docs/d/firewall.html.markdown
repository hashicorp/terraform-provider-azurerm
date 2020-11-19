---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_firewall"
description: |-
  Gets information about an existing Azure Firewall.

---

# Data Source: azurerm_firewall

Use this data source to access information about an existing Azure Firewall.

## Example Usage

```hcl
data "azurerm_firewall" "example" {
  name                = "firewall1"
  resource_group_name = "firewall-RG"
}

output "firewall_private_ip" {
  value = data.azurerm_firewall.example.ip_configuration[0].private_ip_address
}
```

## Argument Reference

* `name` - The name of the Azure Firewall.

* `resource_group_name` - The name of the Resource Group in which the Azure Firewall exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Azure Firewall.

* `location` - The Azure location where the Azure Firewall exists.

* `sku_name` - The sku name of the Azure Firewall.

* `sku_tier` - The sku tier of the Azure Firewall.

* `firewall_policy_id` - The ID of the Firewall Policy applied to the Azure Firewall.

* `ip_configuration` - A `ip_configuration` block as defined below.

* `dns_servers` - The list of DNS servers that the Azure Firewall will direct DNS traffic to the for name resolution.

* `management_ip_configuration` - A `management_ip_configuration` block as defined below, which allows force-tunnelling of traffic to be performed by the firewall.

* `threat_intel_mode` - The operation mode for threat intelligence-based filtering.

* `virtual_hub` - A `virtual_hub` block as defined below.

* `zones` - The availability zones in which the Azure Firewall is created.

* `tags` - A mapping of tags assigned to the Azure Firewall.

---

A `ip_configuration` block exports the following:

* `subnet_id` - The ID of the Subnet where the Azure Firewall is deployed.

* `private_ip_address` - The Private IP Address of the Azure Firewall.

* `public_ip_address_id`- The ID of the Public IP address of the Azure Firewall.

---

A `management_ip_configuration` block exports the following:

* `subnet_id` - The ID of the Subnet where the Azure Firewall is deployed.

* `private_ip_address` - The Private IP Address of the Azure Firewall.

* `public_ip_address_id`- The ID of the Public IP address of the Azure Firewall.

---

A `virtual_hub` block exports the following:

* `virtual_hub_id` - The ID of the Virtual Hub where the Azure Firewall resides in.

* `public_ip_count` - The number of public IPs assigned to the Azure Firewall.

* `public_ip_addresses` - The list of public IP addresses associated with the Azure Firewall.

* `private_ip_address` - The private IP address associated with the Azure Firewall.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Firewall.
