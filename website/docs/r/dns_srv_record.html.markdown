---
subcategory: "DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_srv_record"
description: |-
  Manages a DNS SRV Record.
---

# azurerm_dns_srv_record

Enables you to manage DNS SRV Records within Azure DNS.

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

resource "azurerm_dns_srv_record" "example" {
  name                = "test"
  zone_name           = azurerm_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 300

  record {
    priority = 1
    weight   = 5
    port     = 8080
    target   = "target1.contoso.com"
  }

  tags = {
    Environment = "Production"
  }
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DNS SRV Record.

* `resource_group_name` - (Required) Specifies the resource group where the resource exists. Changing this forces a new resource to be created.

* `zone_name` - (Required) Specifies the DNS Zone where the DNS Zone (parent resource) exists. Changing this forces a new resource to be created.

* `ttl` - (Required) The Time To Live (TTL) of the DNS record in seconds.

* `record` - (Required) A list of values that make up the SRV record. Each `record` block supports fields documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

The `record` block supports:

* `priority` - (Required) Priority of the SRV record.

* `weight` - (Required) Weight of the SRV record.

* `port` - (Required) Port the service is listening on.

* `target` - (Required) FQDN of the service.


## Attributes Reference

The following attributes are exported:

* `id` - The DNS SRV Record ID.
* `fqdn` - The FQDN of the DNS SRV Record.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DNS SRV Record.
* `update` - (Defaults to 30 minutes) Used when updating the DNS SRV Record.
* `read` - (Defaults to 5 minutes) Used when retrieving the DNS SRV Record.
* `delete` - (Defaults to 30 minutes) Used when deleting the DNS SRV Record.

## Import

SRV records can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dns_srv_record.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/dnszones/zone1/SRV/myrecord1
```
