---
subcategory: "DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_zone"
description: |-
  Gets information about an existing DNS Zone.

---

# Data Source: azurerm_dns_zone

Use this data source to access information about an existing DNS Zone.

## Example Usage

```hcl
data "azurerm_dns_zone" "example" {
  name                = "search-eventhubns"
  resource_group_name = "search-service"
}

output "dns_zone_id" {
  value = data.azurerm_dns_zone.example.id
}
```

## Argument Reference

* `name` - The name of the DNS Zone.

* `resource_group_name` - (Optional) The Name of the Resource Group where the DNS Zone exists.
If the Name of the Resource Group is not provided, the first DNS Zone from the list of DNS Zones
in your subscription that matches `name` will be returned.

## Attributes Reference

* `id` - The ID of the DNS Zone.

* `max_number_of_record_sets` - Maximum number of Records in the zone.

* `number_of_record_sets` - The number of records already in the zone.

* `name_servers` - A list of values that make up the NS record for the zone.

* `tags` - A mapping of tags assigned to the DNS Zone.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the DNS Zone.
