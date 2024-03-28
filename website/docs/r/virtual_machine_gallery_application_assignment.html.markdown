---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_gallery_application_assignment"
description: |-
  Manages a Virtual Machine Gallery Application Assignment.
---

# azurerm_virtual_machine_gallery_application_assignment

Manages a Virtual Machine Gallery Application Assignment.

~> **Note:** Gallery Application Assignments can be defined either directly on `azurerm_linux_virtual_machine` and `azurerm_windows_virtual_machine` resources, or using the `azurerm_virtual_machine_gallery_application_assignment` resource - but the two approaches cannot be used together. If both are used with the same Virtual Machine, spurious changes will occur. It's recommended to use `ignore_changes` for the `gallery_application` block on the associated virtual machine resources, to avoid a persistent diff when using this resource.
## Example Usage

```hcl
data "azurerm_virtual_machine" "example" {
  name                = "example-vm"
  resource_group_name = "example-resources-vm"
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_shared_image_gallery" "example" {
  name                = "examplegallery"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_gallery_application" "example" {
  name              = "example-app"
  gallery_id        = azurerm_shared_image_gallery.example.id
  location          = azurerm_resource_group.example.location
  supported_os_type = "Linux"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorage"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "example-container"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "blob"
}

resource "azurerm_storage_blob" "example" {
  name                   = "scripts"
  storage_account_name   = azurerm_storage_account.example.name
  storage_container_name = azurerm_storage_container.example.name
  type                   = "Block"
  source_content         = "[scripts file content]"
}

resource "azurerm_gallery_application_version" "example" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.example.id
  location               = azurerm_gallery_application.example.location

  manage_action {
    install = "[install command]"
    remove  = "[remove command]"
  }

  source {
    media_link = azurerm_storage_blob.example.id
  }

  target_region {
    name                   = azurerm_gallery_application.example.location
    regional_replica_count = 1
  }
}

resource "azurerm_virtual_machine_gallery_application_assignment" "example" {
  gallery_application_version_id = azurerm_gallery_application_version.example.id
  virtual_machine_id             = data.azurerm_virtual_machine.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `gallery_application_version_id` - (Required) The ID of the Gallery Application Version. Changing this forces a new resource to be created.

* `virtual_machine_id` - (Required) The ID of the Virtual Machine. Changing this forces a new resource to be created.

* `configuration_blob_uri` - (Optional) Specifies the URI to an Azure Blob that will replace the default configuration for the package if provided. Changing this forces a new resource to be created.

* `order` - (Optional) Specifies the order in which the packages have to be installed. Possible values are between `0` and `2147483647`. Defaults to `0`.

* `tag` - (Optional) Specifies a passthrough value for more generic context. This field can be any valid `string` value. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Virtual Machine Gallery Application Assignment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Machine Gallery Application Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Machine Gallery Application Assignment.
* `update` - (Defaults to 30 minutes) Used when updating the Virtual Machine Gallery Application Assignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Machine Gallery Application Assignment.

## Import

Virtual Machine Gallery Application Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_machine_gallery_application_assignment.example subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/virtualMachines/machine1|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/galleries/gallery1/applications/galleryApplication1/versions/galleryApplicationVersion1
```
