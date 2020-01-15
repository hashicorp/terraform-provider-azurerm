---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_virtual_machine"
description: |-
  Manage Azure MsSqlVirtualMachine instance.
---

# azurerm_mssql_virtual_machine

Manage Azure MsSqlVirtualMachine instance.


## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-group"
  location = "example-location"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_kind             = "StorageV2"
  account_tier             = "Premium"
  account_replication_type = "LRS"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vn"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "test" {
  name                      = "example-sub"
  resource_group_name       = azurerm_resource_group.test.name
  virtual_network_name      = azurerm_virtual_network.test.name
  address_prefix            = "10.0.0.0/24"
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_public_ip" "vm" {
  name                = "exampleIP"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Dynamic"
}

resource "azurerm_network_security_group" "nsg" {
  name                = "examplensg"
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
  name                      = "exampleni"
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
    name          = "exampleOSDisk"
    vhd_uri       = "${azurerm_storage_account.example.primary_blob_endpoint}vhds/exampleOSDisk.vhd"
    caching       = "ReadOnly"
    create_option = "FromImage"
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
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  virtual_machine_resource_id = azurerm_virtual_machine.example.id
  sql_license_type            = "PAYG"
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group that contains the resource. You can obtain this value from the Azure Resource Manager API or the portal. Changing this forces a new resource to be created.

* `location` - (Required) The resource location. Changing this forces a new resource to be created.

* `virtual_machine_resource_id` - (Required) The ARM Resource id of underlying virtual machine created from SQL marketplace image.

* `sql_virtual_machine_group_resource_id` - (Optional) The ARM resource id of the SQL virtual machine group this SQL virtual machine is or will be part of.

* `sql_license_type` - (Optional) The SQL Server license type. Possible values include: 'PAYG'(Pay As You Go), 'AHUB'(Azure Hybrid Benefit). Defaults to `PAYG`.

* `sql_sku` - (Optional) The SQL Server edition type. Possible values include: 'Developer', 'Express', 'Standard', 'Enterprise', 'Web'. Defaults to `Developer`.

* `auto_patching` - (Optional) The `auto_patching_setting` block defined below.SQL Server Azure VMs can use Automated Patching to schedule a maintenance window for installing important windows and SQL Server updates automatically. Please refer [automated patching](https://docs.microsoft.com/en-us/azure/virtual-machines/windows/sql/virtual-machines-windows-sql-automated-patching) for more information.

* `key_vault_credential` - (Optional) The `key_vault_credential_setting` block defined below. For more information, please refer to [virtual machines windows sql keyvault](https://docs.microsoft.com/en-us/azure/virtual-machines/windows/sql/virtual-machines-windows-ps-sql-keyvault)

* `server_configuration` - (Optional) The `server_configurations_management_setting` block defined below.

* `tags` - (Optional) The Resource tags. Changing this forces a new resource to be created.


The `auto_patching_setting` block supports the following:

* `enable` - (Optional) Enable or disable autopatching on SQL virtual machine.

* `day_of_week` - (Optional) The day of week to apply the patch on. Defaults to `Monday`.

* `maintenance_window_starting_hour` - (Optional) The hour of the day when patching is initiated. Local VM time.

* `maintenance_window_duration_in_minutes` - (Optional) The duration of patching.

* `name` - (Computed) The name of the SQL virtual machine, which is the same with the name of the Virtual Machine provided.
---

The `key_vault_credential_setting` block supports the following:

* `enable` - (Optional) Enable or disable key vault credential setting.

* `credential_name` - (Optional) The credential name.

* `azure_key_vault_url` - (Optional) The azure Key Vault url.

* `service_principal_name` - (Optional) The service principal name to access key vault.

* `service_principal_secret` - (Optional) The service principal name secret to access key vault.

---

The `server_configurations_management_setting` block supports the following:

* `sql_connectivity_type` - (Optional) The SQL Server connectivity option. Defaults to `LOCAL`.

* `sql_connectivity_port` - (Optional) The SQL Server port.

* `sql_connectivity_update_user_name` - (Optional) The SQL Server sysadmin login to create.

* `sql_connectivity_update_password` - (Optional) The SQL Server sysadmin login password.

* `is_r_services_enabled` - (Optional) Enable or disable R services (SQL 2016 onwards).Enables SQL Server Machine Learning Services (In-Database), allowing you to utilize advanced analytics within your SQL Server. SQL Server Machine Learning Services (In-Database) is only supported with SQL Server 2017 Enterprise.

## Attributes Reference

The following attributes are exported:
* `id` - Resource ID.

* `name` - Resource name.

## Import

Sql Virtual Machine can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_sql_virtual_machine.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resource-group/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachines/example-sql-virtual-machine
```
