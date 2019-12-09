---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_public_ip_prefix"
sidebar_current: "docs-azurerm-datasource-public-ip-prefix-x"
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
  value = "${data.azurerm_public_ip_prefix.example.ip_prefix}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the public IP prefix.
* `resource_group_name` - (Required) Specifies the name of the resource group.

## Attributes Reference

* `name` - The name of the Public IP prefix resource.
* `resource_group_name` - The name of the resource group in which to create the public IP.
* `location` - The supported Azure location where the resource exists.
* `sku` - The SKU of the Public IP Prefix.
* `prefix_length` - The number of bits of the prefix.
* `tags` - A mapping of tags to assigned to the resource.
