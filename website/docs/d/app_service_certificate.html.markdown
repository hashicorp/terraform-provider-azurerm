---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_certificate"
description: |-
  Gets information about an existing App Service Certificate.

---

# Data Source: azurerm_app_service_certificate

Use this data source to access information about an App Service Certificate.

## Example Usage

```hcl
data "azurerm_app_service_certificate" "example" {
  name                = "example-app-service-certificate"
  resource_group_name = "example-rg"
}

output "app_service_certificate_id" {
  value = data.azurerm_app_service_certificate.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - Specifies the name of the certificate.

* `resource_group_name` - The name of the resource group in which to create the certificate.

## Attributes Reference

The following attributes are exported:

* `id` - The App Service certificate ID.

* `friendly_name` - The friendly name of the certificate.

* `subject_name` - The subject name of the certificate.

* `host_names` - List of host names the certificate applies to.

* `issuer` - The name of the certificate issuer.

* `issue_date` - The issue date for the certificate.

* `expiration_date` - The expiration date for the certificate.

* `thumbprint` - The thumbprint for the certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Certificate.
