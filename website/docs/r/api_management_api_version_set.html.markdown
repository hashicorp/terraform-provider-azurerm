---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_api_version_set"
description: |-
  Manages an API Version Set within an API Management Service.
---

# azurerm_api_management_api_version_set

Manages an API Version Set within an API Management Service.

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
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_api_version_set" "example" {
  name                = "example-apimapi-1.0.0"
  resource_group_name = azurerm_resource_group.example.name
  api_management_name = azurerm_api_management.example.name
  display_name        = "ExampleAPIVersionSet"
  versioning_scheme   = "Segment"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the API Version Set. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The name of the [API Management Service](api_management.html) in which the API Version Set should exist. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the parent API Management Service exists. Changing this forces a new resource to be created.

* `display_name` - (Required) The display name of this API Version Set.

* `versioning_scheme` - (Required) Specifies where in an Inbound HTTP Request that the API Version should be read from. Possible values are `Header`, `Query` and `Segment`.

---

* `description` - (Optional) The description of API Version Set.

* `version_header_name` - (Optional) The name of the Header which should be read from Inbound Requests which defines the API Version.

-> **NOTE:** This must be specified when `versioning_scheme` is set to `Header`.

* `version_query_name` - (Optional) The name of the Query String which should be read from Inbound Requests which defines the API Version.

-> **NOTE:** This must be specified when `versioning_scheme` is set to `Query`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API Version Set.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management API Version Set.
* `update` - (Defaults to 30 minutes) Used when updating the API Management API Version Set.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management API Version Set.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management API Version Set.

## Import

API Version Set can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_api_version_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/apiVersionSets/set1
```
