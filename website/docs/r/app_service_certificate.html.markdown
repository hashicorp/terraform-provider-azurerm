---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_certificate"
description: |-
  Manages an App Service certificate.

---

# azurerm_app_service_certificate

Manages an App Service certificate.

## Example Usage

This example provisions an App Service Certificate from a Local File. Additional examples of how to use the `azurerm_app_service_certificate` resource can be found [in the `./examples/app-service-certificate` directory within the GitHub Repository](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/app-service-certificate).

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

-> **Note:** The resource group must be the same as that which the app service plan is defined in - otherwise the certificate will not show as available for the app services.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `pfx_blob` - (Optional) The base64-encoded contents of the certificate. Changing this forces a new resource to be created.

-> **Note:** Exactly one of `key_vault_secret_id` or `pfx_blob` must be specified.

* `password` - (Optional) The password to access the certificate's private key. Changing this forces a new resource to be created.

* `app_service_plan_id` - (Optional) The ID of the associated App Service plan. Must be specified when the certificate is used inside an App Service Environment hosted App Service or with Basic and Premium App Service plans. Changing this forces a new resource to be created.

* `key_vault_secret_id` - (Optional) The ID of the Key Vault secret. Changing this forces a new resource to be created.

-> **Note:** Exactly one of `key_vault_secret_id` or `pfx_blob` must be specified.

-> **Note:** If using `key_vault_secret_id`, the WebApp Service Resource Principal ID `abfa0a7c-a6b6-4736-8310-5855508787cd` must have 'Secret -> get' and 'Certificate -> get' permissions on the Key Vault containing the certificate. (Source: [App Service Blog](https://azure.github.io/AppService/2016/05/24/Deploying-Azure-Web-App-Certificate-through-Key-Vault.html)) If you use Terraform to create the access policy you have to specify the Object ID of this Principal. This Object ID can be retrieved via following data reference, since it is different in every AAD Tenant:

```hcl
data "azuread_service_principal" "MicrosoftWebApp" {
  application_id = "abfa0a7c-a6b6-4736-8310-5855508787cd"
}
```

* `key_vault_id` - (Optional) The ID of the Key Vault. Must be specified if the Key Vault of `key_vault_secret_id` is in a different subscription from the App Service Certificate. Changing this forces a new resource to be created.

-> **Note:** `key_vault_id` can only be specified if `key_vault_secret_id` has been set.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The App Service certificate ID.

* `friendly_name` - The friendly name of the certificate.

* `subject_name` - The subject name of the certificate.

* `host_names` - List of host names the certificate applies to.

* `issuer` - The name of the certificate issuer.

* `issue_date` - The issue date for the certificate.

* `expiration_date` - The expiration date for the certificate.

* `thumbprint` - The thumbprint for the certificate.

* `hosting_environment_profile_id` - The ID of the App Service Environment where the certificate is in use.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Service Certificate.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Certificate.
* `update` - (Defaults to 30 minutes) Used when updating the App Service Certificate.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Service Certificate.

## Import

App Service Certificates can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_certificate.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/certificates/certificate1
```
