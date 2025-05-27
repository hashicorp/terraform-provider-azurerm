---
subcategory: "DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_a_record"
description: |-
  Manages a DNS A Record.
---

# azurerm_dns_a_record

Enables you to manage DNS A Records within Azure DNS.

~> **Note:** [The Azure DNS API has a throttle limit of 500 read (GET) operations per 5 minutes](https://docs.microsoft.com/azure/azure-resource-manager/management/request-limits-and-throttling#network-throttling) - whilst the default read timeouts will work for most cases - in larger configurations you may need to set a larger [read timeout](https://www.terraform.io/language/resources/syntax#operation-timeouts) then the default 5min. Although, we'd generally recommend that you split the resources out into smaller Terraform configurations to avoid the problem entirely.

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

resource "azurerm_dns_a_record" "example" {
  name                = "test"
  zone_name           = azurerm_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 300
  records             = ["10.0.180.17"]
}
```

## Example Usage (Alias Record)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dns_zone" "example" {
  name                = "mydomain.com"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_public_ip" "example" {
  name                = "mypublicip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Dynamic"
  ip_version          = "IPv4"
}

resource "azurerm_dns_a_record" "example" {
  name                = "test"
  zone_name           = azurerm_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 300
  target_resource_id  = azurerm_public_ip.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DNS A Record. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the resource group where the DNS Zone (parent resource) exists. Changing this forces a new resource to be created.

* `zone_name` - (Required) Specifies the DNS Zone where the resource exists. Changing this forces a new resource to be created.

~> **Note:** The `zone_name` should be the name of resource `azurerm_dns_zone` instead of `azurerm_private_dns_zone`.

* `ttl` - (Required) The Time To Live (TTL) of the DNS record in seconds.

* `records` - (Optional) List of IPv4 Addresses. Conflicts with `target_resource_id`.

* `target_resource_id` - (Optional) The Azure resource id of the target object. Conflicts with `records`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

~> **Note:** either `records` OR `target_resource_id` must be specified, but not both.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The DNS A Record ID.

* `fqdn` - The FQDN of the DNS A Record.

~> **Note:** The FQDN of the DNS A Record which has a full-stop at the end is by design. Please [see the documentation](https://en.wikipedia.org/wiki/Fully_qualified_domain_name) for more information.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DNS A Record.

* `read` - (Defaults to 5 minutes) Used when retrieving the DNS A Record.

* `update` - (Defaults to 30 minutes) Used when updating the DNS A Record.

* `delete` - (Defaults to 30 minutes) Used when deleting the DNS A Record.

## Import

A records can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dns_a_record.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/dnsZones/zone1/A/myrecord1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network`: 2018-05-01
