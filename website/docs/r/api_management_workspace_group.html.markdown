---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_workspace_group"
description: |-
  Manages an API Management Workspace Group.
---

# azurerm_api_management_workspace_group

Manages an API Management Workspace Group.

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
  sku_name            = "Premium_1"
}

resource "azurerm_api_management_workspace" "example" {
  name              = "example-workspace"
  api_management_id = azurerm_api_management.example.id
  display_name      = "Example Workspace"
}

resource "azurerm_api_management_workspace_group" "example" {
  name                        = "example-group"
  api_management_workspace_id = azurerm_api_management_workspace.example.id
  display_name                = "Example Group"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the API Management Workspace Group. Changing this forces a new resource to be created.

* `api_management_workspace_id` - (Required) Specifies the ID of the API Management Workspace. Changing this forces a new resource to be created.

* `display_name` - (Required) Specifies the display name of the API Management Workspace Group.

---

* `external_id` - (Optional) Specifies the ID of the group from an external identity provider. For example, for Azure Active Directory: `aad://<tenant id>/groups/<group object id>`. Changing this forces a new resource to be created.

* `description` - (Optional) Specifies the description of the API Management Workspace Group.

* `type` - (Optional) Specifies the type of the API Management Workspace Group. Possible values are `custom`, `external`. Defaults to `custom`. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Workspace Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Workspace Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Workspace Group.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Workspace Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Workspace Group.

## Import

API Management Workspace Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_workspace_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/workspaces/workspace1/groups/group1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ApiManagement` - 2024-05-01
