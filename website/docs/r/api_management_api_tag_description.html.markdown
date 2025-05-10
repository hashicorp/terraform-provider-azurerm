---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_api_tag_description"
description: |-
  Manages an API Tag Description within an API Management Service.
---

# azurerm_api_management_api_tag_description

Manages an API Tag Description within an API Management Service.

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

  sku_name = "Developer_1"
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
    content_value  = "https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm/refs/heads/main/internal/services/apimanagement/testdata/api_management_api_swagger.json"
  }
}

resource "azurerm_api_management_tag" "example" {
  api_management_id = azurerm_api_management.example.id
  name              = "example-Tag"
}

resource "azurerm_api_management_api_tag_description" "example" {
  api_tag_id                = azurerm_api_management_tag.example.id
  description               = "This is an example description"
  external_docs_url         = "https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs"
  external_docs_description = "This is an example external docs description"
}
```

## Argument Reference

The following arguments are supported:

* `api_tag_id` - (Required) The The ID of the API Management API Tag. Changing this forces a new API Management API Tag Description to be created.

* `description` - (Optional) The description of the Tag.

* `external_documentation_url` - (Optional) The URL of external documentation resources describing the tag.

* `external_documentation_description` - (Optional) The description of the external documentation resources describing the tag.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management API Schema.

## Timeouts

The `timeouts` block allows you to
specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management API Schema.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management API Schema.
* `update` - (Defaults to 30 minutes) Used when updating the API Management API Schema.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management API Schema.

## Import

API Management API Schema's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_api_tag_description.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/apis/api1/tagDescriptions/tagDescriptionId1
```
