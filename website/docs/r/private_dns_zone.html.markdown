---
subcategory: "Private DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_zone"
description: |-
  Manages a Private DNS Zone.
---

# azurerm_private_dns_zone

Enables you to manage Private DNS zones within Azure DNS. These zones are hosted on Azure's name servers.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_private_dns_zone" "example" {
  name                = "mydomain.com"
  resource_group_name = azurerm_resource_group.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Private DNS Zone. Must be a valid domain name. Changing this forces a new resource to be created.

-> **Note:** If you are going to be using the Private DNS Zone with a Private Endpoint the name of the Private DNS Zone must follow the **Private DNS Zone name** schema in the [product documentation](https://docs.microsoft.com/azure/private-link/private-endpoint-dns#virtual-network-and-on-premises-workloads-using-a-dns-forwarder) in order for the two resources to be connected successfully.

* `resource_group_name` - (Required) Specifies the resource group where the resource exists. Changing this forces a new resource to be created.

* `soa_record` - (Optional) An `soa_record` block as defined below. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `soa_record` block supports:

* `email` - (Required) The email contact for the SOA record.

* `expire_time` - (Optional) The expire time for the SOA record. Defaults to `2419200`.

* `minimum_ttl` - (Optional) The minimum Time To Live for the SOA record. By convention, it is used to determine the negative caching duration. Defaults to `10`.

* `refresh_time` - (Optional) The refresh time for the SOA record. Defaults to `3600`.

* `retry_time` - (Optional) The retry time for the SOA record. Defaults to `300`.

* `ttl` - (Optional) The Time To Live of the SOA Record in seconds. Defaults to `3600`.

* `tags` - (Optional) A mapping of tags to assign to the Record Set.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Private DNS Zone ID.
* `soa_record` - A `soa_record` block as defined below.
* `number_of_record_sets` - The current number of record sets in this Private DNS zone.
* `max_number_of_record_sets` - The maximum number of record sets that can be created in this Private DNS zone.
* `max_number_of_virtual_network_links` - The maximum number of virtual networks that can be linked to this Private DNS zone.
* `max_number_of_virtual_network_links_with_registration` - The maximum number of virtual networks that can be linked to this Private DNS zone with registration enabled.

---

A `soa_record` block exports the following:

* `fqdn` - The fully qualified domain name of the Record Set.

* `host_name` - The domain name of the authoritative name server for the SOA record.

* `serial_number` - The serial number for the SOA record.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Private DNS Zone.
* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS Zone.
* `update` - (Defaults to 30 minutes) Used when updating the Private DNS Zone.
* `delete` - (Defaults to 30 minutes) Used when deleting the Private DNS Zone.

## Import

Private DNS Zones can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_private_dns_zone.zone1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/privateDnsZones/zone1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network`: 2024-06-01
