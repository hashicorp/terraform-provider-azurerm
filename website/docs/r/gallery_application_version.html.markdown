---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_gallery_application_version"
description: |-
  Manages a Gallery Application Version.
---

# azurerm_gallery_application_version

Manages a Gallery Application Version.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_shared_image_gallery" "example" {
  name                = "example-gallery"
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
  name                     = "example-storage"
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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The version name of the Gallery Application Version, such as `1.0.0`. Changing this forces a new resource to be created.

* `gallery_application_id` - (Required) The ID of the Gallery Application. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Gallery Application Version exists. Changing this forces a new resource to be created.

* `manage_action` - (Required) A `manage_action` block as defined below.

* `source` - (Required) A `source` block as defined below.

* `target_region` - (Required) One or more `target_region` blocks as defined below.

---

* `enable_health_check` - (Optional) Should the Gallery Application reports health. Defaults to `false`.

* `end_of_life_date` - (Optional) The end of life date in RFC3339 format of the Gallery Application Version.

* `exclude_from_latest` - (Optional) Should the Gallery Application Version be excluded from the `latest` filter? If set to `true` this Gallery Application Version won't be returned for the `latest` version. Defaults to `false`.

* `tags` - (Optional) A mapping of tags to assign to the Gallery Application Version.

---

A `manage_action` block supports the following:

* `install` - (Required) The command to install the Gallery Application. Changing this forces a new resource to be created.

* `remove` - (Required) The command to remove the Gallery Application. Changing this forces a new resource to be created.

* `update` - (Optional) The command to update the Gallery Application. Changing this forces a new resource to be created.

---

A `source` block supports the following:

* `media_link` - (Required) The Storage Blob URI of the source application package. Changing this forces a new resource to be created.

* `default_configuration_link` - (Optional) The Storage Blob URI of the default configuration. Changing this forces a new resource to be created.

---

A `target_region` block supports the following:

* `name` - (Required) The Azure Region in which the Gallery Application Version exists.

* `regional_replica_count` - (Required) The number of replicas of the Gallery Application Version to be created per region. Possible values are between `1` and `10`.

* `storage_account_type` - (Optional) The storage account type for the Gallery Application Version. Possible values are `Standard_LRS`, `Premium_LRS` and `Standard_ZRS`. Defaults to `Standard_LRS`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Gallery Application Version.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Gallery Application Version.
* `read` - (Defaults to 5 minutes) Used when retrieving the Gallery Application Version.
* `update` - (Defaults to 30 minutes) Used when updating the Gallery Application Version.
* `delete` - (Defaults to 30 minutes) Used when deleting the Gallery Application Version.

## Import

Gallery Application Versions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_gallery_application_version.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/galleries/gallery1/applications/galleryApplication1/versions/galleryApplicationVersion1
```
