---
subcategory: "DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_ns_record"
sidebar_current: "docs-azurerm-resource-dns-ns-record"
description: |-
  Manages a DNS NS Record.
---

# azurerm_dns_ns_record

Enables you to manage DNS NS Records within Azure DNS.

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

resource "azurerm_dns_ns_record" "example" {
  name                = "test"
  zone_name           = "${azurerm_dns_zone.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  ttl                 = 300

  records = ["ns1.contoso.com", "ns2.contoso.com"]

  tags = {
    Environment = "Production"
  }
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DNS NS Record.

* `resource_group_name` - (Required) Specifies the resource group where the resource exists. Changing this forces a new resource to be created.

* `zone_name` - (Required) Specifies the DNS Zone where the DNS Zone (parent resource) exists. Changing this forces a new resource to be created.

* `ttl` - (Required) The Time To Live (TTL) of the DNS record in seconds.

* `records` - (Optional) A list of values that make up the NS record. *WARNING*: Either `records` or `record` is required.

* `record` - (Optional) A list of values that make up the NS record. Each `record` block supports fields documented below. This field has been deprecated and will be removed in a future release.

* `tags` - (Optional) A mapping of tags to assign to the resource.

The `record` block supports:

* `nsdname` - (Required) The value of the record.

## Attributes Reference

The following attributes are exported:

* `id` - The DNS NS Record ID.
* `fqdn` - The FQDN of the DNS NS Record.

## Import

NS records can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dns_ns_record.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/dnszones/zone1/NS/myrecord1
```
