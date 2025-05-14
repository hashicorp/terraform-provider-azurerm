---
subcategory: "Private DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_cname_record"
description: |-
  Manages a Private DNS CNAME Record.
---

# azurerm_private_dns_cname_record

Enables you to manage DNS CNAME Records within Azure Private DNS.

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

resource "azurerm_private_dns_cname_record" "example" {
  name                = "test"
  zone_name           = azurerm_private_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 300
  record              = "contoso.com"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DNS CNAME Record. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the resource group where the resource exists. Changing this forces a new resource to be created.

* `zone_name` - (Required) Specifies the Private DNS Zone where the resource exists. Changing this forces a new resource to be created.

* `ttl` - (Required) The Time To Live (TTL) of the DNS record in seconds. Possible values are between `0` and `2147483647`.

* `record` - (Required) The target of the CNAME.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Private DNS CNAME Record ID.

* `fqdn` - The FQDN of the DNS CNAME Record.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Private DNS CNAME Record.
* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS CNAME Record.
* `update` - (Defaults to 30 minutes) Used when updating the Private DNS CNAME Record.
* `delete` - (Defaults to 30 minutes) Used when deleting the Private DNS CNAME Record.

## Import

Private DNS CName Records can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_private_dns_cname_record.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/privateDnsZones/zone1/CNAME/myrecord1
```
