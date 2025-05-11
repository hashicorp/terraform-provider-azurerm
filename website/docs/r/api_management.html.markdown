---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management"
description: |-
  Manages an API Management Service.
---

# azurerm_api_management

Manages an API Management Service.

## Disclaimers

-> **Note:** When creating a new API Management resource in version 3.0 of the AzureRM Provider and later, please be aware that the AzureRM Provider will now clean up any sample APIs and Products created by the Azure API during the creation of the API Management resource.

-> **Note:** Version 2.77 and later of the Azure Provider include a Feature Toggle which will purge an API Management resource on destroy, rather than the default soft-delete. See [the Features block documentation](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block) for more information on Feature Toggles within Terraform.

~> **Note:** It's possible to define Custom Domains both within [the `azurerm_api_management` resource](api_management.html) via the `hostname_configurations` block and by using [the `azurerm_api_management_custom_domain` resource](api_management_custom_domain.html). However it's not possible to use both methods to manage Custom Domains within an API Management Service, since there'll be conflicts.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_api_management" "example" {
  name                = "example-apim"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  publisher_name      = "My Company"
  publisher_email     = "company@terraform.io"

  sku_name = "Developer_1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the API Management Service. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the API Management Service exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the API Management Service should exist. Changing this forces a new resource to be created.

* `publisher_name` - (Required) The name of publisher/company.

* `publisher_email` - (Required) The email of publisher/company.

* `sku_name` - (Required) `sku_name` is a string consisting of two parts separated by an underscore(\_). The first part is the `name`, valid values include: `Consumption`, `Developer`, `Basic`, `Standard` and `Premium`. The second part is the `capacity` (e.g. the number of deployed units of the `sku`), which must be a positive `integer` (e.g. `Developer_1`).

~> **Note:** Premium SKUs are limited to a default maximum of 12 (i.e. `Premium_12`), this can, however, be increased via support request.

~> **Note:** Consumption SKU capacity should be 0 (e.g. `Consumption_0`) as this tier includes automatic scaling.

---

* `additional_location` - (Optional) One or more `additional_location` blocks as defined below.

* `certificate` - (Optional) One or more `certificate` blocks (up to 10) as defined below.

* `client_certificate_enabled` - (Optional) Enforce a client certificate to be presented on each request to the gateway? This is only supported when SKU type is `Consumption`.

* `delegation` - (Optional) A `delegation` block as defined below.

* `gateway_disabled` - (Optional) Disable the gateway in main region? This is only supported when `additional_location` is set.

* `min_api_version` - (Optional) The version which the control plane API calls to API Management service are limited with version equal to or newer than.

* `zones` - (Optional) Specifies a list of Availability Zones in which this API Management service should be located.

~> **Note:** Availability zones are only supported in the Premium tier.

* `identity` - (Optional) An `identity` block as defined below.

* `hostname_configuration` - (Optional) A `hostname_configuration` block as defined below.

* `notification_sender_email` - (Optional) Email address from which the notification will be sent.

* `protocols` - (Optional) A `protocols` block as defined below.

* `security` - (Optional) A `security` block as defined below.

* `sign_in` - (Optional) A `sign_in` block as defined below.

* `sign_up` - (Optional) A `sign_up` block as defined below.

* `tenant_access` - (Optional) A `tenant_access` block as defined below.

* `public_ip_address_id` - (Optional) ID of a standard SKU IPv4 Public IP.

~> **Note:** Custom public IPs are only supported on the `Premium` and `Developer` tiers when deployed in a virtual network.

* `public_network_access_enabled` - (Optional) Is public access to the service allowed? Defaults to `true`.

~> **Note:** This option is applicable only to the Management plane, not the API gateway or Developer portal. It is required to be `true` on the creation.

* `virtual_network_type` - (Optional) The type of virtual network you want to use, valid values include: `None`, `External`, `Internal`. Defaults to `None`.

