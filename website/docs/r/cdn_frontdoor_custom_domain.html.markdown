---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_custom_domain"
description: |-
  Manages a Front Door (Standard/Premium) Custom Domain.
---

# azurerm_cdn_frontdoor_custom_domain

Manages a Front Door (Standard/Premium) Custom Domain.

-> **Note:** If you are using Terraform to manage your DNS Auth and DNS CNAME records for your Custom Domain you will need to add configuration blocks for both the `azurerm_dns_txt_record` (see the `Example DNS Auth TXT Record Usage` below) and the `azurerm_dns_cname_record` (see the `Example CNAME Record Usage` below) to your configuration file.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-cdn-frontdoor-custom-domain"
  location = "West Europe"
}

resource "azurerm_dns_zone" "example" {
  name                = "example-domain.com"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "example-cdn-frontdoor-profile"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Standard_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_endpoint" "example" {
  name                     = "example-cdn-frontdoor-endpoint"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
}

resource "azurerm_cdn_frontdoor_origin_group" "example" {
  name                     = "example-cdn-frontdoor-origin-group"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id

  load_balancing {
    additional_latency_in_milliseconds = 0
    sample_size                        = 4
    successful_samples_required        = 3
  }
}

resource "azurerm_cdn_frontdoor_origin" "example" {
  name                          = "example-cdn-frontdoor-origin"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.example.id

  enabled                        = true
  certificate_name_check_enabled = false

  host_name          = "www.example.com"
  http_port          = 80
  https_port         = 443
  origin_host_header = "www.example.com"
}

resource "azurerm_cdn_frontdoor_custom_domain" "example" {
  name                     = "example-cdn-frontdoor-custom-domain"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
  dns_zone_id              = azurerm_dns_zone.example.id
  host_name                = join(".", ["example-subdomain", azurerm_dns_zone.example.name])

  tls {
    certificate_type = "ManagedCertificate"
    minimum_version  = "TLS12"
  }
}

resource "azurerm_cdn_frontdoor_route" "example" {
  name                          = "example-cdn-frontdoor-route"
  cdn_frontdoor_endpoint_id     = azurerm_cdn_frontdoor_endpoint.example.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.example.id
  cdn_frontdoor_origin_ids      = [azurerm_cdn_frontdoor_origin.example.id]
  enabled                       = true

  forwarding_protocol    = "HttpsOnly"
  https_redirect_enabled = true
  patterns_to_match      = ["/*"]
  supported_protocols    = ["Http", "Https"]

  cdn_frontdoor_custom_domain_ids = [azurerm_cdn_frontdoor_custom_domain.example.id]
  link_to_default_domain          = false
}

resource "azurerm_cdn_frontdoor_firewall_policy" "example" {
  name                = "examplecdnfdwafpolicy"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = azurerm_cdn_frontdoor_profile.example.sku_name
  mode                = "Prevention"
}

