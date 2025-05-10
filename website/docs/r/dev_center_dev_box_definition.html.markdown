---
subcategory: "Dev Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_center_dev_box_definition"
description: |-
  Manages a Dev Center Dev Box Definition.
---

# azurerm_dev_center_dev_box_definition

Manages a Dev Center Dev Box Definition.

## Example Usage

```hcl
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

resource "azurerm_dev_center_dev_box_definition" "example" {
  name               = "example-dcet"
  location           = azurerm_resource_group.example.location
  dev_center_id      = azurerm_dev_center.example.id
  image_reference_id = "${azurerm_dev_center.example.id}/galleries/default/images/microsoftvisualstudio_visualstudioplustools_vs-2022-ent-general-win10-m365-gen2"
  sku_name           = "general_i_8c32gb256ssd_v2"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Dev Center Dev Box Definition. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Dev Center Dev Box Definition should exist. Changing this forces a new resource to be created.

* `dev_center_id` - (Required) The ID of the associated Dev Center. Changing this forces a new resource to be created.

* `image_reference_id` - (Required) The ID of the image for the Dev Center Dev Box Definition.

* `sku_name` - (Required) The name of the SKU for the Dev Center Dev Box Definition.

* `tags` - (Optional) A mapping of tags which should be assigned to the Dev Center Dev Box Definition.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Center Dev Box Definition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Dev Center Dev Box Definition.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dev Center Dev Box Definition.
* `update` - (Defaults to 30 minutes) Used when updating the Dev Center Dev Box Definition.
* `delete` - (Defaults to 30 minutes) Used when deleting the Dev Center Dev Box Definition.

## Import

An existing Dev Center Dev Box Definition can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_dev_center_dev_box_definition.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DevCenter/devCenters/dc1/devBoxDefinitions/et1
```
