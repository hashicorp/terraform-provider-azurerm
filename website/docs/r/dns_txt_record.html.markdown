---
subcategory: "DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_txt_record"
description: |-
  Manages a DNS TXT Record.
---

# azurerm_dns_txt_record

Enables you to manage DNS TXT Records within Azure DNS.

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

resource "azurerm_dns_txt_record" "example" {
  name                = "test"
  zone_name           = azurerm_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 300

  record {
    value = "google-site-authenticator"
  }

  record {
    value = "more site information here"
  }

  tags = {
    Environment = "Production"
  }
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DNS TXT Record.

* `resource_group_name` - (Required) Specifies the resource group where the resource exists. Changing this forces a new resource to be created.

* `zone_name` - (Required) Specifies the DNS Zone where the DNS Zone (parent resource) exists. Changing this forces a new resource to be created.

* `ttl` - (Required) The Time To Live (TTL) of the DNS record in seconds.

* `record` - (Required) A list of values that make up the txt record. Each `record` block supports fields documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

The `record` block supports:

* `value` - (Required) The value of the record. Max length: 1024 characters

## Attributes Reference

The following attributes are exported:

* `id` - The DNS TXT Record ID.
* `fqdn` - The FQDN of the DNS TXT Record.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DNS TXT Record.
* `update` - (Defaults to 30 minutes) Used when updating the DNS TXT Record.
* `read` - (Defaults to 5 minutes) Used when retrieving the DNS TXT Record.
* `delete` - (Defaults to 30 minutes) Used when deleting the DNS TXT Record.

## Import

TXT records can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dns_txt_record.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/dnszones/zone1/TXT/myrecord1
```
