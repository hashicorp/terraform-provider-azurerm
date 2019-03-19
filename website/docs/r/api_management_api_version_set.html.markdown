---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_version_set"
sidebar_current: "docs-azurerm-resource-api-management-api-version-set"
description: |-
  Manages an API Management API Version Set.
---

# azurerm_api_management_version_set

Manages an API Management API Version Set.


## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West US"
}

resource "azurerm_api_management" "example" {
  name                = "example-apim"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku {
    name     = "Developer"
    capacity = 1
  }
}

resource "azurerm_api_management_version_set" "example" {
  name                = "example-apimapivs"
  resource_group_name = "${azurerm_resource_group.example.name}"
  api_management_name = "${azurerm_api_management.example.name}"
  description         = "ExampleAPIVersionSetDescription"
  display_name        = "ExampleAPIVersionSet"
  versioning_schema   = "Segment"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the API Version Set. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the API Version Set should exist. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The name of the [API Management Service](api_management.html) in which the API Version Set should exist. Changing this forces a new resource to be created.

* `description` - (Required) The description of API Version Set.

* `display_name` - (Required) The display name of this API Version Set.

* `versioning_schema` - (Required) A value that determines where the API Version identifier will be located in a HTTP request. Allowed values include: `Segment`, `Header`, `Query`.

* `version_header_name` - (Optional) Name of HTTP header parameter that indicates the API Version if `versioning_schema` is set to `Header`.

* `version_query_name` - (Optional) Name of query parameter that indicates the API Version if `versioning_schema` is set to `Query`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API Version Set.

## Import

API Version Set can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_version_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.ApiManagement/service/example-apim/api-version-sets/example-apimp
```
