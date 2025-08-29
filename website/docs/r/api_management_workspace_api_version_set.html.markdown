---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_workspace_api_version_set"
description: |-
  Manages an API Version Set within an API Management Workspace.
---

# azurerm_api_management_workspace_api_version_set

Manages an API Version Set within an API Management Workspace.

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
  publisher_name      = "Example Publisher"
  publisher_email     = "publisher@example.com"

  sku_name = "Premium_1"
}

resource "azurerm_api_management_workspace" "example" {
  name              = "example-workspace"
  api_management_id = azurerm_api_management.example.id
  display_name      = "Example Workspace"
  description       = "Example workspace for development"
}

resource "azurerm_api_management_workspace_api_version_set" "example" {
  name                        = "example-version-set"
  api_management_workspace_id = azurerm_api_management_workspace.example.id
  display_name                = "Example API Version Set"
  versioning_scheme           = "Segment"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the API Management Workspace API Version Set. Changing this forces a new resource to be created.

* `api_management_workspace_id` - (Required) Specifies the ID of the API Management Workspace. Changing this forces a new resource to be created.

* `display_name` - (Required) Specifies the display name of the API Management Workspace API Version Set.

* `versioning_scheme` - (Required) Specifies where in a request that the API Version should be read from. Possible values are `Header`, `Query` and `Segment`.

* `description` - (Optional) Specifies the description of the API Management Workspace API Version Set.

* `version_header_name` - (Optional) Specifies the name of the header to read from inbound requests to determine the API version.

* `version_query_name` - (Optional) Specifies the name of the query string parameter to read from inbound requests to determine the API version.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Workspace API Version Set.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Workspace API Version Set.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Workspace API Version Set.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Workspace API Version Set.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Workspace API Version Set.

## Import

API Management Workspace API Version Sets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_workspace_api_version_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/workspaces/workspace1/apiVersionSets/set1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ApiManagement` - 2024-05-01
