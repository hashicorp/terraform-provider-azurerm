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
  publisher_name      = "My Company"
  publisher_email     = "company@terraform.io"
  sku_name            = "Premium_1"
}

resource "azurerm_api_management_workspace" "example" {
  name              = "example-workspace"
  api_management_id = azurerm_api_management.example.id
  display_name      = "ExampleWorkspace"
}

resource "azurerm_api_management_workspace_named_value" "example" {
  name                        = "example-named-value"
  api_management_workspace_id = azurerm_api_management_workspace.example.id
  display_name                = "ExampleProperty"
  value                       = "Example Value"
}
```

## Arguments Reference

The following arguments are supported:

* `api_management_workspace_id` - (Required) The ID of the API Management Workspace. Changing this forces a new resource to be created.

* `display_name` - (Required) The display name of this API Management Workspace Named Value.

* `name` - (Required) The name of the API Management Workspace Named Value. Changing this forces a new resource to be created.

* `secret` - (Optional) Specifies whether the API Management Workspace Named Value is secret. Valid values are `true` or `false`. The default value is `false`.

~> **Note:** Setting the field `secret` to `true` does not make this field sensitive in Terraform, instead it marks the value as secret and encrypts the value in Azure.

* `tags` - (Optional) A list of tags to be applied to the API Management Workspace Named Value.

* `value` - (Optional) The value of this API Management Workspace Named Value.

~> **Note:** Exactly one of `value` or `value_from_key_vault` must be specified.

* `value_from_key_vault` - (Optional) A `value_from_key_vault` block as defined below.

~> **Note:** Exactly one of `value` or `value_from_key_vault` must be specified. If `value_from_key_vault` is specified, `secret` must also be set to `true`.

---

A `value_from_key_vault` block supports the following:

* `identity_client_id` - (Optional) The client ID of User Assigned Identity, for the API Management Service, which will be used to access the key vault secret. The System Assigned Identity will be used in absence.

* `secret_id` - (Required) The resource ID of the Key Vault Secret.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Workspace Named Value.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Workspace Named Value.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Workspace Named Value.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Workspace Named Value.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Workspace Named Value.

## Import

API Management Workspace Named Values can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_workspace_named_value.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/workspaces/workspace1/namedValues/namedValue1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ApiManagement` - 2024-05-01
