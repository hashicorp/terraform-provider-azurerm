---
subcategory: "Dev Test"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_test_virtual_network"
description: |-
  Manages a Virtual Network within a DevTest Lab.
---

# azurerm_dev_test_virtual_network

Manages a Virtual Network within a DevTest Lab.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dev_test_lab" "example" {
  name                = "example-devtestlab"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  tags = {
    "Sydney" = "Australia"
  }
}

resource "azurerm_dev_test_virtual_network" "example" {
  name                = "example-network"
  lab_name            = azurerm_dev_test_lab.example.name
  resource_group_name = azurerm_resource_group.example.name

  subnet {
    use_public_ip_address           = "Allow"
    use_in_virtual_machine_creation = "Allow"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Dev Test Virtual Network. Changing this forces a new resource to be created.

* `lab_name` - (Required) Specifies the name of the Dev Test Lab in which the Virtual Network should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Dev Test Lab resource exists. Changing this forces a new resource to be created.

* `description` - (Optional) A description for the Virtual Network.

* `subnet` - (Optional) A `subnet` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `subnet` block supports the following:

* `use_public_ip_address` - (Optional) Can Virtual Machines in this Subnet use Public IP Addresses? Possible values are `Allow`, `Default` and `Deny`. Defaults to `Allow`.

* `use_in_virtual_machine_creation` - (Optional) Can this subnet be used for creating Virtual Machines? Possible values are `Allow`, `Default` and `Deny`. Defaults to `Allow`.

* `shared_public_ip_address` - (Optional) A `shared_public_ip_address` block as defined below.

---

A `shared_public_ip_address` block supports the following:

* `allowed_ports` - (Optional) A list of `allowed_ports` blocks as defined below.

---

An `allowed_ports` block supports the following:

* `backend_port` - (Optional) The port on the Virtual Machine that the traffic will be sent to.

* `transport_protocol` - (Optional) The transport protocol that the traffic will use. Possible values are `TCP` and `UDP`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Test Virtual Network.

* `subnet` - A `subnet` block as defined below.

* `unique_identifier` - The unique immutable identifier of the Dev Test Virtual Network.

---

A `subnet` block exports the following:

* `name` - The name of the Subnet for this Virtual Network.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DevTest Virtual Network.
* `read` - (Defaults to 5 minutes) Used when retrieving the DevTest Virtual Network.
* `update` - (Defaults to 30 minutes) Used when updating the DevTest Virtual Network.
* `delete` - (Defaults to 30 minutes) Used when deleting the DevTest Virtual Network.

## Import

DevTest Virtual Networks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dev_test_virtual_network.network1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/lab1/virtualNetworks/network1
```
