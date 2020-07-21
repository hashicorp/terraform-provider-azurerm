---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_endpoint_custom_domain"
description: |-
  Manages a CDN Endpoint Custom Domain.
---

# azurerm_cdn_endpoint_custom_domain

Manages a CDN Endpoint Custom Domain.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

variable "domain_name" {
  type = string
}

variable "dns_zone_name" {
  type = string
}

variable "dns_zone_rg" {
  type = string
}

resource "azurerm_resource_group" "example" {
  name     = "example-test"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplecdnsa"
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
    host_name = trimsuffix(trimprefix(trimprefix(azurerm_storage_account.example.primary_blob_endpoint, "https://"), "http://"), "/")
  }
}

data "azurerm_dns_zone" "example" {
  name                = var.dns_zone_name
  resource_group_name = var.dns_zone_rg
}

resource "azurerm_dns_cname_record" "example" {
  name                = var.sub_domain_name
  zone_name           = data.azurerm_dns_zone.example.name
  resource_group_name = data.azurerm_dns_zone.example.resource_group_name
  ttl                 = 3600
  target_resource_id  = azurerm_cdn_endpoint.example.id
}

resource "azurerm_cdn_endpoint_custom_domain" "example" {
  name            = "example-customdomain"
  cdn_endpoint_id = azurerm_cdn_endpoint.example.id
  host_name       = "${azurerm_dns_cname_record.example.name}.${var.domain_name}"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this CDN Endpoint Custom Domain. Changing this forces a new CDN Endpoint Custom Domain to be created.

* `cdn_endpoint_id` - (Required) The ID of the CDN Endpoint. Changing this forces a new CDN Endpoint Custom Domain to be created.

* `host_name` - (Required) The host name of the custom domain. Changing this forces a new CDN Endpoint Custom Domain to be created.

* `cdn_managed_https_settings` (Optional) - A `cdn_managed_https_settings` block as defined below. Only one of `cdn_managed_https_settings` and `user_managed_https_settings` can be specified.

* `user_managed_https_settings` (Optional) - A `user_managed_https_settings` block as defined below. Only one of `cdn_managed_https_settings` and `user_managed_https_settings` can be specified.

!> **Warning** It is allowed to update the HTTPS settings on the CDN Endpoint Custom Domain only by toggling it. It is not allowed to in-place update the HTTPS settings, which means it is not allowed to modify an already enabled http settings with different attributes. This is because setting different HTTPS settings will need a disable-then-enable process. When HTTPS settings got disabled, the service will take 8 hours to clean up your previous enablement request for the same custom domain and there is no way to get notification when that clean up has done.

---

A `cdn_managed_https_settings` block supports the following:

* `certificate_type` - (Required) The type of the HTTPS certificate. Possible values are `Shared` and `Dedicated`.

---

A `user_managed_https_settings` block supports the following:

* `subscription_id` - (Required) The subscription ID where the Key Vault Certificate that contains the HTTPS certificate resides in.

* `resource_group_name` - (Required) The name of Resource Group where the Key Vault Certificate that contains the HTTPS certificate resides in.

* `vault_name` - (Required) The name of Key Vault where the Key Vault Certificate that contains the HTTPS certificate resides in.

* `secret_name` - (Required) The name of Key Vault Certificate that contains the HTTPS certificate.

* `secret_version` - (Required) The version of Key Vault Certificate that contains the HTTPS certificate.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the CDN Endpoint Custom Domain.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 hours) Used when creating the CDN Endpoint Custom Domain.
* `read` - (Defaults to 5 minutes) Used when retrieving the CDN Endpoint Custom Domain.
* `update` - (Defaults to 10 hours) Used when updating the CDN Endpoint Custom Domain.
* `delete` - (Defaults to 10 hours) Used when deleting the CDN Endpoint Custom Domain.

## Import

CDN Endpoint Custom Domains can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_endpoint_custom_domain.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Cdn/profiles/profile1/endpoints/endpoint1/customdomains/domain1
```
