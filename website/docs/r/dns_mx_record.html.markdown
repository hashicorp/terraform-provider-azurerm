---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_mx_record"
sidebar_current: "docs-azurerm-resource-dns-mx-record"
description: |-
  Manages a DNS MX Record.
---

# azurerm_dns_mx_record

Enables you to manage DNS MX Records within Azure DNS.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acceptanceTestResourceGroup1"
  location = "West US"
}

resource "azurerm_dns_zone" "test" {
  name                = "mydomain.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_mx_record" "test" {
  name                = "test"
  zone_name           = "${azurerm_dns_zone.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  ttl                 = 300

  record {
    preference = 10
    exchange   = "mail1.contoso.com"
  }

  record {
    preference = 20
    exchange   = "mail2.contoso.com"
  }

  tags = {
    Environment = "Production"
  }
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DNS MX Record.

* `resource_group_name` - (Required) Specifies the resource group where the resource exists. Changing this forces a new resource to be created.

* `zone_name` - (Required) Specifies the DNS Zone where the resource exists. Changing this forces a new resource to be created.

* `ttl` - (Required) The Time To Live (TTL) of the DNS record in seconds.

* `record` - (Required) A list of values that make up the MX record. Each `record` block supports fields documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

The `record` block supports:

* `preference` - (Required) String representing the "preference” value of the MX records. Records with lower preference value take priority.

* `exchange` - (Required) The mail server responsible for the domain covered by the MX record.

## Attributes Reference

The following attributes are exported:

* `id` - The DNS MX Record ID.

## Import

MX records can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dns_mx_record.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/dnszones/zone1/MX/myrecord1
```
