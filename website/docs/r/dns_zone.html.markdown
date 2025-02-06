---
subcategory: "DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_zone"
description: |-
  Manages a DNS Zone.
---

# azurerm_dns_zone

Enables you to manage DNS zones within Azure DNS. These zones are hosted on Azure's name servers to which you can delegate the zone from the parent domain.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dns_zone" "example-public" {
  name                = "mydomain.com"
  resource_group_name = azurerm_resource_group.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DNS Zone. Must be a valid domain name. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the resource group where the resource exists. Changing this forces a new resource to be created.

* `soa_record` - (Optional) An `soa_record` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `soa_record` block supports:

* `email` - (Required) The email contact for the SOA record.

* `host_name` - (Optional) The domain name of the authoritative name server for the SOA record. If not set, computed value from Azure will be used.

* `expire_time` - (Optional) The expire time for the SOA record. Defaults to `2419200`.

* `minimum_ttl` - (Optional) The minimum Time To Live for the SOA record. By convention, it is used to determine the negative caching duration. Defaults to `300`.

* `refresh_time` - (Optional) The refresh time for the SOA record. Defaults to `3600`.

* `retry_time` - (Optional) The retry time for the SOA record. Defaults to `300`.

* `serial_number` - (Optional) The serial number for the SOA record. Defaults to `1`.

* `ttl` - (Optional) The Time To Live of the SOA Record in seconds. Defaults to `3600`.

* `tags` - (Optional) A mapping of tags to assign to the Record Set.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The DNS Zone ID.

* `max_number_of_record_sets` - (Optional) Maximum number of Records in the zone. Defaults to `1000`.

* `number_of_record_sets` - (Optional) The number of records already in the zone.

* `name_servers` - (Optional) A list of values that make up the NS record for the zone.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DNS Zone.

* `update` - (Defaults to 30 minutes) Used when updating the DNS Zone.

* `read` - (Defaults to 5 minutes) Used when retrieving the DNS Zone.

* `delete` - (Defaults to 30 minutes) Used when deleting the DNS Zone.

## Import

DNS Zones can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dns_zone.zone1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/dnsZones/zone1
```
