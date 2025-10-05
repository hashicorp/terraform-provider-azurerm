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

## Argument Reference

The following arguments are supported:

* `name` - (Required) A substring to match some number of IP Groups.

* `resource_group_name` - (Required) The name of the Resource Group where the IP Groups exist.

## Attribute Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `ids` - A list of IP Group IDs.

* `names` - A list of IP Group Names.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 10 minutes) Used when retrieving the IP Groups.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Network` - 2024-05-01
