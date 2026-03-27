---
subcategory: "DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_cname_record_public_ip_address_association"
description: |-
  Associates a DNS CNAME Record with a Public IP Address by setting the reverse FQDN.
---

# azurerm_dns_cname_record_public_ip_address_association

Associates a [DNS CNAME Record](dns_cname_record.html) with a [Public IP Address](public_ip.html) by setting the `reverse_fqdn` property on the Public IP Address to the FQDN of the CNAME Record.

This resource is useful for breaking circular dependencies between DNS CNAME Records and Public IP Addresses when you need to reference the FQDN of a CNAME Record in a Public IP Address configuration.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dns_zone" "example" {
  name                = "example.com"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_public_ip" "example" {
  name                = "example-pip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"
  domain_name_label   = "example-pip"
}

resource "azurerm_dns_cname_record" "example" {
  name                = "www"
  zone_name           = azurerm_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 300
  record              = azurerm_public_ip.example.fqdn
}

resource "azurerm_dns_cname_record_public_ip_address_association" "example" {
  dns_cname_record_id  = azurerm_dns_cname_record.example.id
  public_ip_address_id = azurerm_public_ip.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `dns_cname_record_id` - (Required) The ID of the DNS CNAME Record. Changing this forces a new resource to be created.

* `public_ip_address_id` - (Required) The ID of the Public IP Address. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the DNS CNAME Record Public IP Address Association.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DNS CNAME Record Public IP Address Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the DNS CNAME Record Public IP Address Association.
* `delete` - (Defaults to 30 minutes) Used when deleting the DNS CNAME Record Public IP Address Association.

## Import

DNS CNAME Record Public IP Address Associations can be imported using the `resource id` of the DNS CNAME Record and the Public IP Address separated by `|`, e.g.

```shell
terraform import azurerm_dns_cname_record_public_ip_address_association.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/dnsZones/example.com/CNAME/www|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/publicIPAddresses/mypip1"
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network` - 2025-01-01, 2018-05-01
