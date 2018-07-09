---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management"
sidebar_current: "docs-azurerm-datasource-azurerm_api_management"
description: |-
  Get information about an API Management service.
---

# Data Source: azurerm_api_management

Use this data source to obtain information about an API Management service.

## Example Usage

```hcl
data "azurerm_api_management" "test" {
  name                = "search-api"
  resource_group_name = "search-service"
}

output "api_management_id" {
  value = "${data.azurerm_api_management.test.id}"
}
```

## Argument Reference

* `name` - (Required) The name of the API Management service.

* `resource_group_name` - (Required) The Name of the Resource Group where the API Management service exists.

## Attributes Reference

* `id` - The ID of the API Management Service.

* `location` - The Azure location where the API Management Service exists.

* `publisher_name` - The name of publisher/company.

* `publisher_email` - The email of publisher/company.

* `sku` - A `sku` block as documented below.

* `notification_sender_email` - Email address from which the notification will be sent.

* `created` - Creation date of the API Management service.

* `gateway_url` - Gateway URL of the API Management service.

* `gateway_regional_url` - Gateway URL of the API Management service in the Default Region.

* `portal_url` - Publisher portal endpoint Url of the API Management service.

* `management_api_url` - Management API endpoint URL of the API Management service.

* `scm_url` - SCM endpoint URL of the API Management service.

* `additional_location` - Additional datacenter locations of the API Management service. The `additional_location` block is documented below.

* `certificate` - List of Certificates that is installed in the API Management service. Max supported certificates that can be installed is 10. The `certificate` block is documented below.

* `custom_properties` - Custom properties of the API Management service. The property `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TripleDes168` means the cipher TLS_RSA_WITH_3DES_EDE_CBC_SHA is disabled for all TLS(1.0, 1.1 and 1.2). The property `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls11` means just TLS 1.1 is disabled and the property `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls10` means TLS 1.0 is disabled on an API Management service.

* `hostname_configuration` - Custom hostname configuration of the API Management service. The `hostname_configuration` block is documented below.

* `tags` - A mapping of tags assigned to the resource.

---

A `sku` block supports the following:

* `name` - Specifies the plan's pricing tier.

* `capacity` - Specifies the number of units associated with this API Management service.


A `additional_location` block supports the following:

* `location` - (Required) The location name of the additional region among Azure Data center regions.

* `sku` - (Required) SKU properties of the API Management service. The `hostname_configuration` block is documented above.

* `gateway_regional_url` - Gateway URL of the API Management service in the Region.

* `static_ips` - Static IP addresses of the location's virtual machines.

A `certificate` block supports the following:

* `store_name` - The local certificate store location. Only Root and CertificateAuthority are valid locations. Possible values include: `CertificateAuthority`, `Root`.

* `certificate_info` - A `certificate_info` block as documented below.

A `certificate_info` block supports the following:

* `expiry` - Expiration date of the certificate.

* `thumbprint` - Thumbprint of the certificate.

* `subject` - Subject of the certificate.

A `hostname_configuration` block supports the following:

* `type` - Hostname type. Possible values include: `Proxy`, `Portal`, `Management` or `Scm`

* `host_name` - Hostname to configure on the Api Management service.

* `certificate` - Base64 Encoded certificate.

* `certificate_password` - Certificate Password.

* `default_ssl_binding` - If set to true the certificate associated with this Hostname is setup as the Default SSL Certificate. If a client does not send the SNI header, then this will be the certificate that will be challenged. The property is useful if a service has multiple custom hostname enabled and it needs to decide on the default ssl certificate. The setting only applied to Proxy Hostname Type.

* `negotiate_client_certificate` - If set to true will always negotiate client certificate on the hostname. Default Value is false.
