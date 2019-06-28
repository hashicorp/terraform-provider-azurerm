---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_media_services_account"
sidebar_current: "docs-azurerm-resource-media-media-services-account"
description: |-
  Manages a Media Services Account.
---

# azurerm_media_services_account

Manages a Media Services Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "media-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "test" {
  name                     = "examplestoracc"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services_account" "test" {
  name                = "examplemediaacc"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  storage_account {
    id         = "${azurerm_storage_account.test.id}"
    is_primary = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Media Services Account. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Media Services Account. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `storage_account` - (Required) One or more `storage_account` blocks as defined below.

---

A `storage_account` block supports the following:

* `id` - (Required) Specifies the ID of the Storage Account that will be associated with the Media Services instance.

* `is_primary` - (Required) Specifies whether the storage account should be the primary account or not. Defaults to `false`.

~> **NOTE:** Whilst multiple `storage_account` blocks can be specified - one of them must be set to the primary

## Attributes Reference

The following attributes are exported:

* `id` - The Resource ID of the Media Services Account.

## Import

Media Services Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_media_services_account.account /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Media/mediaservices/account1
```
