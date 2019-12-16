---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_certificate"
sidebar_current: "docs-azurerm-resource-app-service-certificate"
description: |-
  Manages an App Service certificate.

---

# azurerm_app_service_certificate

Manages an App Service certificate.

## Example Usage

This example provisions an App Service Certificate from a Local File. Additional examples of how to use the `azurerm_app_service_certificate` resource can be found [in the ./examples/app-service/certificate` directory within the Github Repository](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/app-service/certificate).


```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_app_service_certificate" "example" {
  name                = "example-cert"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  pfx_blob            = filebase64("certificate.pfx")
  password            = "terraform"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the certificate. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the certificate. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `pfx_blob` - (Optional) The base64-encoded contents of the certificate. Changing this forces a new resource to be created.

-> **NOTE:** Either `pfx_blob` or `key_vault_secret_id` must be set - but not both.

* `password` - (Optional) The password to access the certificate's private key. Changing this forces a new resource to be created.

* `key_vault_secret_id` - (Optional) The ID of the Key Vault secret. Changing this forces a new resource to be created.

-> **NOTE:** If using `key_vault_secret_id`, the WebApp Service Resource Principal ID `abfa0a7c-a6b6-4736-8310-5855508787cd` with object_id `f8daea97-62e7-4026-becf-13c2ea98e8b4` of  must have 'Secret -> get' and 'Certificate -> get' permissions on the Key Vault containing the certificate.  (Source: [App Service Blog](https://azure.github.io/AppService/2016/05/24/Deploying-Azure-Web-App-Certificate-through-Key-Vault.html))

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

## Import

App Service certificates can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_certificate.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/certificates/certificate1
```
