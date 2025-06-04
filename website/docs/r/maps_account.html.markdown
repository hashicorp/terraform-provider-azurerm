---
subcategory: "Maps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_maps_account"
description: |-
  Manages an Azure Maps Account.
---

# azurerm_maps_account

Manages an Azure Maps Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_maps_account" "example" {
  name                         = "example-maps-account"
  resource_group_name          = azurerm_resource_group.example.name
  sku_name                     = "S1"
  local_authentication_enabled = true

  tags = {
    environment = "Test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Azure Maps Account. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Azure Maps Account should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Location in which the Azure Maps Account should be provisioned. Changing this forces a new resource to be created.

* `sku_name` - (Required) The SKU of the Azure Maps Account. Possible values are `S0`, `S1` and `G2`. Changing this forces a new resource to be created.

~> **Note:** Gen1 SKUs (`S0` and `S1`) are deprecated and can no longer be used for new deployments, which should instead use a Gen2 SKU (`G2`) - more information can be found [in the Azure documentation](https://learn.microsoft.com/azure/azure-maps/how-to-manage-pricing-tier).

* `cors` - (Optional) - A `cors` block as defined below

* `data_store` - (Optional) One or more `data_store` blocks as defined below.

* `identity` - (Optional) An `identity` block as defined below.

* `local_authentication_enabled` - (Optional) Is local authentication enabled for this Azure Maps Account? When `false`, all authentication to the Azure Maps data-plane REST API is disabled, except Azure AD authentication. Defaults to `true`.

* `tags` - (Optional) A mapping of tags to assign to the Azure Maps Account.

---

A `cors` block supports the following:

* `allowed_origins` - (Required) A list of origins that should be allowed to make cross-origin calls.

---

A `data_store` block supports the following:

* `storage_account_id` - (Required) The ID of the Storage Account that should be linked to this Azure Maps Account.

* `unique_name` - (Required) The name given to the linked Storage Account.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Azure Maps Account. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) A list of User Assigned Managed Identity IDs to be assigned to this Azure Maps Account.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Maps Account.

* `identity` - An `identity` block as defined below.

* `primary_access_key` - The primary key used to authenticate and authorize access to the Maps REST APIs.

* `secondary_access_key` - The secondary key used to authenticate and authorize access to the Maps REST APIs.

* `x_ms_client_id` - A unique identifier for the Maps Account.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Maps Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Maps Account.
* `update` - (Defaults to 30 minutes) Used when updating the Maps Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Maps Account.

## Import

A Maps Account can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_maps_account.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Maps/accounts/my-maps-account
```
