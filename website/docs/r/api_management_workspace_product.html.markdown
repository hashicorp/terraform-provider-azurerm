---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_workspace_product"
description: |-
  Manages an API Management Workspace Product.
---

# azurerm_api_management_workspace_product

Manages an API Management Workspace Product.

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

  sku_name = "Premium_1"
}

resource "azurerm_api_management_workspace" "example" {
  name              = "example-workspace"
  api_management_id = azurerm_api_management.example.id
  display_name      = "Example Workspace"
}

resource "azurerm_api_management_workspace_product" "example" {
  name                        = "example-product"
  api_management_workspace_id = azurerm_api_management_workspace.example.id
  display_name                = "Example Product"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the API Management Workspace Product. Changing this forces a new resource to be created.

* `api_management_workspace_id` - (Required) Specifies the ID of the API Management Workspace. Changing this forces a new resource to be created.

* `display_name` - (Required) Specifies the display name of the API Management Workspace Product.

* `description` - (Optional) Specifies a description for the API Management Workspace Product.

* `published_enabled` - (Optional) Specifies whether the API Management Workspace Product is published. Defaults to `false`.

* `require_approval_enabled` - (Optional) Specifies whether new subscriptions require administrator approval before users can access the APIs of the API Management Workspace Product. Defaults to `false`.

-> **Note:** This can only be specified when `require_subscription_enabled` is `true`.

* `require_subscription_enabled` - (Optional) Specifies whether a subscription is required to access the APIs included in the API Management Workspace Product. Defaults to `true`.

* `subscriptions_limit` - (Optional) Specifies the maximum number of subscriptions a user can have for the API Management Workspace Product.

-> **Note:** This can only be specified when `require_subscription_enabled` is `true`.

* `terms` - (Optional) Specifies the terms of the API Management Workspace Product.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Workspace Product.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Workspace Product.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Workspace Product.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Workspace Product.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Workspace Product.

## Import

API Management Workspace Products can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_workspace_product.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/workspaces/workspace1/products/product1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ApiManagement` - 2024-05-01
