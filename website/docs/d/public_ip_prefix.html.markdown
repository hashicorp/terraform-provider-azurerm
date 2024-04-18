---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_public_ip_prefix"
description: |-
  Gets information about an existing Public IP Prefix.

---

# Data Source: azurerm_public_ip_prefix

Use this data source to access information about an existing Public IP Prefix.

## Example Usage (reference an existing)

```hcl
data "azurerm_public_ip_prefix" "example" {
  name                = "name_of_public_ip"
  resource_group_name = "name_of_resource_group"
}

output "public_ip_prefix" {
  value = data.azurerm_public_ip_prefix.example.ip_prefix
}
```

## Argument Reference

* `name` - Specifies the name of the public IP prefix.
* `resource_group_name` - Specifies the name of the resource group.

## Attributes Reference

* `id` - The ID of the Public IP Prefix.
* `ip_prefix` - The Public IP address range, in CIDR notation.
* `location` - The supported Azure location where the resource exists.
* `sku` - The SKU of the Public IP Prefix.
* `prefix_length` - The number of bits of the prefix.
* `tags` - A mapping of tags to assigned to the resource.
* `zones` - A list of Availability Zones in which this Public IP Prefix is located.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Public IP Prefix.
