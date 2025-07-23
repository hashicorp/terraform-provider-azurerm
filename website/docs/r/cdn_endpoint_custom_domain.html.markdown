---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_endpoint_custom_domain"
description: |-
  Manages a Custom Domain for a CDN Endpoint.
---

# azurerm_cdn_endpoint_custom_domain

Manages a Custom Domain for a CDN Endpoint.

!> **Note:** The CDN services from Edgio(formerly Verizon) was shut down on 15 January 2025 and is no longer available .

!> **Note:** Support for CDN services from Akamai was removed on 31 October 2023.

## Example Usage

```hcl
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
  sku                 = "Standard_Microsoft"
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
  name                = "example-domain.com"
  resource_group_name = "domain-rg"
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

* `name` - (Required) The name which should be used for this CDN Endpoint Custom Domain. Changing this forces a new CDN Endpoint Custom Domain to be created.

* `cdn_endpoint_id` - (Required) The ID of the CDN Endpoint. Changing this forces a new CDN Endpoint Custom Domain to be created.

* `host_name` - (Required) The host name of the custom domain. Changing this forces a new CDN Endpoint Custom Domain to be created.

* `cdn_managed_https` - (Optional) A `cdn_managed_https` block as defined below.

* `user_managed_https` - (Optional) A `user_managed_https` block as defined below.

~> **Note:** Only one of `cdn_managed_https` and `user_managed_https` can be specified.

---

A `cdn_managed_https` block supports the following:

* `certificate_type` - (Required) The type of HTTPS certificate. Possible values are `Shared` and `Dedicated`.

* `protocol_type` - (Required) The type of protocol. Possible values are `ServerNameIndication` and `IPBased`.

* `tls_version` - (Optional) The minimum TLS protocol version that is used for HTTPS. Possible values are `TLS10` (representing TLS 1.0/1.1), `TLS12` (representing TLS 1.2) and `None` (representing no minimums). Defaults to `TLS12`.

~> **Note:** Azure Services will require TLS 1.2+ by August 2025, please see this [announcement](https://azure.microsoft.com/en-us/updates/v2/update-retirement-tls1-0-tls1-1-versions-azure-services/) for more.

---

A `user_managed_https` block supports the following:

* `key_vault_secret_id` - (Required) The ID of the Key Vault Secret that contains the HTTPS certificate.

* `tls_version` - (Optional) The minimum TLS protocol version that is used for HTTPS. Possible values are `TLS10` (representing TLS 1.0/1.1), `TLS12` (representing TLS 1.2) and `None` (representing no minimums). Defaults to `TLS12`.

~> **Note:** Azure Services will require TLS 1.2+ by August 2025, please see this [announcement](https://azure.microsoft.com/en-us/updates/v2/update-retirement-tls1-0-tls1-1-versions-azure-services/) for more.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the CDN Endpoint Custom Domain.

## Timeouts

The `timeouts` block allows you to
specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 12 hours) Used when creating the Custom Domain for this CDN Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the CDN Endpoint Custom Domain.
* `update` - (Defaults to 24 hours) Used when updating the CDN Endpoint Custom Domain.
* `delete` - (Defaults to 12 hours) Used when deleting the CDN Endpoint Custom Domain.

## Import

CDN Endpoint Custom Domains can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_endpoint_custom_domain.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Cdn/profiles/profile1/endpoints/endpoint1/customDomains/domain1
```
