---
subcategory: "Media"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_media_services_account"
description: |-
  Manages a Media Services Account.
---

# azurerm_media_services_account

Manages a Media Services Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "media-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestoracc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services_account" "example" {
  name                = "examplemediaacc"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  storage_account {
    id         = azurerm_storage_account.example.id
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

* `identity` - (Optional) An `identity` block is documented below.

* `storage_authentication` - (Optional) Specifies the storage authentication type. 
Possible value is  `ManagedIdentity` or `System`.

* `tags` - (Optional) A mapping of tags assigned to the resource.


---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Media Services Account. Possible value is  `SystemAssigned`. 

---


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Media Services Account.

* `identity` - An `identity` block as defined below.
---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Media Services Account.
* `update` - (Defaults to 30 minutes) Used when updating the Media Services Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Media Services Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Media Services Account.

## Import

Media Services Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_media_services_account.account /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Media/mediaservices/account1
```
