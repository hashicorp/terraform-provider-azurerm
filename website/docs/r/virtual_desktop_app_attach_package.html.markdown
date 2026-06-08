---
subcategory: "Desktop Virtualization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_desktop_app_attach_package"
description: |-
  Manages a Virtual Desktop App Attach Package.
---

# azurerm_virtual_desktop_app_attach_package

Manages a Virtual Desktop App Attach Package.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "example"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_replication_type = "LRS"
  account_tier             = "Standard"
}

resource "azurerm_storage_share" "example" {
  name               = "example"
  quota              = 16
  storage_account_id = azurerm_storage_account.example.id
}

resource "azurerm_storage_share_file" "example" {
  name              = "example"
  source            = "${path.module}/testdata/example"
  storage_share_url = azurerm_storage_share.example.url
}

resource "azurerm_virtual_network" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/24"]
}

resource "azurerm_subnet" "example" {
  name                 = "example"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.0.0/28"]
}

resource "azurerm_network_interface" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  ip_configuration {
    name                          = "example"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.example.id
  }
}

resource "azurerm_nat_gateway" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_desktop_host_pool" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  load_balancer_type  = "BreadthFirst"
  type                = "Pooled"
}

resource "azurerm_virtual_desktop_host_pool_registration_info" "example" {
  expiration_date = "2026-01-01T00:00:00Z"
  hostpool_id     = azurerm_virtual_desktop_host_pool.example.id
}

resource "azurerm_windows_virtual_machine" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  network_interface_ids = [
    azurerm_network_interface.example.id
  ]
  size                = "Standard_F1als_v7"
  admin_password      = "Password1234"
  admin_username      = "adminuser"
  secure_boot_enabled = true
  vtpm_enabled        = true

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  identity {
    type = "SystemAssigned"
  }

  source_image_reference {
    offer     = "office-365"
    publisher = "microsoftwindowsdesktop"
    sku       = "win11-24h2-avd-m365"
    version   = "latest"
  }
}

resource "azurerm_virtual_machine_extension" "example0" {
  name                 = "example0"
  publisher            = "Microsoft.Azure.Security.WindowsAttestation"
  type                 = "GuestAttestation"
  type_handler_version = "1.0"
  virtual_machine_id   = azurerm_windows_virtual_machine.example.id

  depends_on = [
    azurerm_nat_gateway.example
  ]
}

resource "azurerm_virtual_machine_extension" "example1" {
  name                 = "example1"
  publisher            = "Microsoft.Powershell"
  type                 = "DSC"
  type_handler_version = "2.83"
  virtual_machine_id   = azurerm_windows_virtual_machine.example.id

  protected_settings = jsonencode({
    properties = {
      registrationInfoToken = azurerm_virtual_desktop_host_pool_registration_info.example.token
    }
  })

  settings = jsonencode({
    modulesUrl            = "https://wvdportalstorageblob.blob.core.windows.net/galleryartifacts/Configuration_01-20-2022.zip"
    configurationFunction = "Configuration.ps1\\AddSessionHost"
    properties = {
      hostPoolName = azurerm_virtual_desktop_host_pool.example.name
      aadJoin      = true
    }
  })

  depends_on = [
    azurerm_nat_gateway.example
  ]
}

resource "azurerm_virtual_machine_extension" "example2" {
  name                 = "example2"
  publisher            = "Microsoft.Azure.ActiveDirectory"
  type                 = "AADLoginForWindows"
  type_handler_version = "2.2"
  virtual_machine_id   = azurerm_windows_virtual_machine.example.id

  depends_on = [
    azurerm_nat_gateway.example
  ]
}

resource "azurerm_virtual_desktop_app_attach_package" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  display_name        = "example"
  host_pool_ids = [
    azurerm_virtual_desktop_host_pool.example.id
  ]
  msix_package_name     = "example"
  storage_share_file_id = azurerm_storage_share_file.example.id

  depends_on = [
    azurerm_virtual_machine_extension.example0,
    azurerm_virtual_machine_extension.example1,
    azurerm_virtual_machine_extension.example2
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the App Attach Package. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the App Attach Package should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the App Attach Package should exist. Changing this forces a new resource to be created.

* `display_name` - (Required) The user friendly name of the App Attach Package to be displayed in the portal.

* `host_pool_ids` - (Required) A set of Virtual Desktop Host Pool IDs that the App Attach Package is associated with.

~> **Note:** Changing the number of `host_pool_ids` forces a new App Attach Package to be created.

* `msix_package_name` - (Required) The App Attach Package full name from the `appxmanifest.xml` of the corresponding Windows application package.

* `storage_share_file_id` - (Required) The ID of the Storage Share File containing the VHD or CIM image of the App Attach Package on the network share.

* `health_check_status_on_failure` - (Optional) Indicates how the health check should behave if the App Attach Package fails staging. Possible values are `DoNotFail`, `NeedsAssistance`, and `Unhealthy`. Defaults to `NeedsAssistance`.

* `register_at_log_on_enabled` - (Optional) Whether the registration of App Attach Package in feed during log on is enabled. Defaults to `true`.

* `state_enabled` - (Optional) Whether the App Attach Package is active across the Virtual Desktop Host Pool. Defaults to `false`.

* `tags` - (Optional) A mapping of tags which should be assigned to the App Attach Package.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the App Attach Package.

* `last_updated` - The date the App Attach Package was last updated, found in the `appxmanifest.xml` of the corresponding Windows application package.

* `package_applications` - A `package_applications` block as defined below.

* `package_family_name` - The App Attach Package family name from the `appxmanifest.xml` of the corresponding Windows application package. Contains the package name and publisher name.

* `package_name` - The App Attach Package name from the `appxmanifest.xml` of the corresponding Windows application package.

* `package_relative_path` - The relative path to the App Attach Package inside the correponding VHD or CIM image.

* `version` - The App Attach Package version found in the `appxmanifest.xml` of the corresponding Windows application package.

---

A `package_applications` block exports the following:

* `app_id` - The App Attach Package application ID, found in the `appxmanifest.xml` of the corresponding Windows application package.

* `app_user_model_id` - Used to activate the App Attach Package application. Consists of the package name and application ID found in the `appxmanifest.xml` of the corresponding Windows application package.

* `description` - The description of the App Attach Package application.

* `friendly_name` - The user friendly name of the App Attach Package application.

* `icon_image_name` - The icon image name of the App Attach Package application.

* `raw_icon` - The App Attach Package application icon as a 64 bit string in a byte array.

* `raw_png` - The App Attach Package application PNG icon as a 64 bit string in a byte array.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Attach Package.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Attach Package.
* `update` - (Defaults to 30 minutes) Used when updating the App Attach Package.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Attach Package.

## Import

An App Attach Package can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_desktop_app_attach_package.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DesktopVirtualization/appAttachPackages/appAttachPackage1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DesktopVirtualization` - 2025-10-10
