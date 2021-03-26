---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_endpoint_custom_domain"
description: |- Manages a CDN Endpoint Custom Domain.
---

# azurerm_cdn_endpoint_custom_domain

Manages a CDN Endpoint Custom Domain.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

variable "domain_rg" {
  type = string
}

variable "domain_name" {
  type = string
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "west europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "example"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_cdn_profile" "example" {
  name                = "example-profile"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "example" {
  name                = "example-endpoint"
  profile_name        = azurerm_cdn_profile.example.name
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  origin {
    name      = "example"
    host_name = azurerm_storage_account.example.primary_blob_host
  }
}

data "azurerm_dns_zone" "example" {
  name                = var.domain_name
  resource_group_name = var.domain_rg
}

resource "azurerm_dns_cname_record" "example" {
  name                = "example"
  zone_name           = data.azurerm_dns_zone.example.name
  resource_group_name = data.azurerm_dns_zone.example.resource_group_name
  ttl                 = 3600
  target_resource_id  = azurerm_cdn_endpoint.example.id
}

resource "azurerm_cdn_endpoint_custom_domain" "example" {
  name            = "example-domain"
  cdn_endpoint_id = azurerm_cdn_endpoint.example.id
  host_name       = "${azurerm_dns_cname_record.example.name}.${data.azurerm_dns_zone.example.name}"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this CDN Endpoint Custom Domain. Changing this forces a new CDN
  Endpoint Custom Domain to be created.

* `cdn_endpoint_id` - (Required) The ID of the CDN Endpoint. Changing this forces a new CDN Endpoint Custom Domain to be
  created.

* `host_name` - (Required) The host name of the custom domain. Changing this forces a new CDN Endpoint Custom Domain to
  be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the CDN Endpoint Custom Domain.

## Timeouts

The `timeouts` block allows you to
specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 20 hours) Used when creating the CDN Endpoint Custom Domain.
* `read` - (Defaults to 5 minutes) Used when retrieving the CDN Endpoint Custom Domain.
* `delete` - (Defaults to 20 hours) Used when deleting the CDN Endpoint Custom Domain.

## Import

CDN Endpoint Custom Domains can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_endpoint_custom_domain.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Cdn/profiles/profile1/endpoints/endpoint1/customdomains/domain1
```
