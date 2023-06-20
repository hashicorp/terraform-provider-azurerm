---
subcategory: "Private DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_zone_virtual_network_link"
description: |-
  Gets information about an existing Private DNS Zone Virtual Network Link.
---

# Data Source: azurerm_private_dns_zone_virtual_network_link

Use this data source to access information about an existing Private DNS zone Virtual Network Link. These Links enable DNS resolution and registration inside Azure Virtual Networks using Azure Private DNS.

## Example Usage

```hcl
data "azurerm_private_dns_zone_virtual_network_link" "example" {
  name                  = "test"
  resource_group_name   = "test-rg"
  private_dns_zone_name = "test-zone"
}

output "private_dns_a_record_id" {
  value = data.azurerm_private_dns_zone_virtual_network_link.example.id
}
```

## Argument Reference

* `name` - The name of the Private DNS Zone Virtual Network Link.

* `private_dns_zone_name` - The name of the Private DNS zone (without a terminating dot).

* `resource_group_name` - Specifies the resource group where the Private DNS Zone exists.

## Attributes Reference

* `id` - The ID of the Private DNS Zone Virtual Network Link.

* `virtual_network_id` - The ID of the Virtual Network that is linked to the DNS Zone.

* `registration_enabled` - Whether the auto-registration of virtual machine records in the virtual network in the Private DNS zone is enabled or not.

* `tags` - A mapping of tags to assign to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS Zone Virtual Network Link.
