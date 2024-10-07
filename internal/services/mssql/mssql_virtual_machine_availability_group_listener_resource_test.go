// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/availabilitygrouplisteners"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlVirtualMachineAvailabilityGroupListenerResource struct{}

func TestAccMsSqlVirtualMachineAvailabilityGroupListener_loadBalancerConfiguration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_machine_availability_group_listener", "test")
	r := MsSqlVirtualMachineAvailabilityGroupListenerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.configureDomain(data),
		},
		{
			PreConfig: func() { time.Sleep(12 * time.Minute) },
			Config:    r.setDomainUser(data),
		},
		{
			Config: r.loadBalancerConfiguration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlVirtualMachineAvailabilityGroupListener_multiSubnetIpConfiguration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_machine_availability_group_listener", "test")
	r := MsSqlVirtualMachineAvailabilityGroupListenerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.configureDomain(data),
		},
		{
			PreConfig: func() { time.Sleep(12 * time.Minute) },
			Config:    r.setDomainUser(data),
		},
		{
			Config: r.multiSubnetIpConfiguration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (MsSqlVirtualMachineAvailabilityGroupListenerResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {

	id, err := availabilitygrouplisteners.ParseAvailabilityGroupListenerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.VirtualMachinesAvailabilityGroupListenersClient.Get(ctx, *id, availabilitygrouplisteners.GetOperationOptions{Expand: utils.String("*")})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("%s does not exist", *id)
		}
		return nil, fmt.Errorf("reading %s: %v", *id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r MsSqlVirtualMachineAvailabilityGroupListenerResource) loadBalancerConfiguration(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku = "Standard"

  frontend_ip_configuration {
    name                          = "Internal"
    private_ip_address_allocation = "Static"
    private_ip_address_version    = "IPv4"
    private_ip_address            = "10.0.2.10"
    subnet_id                     = azurerm_subnet.domain_clients[0].id
  }

  lifecycle {
    ignore_changes = [
      frontend_ip_configuration
    ]
  }
}

resource "azurerm_mssql_virtual_machine_availability_group_listener" "test" {
  name                         = "acctestli-%[3]s"
  availability_group_name      = "availabilitygroup1"
  port                         = 1433
  sql_virtual_machine_group_id = azurerm_mssql_virtual_machine_group.test.id

  load_balancer_configuration {
    load_balancer_id   = azurerm_lb.test.id
    private_ip_address = "10.0.2.11"
    probe_port         = 51572
    subnet_id          = azurerm_subnet.domain_clients[0].id

    sql_virtual_machine_ids = [
      azurerm_mssql_virtual_machine.test[0].id,
      azurerm_mssql_virtual_machine.test[1].id
    ]
  }

  replica {
    sql_virtual_machine_id = azurerm_mssql_virtual_machine.test[0].id
    role                   = "Secondary"
    commit                 = "Asynchronous_Commit"
    failover_mode          = "Manual"
    readable_secondary     = "No"
  }

  replica {
    sql_virtual_machine_id = azurerm_mssql_virtual_machine.test[1].id
    role                   = "Primary"
    commit                 = "Synchronous_Commit"
    failover_mode          = "Automatic"
    readable_secondary     = "All"
  }
}
`, r.template(data, true), data.RandomInteger, data.RandomString)
}

func (r MsSqlVirtualMachineAvailabilityGroupListenerResource) multiSubnetIpConfiguration(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_virtual_machine_availability_group_listener" "test" {
  name                         = "acctestli-%[2]s"
  availability_group_name      = "default"
  port                         = 1433
  sql_virtual_machine_group_id = azurerm_mssql_virtual_machine_group.test.id

  multi_subnet_ip_configuration {
    private_ip_address     = "10.0.2.11"
    sql_virtual_machine_id = azurerm_mssql_virtual_machine.test[0].id
    subnet_id              = azurerm_subnet.domain_clients[0].id
  }

  multi_subnet_ip_configuration {
    private_ip_address     = "10.0.3.11"
    sql_virtual_machine_id = azurerm_mssql_virtual_machine.test[1].id
    subnet_id              = azurerm_subnet.domain_clients[1].id
  }

  replica {
    sql_virtual_machine_id = azurerm_mssql_virtual_machine.test[0].id
    role                   = "Primary"
    commit                 = "Synchronous_Commit"
    failover_mode          = "Automatic"
    readable_secondary     = "All"
  }

  replica {
    sql_virtual_machine_id = azurerm_mssql_virtual_machine.test[1].id
    role                   = "Secondary"
    commit                 = "Asynchronous_Commit"
    failover_mode          = "Manual"
    readable_secondary     = "Read_Only"
  }
}
`, r.template(data, false), data.RandomString)
}

