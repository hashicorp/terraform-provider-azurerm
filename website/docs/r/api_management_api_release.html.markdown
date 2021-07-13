---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_api_release"
description: |-
  Manages a API Management API Release.
---

# azurerm_api_management_api_release

Manages a API Management API Release.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

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
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_api" "example" {
  name                = "example-api"
  resource_group_name = azurerm_resource_group.example.name
  api_management_name = azurerm_api_management.example.name
  revision            = "1"
  display_name        = "Example API"
  path                = "example"
  protocols           = ["https"]

  import {
    content_format = "swagger-link-json"
    content_value  = "http://conferenceapi.azurewebsites.net/?format=json"
  }
}

resource "azurerm_api_management_api_release" "example" {
  name   = "example-Api-Release"
  api_id = azurerm_api_management_api.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this API Management API Release. Changing this forces a new API Management API Release to be created.

* `api_id` - (Required) The ID of the API Management API. Changing this forces a new API Management API Release to be created.

---

* `notes` - (Optional) The Release Notes.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the API Management API Release.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management API Release.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management API Release.
* `update` - (Defaults to 30 minutes) Used when updating the API Management API Release.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management API Release.

## Import

API Management API Releases can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_api_release.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/apis/api1/releases/release1
```
