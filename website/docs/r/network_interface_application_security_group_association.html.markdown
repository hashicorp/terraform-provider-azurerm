---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_interface_application_security_group_association"
description: |-
  Manages the association between a Network Interface and a Application Security Group

---

# azurerm_network_interface_application_security_group_association

Manages the association between a Network Interface and a Application Security Group.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_application_security_group" "example" {
  name                = "example-asg"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_network_interface" "example" {
  name                = "example-nic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_network_interface_application_security_group_association" "example" {
  network_interface_id          = azurerm_network_interface.example.id
  application_security_group_id = azurerm_application_security_group.example.id
}
```

## Argument Reference

The following arguments are supported:

* `network_interface_id` - (Required) The ID of the Network Interface. Changing this forces a new resource to be created.

* `application_security_group_id` - (Required) The ID of the Application Security Group which this Network Interface which should be connected to. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The (Terraform specific) ID of the Association between the Network Interface and the Application Security Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the association between the Network Interface and the Application Security Group.
* `update` - (Defaults to 30 minutes) Used when updating the association between the Network Interface and the Application Security Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the association between the Network Interface and the Application Security Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the association between the Network Interface and the Application Security Group.

## Import

Associations between Network Interfaces and Application Security Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_interface_application_security_group_association.association1 "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.network/networkInterfaces/nic1|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/applicationSecurityGroups/securityGroup1"
```

-> **NOTE:** This ID is specific to Terraform - and is of the format `{networkInterfaceId}|{applicationSecurityGroupId}`.