resource "azurerm_cdn_frontdoor_security_policy" "example" {
  name                     = "example-cdn-frontdoor-security-policy"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id

  security_policies {
    firewall {
      cdn_frontdoor_firewall_policy_id = azurerm_cdn_frontdoor_firewall_policy.example.id

      association {
        domain {
          cdn_frontdoor_domain_id = azurerm_cdn_frontdoor_custom_domain.example.id
        }

        patterns_to_match = ["/*"]
      }
    }
  }
}
```

## Example DNS Auth TXT Record Usage

The name of your DNS TXT record should be in the format of `_dnsauth.<your_subdomain>`. So, for example, if we use the `host_name` in the example usage above you would create a DNS TXT record with the name of `_dnsauth.example-subdomain` which contains the value of the Front Door Custom Domain `validation_token` field. See the [product documentation](https://learn.microsoft.com/azure/frontdoor/standard-premium/how-to-add-custom-domain) for more information.

-> **Note:** Domain ownership validation is performed asynchronously by the Azure Front Door service (the domain typically transitions through states like `Submitting` and `Pending` before becoming `Approved`). If validation appears to be taking longer than expected, refer to the Azure Front Door documentation on [domain validation](https://learn.microsoft.com/azure/frontdoor/domain#domain-validation) and [domain validation states](https://learn.microsoft.com/azure/frontdoor/domain#domain-validation).

```hcl
resource "azurerm_dns_txt_record" "example" {
  name                = join(".", ["_dnsauth", "example-subdomain"])
  zone_name           = azurerm_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 3600

  record {
    value = azurerm_cdn_frontdoor_custom_domain.example.validation_token
  }
}
```

## Example CNAME Record Usage

~> **Note:** You **must** include the `depends_on` meta-argument which references both the `azurerm_cdn_frontdoor_route` and the `azurerm_cdn_frontdoor_security_policy` that are associated with your Custom Domain. The reason for these `depends_on` meta-arguments is because all of the resources for the Custom Domain need to be associated within Front Door before the CNAME record can be written to the domains DNS, else the CNAME validation will fail and Front Door will not enable traffic to the Domain.

```hcl
resource "azurerm_dns_cname_record" "example" {
  depends_on = [azurerm_cdn_frontdoor_route.example, azurerm_cdn_frontdoor_security_policy.example]

  name                = "example-subdomain"
  zone_name           = azurerm_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 3600
  record              = azurerm_cdn_frontdoor_endpoint.example.host_name
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Front Door Custom Domain. Changing this forces a new resource to be created.

* `cdn_frontdoor_profile_id` - (Required) The ID of the Front Door Profile. Changing this forces a new resource to be created.

* `host_name` - (Required) The host name of the domain. The `host_name` field must be the FQDN of your domain (e.g. `contoso.fabrikam.com`). Changing this forces a new resource to be created.

* `tls` - (Required) A `tls` block as defined below.

* `dns_zone_id` - (Optional) The ID of the Azure DNS Zone which should be used for this Front Door Custom Domain. If you are using Azure to host your [DNS domains](https://learn.microsoft.com/azure/dns/dns-overview), you must delegate the domain provider's domain name system (DNS) to an Azure DNS Zone. For more information, see [Delegate a domain to Azure DNS](https://learn.microsoft.com/azure/dns/dns-delegate-domain-azure-dns). Otherwise, if you're using your own domain provider to handle your DNS, you must validate the Front Door Custom Domain by creating the DNS TXT records manually.

<!-- * `pre_validated_cdn_frontdoor_custom_domain_id` - (Optional) The resource ID of the pre-validated Front Door Custom Domain. This domain type is used when you wish to onboard a validated Azure service domain, and then configure the Azure service behind an Azure Front Door.

-> **Note:** Currently `pre_validated_cdn_frontdoor_custom_domain_id` only supports domains validated by Static Web App. -->

---

A `tls` block supports the following:

* `certificate_type` - (Optional) Defines the source of the SSL certificate. Possible values include `CustomerCertificate` and `ManagedCertificate`. Defaults to `ManagedCertificate`.

-> **Note:** It may take up to 15 minutes for the Front Door service to validate the state and domain ownership of the custom domain.

* `cdn_frontdoor_secret_id` - (Optional) The resource ID of the Front Door Secret.

~> **Note:** `cdn_frontdoor_secret_id` must be specified when `certificate_type` is `CustomerCertificate`, and must not be specified when `certificate_type` is `ManagedCertificate`.

* `cipher_suite` - (Optional) A `cipher_suite` block as defined below.

-> **Note:** Azure Front Door returns a default cipher suite even when not configured. The provider only persists the `cipher_suite` block in state when it is explicitly configured or when Azure returns a non-default/customized configuration.

* `minimum_version` - (Optional) TLS protocol version that will be used for HTTPS. Possible values include `TLS12`. Defaults to `TLS12`.

~> **Note:** As of March 1, 2025, support for Transport Layer Security (TLS) 1.0 and 1.1 has been retired for Azure Front Door, and all connections to Azure Front Door must employ TLS 1.2 or later. For more information, see the product [announcement](https://azure.microsoft.com/updates/v2/update-retirement-tls1-0-tls1-1-versions-azure-services/) for more details.

---

A `cipher_suite` block supports the following:

* `type` - (Required) The cipher suite policy type. Possible values include `Customized`, `TLS12_2022`, and `TLS12_2023`.

* `custom_ciphers` - (Optional) A `custom_ciphers` block as defined below.

~> **Note:** `custom_ciphers` can only be specified when `type` is `Customized`. At least one cipher suite must be selected in `custom_ciphers` when `type` is `Customized`.

---

A `custom_ciphers` block supports the following:

* `tls12` - (Optional) A set of TLS 1.2 cipher suites. Possible values include `DHE_RSA_AES128_GCM_SHA256`, `DHE_RSA_AES256_GCM_SHA384`, `ECDHE_RSA_AES128_GCM_SHA256`, `ECDHE_RSA_AES128_SHA256`, `ECDHE_RSA_AES256_GCM_SHA384`, and `ECDHE_RSA_AES256_SHA384`.

~> **Note:** When `minimum_version` is `TLS12` (default), at least one cipher suite must be specified in `custom_ciphers.tls12`.

* `tls13` - (Optional) A set of TLS 1.3 cipher suites. Possible values include `TLS_AES_128_GCM_SHA256` and `TLS_AES_256_GCM_SHA384`.

~> **Note:** When `tls13` is specified, it must contain both `TLS_AES_128_GCM_SHA256` and `TLS_AES_256_GCM_SHA384`.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Front Door Custom Domain.

* `expiration_date` - The date time that the token expires.

* `validation_token` - Challenge used for DNS TXT record or file based validation.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

~> **Note:** Deleting a Front Door Custom Domain can take a significant amount of time while the Azure Front Door service performs backend synchronization. During this period, the domain may remain visible in the Azure Portal with a provisioning state of `Deleting`. If you encounter `context deadline exceeded` during deletion, increase the `delete` timeout accordingly.

* `create` - (Defaults to 12 hours) Used when creating the Front Door Custom Domain.
* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Custom Domain.
* `update` - (Defaults to 24 hours) Used when updating the Front Door Custom Domain.
* `delete` - (Defaults to 12 hours) Used when deleting the Front Door Custom Domain.

## Import

Front Door Custom Domains can be imported using the resource id, e.g.

```shell
terraform import azurerm_cdn_frontdoor_custom_domain.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/customDomains/customDomain1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Cdn` - 2025-04-15
