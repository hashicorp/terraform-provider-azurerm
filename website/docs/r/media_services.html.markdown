---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_media_services"
sidebar_current: "docs-azurerm-resource-media-media-services"
description: |-
  Manages a Media Services account.
---

# azurerm_media_services

Manages a Media Services account.

## Example Usage

```hcl
resource "azurerm_resource_group" "testrg" {
  name     = "amstestrg"
  location = "westus"
}

resource "azurerm_storage_account" "testsa" {
  name                     = "amstestsa"
  resource_group_name      = "${azurerm_resource_group.testrg.name}"
  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services" "ams" {

  name                = "amstest"
  location            = "${azurerm_resource_group.testrg.location}"
  resource_group_name = "${azurerm_resource_group.testrg.name}"
		
  storage_account {
		id          = "${azurerm_storage_account.testsa.id}"
		is_primary  = true
	}
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the MariaDB Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the MariaDB Server. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `storage_account` - (Required) A `storage_account` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `storage_account` block supports the following:

* `id` - (Required) Specifies the Id for the storage account that will be associated with the Media Services instance.

* `is_primary` - (Required) Specifies whether the storage account should be the primary account or not. Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The `resource id` of the Media Services account.

## Import

MariaDB Server's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_media_services.ams1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Media/mediaservices/ams1
```
## Note

Multiple `storage_account` arguments can be provided to `azurerm_media_services`, but only one can be marked as `is_primary`=`true`.