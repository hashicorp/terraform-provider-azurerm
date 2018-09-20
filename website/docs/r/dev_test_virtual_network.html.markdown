---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_test_virtual_network"
sidebar_current: "docs-azurerm-resource-dev-test-virtual-network"
description: |-
  Manages a Virtual Network within a Dev Test Lab.
---

# azurerm_dev_test_virtual_network

Manages a Virtual Network within a Dev Test Lab.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West US"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "example-devtestlab"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags {
    "Sydney" = "Australia"
  }
}


resource "azurerm_dev_test_virtual_network" "test" {
  name                = "example-network"
  lab_name            = "${azurerm_dev_test_lab.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Dev Test Lab. Changing this forces a new resource to be created.

* `lab_name` - (Required) Specifies the name of the Dev Test Lab in which the Virtual Network should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Dev Test Lab resource exists. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the Dev Test Lab exists. Changing this forces a new resource to be created.

* `description` - (Optional) A description for the Virtual Network.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Dev Test Virtual Network.

* `unique_identifier` - The unique immutable identifier of the Dev Test Virtual Network.

## Import

Dev Test Virtual Networks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dev_test_virtual_network.network1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/lab1/virtualnetworks/network1
```
