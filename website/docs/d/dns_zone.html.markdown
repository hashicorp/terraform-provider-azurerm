---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_zone"
sidebar_current: "docs-azurerm-datasource-dns-zone"
description: |-
  Gets information about an existing DNS Zone.

---

# Data Source: azurerm_dns_zone

Use this data source to access information about an existing DNS Zone.

## Example Usage

```hcl
data "azurerm_dns_zone" "test" {
  name                = "search-eventhubns"
  resource_group_name = "search-service"
}

output "dns_zone_id" {
  value = "${data.azurerm_dns_zone.test.id}"
}
```

## Argument Reference

* `name` - (Required) The name of the DNS Zone.
* `resource_group_name` - (Optional) The Name of the Resource Group where the DNS Zone exists.
If the Name of the Resource Group is not provided, the first DNS Zone from the list of DNS Zones
in your subscription that matches `name` will be returned.

## Attributes Reference

* `id` - The ID of the DNS Zone.

* `max_number_of_record_sets` - Maximum number of Records in the zone.
* `number_of_record_sets` - The number of records already in the zone.
* `name_servers` - A list of values that make up the NS record for the zone.
* `registration_virtual_network_ids` - A list of Virtual Network ID's that register hostnames in this DNS zone.
* `resolution_virtual_network_ids` - A list of Virtual Network ID's that resolve records in this DNS zone.
* `tags` - A mapping of tags to assign to the EventHub Namespace.
* `zone_type` - (**Deprecated**) The type of this DNS zone, such as `Public` or `Private`.

~> **NOTE:** This field is deprecated since Private DNS is now a separate resource and will be removed in 2.0 of the Azure Provider.
