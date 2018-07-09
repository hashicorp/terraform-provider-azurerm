---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management"
sidebar_current: "docs-azurerm-resource-api-management-x"
description: |-
  Create a Api Management service.
---

# azurerm_api_management

Create a Api Management component.

## Example Usage (Developer)

```hcl
resource "azurerm_resource_group" "test" {
  name     = "api-rg-dev"
  location = "West Europe"
}

resource "azurerm_api_management" "test" {
  name                = "api-mngmnt-dev"
  publisher_name      = "My Company"
  publisher_email     = "company@terraform.io"
  sku {
    name = "Developer"
  }
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
```

## Example Usage (Complete)

```hcl
resource "azurerm_resource_group" "west" {
  name     = "api-rg-premium-west"
  location = "West Europe"
}

resource "azurerm_resource_group" "north" {
  name     = "api-rg-premium-north"
  location = "North Europe"
}

resource "azurerm_api_management" "test" {
  name                          = "api-mngmnt"
  publisher_name                = "My Company"
  publisher_email               = "company1@terraform.io"
  notification_sender_email     = "api@terraform.io"

  additional_location {
	location = "${azurerm_resource_group.north.location}"
	sku {
		name = "Premium"
	}
  }

  certificate {
	certificate            = "${base64encode(file("testdata/api_management_api_test.pfx"))}"
	certificate_password   = "terraform"
	store_name             = "CertificateAuthority"
  }

  certificate {
	certificate            = "${base64encode(file("testdata/api_management_api_test.pfx"))}"
	certificate_password   = "terraform"
	store_name             = "Root"
  }

  custom_properties {
	Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TripleDes168 = "true"
  }

  hostname_configuration {
	type                         = "Proxy"
	host_name                    = "api.terraform.io"
	certificate                  = "${base64encode(file("testdata/api_management_api_test.pfx"))}"
	certificate_password         = "terraform"
	default_ssl_binding          = true
	negotiate_client_certificate = false
  }

  hostname_configuration {
	type                         = "Proxy"
	host_name                    = "api2.terraform.io"
	certificate                  = "${base64encode(file("testdata/api_management_api2_test.pfx"))}"
	certificate_password         = "terraform"
	negotiate_client_certificate = true
  }

  hostname_configuration {
	type                         = "Portal"
	host_name                    = "portal.terraform.io"
	certificate                  = "${base64encode(file("testdata/api_management_portal_test.pfx"))}"
	certificate_password         = "terraform"
  }

  sku {
    name = "Premium"
  }

  tags {
	test = "true"
  }

  location            = "${azurerm_resource_group.west.location}"
  resource_group_name = "${azurerm_resource_group.west.name}"
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the API Management service.

* `location` - (Required) The Azure location where the API Management Service exists.

* `resource_group_name` - (Required) The Name of the Resource Group where the API Management service exists.

* `publisher_name` - (Required) The name of publisher/company.

* `publisher_email` - (Required) The email of publisher/company.

* `sku` - (Required) A `sku` block as documented below.

* `notification_sender_email` - (Optional) Email address from which the notification will be sent.

* `additional_location` - (Optional) Additional datacenter locations of the API Management service. The `additional_location` block is documented below.

* `certificate` - (Optional) List of Certificates that is installed in the API Management service. Max supported certificates that can be installed is 10. The `certificate` block is documented below.

* `custom_properties` - (Optional) Custom properties of the API Management service. The property `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TripleDes168` means the cipher TLS_RSA_WITH_3DES_EDE_CBC_SHA is disabled for all TLS(1.0, 1.1 and 1.2). The property `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls11` means just TLS 1.1 is disabled and the property `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls10` means TLS 1.0 is disabled on an API Management service.

* `hostname_configuration` - (Optional) Custom hostname configuration of the API Management service. The `hostname_configuration` block is documented below.

* `tags` - (Optional) A mapping of tags assigned to the resource.

---

`sku` block supports the following:

* `name` - (Required) Specifies the plan's pricing tier.

* `capacity` - (Optional) Specifies the number of units associated with this API Management service.

`additional_location` block supports the following:

* `location` - (Required) The location name of the additional region among Azure Data center regions.

* `sku` - (Required) SKU properties of the API Management service. The `hostname_configuration` block is documented above.

`certificate` block supports the following:

* `encoded_certificate` - (Required) Base64 Encoded PFX certificate.

* `certificate_password` - (Required) Certificate password.

* `store_name` - (Required) The local certificate store location. Only Root and CertificateAuthority are valid locations. Possible values include: `CertificateAuthority`, `Root`.

`hostname_configuration` block supports the following:

* `type` - (Required) Hostname type. Possible values include: `Proxy`, `Portal`, `Management` or `Scm`

* `host_name` - (Required) Hostname to configure on the Api Management service.

* `certificate` - (Required) Base64 Encoded certificate.

* `certificate_password` - (Required) Certificate Password.

* `default_ssl_binding` - (Optional) If set to true the certificate associated with this Hostname is setup as the Default SSL Certificate. If a client does not send the SNI header, then this will be the certificate that will be challenged. The property is useful if a service has multiple custom hostname enabled and it needs to decide on the default ssl certificate. The setting only applied to Proxy Hostname Type.

* `negotiate_client_certificate` - (Optional) If set to true will always negotiate client certificate on the hostname. Default Value is false.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the App Service Plan component.

* `created` - Creation date of the API Management service.

* `gateway_url` - Gateway URL of the API Management service.

* `gateway_regional_url` - Gateway URL of the API Management service in the Default Region.

* `portal_url` - Publisher portal endpoint Url of the API Management service.

* `management_api_url` - Management API endpoint URL of the API Management service.

* `scm_url` - SCM endpoint URL of the API Management service.

* `additional_location` - Additional datacenter locations of the API Management service. The `additional_location` block is documented below.

---

`additional_location` block exports the following:

* `gateway_regional_url` - Gateway URL of the API Management service in the Region.

* `static_ips` - Static IP addresses of the location's virtual machines.

## Import

Api Management services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1
```
