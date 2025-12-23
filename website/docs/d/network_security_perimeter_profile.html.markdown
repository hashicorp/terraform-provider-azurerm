---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_network_security_perimeter_profile"
description: |-
  Gets information about an existing Network Security Perimeter Profile.
---

# Data Source: azurerm_network_security_perimeter_profile

Use this data source to access information about an existing Network Security Perimeter Profile.

## Example Usage

```hcl

data "azurerm_network_security_perimeter" "example" {
  name = "existing"
  resource_group_name = "existing"
}

data "azurerm_network_security_perimeter_profile" "example" {
  name = "existing"
  perimeter_id = data.azurerm_network_security_perimeter.example.id
}

output "id" {
  value = data.azurerm_network_security_perimeter_profile.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Network Security Perimeter Profile.

* `perimeter_id` - (Required) The ID of the Network Security Perimeter.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Network Security Perimeter Profile.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Network Security Perimeter Profile.