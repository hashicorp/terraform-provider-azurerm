---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_static_web_app_build"
description: |-
  Manages the Environment Variables of a Static Web App Build.
---

# azurerm_static_web_app_build

Manages the Environment Variables of a Static Web App Build, such as a preview environment.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_static_web_app" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_static_web_app_build" "example" {
  name                = azurerm_static_web_app.example.name
  resource_group_name = azurerm_resource_group.example.name
  build               = "preview"

  environment_variables = {
    "X_API_KEY" = "example"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Static Web App. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Static Web App exists. Changing this forces a new resource to be created.

* `build` - (Required) The name of the Static Web App Build, such as the name of a preview environment.

* `environment_variables` - (Required) A map of Environment Variables to set on the Static Web App Build.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Static Web App Build.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Static Web App Build.
* `read` - (Defaults to 5 minutes) Used when retrieving the Static Web App Build.
* `update` - (Defaults to 30 minutes) Used when updating the Static Web App Build.
* `delete` - (Defaults to 30 minutes) Used when deleting the Static Web App Build.

## Import

Static Web App Builds can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_static_web_app_build.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Web/staticSites/my-static-site1/builds/preview
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Web` - 2023-01-01
