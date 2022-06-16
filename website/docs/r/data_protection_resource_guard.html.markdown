---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_resource_guard"
description: |-
  Manages a Resource Guard.
---

# azurerm_data_protection_resource_guard

Manages a Resource Guard.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_protection_resource_guard" "example" {
  name                = "example-resourceguard"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Resource Guard. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Resource Guard should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Resource Guard should exist. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `vault_critical_operation_exclusion_list` - (Optional) A list of the critical operations which are not protected by this Resource Guard.

* `tags` - (Optional) A mapping of tags which should be assigned to the Resource Guard.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Resource Guard. The only possible value is `SystemAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Resource Guard.

* `identity` - An `identity` block as defined below, which contains the Identity information for this Resource Guard.

---

`identity` exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this Resource Guard.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this Resource Guard.

-> You can access the Principal ID via `${azurerm_data_protection_resource_guard.example.identity.0.principal_id}` and the Tenant ID via `${azurerm_data_protection_resource_guard.example.identity.0.tenant_id}`

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Resource Guard.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Guard.
* `update` - (Defaults to 30 minutes) Used when updating the Resource Guard.
* `delete` - (Defaults to 30 minutes) Used when deleting the Resource Guard.

## Import

Resource Guards can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_resource_guard.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataProtection/resourceGuards/resourceGuard1
```
