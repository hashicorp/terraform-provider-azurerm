---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_public_ips"
sidebar_current: "docs-azurerm-datasource-public-ips"
description: |-
  Provides a list of public IP addresses.
---

# azurerm\_public\_ips

Use this data source to get a list of associated or unassociated public IP addresses
in a resource group, optionally specifying a minimum required number.

## Example Usage

```hcl
data "azurerm_public_ips" "test" {
  resource_group_name = "pip-test"
  attached            = false
  minimum_count       = 2
}

resource "azurerm_lb" "load_balancer" {
  count               = 2
  name                = "load_balancer-${count.index}"
  location            = "northeurope"
  resource_group_name = "acctestRG"

  frontend_ip_configuration {
    name                 = "frontend"
    public_ip_address_id = "${lookup(data.azurerm_public_ips.test.public_ips[count.index], "public_ip_address_id")}"
  }
}
```

## Argument Reference

* `resource_group_name` - (Required) Specifies the name of the resource group.
* `attached` - (Required) Whether to return public IPs that are attached or not.
* `minimum_count` - (Optional) Specifies the minimum number of IP addresses that
must be available, otherwise an error will be raised.


## Attributes Reference

* `public_ips` - A list of public IP addresses. Each public IP is represented by a
map containing the following keys; public_ip_address_id, name, fqdn, domain_name_label,
ip_address. Note that if the public IP is unassigned then some values may be empty.