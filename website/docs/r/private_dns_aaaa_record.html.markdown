---
subcategory: "Private DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_aaaa_record"
description: |-
  Manages a Private DNS AAAA Record.
---

# azurerm_private_dns_aaaa_record

Enables you to manage DNS AAAA Records within Azure Private DNS.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "mydomain.com"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_private_dns_aaaa_record" "test" {
  name                = "test"
  zone_name           = azurerm_private_dns_zone.test.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 300
  records             = ["fd5d:70bc:930e:d008:0000:0000:0000:7334", "fd5d:70bc:930e:d008::7335"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DNS A Record. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the resource group where the resource exists. Changing this forces a new resource to be created.

* `zone_name` - (Required) Specifies the Private DNS Zone where the resource exists. Changing this forces a new resource to be created.

* `ttl` - (Required) The Time To Live (TTL) of the DNS record in seconds.

* `records` - (Required) A list of IPv6 Addresses.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Private DNS AAAA Record ID.

* `fqdn` - The FQDN of the DNS AAAA Record.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Private DNS AAAA Record.
* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS AAAA Record.
* `update` - (Defaults to 30 minutes) Used when updating the Private DNS AAAA Record.
* `delete` - (Defaults to 30 minutes) Used when deleting the Private DNS AAAA Record.

## Import

Private DNS AAAA Records can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_private_dns_aaaa_record.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/privateDnsZones/zone1/AAAA/myrecord1
```