~> **Note:** Please ensure that in the subnet, inbound port 3443 is open when `virtual_network_type` is `Internal` or `External`. Additionally, please ensure other necessary ports are open according to [api management network configuration](https://learn.microsoft.com/azure/api-management/virtual-network-reference).

* `virtual_network_configuration` - (Optional) A `virtual_network_configuration` block as defined below. Required when `virtual_network_type` is `External` or `Internal`.

* `tags` - (Optional) A mapping of tags assigned to the resource.

---

A `additional_location` block supports the following:

* `location` - (Required) The name of the Azure Region in which the API Management Service should be expanded to.

* `capacity` - (Optional) The number of compute units in this region. Defaults to the capacity of the main region.

* `zones` - (Optional) A list of availability zones.

* `public_ip_address_id` - (Optional) ID of a standard SKU IPv4 Public IP.

~> **Note:** Availability zones and custom public IPs are only supported in the Premium tier.

* `virtual_network_configuration` - (Optional) A `virtual_network_configuration` block as defined below. Required when `virtual_network_type` is `External` or `Internal`.

* `gateway_disabled` - (Optional) Only valid for an Api Management service deployed in multiple locations. This can be used to disable the gateway in this additional location.

---

A `certificate` block supports the following:

* `encoded_certificate` - (Required) The Base64 Encoded PFX or Base64 Encoded X.509 Certificate.

* `store_name` - (Required) The name of the Certificate Store where this certificate should be stored. Possible values are `CertificateAuthority` and `Root`.

* `certificate_password` - (Optional) The password for the certificate.

---

A `delegation` block supports the following:

* `subscriptions_enabled` - (Optional) Should subscription requests be delegated to an external url? Defaults to `false`.

* `user_registration_enabled` - (Optional) Should user registration requests be delegated to an external url? Defaults to `false`.

* `url` - (Optional) The delegation URL.

* `validation_key` - (Optional) A base64-encoded validation key to validate, that a request is coming from Azure API Management.

---

A `hostname_configuration` block supports the following:

* `management` - (Optional) One or more `management` blocks as documented below.

* `portal` - (Optional) One or more `portal` blocks as documented below.

* `developer_portal` - (Optional) One or more `developer_portal` blocks as documented below.

* `proxy` - (Optional) One or more `proxy` blocks as documented below.

* `scm` - (Optional) One or more `scm` blocks as documented below.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this API Management Service. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) A list of User Assigned Managed Identity IDs to be assigned to this API Management Service.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `management`, `portal`, `developer_portal` and `scm` block supports the following:

* `host_name` - (Required) The Hostname to use for the Management API.

* `key_vault_certificate_id` - (Optional) The ID of the Key Vault Secret containing the SSL Certificate, which must be of the type `application/x-pkcs12`.

-> **Note:** Setting this field requires the `identity` block to be specified, since this identity is used for to retrieve the Key Vault Certificate. Possible values are versioned or versionless secret ID. Auto-updating the Certificate from the Key Vault requires the Secret version isn't specified.

* `certificate` - (Optional) The Base64 Encoded Certificate.

* `certificate_password` - (Optional) The password associated with the certificate provided above.

-> **Note:** Either `key_vault_certificate_id` or `certificate` and `certificate_password` must be specified.

* `negotiate_client_certificate` - (Optional) Should Client Certificate Negotiation be enabled for this Hostname? Defaults to `false`.

* `ssl_keyvault_identity_client_id` - (Optional) System or User Assigned Managed identity clientId as generated by Azure AD, which has `GET` access to the keyVault containing the SSL certificate.

-> **Note:** If a User Assigned Managed identity is specified for `ssl_keyvault_identity_client_id` then this identity must be associated to the `azurerm_api_management` within an `identity` block.

---

A `proxy` block supports the following:

* `default_ssl_binding` - (Optional) Is the certificate associated with this Hostname the Default SSL Certificate? This is used when an SNI header isn't specified by a client. Defaults to `false`.

