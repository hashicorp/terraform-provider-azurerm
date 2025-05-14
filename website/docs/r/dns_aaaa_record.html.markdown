---
subcategory: "DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_aaaa_record"
description: |-
  Manages a DNS AAAA Record.
---

# azurerm_dns_aaaa_record

Enables you to manage DNS AAAA Records within Azure DNS.

~> **Note:** [The Azure DNS API has a throttle limit of 500 read (GET) operations per 5 minutes](https://docs.microsoft.com/azure/azure-resource-manager/management/request-limits-and-throttling#network-throttling) - whilst the default read timeouts will work for most cases - in larger configurations you may need to set a larger [read timeout](https://www.terraform.io/language/resources/syntax#operation-timeouts) then the default 5min. Although, we'd generally recommend that you split the resources out into smaller Terraform configurations to avoid the problem entirely.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dns_zone" "example" {
  name                = "mydomain.com"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_dns_aaaa_record" "example" {
  name                = "test"
  zone_name           = azurerm_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 300
  records             = ["2001:db8::1:0:0:1"]
}
```

## Example Usage (Alias Record)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dns_zone" "example" {
  name                = "mydomain.com"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_public_ip" "example" {
  name                = "mypublicip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Dynamic"
  ip_version          = "IPv6"
}

resource "azurerm_dns_aaaa_record" "example" {
  name                = "test"
  zone_name           = azurerm_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 300
  target_resource_id  = azurerm_public_ip.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DNS AAAA Record. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the resource group where the DNS Zone (parent resource) exists. Changing this forces a new resource to be created.

* `zone_name` - (Required) Specifies the DNS Zone where the resource exists. Changing this forces a new resource to be created.

* `ttl` - (Required) The Time To Live (TTL) of the DNS record in seconds.

* `records` - (Optional) List of IPv6 Addresses. Conflicts with `target_resource_id`.

* `target_resource_id` - (Optional) The Azure resource id of the target object. Conflicts with `records`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

~> **Note:** either `records` OR `target_resource_id` must be specified, but not both.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The DNS AAAA Record ID.

* `fqdn` - The FQDN of the DNS AAAA Record.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DNS AAAA Record.

* `read` - (Defaults to 5 minutes) Used when retrieving the DNS AAAA Record.

* `update` - (Defaults to 30 minutes) Used when updating the DNS AAAA Record.

* `delete` - (Defaults to 30 minutes) Used when deleting the DNS AAAA Record.

## Import

AAAA records can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dns_aaaa_record.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/dnsZones/zone1/AAAA/myrecord1
```
