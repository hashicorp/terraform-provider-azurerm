---
subcategory: "Private DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_srv_record"
description: |-
  Manages a Private DNS SRV Record.
---

# azurerm_private_dns_srv_record

Enables you to manage DNS SRV Records within Azure Private DNS.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_private_dns_zone" "example" {
  name                = "contoso.com"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_private_dns_srv_record" "example" {
  name                = "test"
  resource_group_name = azurerm_resource_group.example.name
  zone_name           = azurerm_private_dns_zone.example.name
  ttl                 = 300

  record {
    priority = 1
    weight   = 5
    port     = 8080
    target   = "target1.contoso.com"
  }

  record {
    priority = 10
    weight   = 10
    port     = 8080
    target   = "target2.contoso.com"
  }

  tags = {
    Environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DNS SRV Record. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the resource group where the resource exists. Changing this forces a new resource to be created.

* `zone_name` - (Required) Specifies the Private DNS Zone where the resource exists. Changing this forces a new resource to be created.

* `record` - (Required) One or more `record` blocks as defined below.

* `ttl` - (Required) The Time To Live (TTL) of the DNS record in seconds.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `record` block supports the following:

* `priority` - (Required) The priority of the SRV record.

* `weight` - (Required) The Weight of the SRV record.

* `port` - (Required) The Port the service is listening on.

* `target` - (Required) The FQDN of the service.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Private DNS SRV Record ID.

* `fqdn` - The FQDN of the DNS SRV Record.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Private DNS SRV Record.
* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS SRV Record.
* `update` - (Defaults to 30 minutes) Used when updating the Private DNS SRV Record.
* `delete` - (Defaults to 30 minutes) Used when deleting the Private DNS SRV Record.

## Import

Private DNS SRV Records can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_private_dns_srv_record.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/privateDnsZones/contoso.com/SRV/test
```
