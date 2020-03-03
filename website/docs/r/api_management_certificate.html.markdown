---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_certificate"
description: |-
  Manages an Certificate within an API Management Service.
---

# azurerm_api_management_certificate

Manages an Certificate within an API Management Service.

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

resource "azurerm_api_management_certificate" "example" {
  name                = "example-cert"
  api_management_name = azurerm_api_management.example.name
  resource_group_name = azurerm_resource_group.example.name
  data                = filebase64("example.pfx")
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the API Management Certificate. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The Name of the API Management Service where this Service should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The Name of the Resource Group where the API Management Service exists. Changing this forces a new resource to be created.

* `data` - (Required) The base-64 encoded certificate data, which must be a PFX file. Changing this forces a new resource to be created.

* `password` - (Optional) The password used for this certificate. Changing this forces a new resource to be created.

---

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API Management Certificate.

* `expiration` - The Expiration Date of this Certificate, formatted as an RFC3339 string.

* `subject` - The Subject of this Certificate.

* `thumbprint` - The Thumbprint of this Certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Certificate.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Certificate.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Certificate.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Certificate.

## Import

API Management Certificates can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_certificate.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/certificates/certificate1
```