func (r MsSqlVirtualMachineAvailabilityGroupListenerResource) template(data acceptance.TestData, isSingleSubnet bool) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_subnet" "domain_clients" {
  count = %[4]t ? 1 : 2

  name                 = "domain-clients-${count.index}"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.${count.index + 2}.0/24"]
}

resource "azurerm_network_interface" "client_single_subnet" {
  count = %[4]t ? 2 : 0

  name                = "acctestnic-client-${count.index}-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.domain_clients[0].id
  }
}

resource "azurerm_network_interface" "client_multi_subnet" {
  count = %[4]t ? 0 : 2

  name                = "acctestnic-client-${count.index}-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.domain_clients[count.index].id
  }

  lifecycle {
    ignore_changes = [
      ip_configuration
    ]
  }
}

resource "azurerm_availability_set" "test" {
  name                = "acctestavset-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_windows_virtual_machine" "client" {
  count = 2

  name                = "acctest-${count.index}-%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  size                = "Standard_F2"
  admin_username      = local.admin_username
  admin_password      = local.admin_password
  custom_data         = local.custom_data
  availability_set_id = azurerm_availability_set.test.id

  network_interface_ids = [
    %[4]t ? azurerm_network_interface.client_single_subnet[count.index].id : azurerm_network_interface.client_multi_subnet[count.index].id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftSQLServer"
    offer     = "SQL2019-WS2019"
    sku       = "SQLDEV"
    version   = "latest"
  }
}

resource "azurerm_virtual_machine_extension" "join_domain" {
  count = 2

  name                 = "join-domain-${count.index}"
  virtual_machine_id   = azurerm_windows_virtual_machine.client[count.index].id
  publisher            = "Microsoft.Compute"
  type                 = "JsonADDomainExtension"
  type_handler_version = "1.3"

  settings = jsonencode({
    Name    = local.active_directory_domain_name,
    OUPath  = "",
    User    = "${local.active_directory_domain_name}\\${local.admin_username}",
    Restart = "true",
    Options = "3"
  })

  protected_settings = jsonencode({
    Password = local.admin_password
  })
}

resource "azurerm_mssql_virtual_machine" "test" {
  count = 2

  virtual_machine_id           = azurerm_windows_virtual_machine.client[count.index].id
  sql_license_type             = "PAYG"
  sql_virtual_machine_group_id = azurerm_mssql_virtual_machine_group.test.id

  wsfc_domain_credential {
    cluster_bootstrap_account_password = local.admin_password
    cluster_operator_account_password  = local.admin_password
    sql_service_account_password       = local.admin_password
  }

  depends_on = [
    azurerm_virtual_machine_extension.join_domain
  ]
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_virtual_machine_group" "test" {
  name                = "acctestgr-%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sql_image_offer     = "SQL2019-WS2019"
  sql_image_sku       = "Developer"

  wsfc_domain_profile {
    fqdn = local.active_directory_domain_name

    cluster_bootstrap_account_name = "${local.admin_username}@${local.active_directory_domain_name}"
    cluster_operator_account_name  = "${local.admin_username}@${local.active_directory_domain_name}"
    sql_service_account_name       = "${local.admin_username}@${local.active_directory_domain_name}"
    storage_account_url            = azurerm_storage_account.test.primary_blob_endpoint
    storage_account_primary_key    = azurerm_storage_account.test.primary_access_key
    cluster_subnet_type            = %[4]t ? "SingleSubnet" : "MultiSubnet"
  }
}
`, r.setDomainUser(data), data.RandomInteger, data.RandomString, isSingleSubnet)
}

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
`, r.domainDependencies(data))
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
`, r.domainDependencies(data))
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
