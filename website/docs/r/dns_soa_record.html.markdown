---
subcategory: "DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_soa_record"
description: |-
  Manages a DNS SOA Record.
---

# azurerm_dns_soa_record

Enables you to manage DNS SOA Records within Azure DNS.

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

resource "azurerm_dns_soa_record" "example" {
  name                = "mysoarecord"
  zone_name           = azurerm_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 300

  record {
    email         = "mydomain.com"
    expire_time   = 2419200
    host_name     = "target.contoso.com"
    minimum_ttl   = 300
  	refresh_time  = 3600
    retry_time    = 300
  	serial_number = 1
  }

  tags = {
    Environment = "Production"
  }
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DNS SOA Record.

* `resource_group_name` - (Required) Specifies the resource group where the resource exists. Changing this forces a new resource to be created.

* `zone_name` - (Required) Specifies the DNS Zone where the DNS Zone (parent resource) exists.

* `ttl` - (Required) The Time To Live (TTL) of the DNS record in seconds.

* `record` - (Required) An `record` block as defined below. 

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `record` block supports:

* `email` - (Required) The email of the SOA record.

* `expire_time` - (Required) The expire time of the SOA record.

* `host_name` - (Required) The domain name of the authoritative name server of the SOA record.

* `minimum_ttl` - (Required) The minimum value of the SOA record. By convention this is used to determine the negative caching duration.

* `refresh_time` - (Required) The refresh value of the SOA record.

* `retry_time` - (Required) The retry time of the SOA record.

* `serial_number` - (Required) The serial number of the SOA record.

## Attributes Reference

The following attributes are exported:

* `id` - The DNS SOA Record ID.

* `fqdn` - The FQDN of the DNS SOA Record.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DNS SOA Record.
* `update` - (Defaults to 30 minutes) Used when updating the DNS SOA Record.
* `read` - (Defaults to 5 minutes) Used when retrieving the DNS SOA Record.
* `delete` - (Defaults to 30 minutes) Used when deleting the DNS SOA Record.

## Import

SOA records can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dns_soa_record.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/dnszones/zone1/SOA/myrecord1
```
