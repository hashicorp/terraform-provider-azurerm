---
layout: "azurerm"
page_title: "Resource: azurerm_api_center_api_definition"
subcategory: "ApiCenter"
description: |-
  Documentation for azurerm_api_center_api_definition.
---

# azurerm_api_center_api_definition

Documentation for azurerm_api_center_api_definition.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_api_center_service" "example" {
  name                = "example-service"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_api_center_api" "example" {
  name       = "example-api"
  service_id = azurerm_api_center_service.example.id
  kind       = "rest"
}

resource "azurerm_api_center_api_version" "example" {
  name            = "example-version"
  api_id          = azurerm_api_center_api.example.id
  lifecycle_stage = "development"
}

resource "azurerm_api_center_api_definition" "example" {
  name           = "example-definition"
  api_version_id = azurerm_api_center_api_version.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the API definition. Changing this value forces a new resource to be created.

* `api_version_id` - (Required) The ID of the API version to which this API definition belongs. Changing this value forces a new resource to be created.

* `description` - (Optional) The description of the API definition.

* `display_name` - (Optional) The display name of the API definition. String length must be between `1` and `50` characters inclusive.

* `specification` - (Optional) (see [nested schema](#schema-6e6573746564-696e707574-73706563696669636174696f6e))

---

<a id="schema-6e6573746564-696e707574-73706563696669636174696f6e"></a>
The `specification` block supports the following:

* `name` - (Required) The name of the Specification.

* `content` - (Optional) The content of the Specification.

* `url` - (Optional) The URL containing the Specification content.

* `version` - (Optional) The version of the Specification.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the resource.
* `update` - (Defaults to 30 minutes) Used when updating the resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the resource.

## Import

API Center API Definitions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_center_api_definition.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.ApiCenter/services/example-service/workspaces/default/apis/example-api/versions/example-version/definitions/example-definition
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ApiCenter` - 2024-03-01
