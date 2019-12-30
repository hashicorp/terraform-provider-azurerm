---
subcategory: "DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_aaaa_record"
sidebar_current: "docs-azurerm-resource-dns-aaaa-record"
description: |-
  Manages a DNS AAAA Record.
---

# azurerm_dns_aaaa_record

Enables you to manage DNS AAAA Records within Azure DNS.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "acceptanceTestResourceGroup1"
  location = "West US"
}

resource "azurerm_dns_zone" "example" {
  name                = "mydomain.com"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_dns_aaaa_record" "example" {
  name                = "test"
  zone_name           = "${azurerm_dns_zone.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  ttl                 = 300

}
```

## Example Usage (Alias Record)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "acceptanceTestResourceGroup1"
  location = "West US"
}

resource "azurerm_dns_zone" "example" {
  name                = "mydomain.com"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_public_ip" "example" {
  name                = "mypublicip"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  allocation_method   = "Dynamic"
  ip_version          = "IPv6"
}

resource "azurerm_dns_aaaa_record" "example" {
  name                = "test"
  zone_name           = "${azurerm_dns_zone.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  ttl                 = 300
  target_resource_id  = "${azurerm_public_ip.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DNS AAAA Record.

* `resource_group_name` - (Required) Specifies the resource group where the DNS Zone (parent resource) exists. Changing this forces a new resource to be created.

* `zone_name` - (Required) Specifies the DNS Zone where the resource exists. Changing this forces a new resource to be created.

* `TTL` - (Required) The Time To Live (TTL) of the DNS record in seconds.

* `records` - (Optional) List of IPv4 Addresses. Conflicts with `target_resource_id`.

* `target_resource_id` - (Optional) The Azure resource id of the target object. Conflicts with `records`

* `tags` - (Optional) A mapping of tags to assign to the resource.

~> **Note:** either `records` OR `target_resource_id` must be specified, but not both.

## Attributes Reference

The following attributes are exported:

* `id` - The DNS AAAA Record ID.
* `fqdn` - The FQDN of the DNS AAAA Record.

## Import

AAAA records can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dns_aaaa_record.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/dnszones/zone1/AAAA/myrecord1
```
