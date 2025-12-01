---
subcategory: "Desktop Virtualization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_desktop_application"
description: |-
  Manages a Virtual Desktop Application.
---

# azurerm_virtual_desktop_application

Manages a Virtual Desktop Application.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "rg-example-virtualdesktop"
  location = "West Europe"
}

resource "azurerm_virtual_desktop_host_pool" "pooledbreadthfirst" {
  name                = "pooledbreadthfirst"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  type               = "Pooled"
  load_balancer_type = "BreadthFirst"
}

resource "azurerm_virtual_desktop_host_pool" "personalautomatic" {
  name                = "personalautomatic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  type                             = "Personal"
  personal_desktop_assignment_type = "Automatic"
  load_balancer_type               = "BreadthFirst"
}

resource "azurerm_virtual_desktop_application_group" "remoteapp" {
  name                = "acctag"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  type          = "RemoteApp"
  host_pool_id  = azurerm_virtual_desktop_host_pool.pooledbreadthfirst.id
  friendly_name = "TestAppGroup"
  description   = "Acceptance Test: An application group"
}

resource "azurerm_virtual_desktop_application" "chrome" {
  name                         = "googlechrome"
  application_group_id         = azurerm_virtual_desktop_application_group.remoteapp.id
  friendly_name                = "Google Chrome"
  description                  = "Chromium based web browser"
  path                         = "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
  command_line_argument_policy = "DoNotAllow"
  command_line_arguments       = "--incognito"
  show_in_portal               = false
  icon_path                    = "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
  icon_index                   = 0
}
```

With MSIX application type, associated to a MSIX package.

```hcl
resource "azurerm_virtual_desktop_msix_package" "example" {
  name                = "example-msix-package"
  host_pool_name      = azurerm_virtual_desktop_host_pool.pooledbreadthfirst.name
  resource_group_name = azurerm_resource_group.example.name
  image_path          = "\\\\path\\to\\image.vhd"
  last_updated_in_utc = "2021-09-01T00:00:00"

  package_application {
    app_id            = "my-app-1"
    app_user_model_id = "user-model-id"
    description       = "Description of my app 1"
    friendly_name     = "App1"
    icon_image_name   = "icon.ico"
    raw_png           = "VGhpcyBpcyBhIHN0cmluZyB0byBoYXNo"
    raw_icon          = "VGhpcyBpcyBhIHN0cmluZyB0byBoYXNo"
  }

  package_family_name   = "example-package-family-name"
  package_name          = "example-package-name"
  package_relative_path = "relative\\path"
  version               = "0.0.0.1"
}

resource "azurerm_virtual_desktop_application" "example" {
  name                         = "example-msix-app"
  application_group_id         = azurerm_virtual_desktop_application_group.remoteapp.id
  command_line_argument_policy = "DoNotAllow"
  msix_package_app_id          = azurerm_virtual_desktop_msix_package.example.package_application[0].app_id
  msix_package_family_name     = azurerm_virtual_desktop_msix_package.example.package_family_name
  application_type             = "MsixApplication"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Virtual Desktop Application. Changing the name forces a new resource to be created.

* `application_group_id` - (Required) Resource ID for a Virtual Desktop Application Group to associate with the Virtual Desktop Application. Changing this forces a new resource to be created.

* `friendly_name` - (Optional) Option to set a friendly name for the Virtual Desktop Application.

* `description` - (Optional) Option to set a description for the Virtual Desktop Application.

* `path` - (Optional) The file path location of the app on the Virtual Desktop OS. Required for `InBuilt` application type, but should not be specified for `MsixApplication`.

* `command_line_argument_policy` - (Required) Specifies whether this published application can be launched with command line arguments provided by the client, command line arguments specified at publish time, or no command line arguments at all. Possible values include: `DoNotAllow`, `Allow`, `Require`.

* `command_line_arguments` - (Optional) Command Line Arguments for Virtual Desktop Application.

* `show_in_portal` - (Optional) Specifies whether to show the RemoteApp program in the RD Web Access server.

* `icon_path` - (Optional) Specifies the path for an icon which will be used for this Virtual Desktop Application.

* `icon_index` - (Optional) The index of the icon you wish to use.

* `application_type` - (Optional) The remote application type. Possible values include: `InBuilt`, `MsixApplication`. Defaults to `InBuilt`.

* `msix_package_app_id` - (Optional) The package application id for MSIX applications. This corresponds to the `app_id` argument in `package_application` block in `azurerm_virtual_desktop_msix_package`.

* `msix_package_family_name` - (Optional) The package family name for MSIX applications. This corresponds to the `package_family_name` argument in `azurerm_virtual_desktop_msix_package`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Desktop Application.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Virtual Desktop Application.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Desktop Application.
* `update` - (Defaults to 1 hour) Used when updating the Virtual Desktop Application.
* `delete` - (Defaults to 1 hour) Used when deleting the Virtual Desktop Application.

## Import

Virtual Desktop Application can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_desktop_application.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/Microsoft.DesktopVirtualization/applicationGroups/myapplicationgroup/applications/myapplication
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DesktopVirtualization` - 2024-04-03
