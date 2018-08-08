---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_zone"
sidebar_current: "docs-azurerm-resource-dns-zone"
description: |-
  Manages a DNS Zone.
---

# azurerm_dns_zone

Enables you to manage DNS zones within Azure DNS. These zones are hosted on Azure's name servers to which you can delegate the zone from the parent domain.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_dns_zone" "example" {
  name                = "mydomain.com"
  resource_group_name = "${azurerm_resource_group.example.name}"
  zone_type           = "Public"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DNS Zone. Must be a valid domain name.

* `resource_group_name` - (Required) Specifies the resource group where the resource exists. Changing this forces a new resource to be created.

* `zone_type` - (Required) Specifies the type of this DNS zone. Possible values are `Public` or `Private` (Defaults to `Public`).

* `registration_virtual_network_ids` - (Optional) A list of Virtual Network ID's that register hostnames in this DNS zone. This field can only be set when `zone_type` is set to `Private`.

* `resolution_virtual_network_ids` - (Optional) A list of Virtual Network ID's that resolve records in this DNS zone. This field can only be set when `zone_type` is set to `Private`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The DNS Zone ID.
* `max_number_of_record_sets` - (Optional) Maximum number of Records in the zone. Defaults to `1000`.
* `number_of_record_sets` - (Optional) The number of records already in the zone.
* `name_servers` - (Optional) A list of values that make up the NS record for the zone.


## Import

DNS Zones can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dns_zone.zone1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/dnsZones/zone1
```
