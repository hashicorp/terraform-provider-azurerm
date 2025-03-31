---
subcategory: "Dev Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_center_project"
description: |-
  Manages a Dev Center Project.
---

# azurerm_dev_center_project

Manages a Dev Center Project.

## Example Usage

```hcl
resource "azurerm_dev_center" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  identity {
    type = "example-value"
  }
}
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_dev_center_project" "example" {
  dev_center_id       = azurerm_dev_center.example.id
  location            = azurerm_resource_group.example.location
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `dev_center_id` - (Required) Resource Id of an associated DevCenter. Changing this forces a new Dev Center Project to be created.

* `location` - (Required) The Azure Region where the Dev Center Project should exist. Changing this forces a new Dev Center Project to be created.

* `name` - (Required) Specifies the name of this Dev Center Project. Changing this forces a new Dev Center Project to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group within which this Dev Center Project should exist. Changing this forces a new Dev Center Project to be created.

* `description` - (Optional) Description of the project. Changing this forces a new Dev Center Project to be created.

* `maximum_dev_boxes_per_user` - (Optional) When specified, limits the maximum number of Dev Boxes a single user can create across all pools in the project.

* `tags` - (Optional) A mapping of tags which should be assigned to the Dev Center Project.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Center Project.

* `dev_center_uri` - The URI of the Dev Center resource this project is associated with.

---



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating this Dev Center Project.
* `delete` - (Defaults to 30 minutes) Used when deleting this Dev Center Project.
* `read` - (Defaults to 5 minutes) Used when retrieving this Dev Center Project.
* `update` - (Defaults to 30 minutes) Used when updating this Dev Center Project.

## Import

An existing Dev Center Project can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_dev_center_project.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevCenter/projects/{projectName}
```

* Where `{subscriptionId}` is the ID of the Azure Subscription where the Dev Center Project exists. For example `12345678-1234-9876-4563-123456789012`.
* Where `{resourceGroupName}` is the name of Resource Group where this Dev Center Project exists. For example `example-resource-group`.
* Where `{projectName}` is the name of the Project. For example `projectValue`.
