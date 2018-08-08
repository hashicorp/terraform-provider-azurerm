---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_cname_record"
sidebar_current: "docs-azurerm-resource-dns-cname-record"
description: |-
  Manages a DNS CNAME Record.
---

# azurerm_dns_cname_record

Enables you to manage DNS CNAME Records within Azure DNS.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_dns_zone" "example" {
  # ...
}

resource "azurerm_dns_cname_record" "example" {
  name                = "example"
  zone_name           = "${azurerm_dns_zone.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  ttl                 = 300
  record              = "contoso.com"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DNS CNAME Record.

* `resource_group_name` - (Required) Specifies the resource group where the resource exists. Changing this forces a new resource to be created.

* `zone_name` - (Required) Specifies the DNS Zone where the resource exists. Changing this forces a new resource to be created.

* `TTL` - (Required) The Time To Live (TTL) of the DNS record.

* `record` - (Required) The target of the CNAME.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The DNS CName Record ID.

## Import

CNAME records can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dns_cname_record.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/dnsZones/zone1/CNAME/myrecord1
```
