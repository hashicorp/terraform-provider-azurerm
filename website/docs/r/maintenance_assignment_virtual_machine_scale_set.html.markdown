---
subcategory: "Maintenance"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_maintenance_assignment_virtual_machine_scale_set"
description: |-
  Manages a Maintenance Assignment.
---

# azurerm_maintenance_assignment_virtual_machine_scale_set

Manages a maintenance assignment to a virtual machine scale set.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

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
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "example" {
  name                = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "example" {
  name                = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  frontend_ip_configuration {
    name                 = "internal"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  loadbalancer_id     = azurerm_lb.test.id
}

resource "azurerm_lb_probe" "example" {
  resource_group_name = azurerm_resource_group.example.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "example"
  port                = 22
  protocol            = "Tcp"
}

resource "azurerm_lb_rule" "example" {
  name                           = "example"
  resource_group_name            = azurerm_resource_group.example.name
  loadbalancer_id                = azurerm_lb.test.id
  probe_id                       = azurerm_lb_probe.test.id
  backend_address_pool_id        = azurerm_lb_backend_address_pool.test.id
  frontend_ip_configuration_name = "internal"
  protocol                       = "Tcp"
  frontend_port                  = 22
  backend_port                   = 22
}

resource "azurerm_maintenance_configuration" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  scope               = "OSImage"
  visibility          = "Custom"

  window {
    start_date_time      = "2021-12-31 00:00"
    expiration_date_time = "9999-12-31 00:00"
    duration             = "06:00"
    time_zone            = "Pacific Standard Time"
    recur_every          = "1Days"
  }
}

resource "azurerm_linux_virtual_machine_scale_set" "example" {
  name                            = "example"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.example.location
  sku                             = "Standard_F2"
  instances                       = 1
  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  upgrade_mode                    = "Automatic"
  health_probe_id                 = azurerm_lb_probe.example.id
  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name                                   = "internal"
      primary                                = true
      subnet_id                              = azurerm_subnet.example.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.example.id]
    }
  }

  automatic_os_upgrade_policy {
    disable_automatic_rollback  = true
    enable_automatic_os_upgrade = true
  }

  rolling_upgrade_policy {
    max_batch_instance_percent              = 20
    max_unhealthy_instance_percent          = 20
    max_unhealthy_upgraded_instance_percent = 20
    pause_time_between_batches              = "PT0S"
  }

  depends_on = ["azurerm_lb_rule.example"]
}

resource "azurerm_maintenance_assignment_virtual_machine_scale_set" "example" {
  location                     = azurerm_resource_group.example.location
  maintenance_configuration_id = azurerm_maintenance_configuration.example.id
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine.example.id
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `maintenance_configuration_id` - (Required) Specifies the ID of the Maintenance Configuration Resource. Changing this forces a new resource to be created.

* `virtual_machine_scale_set_id` - (Required) Specifies the Virtual Machine Scale Set ID to which the Maintenance Configuration will be assigned. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Maintenance Assignment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Maintenance Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Maintenance Assignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Maintenance Assignment.

## Import

Maintenance Assignment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_maintenance_assignment_virtual_machine_scale_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/virtualMachineScaleSets/vmss1/providers/Microsoft.Maintenance/configurationAssignments/assign1
```
