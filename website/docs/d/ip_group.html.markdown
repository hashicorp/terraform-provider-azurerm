---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_ip_group"
description: |-
  Gets information about an existing IP Group.
---

# Data Source: azurerm_ip_group

Use this data source to access information about an existing IP Group.

## Example Usage

```hcl
data "azurerm_ip_group" "example" {
  name                = "example1-ipgroup"
  resource_group_name = "example-rg"
}

output "cidrs" {
  value = data.azurerm_ip_group.example.cidrs
}
```

## Argument Reference

* `name` - Specifies the Name of the IP Group.

* `resource_group_name` - Specifies the Name of the Resource Group within which the IP Group exists


## Attributes Reference

* `id` - The ID of the IP Group.

* `location` - The supported Azure location where the resource exists.

* `cidrs` - A list of CIDRs or IP addresses.

* `tags` - A mapping of tags assigned to the resource.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the IP Group.
