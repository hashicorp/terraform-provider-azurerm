---
subcategory: "Private DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_txt_record"
description: |-
  Manages a Private DNS TXT Record.
---

# azurerm_private_dns_txtrecord

Enables you to manage DNS TXT Records within Azure Private DNS.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "acceptanceTestResourceGroup1"
  location = "West US"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "contoso.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_txt_record" "test" {
  name                = "test"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
  ttl                 = 300

  record {
    value = "v=spf1 mx ~all"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DNS TXT Record. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the resource group where the resource exists. Changing this forces a new resource to be created.

* `zone_name` - (Required) Specifies the Private DNS Zone where the resource exists. Changing this forces a new resource to be created.

* `record` - (Required) One or more `record` blocks as defined below.

* `ttl ` - (Required) The Time To Live (TTL) of the DNS record in seconds.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `record` block supports the following:

* `value` - (Required) The value of the TXT record. Max length: 1024 characters


## Attributes Reference

The following attributes are exported:

* `id` - The Private DNS TXT Record ID.

* `fqdn` - The FQDN of the DNS TXT Record.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Private DNS TXT Record.
* `update` - (Defaults to 30 minutes) Used when updating the Private DNS TXT Record.
* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS TXT Record.
* `delete` - (Defaults to 30 minutes) Used when deleting the Private DNS TXT Record.

## Import

Private DNS TXT Records can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_private_dns_txt_record.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/privateDnsZones/contoso.com/TXT/test
```
