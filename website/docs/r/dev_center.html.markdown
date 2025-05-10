---
subcategory: "Dev Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_center"
description: |-
  Manages a Dev Center.
---

# azurerm_dev_center

Manages a Dev Center.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_user_assigned_identity" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
resource "azurerm_dev_center" "example" {
  location            = azurerm_resource_group.example.location
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Dev Center should exist. Changing this forces a new Dev Center to be created.

* `name` - (Required) Specifies the name of this Dev Center. Changing this forces a new Dev Center to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group within which this Dev Center should exist. Changing this forces a new Dev Center to be created.

* `identity` - (Optional) An `identity` block as defined below. Specifies the Managed Identity which should be assigned to this Dev Center.

* `tags` - (Optional) A mapping of tags which should be assigned to the Dev Center.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Center.

* `dev_center_uri` - The URI of the Dev Center.

---

## Blocks Reference

### `identity` Block


The `identity` block supports the following arguments:

* `type` - (Required) Specifies the type of Managed Identity that should be assigned to this Dev Center. Possible values are `SystemAssigned`, `SystemAssigned, UserAssigned` and `UserAssigned`.
* `identity_ids` - (Optional) A list of the User Assigned Identity IDs that should be assigned to this Dev Center.


In addition to the arguments defined above, the `identity` block exports the following attributes:

* `principal_id` - The Principal ID for the System-Assigned Managed Identity assigned to this Dev Center.
* `tenant_id` - The Tenant ID for the System-Assigned Managed Identity assigned to this Dev Center.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Dev Center.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dev Center.
* `update` - (Defaults to 30 minutes) Used when updating the Dev Center.
* `delete` - (Defaults to 30 minutes) Used when deleting the Dev Center.

## Import

An existing Dev Center can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_dev_center.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevCenter/devCenters/{devCenterName}
```

* Where `{subscriptionId}` is the ID of the Azure Subscription where the Dev Center exists. For example `12345678-1234-9876-4563-123456789012`.
* Where `{resourceGroupName}` is the name of Resource Group where this Dev Center exists. For example `example-resource-group`.
* Where `{devCenterName}` is the name of the Dev Center. For example `devCenterValue`.
