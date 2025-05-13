---
subcategory: "Dev Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_center_environment_type"
description: |-
  Manages a Dev Center Environment Type.
---

# azurerm_dev_center_environment_type

Manages a Dev Center Environment Type.

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

resource "azurerm_dev_center_environment_type" "example" {
  name          = "example-dcet"
  dev_center_id = azurerm_dev_center.example.id

  tags = {
    Env = "Test"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Dev Center Environment Type. Changing this forces a new resource to be created.

* `dev_center_id` - (Required) The ID of the associated Dev Center. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Dev Center Environment Type.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Center Environment Type.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Dev Center Environment Type.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dev Center Environment Type.
* `update` - (Defaults to 30 minutes) Used when updating the Dev Center Environment Type.
* `delete` - (Defaults to 30 minutes) Used when deleting the Dev Center Environment Type.

## Import

An existing Dev Center Environment Type can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_dev_center_environment_type.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DevCenter/devCenters/dc1/environmentTypes/et1
```
