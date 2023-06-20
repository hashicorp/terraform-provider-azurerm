---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_public_certificate"
description: |-
 Manages an App Service Public Certificate.
---

# azurerm_app_service_public_certificate

Manages an App Service Public Certificate.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_app_service_plan" "example" {
  name                = "example-app-service-plan"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "example" {
  name                = "example-app-service"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  app_service_plan_id = azurerm_app_service_plan.example.id
}

resource "azurerm_app_service_public_certificate" "example" {
  resource_group_name  = azurerm_resource_group.example.name
  app_service_name     = azurerm_app_service.example.name
  certificate_name     = "example-public-certificate"
  certificate_location = "Unknown"
  blob                 = filebase64("app_service_public_certificate.cer")
}
```

## Arguments Reference

The following arguments are supported:

* `app_service_name` - (Required) The name of the App Service. Changing this forces a new App Service Public Certificate to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the App Service Public Certificate should exist. Changing this forces a new App Service Public Certificate to be created.

* `certificate_name` - (Required) The name of the public certificate. Changing this forces a new App Service Public Certificate to be created.

* `certificate_location` - (Required) The location of the certificate. Possible values are `CurrentUserMy`, `LocalMachineMy` and `Unknown`. Changing this forces a new App Service Public Certificate to be created.

* `blob` - (Required) The base64-encoded contents of the certificate. Changing this forces a new App Service Public Certificate to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the App Service Public Certificate.

* `thumbprint` - The thumbprint of the public certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Service Public Certificate.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Public Certificate.
* `update` - (Defaults to 30 minutes) Used when updating the App Service Public Certificate.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Service Public Certificate.

## Import

App Service Public Certificates can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_public_certificate.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Web/sites/site1/publicCertificates/publicCertificate1
```
