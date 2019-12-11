---
subcategory: "Private DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_mx_record"
sidebar_current: "docs-azurerm-resource-private-dns-mx-record"
description: |-
  Manages a Private DNS MX Record.
---

# azurerm_private_dns_mx_record

Enables you to manage DNS MX Records within Azure Private DNS.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "acceptanceTestResourceGroup1"
  location = "West US"
}

resource "azurerm_private_dns_zone" "example" {
  name                = "contoso.com"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_private_dns_mx_record" "example" {
  name                = "example"
  resource_group_name = "${azurerm_resource_group.example.name}"
  zone_name           = "${azurerm_private_dns_zone.example.name}"
  ttl                 = 300

  record {
    preference = 10
    exchange   = "mx1.contoso.com"
  }

  record {
    preference = 20
    exchange   = "backupmx.contoso.com"
  }

  tags = {
    Environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the DNS MX Record. Changing this forces a new resource to be created. Default to '@' for root zone entry.

* `resource_group_name` - (Required) Specifies the resource group where the resource exists. Changing this forces a new resource to be created.

* `zone_name` - (Required) Specifies the Private DNS Zone where the resource exists. Changing this forces a new resource to be created.

* `record` - (Required) One or more `record` blocks as defined below.

* `ttl ` - (Required) The Time To Live (TTL) of the DNS record in seconds.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `record` block supports the following:

* `preference` - (Required) The preference of the MX record.

* `exchange` - (Required) The FQDN of the exchange to MX record points to.

## Attributes Reference

The following attributes are exported:

* `id` - The Private DNS MX Record ID.

## Import

Private DNS MX Records can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_private_dns_srv_record.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/privateDnsZones/contoso.com/MX/@
```
