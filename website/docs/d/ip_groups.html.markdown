---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_ip_groups"
description: |-
  Gets information about existing IP Groups.
---

# Data Source: azurerm_ip_groups

Use this data source to access information about existing IP Groups.

## Example Usage

```hcl
data "azurerm_ip_groups" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "ids" {
  value = data.azurerm_ip_groups.example.ids
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) A substring to match some number of IP Groups.

* `resource_group_name` - (Required) The name of the Resource Group where the IP Groups exist.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `ids` - A list of IP Group IDs.

* `names` - A list of IP Group Names.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 10 minutes) Used when retrieving the IP Groups.
