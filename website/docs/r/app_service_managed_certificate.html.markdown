---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_managed_certificate"
description: |-
  Manages an App Service Managed Certificate.
---

# azurerm_app_service_managed_certificate

This certificate can be used to secure custom domains on App Services (Windows and Linux) hosted on an App Service Plan of Basic and above (free and shared tiers are not supported).

~> **Note:** A certificate is valid for six months, and about a month before the certificate’s expiration date, App Services renews/rotates the certificate. This is managed by Azure and doesn't require this resource to be changed or reprovisioned. It will change the `thumbprint` computed attribute the next time the resource is refreshed after rotation occurs, so keep that in mind if you have any dependencies on this attribute directly.

## Example Usage

```hcl
data "azurerm_dns_zone" "example" {
  name                = "mydomain.com"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_app_service_plan" "example" {
  name                = "example-plan"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Basic"
    size = "B1"
  }
}

resource "azurerm_app_service" "example" {
  name                = "example-app"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  app_service_plan_id = azurerm_app_service_plan.example.id
}

resource "azurerm_dns_txt_record" "example" {
  name                = "asuid.mycustomhost.contoso.com"
  zone_name           = data.azurerm_dns_zone.example.name
  resource_group_name = data.azurerm_dns_zone.example.resource_group_name
  ttl                 = 300

  record {
    value = azurerm_app_service.example.custom_domain_verification_id
  }
}

resource "azurerm_dns_cname_record" "example" {
  name                = "example-adcr"
  zone_name           = data.azurerm_dns_zone.example.name
  resource_group_name = data.azurerm_dns_zone.example.resource_group_name
  ttl                 = 300
  record              = azurerm_app_service.example.default_site_hostname
}

resource "azurerm_app_service_custom_hostname_binding" "example" {
  hostname            = join(".", [azurerm_dns_cname_record.example.name, azurerm_dns_cname_record.example.zone_name])
  app_service_name    = azurerm_app_service.example.name
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_app_service_managed_certificate" "example" {
  custom_hostname_binding_id = azurerm_app_service_custom_hostname_binding.example.id
}

resource "azurerm_app_service_certificate_binding" "example" {
  hostname_binding_id = azurerm_app_service_custom_hostname_binding.example.id
  certificate_id      = azurerm_app_service_managed_certificate.example.id
  ssl_state           = "SniEnabled"
}
```

## Arguments Reference

The following arguments are supported:

* `custom_hostname_binding_id` - (Required) The ID of the App Service Custom Hostname Binding for the Certificate. Changing this forces a new App Service Managed Certificate to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the App Service Managed Certificate.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the App Service Managed Certificate.

* `canonical_name` - The Canonical Name of the Certificate.

* `expiration_date` - The expiration date of the Certificate.

* `friendly_name` - The friendly name of the Certificate.

* `host_names` - The list of Host Names for the Certificate.

* `issue_date` - The Start date for the Certificate.

* `issuer` - The issuer of the Certificate.

* `subject_name` - The Subject Name for the Certificate.

* `thumbprint` - The Certificate Thumbprint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Service Managed Certificate.

* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Managed Certificate.

* `update` - (Defaults to 30 minutes) Used when updating the App Service Managed Certificate.

* `delete` - (Defaults to 30 minutes) Used when deleting the App Service Managed Certificate.

## Import

App Service Managed Certificates can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_managed_certificate.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Web/certificates/customhost.contoso.com
```
