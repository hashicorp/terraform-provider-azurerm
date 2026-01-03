---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_network_security_perimeter"
description: |-
  Gets information about an existing Network Security Perimeter.
---

# Data Source: azurerm_network_security_perimeter

Use this data source to access information about an existing Network Security Perimeter.

## Example Usage

```hcl
data "azurerm_network_security_perimeter" "example" {
  name = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_network_security_perimeter.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Network Security Perimeter.

* `resource_group_name` - (Required) The name of the Resource Group where the Network Security Perimeter exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Network Security Perimeter.

* `location` - The Azure Region where the Network Security Perimeter exists.

* `tags` - A mapping of tags assigned to the Network Security Perimeter.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Network Security Perimeter.