* `host_name` - (Required) The Hostname to use for the Management API.

* `key_vault_certificate_id` - (Optional) The ID of the Key Vault Secret containing the SSL Certificate, which must be of the type `application/x-pkcs12`.

-> **Note:** Setting this field requires the `identity` block to be specified, since this identity is used for to retrieve the Key Vault Certificate. Auto-updating the Certificate from the Key Vault requires the Secret version isn't specified.

* `certificate` - (Optional) The Base64 Encoded Certificate.

* `certificate_password` - (Optional) The password associated with the certificate provided above.

-> **Note:** Either `key_vault_certificate_id` or `certificate` and `certificate_password` must be specified.

* `negotiate_client_certificate` - (Optional) Should Client Certificate Negotiation be enabled for this Hostname? Defaults to `false`.

* `ssl_keyvault_identity_client_id` - (Optional) The Managed Identity Client ID to use to access the Key Vault. This Identity must be specified in the `identity` block to be used.

---

A `protocols` block supports the following:

* `http2_enabled` - (Optional) Should HTTP/2 be supported by the API Management Service? Defaults to `false`.

---

A `security` block supports the following:

* `backend_ssl30_enabled` - (Optional) Should SSL 3.0 be enabled on the backend of the gateway? Defaults to `false`.

-> **Note:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Ssl30` field

* `backend_tls10_enabled` - (Optional) Should TLS 1.0 be enabled on the backend of the gateway? Defaults to `false`.

-> **Note:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls10` field

* `backend_tls11_enabled` - (Optional) Should TLS 1.1 be enabled on the backend of the gateway? Defaults to `false`.

-> **Note:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls11` field

* `frontend_ssl30_enabled` - (Optional) Should SSL 3.0 be enabled on the frontend of the gateway? Defaults to `false`.

-> **Note:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Ssl30` field

* `frontend_tls10_enabled` - (Optional) Should TLS 1.0 be enabled on the frontend of the gateway? Defaults to `false`.

-> **Note:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls10` field

* `frontend_tls11_enabled` - (Optional) Should TLS 1.1 be enabled on the frontend of the gateway? Defaults to `false`.

-> **Note:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls11` field

* `tls_ecdhe_ecdsa_with_aes128_cbc_sha_ciphers_enabled` - (Optional) Should the `TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA` cipher be enabled? Defaults to `false`.

-> **Note:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA` field

* `tls_ecdhe_ecdsa_with_aes256_cbc_sha_ciphers_enabled` - (Optional) Should the `TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA` cipher be enabled? Defaults to `false`.

-> **Note:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA` field

* `tls_ecdhe_rsa_with_aes128_cbc_sha_ciphers_enabled` - (Optional) Should the `TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA` cipher be enabled? Defaults to `false`.

-> **Note:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA` field

* `tls_ecdhe_rsa_with_aes256_cbc_sha_ciphers_enabled` - (Optional) Should the `TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA` cipher be enabled? Defaults to `false`.

-> **Note:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA` field

* `tls_rsa_with_aes128_cbc_sha256_ciphers_enabled` - (Optional) Should the `TLS_RSA_WITH_AES_128_CBC_SHA256` cipher be enabled? Defaults to `false`.

-> **Note:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_128_CBC_SHA256` field

* `tls_rsa_with_aes128_cbc_sha_ciphers_enabled` - (Optional) Should the `TLS_RSA_WITH_AES_128_CBC_SHA` cipher be enabled? Defaults to `false`.

-> **Note:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_128_CBC_SHA` field

* `tls_rsa_with_aes128_gcm_sha256_ciphers_enabled` - (Optional) Should the `TLS_RSA_WITH_AES_128_GCM_SHA256` cipher be enabled? Defaults to `false`.

-> **Note:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_128_GCM_SHA256` field

* `tls_rsa_with_aes256_gcm_sha384_ciphers_enabled` - (Optional) Should the `TLS_RSA_WITH_AES_256_GCM_SHA384` cipher be enabled? Defaults to `false`.

-> **Note:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_256_GCM_SHA384` field

