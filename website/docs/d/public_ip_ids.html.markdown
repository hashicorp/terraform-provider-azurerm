---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_public_ips"
sidebar_current: "docs-azurerm-datasource-public-ip-ids"
description: |-
  Provides a list of unassociated public IP address IDs.
---

# azurerm\_public\_ip\_ids

Use this data source to get a list of unassociated public IP address IDs
in a resource group, optionally specifying a minimum required number.

## Example Usage

```hcl
data "azurerm_public_ips" "datasourceips" {
  resource_group_name = "pipRG"
  minimum_count       = 2
}

resource "azurerm_lb" "load_balancer" {
  count               = 2
  name                = "load_balancer-${count.index}"
  location            = "northeurope"
  resource_group_name = "acctestRG"

  frontend_ip_configuration {
    name                 = "frontend"
    public_ip_address_id = "${data.azurerm_public_ips.datasourceips.ids[count.index]}"
  }
}
```

## Argument Reference

* `resource_group_name` - (Required) Specifies the name of the resource group.
* `minimum_count` - (Optional) Specifies the minimum number of IP addresses that
must be available, otherwise an error will be raised.


## Attributes Reference

* `ids` - A list of public IP address resource IDs.