---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_gallery_application"
description: |-
  Manages a Gallery Application.
---

# azurerm_gallery_application

Manages a Gallery Application.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Gallery Application. Changing this forces a new resource to be created.

* `gallery_id` - (Required) The ID of the Shared Image Gallery. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Gallery Application exists. Changing this forces a new resource to be created.

* `supported_os_type` - (Required) The type of the Operating System supported for the Gallery Application. Possible values are `Linux` and `Windows`. Changing this forces a new resource to be created.

---

* `description` - (Optional) A description of the Gallery Application.

* `end_of_life_date` - (Optional) The end of life date in RFC3339 format of the Gallery Application.

* `eula` - (Optional) The End User Licence Agreement of the Gallery Application.

* `privacy_statement_uri` - (Optional) The URI containing the Privacy Statement associated with the Gallery Application.

* `release_note_uri` - (Optional) The URI containing the Release Notes associated with the Gallery Application.

* `tags` - (Optional) A mapping of tags to assign to the Gallery Application.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Gallery Application.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Gallery Application.
* `read` - (Defaults to 5 minutes) Used when retrieving the Gallery Application.
* `update` - (Defaults to 30 minutes) Used when updating the Gallery Application.
* `delete` - (Defaults to 30 minutes) Used when deleting the Gallery Application.

## Import

`azurerm_gallery_application` resources can be imported using one of the following methods:

* The `terraform import` CLI command with an `id` string:

  ```shell
  terraform import azurerm_gallery_application.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/applications/{applicationName}
  ```

* An `import` block with an `id` argument:
  
  ```hcl
  import {
    to = azurerm_gallery_application.example
    id = "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/applications/{applicationName}"
  }
  ```

* An `import` block with an `identity` argument:

  ```hcl
  import {
    to       = azurerm_gallery_application.example
    identity = {
      TODO Resource Identity Format
    }
  }
  ```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Compute` - 2022-03-03
