// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_virtual_machine_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type MsSqlVirtualMachineAvailabilityGroupListenerResource struct{}

func (r MsSqlVirtualMachineAvailabilityGroupListenerResource) configureDomain(data acceptance.TestData) string {
	return fmt.Sprintf(
		`
%[1]s

resource "azurerm_virtual_machine_extension" "custom_script" {
  name                 = "create-active-directory-forest"
  virtual_machine_id   = azurerm_windows_virtual_machine.domain_controller.id
  publisher            = "Microsoft.Compute"
  type                 = "CustomScriptExtension"
  type_handler_version = "1.9"
  settings             = <<SETTINGS
  {
    "commandToExecute": "powershell.exe -Command \"${local.configure_domain_command}\""
  }
SETTINGS
}
`, r.domainDependencies(data),
	)
}

func (MsSqlVirtualMachineAvailabilityGroupListenerResource) domainDependencies(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

locals {
  domain_controller_vm_name     = "acctest%[3]sdc"
  active_directory_netbios_name = "acctest%[3]s"
  active_directory_domain_name  = "acctest%[3]s.local"
  domain_controller_vm_fqdn     = join(".", [local.domain_controller_vm_name, local.active_directory_domain_name])
  admin_username                = "adminuser"
  admin_password                = "P@ssw0rd1234!"

  auto_logon_data    = "<AutoLogon><Password><Value>${local.admin_password}</Value></Password><Enabled>true</Enabled><LogonCount>1</LogonCount><Username>${local.admin_username}</Username></AutoLogon>"
  first_logon_data   = file("testdata/FirstLogonCommands.xml")
  custom_data_params = "Param($RemoteHostName = \"${local.domain_controller_vm_fqdn}\", $ComputerName = \"${local.domain_controller_vm_name}\")"
  custom_data        = base64encode(join(" ", [local.custom_data_params, file("testdata/winrm.ps1")]))

  import_command           = "Import-Module ADDSDeployment"
  password_command         = "$password = ConvertTo-SecureString ${local.admin_password} -AsPlainText -Force"
  install_ad_command       = "Add-WindowsFeature -Name AD-Domain-Services -IncludeManagementTools"
  configure_ad_command     = "Install-ADDSForest -CreateDnsDelegation:$false -DomainMode Win2012R2 -DomainName ${local.active_directory_domain_name} -DomainNetbiosName ${local.active_directory_netbios_name} -ForestMode Win2012R2 -InstallDns:$true -SafeModeAdministratorPassword $password -Force:$true"
  shutdown_command         = "shutdown -r -t 10"
  exit_code_hack           = "exit 0"
  configure_domain_command = "${local.import_command}; ${local.password_command}; ${local.install_ad_command}; ${local.configure_ad_command}; ${local.shutdown_command}; ${local.exit_code_hack}"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_servers         = ["10.0.1.4", "8.8.8.8"]
}

resource "azurerm_subnet" "domain_controllers" {
  name                 = "domain-controllers"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_network_interface" "domain_controller" {
  name                = "acctestnic-%[2]d-dc"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    private_ip_address_allocation = "Static"
    private_ip_address            = "10.0.1.4"
    subnet_id                     = azurerm_subnet.domain_controllers.id
  }
}

resource "azurerm_windows_virtual_machine" "domain_controller" {
  name                = local.domain_controller_vm_name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  size                = "Standard_F2"
  admin_username      = local.admin_username
  admin_password      = local.admin_password
  custom_data         = local.custom_data

  network_interface_ids = [
    azurerm_network_interface.domain_controller.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  additional_unattend_content {
    content = local.auto_logon_data
    setting = "AutoLogon"
  }

  additional_unattend_content {
    content = local.first_logon_data
    setting = "FirstLogonCommands"
  }
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (r MsSqlVirtualMachineAvailabilityGroupListenerResource) setDomainUser(data acceptance.TestData) string {
	return fmt.Sprintf(
		`
%[1]s

resource "azurerm_virtual_machine_extension" "custom_script" {
  name                 = "set-ad-user"
  virtual_machine_id   = azurerm_windows_virtual_machine.domain_controller.id
  publisher            = "Microsoft.Compute"
  type                 = "CustomScriptExtension"
  type_handler_version = "1.9"
  settings             = <<SETTINGS
  {
    "commandToExecute": "powershell.exe -Command \"Get-ADUser ${local.admin_username} | Set-ADUser -UserPrincipalName ${local.admin_username}@${local.active_directory_domain_name}\""
  }
SETTINGS
}
`, r.domainDependencies(data),
	)
}
