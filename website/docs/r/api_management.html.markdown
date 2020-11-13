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

  policy {
    xml_content = <<XML
    <policies>
      <inbound />
      <backend />
      <outbound />
      <on-error />
    </policies>
XML

  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the API Management Service. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the API Management Service exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the API Management Service should be exist. Changing this forces a new resource to be created.

* `publisher_name` - (Required) The name of publisher/company.

* `publisher_email` - (Required) The email of publisher/company.

* `sku_name` - (Required) `sku_name` is a string consisting of two parts separated by an underscore(\_). The first part is the `name`, valid values include: `Consumption`, `Developer`, `Basic`, `Standard` and `Premium`. The second part is the `capacity` (e.g. the number of deployed units of the `sku`), which must be a positive `integer` (e.g. `Developer_1`).

---

* `additional_location` - (Optional) One or more `additional_location` blocks as defined below.

* `certificate` - (Optional) One or more (up to 10) `certificate` blocks as defined below.

* `identity` - (Optional) An `identity` block is documented below.

* `hostname_configuration` - (Optional) A `hostname_configuration` block as defined below.

* `notification_sender_email` - (Optional) Email address from which the notification will be sent.

* `policy` - (Optional) A `policy` block as defined below.

* `protocols` - (Optional) A `protocols` block as defined below.

* `security` - (Optional) A `security` block as defined below.

* `sign_in` - (Optional) A `sign_in` block as defined below.

* `sign_up` - (Optional) A `sign_up` block as defined below.

* `virtual_network_type` - (Optional) The type of virtual network you want to use, valid values include: `None`, `External`, `Internal`. 
> **NOTE:** Please ensure that in the subnet, inbound port 3443 is open when `virtual_network_type` is `Internal` or `External`. And please ensure other necessary ports are open according to [api management network configuration](https://docs.microsoft.com/en-us/azure/api-management/api-management-using-with-vnet#-common-network-configuration-issues).

* `virtual_network_configuration` - (Optional) A `virtual_network_configuration` block as defined below. Required when `virtual_network_type` is `External` or `Internal`.

* `tags` - (Optional) A mapping of tags assigned to the resource.

---

A `additional_location` block supports the following:

* `location` - (Required) The name of the Azure Region in which the API Management Service should be expanded to.

* `virtual_network_configuration` - (Optional) A `virtual_network_configuration` block as defined below.  Required when `virtual_network_type` is `External` or `Internal`.

---

A `certificate` block supports the following:

* `encoded_certificate` - (Required) The Base64 Encoded PFX Certificate.

* `certificate_password` - (Required) The password for the certificate.

* `store_name` - (Required) The name of the Certificate Store where this certificate should be stored. Possible values are `CertificateAuthority` and `Root`.


---

A `hostname_configuration` block supports the following:

* `management` - (Optional) One or more `management` blocks as documented below.

* `portal` - (Optional) One or more `portal` blocks as documented below.

* `developer_portal` - (Optional) One or more `developer_portal` blocks as documented below.

* `proxy` - (Optional) One or more `proxy` blocks as documented below.

* `scm` - (Optional) One or more `scm` blocks as documented below.

---

A `identity` block supports the following:

~> **Note:** User Assigned Managed Identities are in Preview

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this API Management Service. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) A list of IDs for User Assigned Managed Identity resources to be assigned.

~> **NOTE:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `management`, `portal`, `developer_portal` and `scm` block supports the following:

* `host_name` - (Required) The Hostname to use for the Management API.

* `key_vault_id` - (Optional) The ID of the Key Vault Secret containing the SSL Certificate, which must be should be of the type `application/x-pkcs12`.

-> **NOTE:** Setting this field requires the `identity` block to be specified, since this identity is used for to retrieve the Key Vault Certificate. Auto-updating the Certificate from the Key Vault requires the Secret version isn't specified.

* `certificate` - (Optional) The Base64 Encoded Certificate.

* `certificate_password` - (Optional) The password associated with the certificate provided above.

-> **NOTE:** Either `key_vault_id` or `certificate` and `certificate_password` must be specified.

* `negotiate_client_certificate` - (Optional) Should Client Certificate Negotiation be enabled for this Hostname? Defaults to `false`.

---

A `policy` block supports the following:

* `xml_content` - (Optional) The XML Content for this Policy.

* `xml_link` - (Optional) A link to an API Management Policy XML Document, which must be publicly available.

---

A `proxy` block supports the following:

* `default_ssl_binding` - (Optional) Is the certificate associated with this Hostname the Default SSL Certificate? This is used when an SNI header isn't specified by a client. Defaults to `false`.

* `host_name` - (Required) The Hostname to use for the Management API.

* `key_vault_id` - (Optional) The ID of the Key Vault Secret containing the SSL Certificate, which must be should be of the type `application/x-pkcs12`.

-> **NOTE:** Setting this field requires the `identity` block to be specified, since this identity is used for to retrieve the Key Vault Certificate. Auto-updating the Certificate from the Key Vault requires the Secret version isn't specified.

* `certificate` - (Optional) The Base64 Encoded Certificate.

* `certificate_password` - (Optional) The password associated with the certificate provided above.

-> **NOTE:** Either `key_vault_id` or `certificate` and `certificate_password` must be specified.

* `negotiate_client_certificate` - (Optional) Should Client Certificate Negotiation be enabled for this Hostname? Defaults to `false`.

---

A `protocols` block supports the following:

* `enable_http2` - (Optional) Should HTTP/2 be supported by the API Management Service? Defaults to `false`.

---

A `security` block supports the following:

* `enable_backend_ssl30` - (Optional) Should SSL 3.0 be enabled on the backend of the gateway? Defaults to `false`.

-> **info:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Ssl30` field

* `enable_backend_tls10` - (Optional) Should TLS 1.0 be enabled on the backend of the gateway? Defaults to `false`.

-> **info:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls10` field

* `enable_backend_tls11` - (Optional) Should TLS 1.1 be enabled on the backend of the gateway? Defaults to `false`.

-> **info:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls11` field

* `enable_frontend_ssl30` - (Optional) Should SSL 3.0 be enabled on the frontend of the gateway? Defaults to `false`.

-> **info:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Ssl30` field

* `enable_frontend_tls10` - (Optional) Should TLS 1.0 be enabled on the frontend of the gateway? Defaults to `false`.

-> **info:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls10` field

* `enable_frontend_tls11` - (Optional) Should TLS 1.1 be enabled on the frontend of the gateway? Defaults to `false`.

-> **info:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls11` field

* `enable_triple_des_ciphers` - (Optional) Should the `TLS_RSA_WITH_3DES_EDE_CBC_SHA` cipher be enabled for alL TLS versions (1.0, 1.1 and 1.2)? Defaults to `false`.

-> **info:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TripleDes168` field

* `disable_backend_ssl30` - (Optional) Should SSL 3.0 be disabled on the backend of the gateway? This property was mistakenly inverted and `true` actually enables it. Defaults to `false`.

-> **Note:** This property has been deprecated in favour of the `enable_backend_ssl30` property and will be removed in version 2.0 of the provider.

* `disable_backend_tls10` - (Optional) Should TLS 1.0 be disabled on the backend of the gateway? This property was mistakenly inverted and `true` actually enables it. Defaults to `false`.

-> **Note:** This property has been deprecated in favour of the `enable_backend_tls10` property and will be removed in version 2.0 of the provider.

* `disable_backend_tls11` - (Optional) Should TLS 1.1 be disabled on the backend of the gateway? This property was mistakenly inverted and `true` actually enables it. Defaults to `false`.

-> **Note:** This property has been deprecated in favour of the `enable_backend_tls11` property and will be removed in version 2.0 of the provider.

* `disable_frontend_ssl30` - (Optional) Should SSL 3.0 be disabled on the frontend of the gateway? This property was mistakenly inverted and `true` actually enables it. Defaults to `false`.

-> **Note:** This property has been deprecated in favour of the `enable_frontend_ssl30` property and will be removed in version 2.0 of the provider.

* `disable_frontend_tls10` - (Optional) Should TLS 1.0 be disabled on the frontend of the gateway? This property was mistakenly inverted and `true` actually enables it. Defaults to `false`.

-> **Note:** This property has been deprecated in favour of the `enable_frontend_tls10` property and will be removed in version 2.0 of the provider.

* `disable_frontend_tls11` - (Optional) Should TLS 1.1 be disabled on the frontend of the gateway? This property was mistakenly inverted and `true` actually enables it. Defaults to `false`.

-> **Note:** This property has been deprecated in favour of the `enable_frontend_tls11` property and will be removed in version 2.0 of the provider.

* `disable_triple_des_ciphers` - (Optional) Should the `TLS_RSA_WITH_3DES_EDE_CBC_SHA` cipher be disabled for alL TLS versions (1.0, 1.1 and 1.2)? This property was mistakenly inverted and `true` actually enables it. Defaults to `false`.

-> **Note:** This property has been deprecated in favour of the `enable_triple_des_ciphers` property and will be removed in version 2.0 of the provider.

---

A `sign_in` block supports the following:

* `enabled` - (Required) Should anonymous users be redirected to the sign in page?

---

A `sign_up` block supports the following:

* `enabled` - (Required) Can users sign up on the development portal?

* `terms_of_service` - (Required) A `terms_of_service` block as defined below.

---

A `virtual_network_configuration` block supports the following:

* `subnet_id` - (Required) The id of the subnet that will be used for the API Management.

---

A `terms_of_service` block supports the following:

* `consent_required` - (Required) Should the user be asked for consent during sign up?

* `enabled` - (Required) Should Terms of Service be displayed during sign up?.

* `text` - (Required) The Terms of Service which users are required to agree to in order to sign up.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API Management Service.

* `additional_location` - Zero or more `additional_location` blocks as documented below.

* `gateway_url` - The URL of the Gateway for the API Management Service.

* `gateway_regional_url` - The Region URL for the Gateway of the API Management Service.

* `identity` - An `identity` block as defined below.

* `management_api_url` - The URL for the Management API associated with this API Management service.

* `portal_url` - The URL for the Publisher Portal associated with this API Management service.

* `developer_portal_url` - The URL for the Developer Portal associated with this API Management service.

* `public_ip_addresses` - The Public IP addresses of the API Management Service.

* `private_ip_addresses` - The Private IP addresses of the API Management Service.

* `scm_url` - The URL for the SCM (Source Code Management) Endpoint associated with this API Management service.

---

An `additional_location` block exports the following:

* `gateway_regional_url` - The URL of the Regional Gateway for the API Management Service in the specified region.

* `public_ip_addresses` - Public Static Load Balanced IP addresses of the API Management service in the additional location. Available only for Basic, Standard and Premium SKU.

* `private_ip_addresses` - The Private IP addresses of the API Management Service.  Available only when the API Manager instance is using Virtual Network mode.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the API Management Service.
* `update` - (Defaults to 3 hours) Used when updating the API Management Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Service.
* `delete` - (Defaults to 3 hours) Used when deleting the API Management Service.

## Import

API Management Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1
```
