---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_virtual_machine"
description: |-
  Manages a Microsoft SQL Virtual Machine
---

# azurerm_mssql_virtual_machine

Manages a Microsoft SQL Virtual Machine

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vn"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "test" {
  name                      = "example-sn"
  resource_group_name       = azurerm_resource_group.test.name
  virtual_network_name      = azurerm_virtual_network.test.name
  address_prefix            = "10.0.0.0/24"
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_public_ip" "vm" {
  name                = "example-pip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Dynamic"
}

resource "azurerm_network_security_group" "nsg" {
  name                = "example-nsg"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_network_security_rule" "RDPRule" {
  name                        = "RDPRule"
  resource_group_name         = azurerm_resource_group.example.name
  priority                    = 1000
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = 3389
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  network_security_group_name = azurerm_network_security_group.nsg.name
}

resource "azurerm_network_security_rule" "MSSQLRule" {
  name                        = "MSSQLRule"
  resource_group_name         = azurerm_resource_group.example.name
  priority                    = 1001
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = 1433
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  network_security_group_name = azurerm_network_security_group.nsg.name
}

resource "azurerm_network_interface" "example" {
  name                      = "example-nic"
  location                  = azurerm_resource_group.example.location
  resource_group_name       = azurerm_resource_group.example.name
  network_security_group_id = azurerm_network_security_group.nsg.id

  ip_configuration {
    name                          = "exampleconfiguration1"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.vm.id
  }
}

resource "azurerm_virtual_machine" "example" {
  name                  = "examplevm"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  network_interface_ids = [azurerm_network_interface.example.id]
  vm_size               = "Standard_DS13"

  storage_image_reference {
    publisher = "MicrosoftSQLServer"
    offer     = "SQL2019-WS2019"
    sku       = "SQLDEV"
    version   = "laexample"
  }

  storage_os_disk {
    name              = "exampleOSDisk"
    caching           = "ReadOnly"
    create_option     = "FromImage"
    managed_disk_type = "Premium_LRS"
  }

  os_profile {
    computer_name  = "winhost01"
    admin_username = "exampleadmin"
    admin_password = "Password1234!"
  }

  os_profile_windows_config {
    timezone                  = "Pacific Standard Time"
    provision_vm_agent        = true
    enable_automatic_upgrades = true
  }
}

resource "azurerm_mssql_virtual_machine" "example" {
  virtual_machine_id = azurerm_virtual_machine.test.id
  sql_license_type   = "PAYG"
}
```

## Argument Reference

The following arguments are supported:

* `virtual_machine_id` - (Required) The ID of the Virtual Machine. Changing this forces a new resource to be created.

* `sql_virtual_machine_group_id` - (Optional) The ARM resource id of the SQL virtual machine group this SQL virtual machine is or will be part of.

* `sql_license_type` - (Optional) The SQL Server license type. Possible values are `AHUB` (Azure Hybrid Benefit) and `PAYG` (Pay-As-You-Go). Changing this forces a new resource to be created.

* `auto_patching` - (Optional) An `auto_patching` block as defined below.

* `key_vault_credential` - (Optional) (Optional) An `key_vault_credential` block as defined below.

* `server_configuration` - (Optional) An `server_configuration` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `auto_patching` block supports the following:

* `day_of_week` - (Required) The day of week to apply the patch on.

* `maintenance_window_starting_hour` - (Required) The Hour, in the Virtual Machine Time-Zone when the patching maintenance window should begin.

* `maintenance_window_duration_in_minutes` - (Required) The size of the Maintenance Window in minutes.

---

The `key_vault_credential` block supports the following:

* `name` - (Required) The credential name.

* `azure_key_vault_url` - (Required) The azure Key Vault url. Changing this forces a new resource to be created.

* `service_principal_name` - (Required) The service principal name to access key vault. Changing this forces a new resource to be created.

* `service_principal_secret` - (Required) The service principal name secret to access key vault. Changing this forces a new resource to be created.

---

The `server_configuration` block supports the following:

* `is_r_services_enabled` - (Optional) Should R Services be enabled?

* `sql_connectivity_type` - (Optional) The connectivity type used for this SQL Server. Defaults to `PRIVATE`.

* `sql_connectivity_port` - (Optional) The SQL Server port. Defaults to `1433`.

* `sql_connectivity_update_username` - (Optional) The SQL Server sysadmin login to create.

* `sql_connectivity_update_password` - (Optional) The SQL Server sysadmin login password.



## Attributes Reference

The following attributes are exported:
* `id` - The ID of the SQL Virtual Machine.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the MSSQL Virtual Machine.
* `update` - (Defaults to 60 minutes) Used when updating the MSSQL Virtual Machine.
* `read` - (Defaults to 5 minutes) Used when retrieving the MSSQL Virtual Machine.
* `delete` - (Defaults to 60 minutes) Used when deleting the MSSQL Virtual Machine.


## Import

Sql Virtual Machines can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_mssql_virtual_machine.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachines/example1
```
