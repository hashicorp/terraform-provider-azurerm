---
subcategory: "Dev Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_center_project_environment_type"
description: |-
  Manages a Dev Center Project Environment Type.
---

# azurerm_dev_center_project_environment_type

Manages a Dev Center Project Environment Type.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dev_center" "example" {
  name                = "example-dc"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_dev_center_environment_type" "example" {
  name          = "example-et"
  dev_center_id = azurerm_dev_center.example.id
}

resource "azurerm_dev_center_project" "example" {
  name                = "example-dcp"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  dev_center_id       = azurerm_dev_center.example.id

  depends_on = [azurerm_dev_center_environment_type.example]
}

resource "azurerm_dev_center_project_environment_type" "example" {
  name                  = "example-et"
  location              = azurerm_resource_group.example.location
  dev_center_project_id = azurerm_dev_center_project.example.id
  deployment_target_id  = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"

  identity {
    type = "SystemAssigned"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Dev Center Project Environment Type. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Dev Center Project Environment Type should exist. Changing this forces a new resource to be created.

* `dev_center_project_id` - (Required) The ID of the associated Dev Center Project. Changing this forces a new resource to be created.

* `deployment_target_id` - (Required) The ID of the subscription that the Environment Type will be mapped to. The environment's resources will be deployed into this subscription.

* `identity` - (Required) An `identity` block as defined below.

* `creator_role_assignment_roles` - (Optional) A list of roles to assign to the environment creator.

* `user_role_assignment` - (Optional) A `user_role_assignment` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Dev Center Project Environment Type.

---

An `identity` block supports the following:

* `type` - (Required) The type of identity used for this Dev Center Project Environment Type. Possible values are `SystemAssigned`, `UserAssigned` and `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) The ID of the User Assigned Identity which should be assigned to this Dev Center Project Environment Type.

-> **Note:** `identity_ids` is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `user_role_assignment` block supports the following:

* `user_id` - (Required) The user object ID that is assigned roles.

* `roles` - (Required) A list of roles to assign to the `user_id`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Center Project Environment Type.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Dev Center Project Environment Type.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dev Center Project Environment Type.
* `update` - (Defaults to 30 minutes) Used when updating the Dev Center Project Environment Type.
* `delete` - (Defaults to 30 minutes) Used when deleting the Dev Center Project Environment Type.

## Import

An existing Dev Center Project Environment Type can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_dev_center_project_environment_type.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DevCenter/projects/project1/environmentTypes/et1
```
