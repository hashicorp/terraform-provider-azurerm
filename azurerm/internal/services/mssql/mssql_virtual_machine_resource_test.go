package mssql_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MsSqlVirtualMachineResource struct{}

func TestAccMsSqlVirtualMachine_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_machine", "test")
	r := MsSqlVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlVirtualMachine_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_machine", "test")
	r := MsSqlVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMsSqlVirtualMachine_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_machine", "test")
	r := MsSqlVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sql_connectivity_update_password", "sql_connectivity_update_username"),
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sql_connectivity_update_password", "sql_connectivity_update_username"),
	})
}

func TestAccMsSqlVirtualMachine_autoBackup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_machine", "test")
	r := MsSqlVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withAutoBackup(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withAutoBackupUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withAutoBackup(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlVirtualMachine_autoPatching(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_machine", "test")
	r := MsSqlVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withAutoPatching(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withAutoPatchingUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlVirtualMachine_keyVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_machine", "test")
	r := MsSqlVirtualMachineResource{}
	value, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withKeyVault(data, value),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("r_services_enabled").MatchesRegex(regexp.MustCompile("/*:acctestkv")),
			),
		},
		data.ImportStep("key_vault_credential.0.key_vault_url", "key_vault_credential.0.service_principal_name", "key_vault_credential.0.service_principal_secret"),

		{
			Config: r.withKeyVaultUpdated(data, value),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("r_services_enabled").MatchesRegex(regexp.MustCompile("/*:acctestkv2")),
			),
		},
		data.ImportStep("key_vault_credential.0.key_vault_url", "key_vault_credential.0.service_principal_name", "key_vault_credential.0.service_principal_secret"),
	})
}

func TestAccMsSqlVirtualMachine_storageConfiguration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_machine", "test")
	r := MsSqlVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.storageConfiguration(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (MsSqlVirtualMachineResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.SqlVirtualMachineID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.VirtualMachinesClient.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil, fmt.Errorf("SQL Virtual Machine %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
		}
		return nil, fmt.Errorf("reading SQL Virtual Machine %q (Resource Group %q): %v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (MsSqlVirtualMachineResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VN-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-SN-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.0.0/24"
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_public_ip" "vm" {
  name                = "acctest-PIP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctest-NSG-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_rule" "RDPRule" {
  name                        = "RDPRule"
  resource_group_name         = azurerm_resource_group.test.name
  priority                    = 1000
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = 3389
  source_address_prefix       = "167.220.255.0/25"
  destination_address_prefix  = "*"
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_network_security_rule" "MSSQLRule" {
  name                        = "MSSQLRule"
  resource_group_name         = azurerm_resource_group.test.name
  priority                    = 1001
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = 1433
  source_address_prefix       = "167.220.255.0/25"
  destination_address_prefix  = "*"
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_network_interface" "test" {
  name                = "acctest-NIC-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.vm.id
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctest-VM-%[1]d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_F2s"

  storage_image_reference {
    publisher = "MicrosoftSQLServer"
    offer     = "SQL2017-WS2016"
    sku       = "SQLDEV"
    version   = "latest"
  }

  storage_os_disk {
    name              = "acctvm-%[1]dOSDisk"
    caching           = "ReadOnly"
    create_option     = "FromImage"
    managed_disk_type = "Premium_LRS"
  }

  os_profile {
    computer_name  = "winhost01"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_windows_config {
    timezone                  = "Pacific Standard Time"
    provision_vm_agent        = true
    enable_automatic_upgrades = true
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r MsSqlVirtualMachineResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_virtual_machine" "test" {
  virtual_machine_id = azurerm_virtual_machine.test.id
  sql_license_type   = "PAYG"
}
`, r.template(data))
}

func (r MsSqlVirtualMachineResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_virtual_machine" "import" {
  virtual_machine_id = azurerm_mssql_virtual_machine.test.virtual_machine_id
  sql_license_type   = azurerm_mssql_virtual_machine.test.sql_license_type
}
`, r.basic(data))
}

func (r MsSqlVirtualMachineResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_virtual_machine" "test" {
  virtual_machine_id               = azurerm_virtual_machine.test.id
  sql_license_type                 = "PAYG"
  r_services_enabled               = true
  sql_connectivity_port            = 1433
  sql_connectivity_type            = "PRIVATE"
  sql_connectivity_update_password = "Password1234!"
  sql_connectivity_update_username = "sqllogin"
}
`, r.template(data))
}

func (r MsSqlVirtualMachineResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_virtual_machine" "test" {
  virtual_machine_id               = azurerm_virtual_machine.test.id
  sql_license_type                 = "PAYG"
  r_services_enabled               = false
  sql_connectivity_port            = 1533
  sql_connectivity_type            = "PUBLIC"
  sql_connectivity_update_password = "Password12344321!"
  sql_connectivity_update_username = "sqlloginupdate"
}
`, r.template(data))
}

func (r MsSqlVirtualMachineResource) withAutoPatching(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_virtual_machine" "test" {
  virtual_machine_id = azurerm_virtual_machine.test.id
  sql_license_type   = "PAYG"

  auto_patching {
    day_of_week                            = "Sunday"
    maintenance_window_duration_in_minutes = 60
    maintenance_window_starting_hour       = 2
  }
}
`, r.template(data))
}

func (r MsSqlVirtualMachineResource) withAutoPatchingUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_virtual_machine" "test" {
  virtual_machine_id = azurerm_virtual_machine.test.id
  sql_license_type   = "PAYG"

  auto_patching {
    day_of_week                            = "Monday"
    maintenance_window_duration_in_minutes = 90
    maintenance_window_starting_hour       = 4
  }
}
`, r.template(data))
}

