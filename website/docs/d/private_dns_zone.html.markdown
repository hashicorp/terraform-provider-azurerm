---
subcategory: "Private DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_zone"
description: |-
  Gets information about an existing Private DNS Zone.

---

# Data Source: azurerm_private_dns_zone

Use this data source to access information about an existing Private DNS Zone.

## Example Usage

```hcl
data "azurerm_private_dns_zone" "example" {
  name                = "contoso.internal"
  resource_group_name = "contoso-dns"
}

output "private_dns_zone_id" {
  value = data.azurerm_private_dns_zone.example.id
}
```

## Argument Reference

* `name` - The name of the Private DNS Zone.
* `resource_group_name` - (Optional) The Name of the Resource Group where the Private DNS Zone exists.
If the Name of the Resource Group is not provided, the first Private DNS Zone from the list of Private
DNS Zones in your subscription that matches `name` will be returned.

## Attributes Reference

* `id` - The ID of the Private DNS Zone.

* `max_number_of_record_sets` - Maximum number of recordsets that can be created in this Private Zone.
* `max_number_of_virtual_network_links` - Maximum number of Virtual Networks that can be linked to this Private Zone.
* `max_number_of_virtual_network_links_with_registration` - Maximum number of Virtual Networks that can be linked to this Private Zone with registration enabled.
* `number_of_record_sets` - The number of recordsets currently in the zone.
* `tags` - A mapping of tags for the zone.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS Zone.
