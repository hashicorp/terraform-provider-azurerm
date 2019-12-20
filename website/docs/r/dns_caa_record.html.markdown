---
subcategory: "DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_caa_record"
sidebar_current: "docs-azurerm-resource-dns-caa-record"
description: |-
  Manages a DNS CAA Record.
---

# azurerm_dns_caa_record

Enables you to manage DNS CAA Records within Azure DNS.

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

resource "azurerm_dns_caa_record" "example" {
  name                = "test"
  zone_name           = "${azurerm_dns_zone.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  ttl                 = 300

  record {
    flags = 0
    tag   = "issue"
    value = "example.com"
  }

  record {
    flags = 0
    tag   = "issue"
    value = "example.net"
  }

  record {
    flags = 0
    tag   = "issuewild"
    value = ";"
  }

  record {
    flags = 0
    tag   = "iodef"
    value = "mailto:terraform@nonexisting.tld"
  }

  tags = {
    Environment = "Production"
  }
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DNS CAA Record.

* `resource_group_name` - (Required) Specifies the resource group where the DNS Zone (parent resource) exists. Changing this forces a new resource to be created.

* `zone_name` - (Required) Specifies the DNS Zone where the resource exists. Changing this forces a new resource to be created.

* `ttl` - (Required) The Time To Live (TTL) of the DNS record in seconds.

* `record` - (Required) A list of values that make up the CAA record. Each `record` block supports fields documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

The `record` block supports:

* `flags` - (Required) Extensible CAA flags, currently only 1 is implemented to set the issuer critical flag.

* `tag` - (Required) A property tag, options are issue, issuewild and iodef.

* `value` - (Required) A property value such as a registrar domain.

## Attributes Reference

The following attributes are exported:

* `id` - The DNS CAA Record ID.
* `fqdn` - The FQDN of the DNS CAA Record.

## Import

CAA records can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dns_caa_record.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/dnszones/zone1/CAA/myrecord1
```
