package tests

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func testCheckAzureRMMsSqlVirtualMachineExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Sql Virtual Machine not found: %s", resourceName)
		}

		resourceGroupName := rs.Primary.Attributes["resource_group_name"]
		id, _ := azure.ParseAzureResourceID(rs.Primary.Attributes["virtual_machine_resource_id"])
		name := id.Path["virtualMachines"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQLVM.SQLVirtualMachinesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, resourceGroupName, name, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q) does not exist", name, resourceGroupName)
			}
			return fmt.Errorf("Bad: Get on sqlVirtualMachinesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMMsSqlVirtualMachineDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQLVM.SQLVirtualMachinesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_sql_virtual_machine" {
			continue
		}

		resourceGroupName := rs.Primary.Attributes["resource_group_name"]
		id, _ := azure.ParseAzureResourceID(rs.Primary.Attributes["virtual_machine_resource_id"])
		name := id.Path["virtualMachines"]

		if resp, err := client.Get(ctx, resourceGroupName, name, ""); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on sqlVirtualMachinesClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func TestAccAzureRMMsSqlVirtualMachine_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlVirtualMachine_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlVirtualMachine_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlVirtualMachine_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlVirtualMachine_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sql_license_type", "PAYG"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlVirtualMachine_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlVirtualMachine_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sql_license_type", "PAYG"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMVirtualMachine_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssqlvm-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "StorageV2"
  account_tier             = "Premium"
  account_replication_type = "LRS"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VN-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                      = "acctest-SN-%[1]d"
  resource_group_name       = azurerm_resource_group.test.name
  virtual_network_name      = azurerm_virtual_network.test.name
  address_prefix            = "10.0.0.0/24"
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_public_ip" "vm" {
  name                = "acctest-PIP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_network_security_group" "nsg" {
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
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  network_security_group_name = azurerm_network_security_group.nsg.name
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
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  network_security_group_name = azurerm_network_security_group.nsg.name
}

resource "azurerm_network_interface" "test" {
  name                      = "acctest-NIC-%[1]d"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  network_security_group_id = azurerm_network_security_group.nsg.id

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
  vm_size               = "Standard_DS13"

  storage_image_reference {
    publisher = "MicrosoftSQLServer"
    offer     = "SQL2019-WS2019"
    sku       = "SQLDEV"
    version   = "latest"
  }

  storage_os_disk {
    name          = "acctvm-%[1]dOSDisk"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}vhds/acctvm-%[1]dOSDisk.vhd"
    caching       = "ReadOnly"
    create_option = "FromImage"
  }

  storage_data_disk {
    name          = "datadisk1"
    create_option = "Empty"
    disk_size_gb  = "1"
    lun           = 0
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}vhds/datadisk1.vhd"
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

func testAccAzureRMMsSqlVirtualMachine_basic(data acceptance.TestData) string {
	vmconfig := testAccAzureRMVirtualMachine_template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_mssql_virtual_machine" "test" {
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  virtual_machine_resource_id = azurerm_virtual_machine.test.id
  sql_license_type     = "PAYG"
}
`, vmconfig)
}

func testAccAzureRMMsSqlVirtualMachine_complete(data acceptance.TestData) string {
	vmconfig := testAccAzureRMVirtualMachine_template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_mssql_virtual_machine" "test" {
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  virtual_machine_resource_id = azurerm_virtual_machine.test.id
  sql_license_type            = "PAYG"
  sql_sku                     = "Developer"

   auto_patching {
    day_of_week                            = "Sunday"
    enable                                 = true
    maintenance_window_duration_in_minutes = 60
    maintenance_window_starting_hour       = 2
  }

  key_vault_credential {
    enable = false
  }
  server_configuration {
    is_r_services_enabled                  = false
    sql_connectivity_type                  = "PRIVATE"
    sql_connectivity_port                  = 1433
    sql_connectivity_update_password  = "<password>"
    sql_connectivity_update_user_name = "sqllogin"
  }
}
`, vmconfig)
}
