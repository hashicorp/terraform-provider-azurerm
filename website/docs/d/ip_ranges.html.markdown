---
layout: "azurerm"
page_title: "Azure IP Ranges: azurerm_ip_ranges"
sidebar_current: "docs-azurerm-datasource-ip-ranges"
description: |-
  Get information on Azure Datacenter IP Ranges
---

# Data Source: azurerm_ip_ranges

Use this data source to get the public IP addresses of Azure Datacenters.

## Example Usage

```hcl
data "azurerm_ip_ranges" "test" {
  regions = ["australiaeast", "australiasouth"]
}

resource "azurerm_network_security_rule" "test" {
  count                       = "${length(data.azurerm_ip_ranges.test.subnets)}"
  name                        = "aksWorkerSubnetSecurityGroup-${count.index}"
  priority                    = "${(count.index + 100)}"
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "${data.azurerm_ip_ranges.test.subnets[count.index]}"
  resource_group_name         = "${azurerm_resource_group.test.name}"
  network_security_group_name = "${azurerm_network_security_group.test.name}"
}
```

## Argument Reference

* `regions` - (Optional) Filter IP ranges by Azure region (or include all regions if omitted). Valid items are all Azure regions (eg. `australiaeast`).

## Attributes Reference

* `subnets` - The lexically ordered list of CIDR blocks
