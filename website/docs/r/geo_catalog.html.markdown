---
subcategory: "Orbital Planetary Computer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_geo_catalog"
description: |-
  Manages a GeoCatalog.
---

# azurerm_geo_catalog

Manages a GeoCatalog.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "southeastasia"
}
resource "azurerm_geo_catalog" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the GeoCatalog. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the GeoCatalog should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the GeoCatalog should exist. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the GeoCatalog.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this GeoCatalog. Possible values are `SystemAssigned` and `UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this GeoCatalog.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the GeoCatalog.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this GeoCatalog.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this GeoCatalog.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the GeoCatalog.
* `read` - (Defaults to 5 minutes) Used when retrieving the GeoCatalog.
* `update` - (Defaults to 90 minutes) Used when updating the GeoCatalog.
* `delete` - (Defaults to 1 hour) Used when deleting the GeoCatalog.

## Import

GeoCatalogs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_geo_catalog.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Orbital/geoCatalogs/geoCatalog1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Orbital` - 2026-04-15
