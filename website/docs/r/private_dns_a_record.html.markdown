---
subcategory: "Private DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_a_record"
description: |-
  Manages a Private DNS A Record.
---

# azurerm_private_dns_a_record

Enables you to manage DNS A Records within Azure Private DNS.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "acceptanceTestResourceGroup1"
  location = "West US"
}

resource "azurerm_private_dns_zone" "example" {
  name                = "mydomain.com"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_private_dns_a_record" "example" {
  name                = "test"
  zone_name           = azurerm_private_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 300
  records             = ["10.0.180.17"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DNS A Record.

* `resource_group_name` - (Required) Specifies the resource group where the Private DNS Zone exists. Changing this forces a new resource to be created.

* `zone_name` - (Required) Specifies the Private DNS Zone where the resource exists. Changing this forces a new resource to be created.

* `TTL` - (Required) The Time To Live (TTL) of the DNS record in seconds.

* `records` - (Required) List of IPv4 Addresses.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Private DNS A Record ID.

* `fqdn` - The FQDN of the DNS A Record.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Private DNS A Record.
* `update` - (Defaults to 30 minutes) Used when updating the Private DNS A Record.
* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS A Record.
* `delete` - (Defaults to 30 minutes) Used when deleting the Private DNS A Record.

## Import

Private DNS A Records can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_private_dns_a_record.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/privateDnsZones/zone1/A/myrecord1
```
