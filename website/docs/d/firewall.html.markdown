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

* `ip_configuration` - A `ip_configuration` block as defined below.

---

A `ip_configuration` block exports the following:

* `subnet_id` - The ID of the Subnet where the Azure Firewall is deployed.

* `private_ip_address` - The Private IP Address of the Azure Firewall.

* `public_ip_address_id`- The ID of the Public IP address of the Azure Firewall.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Firewall.
