---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_workspace"
description: |-
  Manages an API Management Workspace.
---

# azurerm_api_management_workspace

Manages an API Management Workspace.

~> **Note:** This resource is currently available only when using the Classic Premium SKU of `azurerm_api_management`. For more details, refer to [Federated API Management with Workspaces](https://learn.microsoft.com/en-us/azure/api-management/workspaces-overview).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_api_management" "example" {
  name                = "example-apimanagement"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Premium_1"
}

resource "azurerm_api_management_workspace" "example" {
  name              = "example-workspace"
  api_management_id = azurerm_api_management.example.id
  display_name      = "my workspace"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this API Management Workspace. Changing this forces a new resource to be created.

* `api_management_id` - (Required) Specifies the ID of the API Management Service in which the API Management Workspace should be created. Changing this forces a new resource to be created.

* `display_name` - (Required) The display name of the API Management Workspace.

* `description` - (Optional) The description of the API Management Workspace.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Workspace.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Workspace.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Workspace.

## Import

API Management Workspace can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_workspace.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/workspaces/workspace1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ApiManagement`: 2024-05-01
