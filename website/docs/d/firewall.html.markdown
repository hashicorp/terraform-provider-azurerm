---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_firewall"
sidebar_current: "docs-azurerm-datasource-firewall"
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
  value = "${data.azurerm_firewall.example.ip_configuration.0.private_ip_address}"
}
```

## Argument Reference

* `name` - (Required) The name of the Azure Firewall.

* `resource_group_name` - (Required) The name of the Resource Group in which the Azure Firewall exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Resource ID of the Azure Firewall.

* `ip_configuration` - A `ip_configuration` block as defined below.

---

A `ip_configuration` block exports the following:

* `subnet_id` - The Resource ID of the subnet where the Azure Firewall is deployed.

* `private_ip_address` - The private IP address of the Azure Firewall.

* `public_ip_address_id`- The Resource ID of the public IP address of the Azure Firewall.