* `tls_rsa_with_aes256_cbc_sha256_ciphers_enabled` - (Optional) Should the `TLS_RSA_WITH_AES_256_CBC_SHA256` cipher be enabled? Defaults to `false`.

-> **Note:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_256_CBC_SHA256` field

* `tls_rsa_with_aes256_cbc_sha_ciphers_enabled` - (Optional) Should the `TLS_RSA_WITH_AES_256_CBC_SHA` cipher be enabled? Defaults to `false`.

-> **Note:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_256_CBC_SHA` field

* `triple_des_ciphers_enabled` - (Optional) Should the `TLS_RSA_WITH_3DES_EDE_CBC_SHA` cipher be enabled for alL TLS versions (1.0, 1.1 and 1.2)? 

-> **Note:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TripleDes168` field

---

A `sign_in` block supports the following:

* `enabled` - (Required) Should anonymous users be redirected to the sign in page?

---

A `sign_up` block supports the following:

* `enabled` - (Required) Can users sign up on the development portal?

* `terms_of_service` - (Required) A `terms_of_service` block as defined below.

---

A `tenant_access` block supports the following:

* `enabled` - (Required) Should the access to the management API be enabled?

---

A `virtual_network_configuration` block supports the following:

* `subnet_id` - (Required) The id of the subnet that will be used for the API Management.

---

A `terms_of_service` block supports the following:

* `consent_required` - (Required) Should the user be asked for consent during sign up?

* `enabled` - (Required) Should Terms of Service be displayed during sign up?.

* `text` - (Optional) The Terms of Service which users are required to agree to in order to sign up.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Service.

* `additional_location` - Zero or more `additional_location` blocks as documented below.

* `gateway_url` - The URL of the Gateway for the API Management Service.

* `gateway_regional_url` - The Region URL for the Gateway of the API Management Service.

* `identity` - An `identity` block as defined below.

* `hostname_configuration` - A `hostname_configuration` block as defined below.

* `management_api_url` - The URL for the Management API associated with this API Management service.

* `portal_url` - The URL for the Publisher Portal associated with this API Management service.

* `developer_portal_url` - The URL for the Developer Portal associated with this API Management service.

* `public_ip_addresses` - The Public IP addresses of the API Management Service.

* `private_ip_addresses` - The Private IP addresses of the API Management Service.

* `scm_url` - The URL for the SCM (Source Code Management) Endpoint associated with this API Management service.

* `tenant_access` - The `tenant_access` block as documented below.

---

An `additional_location` block exports the following:

* `gateway_regional_url` - The URL of the Regional Gateway for the API Management Service in the specified region.

* `public_ip_addresses` - Public Static Load Balanced IP addresses of the API Management service in the additional location. Available only for Basic, Standard and Premium SKU.

* `private_ip_addresses` - The Private IP addresses of the API Management Service. Available only when the API Manager instance is using Virtual Network mode.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

---

A `tenant_access` block exports the following:

* `tenant_id` - The identifier for the tenant access information contract.

* `primary_key` - Primary access key for the tenant access information contract.

* `secondary_key` - Secondary access key for the tenant access information contract.

---

The `certificate` block exports the following:

* `expiry` - The expiration date of the certificate in RFC3339 format: `2000-01-02T03:04:05Z`.

* `thumbprint` - The thumbprint of the certificate.

* `subject` - The subject of the certificate.

---

The `hostname_configuration` block exports the following:

* `proxy` - A `proxy` block as defined below.

---

The `proxy` block exports the following:

* `certificate_source` - The source of the certificate.

* `certificate_status` - The status of the certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the API Management Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Service.
* `update` - (Defaults to 3 hours) Used when updating the API Management Service.
* `delete` - (Defaults to 3 hours) Used when deleting the API Management Service.

## Import

API Management Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1
```
