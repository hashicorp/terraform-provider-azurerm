---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_workspace_named_value"
description: |-
  Manages an API Management Workspace Named Value.
---

# azurerm_api_management_workspace_named_value

Manages an API Management Workspace Named Value.

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

  sku_name = "Premium_1"
}

resource "azurerm_api_management_workspace" "example" {
  name              = "example-workspace"
  api_management_id = azurerm_api_management.example.id
  display_name      = "Example Workspace"
  description       = "Example workspace description"
}

resource "azurerm_api_management_workspace_named_value" "example" {
  name                        = "example-property"
  api_management_workspace_id = azurerm_api_management_workspace.example.id
  display_name                = "ExampleProperty"
  value                       = "Example Value"
}

```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the API Management Workspace Named Value. Changing this forces a new resource to be created.

* `api_management_workspace_id` - (Required) Specifies the ID of the API Management Workspace. Changing this forces a new resource to be created.

* `display_name` - (Required) Specifies the display name of the API Management Workspace Named Value.

* `value` - (Optional) Specifies the value of the API Management Workspace Named Value.

* `value_from_key_vault` - (Optional) A `value_from_key_vault` block as defined below. 

~> **Note:** One and only one of `value` or `value_from_key_vault` must be specified.

* `secret_enabled` - (Optional) Specifies whether encryption is enabled for the API Management Named Value. Defaults to `false`.

~> **Note:** When `value_from_key_vault` is specified, `secret_enabled` must be set to `true`.

* `tags` - (Optional) Specifies a list of tags to be applied to the API Management Workspace Named Value.

---

A `value_from_key_vault` block supports the following:

* `key_vault_secret_id` - (Required) Specifies the ID of the key vault secret.

* `user_assigned_identity_client_id` - (Optional) Specifies the client ID of user-assigned identity to be used for accessing the `key_vault_secret_id`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Workspace Named Value.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Workspace Named Value.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Workspace Named Value.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Workspace Named Value.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Workspace Named Value.

## Import

API Management Workspace Named Values can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_workspace_named_value.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/workspaces/workspace1/namedValues/value1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ApiManagement` - 2024-05-01
