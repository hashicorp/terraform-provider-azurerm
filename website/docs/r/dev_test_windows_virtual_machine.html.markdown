---
subcategory: "Dev Test"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_test_windows_virtual_machine"
description: |-
  Manages a Windows Virtual Machine within a Dev Test Lab.
---

# azurerm_dev_test_windows_virtual_machine

Manages a Windows Virtual Machine within a Dev Test Lab.

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

resource "azurerm_dev_test_windows_virtual_machine" "example" {
  name                   = "example-vm03"
  lab_name               = azurerm_dev_test_lab.example.name
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  size                   = "Standard_DS2"
  username               = "exampleuser99"
  password               = "Pa$w0rd1234!"
  lab_virtual_network_id = azurerm_dev_test_virtual_network.example.id
  lab_subnet_name        = azurerm_dev_test_virtual_network.example.subnet[0].name
  storage_type           = "Premium"
  notes                  = "Some notes about this Virtual Machine."

  gallery_image_reference {
    offer     = "WindowsServer"
    publisher = "MicrosoftWindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Dev Test Machine. Changing this forces a new resource to be created.

-> **NOTE:** The validation requirements for the Name change based on the `os_type` used in this Virtual Machine. For a Linux VM the name must be between 1-62 characters, and for a Windows VM the name must be between 1-15 characters. It must begin and end with a letter or number, and cannot be all numbers.

* `lab_name` - (Required) Specifies the name of the Dev Test Lab in which the Virtual Machine should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Dev Test Lab resource exists. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the Dev Test Lab exists. Changing this forces a new resource to be created.

* `gallery_image_reference` - (Required) A `gallery_image_reference` block as defined below.

* `lab_subnet_name` - (Required) The name of a Subnet within the Dev Test Virtual Network where this machine should exist. Changing this forces a new resource to be created.

* `lab_virtual_network_id` - (Required) The ID of the Dev Test Virtual Network where this Virtual Machine should be created. Changing this forces a new resource to be created.

* `password` - (Required) The Password associated with the `username` used to login to this Virtual Machine. Changing this forces a new resource to be created.

* `size` - (Required) The Machine Size to use for this Virtual Machine, such as `Standard_F2`. Changing this forces a new resource to be created.

* `storage_type` - (Required) The type of Storage to use on this Virtual Machine. Possible values are `Standard` and `Premium`.

* `username` - (Required) The Username associated with the local administrator on this Virtual Machine. Changing this forces a new resource to be created.

---

* `allow_claim` - (Optional) Can this Virtual Machine be claimed by users? Defaults to `true`.

* `disallow_public_ip_address` - (Optional) Should the Virtual Machine be created without a Public IP Address? Changing this forces a new resource to be created.

* `inbound_nat_rule` - (Optional) One or more `inbound_nat_rule` blocks as defined below. Changing this forces a new resource to be created.

-> **NOTE:** If any `inbound_nat_rule` blocks are specified then `disallow_public_ip_address` must be set to `true`.

* `notes` - (Optional) Any notes about the Virtual Machine.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `gallery_image_reference` block supports the following:

* `offer` - (Required) The Offer of the Gallery Image. Changing this forces a new resource to be created.

* `publisher` - (Required) The Publisher of the Gallery Image. Changing this forces a new resource to be created.

* `sku` - (Required) The SKU of the Gallery Image. Changing this forces a new resource to be created.

* `version` - (Required) The Version of the Gallery Image. Changing this forces a new resource to be created.

---

A `inbound_nat_rule` block supports the following:

* `protocol` - (Required) The Protocol used for this NAT Rule. Possible values are `Tcp` and `Udp`. Changing this forces a new resource to be created.

* `backend_port` - (Required) The Backend Port associated with this NAT Rule. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Machine.

* `fqdn` - The FQDN of the Virtual Machine.

* `inbound_nat_rule` - One or more `inbound_nat_rule` blocks as defined below.

* `unique_identifier` - The unique immutable identifier of the Virtual Machine.

---

A `inbound_nat_rule` block exports the following:

* `frontend_port` - The frontend port associated with this Inbound NAT Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DevTest Windows Virtual Machine.
* `update` - (Defaults to 30 minutes) Used when updating the DevTest Windows Virtual Machine.
* `read` - (Defaults to 5 minutes) Used when retrieving the DevTest Windows Virtual Machine.
* `delete` - (Defaults to 30 minutes) Used when deleting the DevTest Windows Virtual Machine.

## Import

DevTest Windows Virtual Machines can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dev_test_windows_virtual_machine.machine1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/lab1/virtualmachines/machine1
```
