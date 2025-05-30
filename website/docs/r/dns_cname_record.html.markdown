---
subcategory: "DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_cname_record"
description: |-
  Manages a DNS CNAME Record.
---

# azurerm_dns_cname_record

Enables you to manage DNS CNAME Records within Azure DNS.

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

resource "azurerm_dns_cname_record" "example" {
  name                = "test"
  zone_name           = azurerm_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 300
  record              = "contoso.com"
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

resource "azurerm_dns_cname_record" "target" {
  name                = "target"
  zone_name           = azurerm_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 300
  record              = "contoso.com"
}

resource "azurerm_dns_cname_record" "example" {
  name                = "test"
  zone_name           = azurerm_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 300
  target_resource_id  = azurerm_dns_cname_record.target.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DNS CNAME Record. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the resource group where the DNS Zone (parent resource) exists. Changing this forces a new resource to be created.

* `zone_name` - (Required) Specifies the DNS Zone where the resource exists. Changing this forces a new resource to be created.

* `ttl` - (Required) The Time To Live (TTL) of the DNS record in seconds.

* `record` - (Optional) The target of the CNAME.

* `target_resource_id` - (Optional) The Azure resource id of the target object. Conflicts with `record`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

~> **Note:** either `record` OR `target_resource_id` must be specified, but not both.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The DNS CName Record ID.

* `fqdn` - The FQDN of the DNS CName Record.

~> **Note:** The FQDN of the DNS CNAME Record which has a full-stop at the end is by design. Please see the documentation for more information.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DNS CNAME Record.

* `read` - (Defaults to 5 minutes) Used when retrieving the DNS CNAME Record.

* `update` - (Defaults to 30 minutes) Used when updating the DNS CNAME Record.

* `delete` - (Defaults to 30 minutes) Used when deleting the DNS CNAME Record.

## Import

CNAME records can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dns_cname_record.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/dnsZones/zone1/CNAME/myrecord1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network`: 2018-05-01
