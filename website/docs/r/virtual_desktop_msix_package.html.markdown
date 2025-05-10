---
subcategory: "Desktop Virtualization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_desktop_msix_package"
description: |-
  Manages a Virtual Desktop MSIX package.
---

# azurerm_virtual_desktop_msix_package

Manages a Virtual Desktop MSIX package.

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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The version specific package full name of the MSIX package within the specified host pool. Changing this forces a new Virtual Desktop MSIX package to be created.

* `host_pool_name` - (Required) The name of the host pool within the specified resource group. Changing this forces a new Virtual Desktop MSIX package to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Virtual Desktop MSIX package should exist. Changing this forces a new Virtual Desktop MSIX package to be created.

* `image_path` - (Required) VHD/CIM image path on the Network Share. Changing this forces a new Virtual Desktop MSIX package to be created.

* `last_updated_in_utc` - (Required) Date Package was last updated, found in the appxmanifest.xml. RFC3339 format without the time offset, for example: "2021-09-01T00:00:00". Changing this forces a new Virtual Desktop MSIX package to be created.

* `package_application` - (Required) One or more `package_application` blocks as defined below. Changing this forces a new Virtual Desktop MSIX package to be created.

* `package_family_name` - (Required) Package Family Name from appxmanifest.xml. Contains Package Name and Publisher name. Changing this forces a new Virtual Desktop MSIX package to be created.

* `package_name` - (Required) Package Name from appxmanifest.xml. Changing this forces a new Virtual Desktop MSIX package to be created.

* `package_relative_path` - (Required) Relative Path to the package inside the image. Changing this forces a new Virtual Desktop MSIX package to be created.

* `version` - (Required) Package version found in the appxmanifest.xml. Example: "1.0.0.1". Changing this forces a new Virtual Desktop MSIX package to be created.

---

* `display_name` - (Optional) User friendly Name to be displayed in the portal.

* `enabled` - (Optional) Make this version of the package the active one across the host pool. Defaults to `true`.

* `package_dependency` - (Optional) One or more `package_dependency` blocks as defined below.

* `regular_registration_enabled` - (Optional) Specifies how to register the package in feed. Defaults to `true`.

---

A `package_application` block supports the following:

* `app_id` - (Required) The package application id, found in appxmanifest.xml. Changing this forces a new Virtual Desktop MSIX package to be created.

* `app_user_model_id` - (Required) Used to activate the package application. Found in appxmanifest.xml. Changing this forces a new Virtual Desktop MSIX package to be created.

* `description` - (Required) The description of the package application. Changing this forces a new Virtual Desktop MSIX package to be created.

* `friendly_name` - (Required) User friendly name of the package application. Changing this forces a new Virtual Desktop MSIX package to be created.

* `icon_image_name` - (Required) Icon image name. Changing this forces a new Virtual Desktop MSIX package to be created.

* `raw_icon` - (Required) Base64-encoded byte array of the icon file. Changing this forces a new Virtual Desktop MSIX package to be created.

* `raw_png` - (Required) Base64-encoded byte array of the png file. Changing this forces a new Virtual Desktop MSIX package to be created.

---

A `package_dependency` block supports the following:

* `dependency_name` - (Required) The name of the package dependency. Changing this forces a new Virtual Desktop MSIX package to be created.

* `min_version` - (Required) Required version of the dependency. Changing this forces a new Virtual Desktop MSIX package to be created.

* `publisher` - (Required) The package publisher name. Changing this forces a new Virtual Desktop MSIX package to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Desktop MSIX package.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Desktop MSIX package.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Desktop MSIX package.
* `update` - (Defaults to 30 minutes) Used when updating the Virtual Desktop MSIX package.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Desktop MSIX package.

## Import

Virtual Desktop MSIX packages can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_desktop_msix_package.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/Microsoft.DesktopVirtualization/hostPools/myHostPool1/msixPackages/myMsixPackage1
```
