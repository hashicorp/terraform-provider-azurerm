---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_group"
description: |-
  Gets information about an existing Container Group instance.
---

# Data Source: azurerm_container_group

Use this data source to access information about an existing Container Group instance.

## Example Usage

```hcl
data "azurerm_container_group" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_container_group.example.id
}

output "ip_address" {
  value = data.azurerm_container_group.example.ip_address
}

output "fqdn" {
  value = data.azurerm_container_group.example.fqdn
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Container Group instance.

* `resource_group_name` - (Required) The name of the Resource Group where the Container Group instance exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Container Group instance.

* `ip_address` - The IP address allocated to the Container Group instance.

* `fqdn` - The FQDN of the Container Group instance derived from `dns_name_label`.

* `location` - The Azure Region where the Container Group instance exists.

* `tags` - A mapping of tags assigned to the Container Group instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Container Group instance.
