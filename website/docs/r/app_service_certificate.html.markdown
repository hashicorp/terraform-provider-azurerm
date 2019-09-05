---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_certificate"
sidebar_current: "docs-azurerm-resource-app-service-certificate"
description: |-
  Manages an App Service certificate.

---

# azurerm_app_service_certificate

Manages an App Service certificate.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_app_service_certificate" "example" {
  name                = "example-cert"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  pfx_blob            = "${filebase64("certificate.pfx")}"
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
terraform import azurerm_app_certificate.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/certificates/certificate1
```