func (r MsSqlVirtualMachineResource) withAutoBackup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_virtual_machine" "test" {
  virtual_machine_id = azurerm_virtual_machine.test.id
  sql_license_type   = "PAYG"

  auto_backup {
    full_backup_frequency      = "Weekly"
    retention_period           = 21
    storage_account_url        = azurerm_storage_account.test.primary_blob_endpoint
    storage_account_access_key = azurerm_storage_account.test.primary_access_key
  }
}
`, r.template(data), data.RandomString)
}

func (r MsSqlVirtualMachineResource) withAutoBackupUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_virtual_machine" "test" {
  virtual_machine_id = azurerm_virtual_machine.test.id
  sql_license_type   = "PAYG"

  auto_backup {
    full_backup_frequency      = "Daily"
    retention_period           = 14
    storage_account_url        = azurerm_storage_account.test.primary_blob_endpoint
    storage_account_access_key = azurerm_storage_account.test.primary_access_key

	backup_system_databases   = true
    backup_schedule_automated = false
    full_backup_start_hour    = 3
    full_backup_window_hours  = 4
    log_backup_frequency      = 40
  }
}
`, r.template(data), data.RandomString)
}

func (r MsSqlVirtualMachineResource) withKeyVault(data acceptance.TestData, value string) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "acckv-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "create",
      "delete",
      "get",
      "purge",
      "update",
    ]

    secret_permissions = [
      "get",
      "delete",
      "set",
    ]
  }

  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_key" "generated" {
  name         = "key-%[2]d"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azuread_application" "test" {
  name = "acctestspa%[2]d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

resource "azuread_service_principal_password" "test" {
  service_principal_id = azuread_service_principal.test.id
  value                = "%[3]s"
  end_date             = "2021-01-01T01:02:03Z"
}

resource "azurerm_mssql_virtual_machine" "test" {
  virtual_machine_id = azurerm_virtual_machine.test.id
  sql_license_type   = "PAYG"
  key_vault_credential {
    name                     = "acctestkv"
    key_vault_url            = azurerm_key_vault_key.generated.id
    service_principal_name   = azuread_service_principal.test.display_name
    service_principal_secret = azuread_service_principal_password.test.value
  }
}
`, r.template(data), data.RandomInteger, value)
}

func (r MsSqlVirtualMachineResource) withKeyVaultUpdated(data acceptance.TestData, value string) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "acckv-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "create",
      "delete",
      "get",
      "update",
    ]

    secret_permissions = [
      "get",
      "delete",
      "set",
    ]
  }

  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_key" "generated" {
  name         = "key-%[2]d"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azuread_application" "test" {
  name = "acctestspa%[2]d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

resource "azuread_service_principal_password" "test" {
  service_principal_id = azuread_service_principal.test.id
  value                = "%[3]s"
  end_date             = "2021-01-01T01:02:03Z"
}

resource "azurerm_mssql_virtual_machine" "test" {
  virtual_machine_id = azurerm_virtual_machine.test.id
  sql_license_type   = "PAYG"
  key_vault_credential {
    name                     = "acctestkv2"
    key_vault_url            = azurerm_key_vault_key.generated.id
    service_principal_name   = azuread_service_principal.test.display_name
    service_principal_secret = azuread_service_principal_password.test.value
  }
}
`, r.template(data), data.RandomInteger, value)
}

func (r MsSqlVirtualMachineResource) storageConfiguration(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_managed_disk" "test" {
  name                 = "accmd-sqlvm-%[2]d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "0"
  caching            = "None"
}

resource "azurerm_mssql_virtual_machine" "test" {
  virtual_machine_id = azurerm_virtual_machine.test.id
  sql_license_type   = "PAYG"

  storage_configuration {
    disk_type             = "NEW"
    storage_workload_type = "OLTP"

    data_settings {
      luns              = [0]
      default_file_path = "F:\\SQLData"
    }

    log_settings {
      luns              = [0]
      default_file_path = "F:\\SQLLog"
    }

    temp_db_settings {
      luns              = [0]
      default_file_path = "F:\\SQLTemp"
    }
  }
}
`, r.template(data), data.RandomInteger)
}
