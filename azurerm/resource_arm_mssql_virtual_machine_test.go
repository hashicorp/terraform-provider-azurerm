package azurerm

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMsSqlVirtualMachine_basic(t *testing.T) {
	resourceName := "azurerm_mssql_virtual_machine.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMsSqlVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlVirtualMachine_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlVirtualMachineExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"sql_image_sku"},
			},
			{
				ResourceName:      "azurerm_subnet.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"network_security_group_id"},
			},
		},
	})
}

func TestAccAzureRMMsSqlVirtualMachine_complete(t *testing.T) {
	resourceName := "azurerm_mssql_virtual_machine.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMsSqlVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlVirtualMachine_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sql_server_license_type", "PAYG"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sql_image_sku","server_configurations_management_settings.0.sql_connectivity_auth_update_password"},
			},
			{
				ResourceName:      "azurerm_subnet.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"network_security_group_id"},
			},
		},
	})
}

func TestAccAzureRMMsSqlVirtualMachine_withStorage(t *testing.T) {
	resourceName := "azurerm_mssql_virtual_machine.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMsSqlVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlVirtualMachine_withStorage(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sql_server_license_type", "PAYG"),
					resource.TestCheckResourceAttr(resourceName, "storage_configuration_settings.0.storage_workload_type", "OLTP"),
					resource.TestCheckResourceAttr(resourceName, "storage_configuration_settings.0.sql_data_default_file_path", "F:\\folderpath\\"),
					resource.TestCheckResourceAttr(resourceName, "storage_configuration_settings.0.sql_log_default_file_path", "G:\\folderpath\\"),
					resource.TestCheckResourceAttr(resourceName, "storage_configuration_settings.0.sql_temp_db_default_file_path", "D:\\TEMP"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sql_image_sku","server_configurations_management_settings.0.sql_connectivity_auth_update_password"},
			},
			{
				ResourceName:      "azurerm_subnet.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"network_security_group_id"},
			},
		},
	})
}


func testCheckAzureRMMsSqlVirtualMachineExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Sql Virtual Machine not found: %s", resourceName)
		}

		resourceGroupName := rs.Primary.Attributes["resource_group_name"]
		id, _ := azure.ParseAzureResourceID(rs.Primary.Attributes["virtual_machine_resource_id"])
		name := id.Path["virtualMachines"]

		client := testAccProvider.Meta().(*ArmClient).MSSQLVM.SQLVirtualMachinesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
	client := testAccProvider.Meta().(*ArmClient).MSSQLVM.SQLVirtualMachinesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAccAzureRMVirtualMachineConfig(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "StorageV2"
  account_tier             = "Premium"
  account_replication_type = "LRS"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.0.0/24"
}

resource "azurerm_public_ip" "vm" {
  name                = "acctpIP%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_network_security_group" "nsg" {
  name                = "accnsg%d"
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
  name                      = "acctni-%d"
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
  name                  = "acctvm-%d"
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
    name          = "acctvm-%dOSDisk"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}vhds/acctvm-%dOSDisk.vhd"
    caching       = "ReadOnly"
    create_option = "FromImage"
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

`, rInt, location, rInt, rInt, rInt, rInt, rInt,rInt,rInt,rInt,rInt)
}

func testAccAzureRMVirtualMachineConfig_withDisk(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "StorageV2"
  account_tier             = "Premium"
  account_replication_type = "LRS"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.0.0/24"
}

resource "azurerm_public_ip" "vm" {
  name                = "acctpIP%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_network_security_group" "nsg" {
  name                = "accnsg%d"
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
  name                      = "acctni-%d"
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
  name                  = "acctvm-%d"
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
    name          = "acctvm-%dOSDisk"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}vhds/acctvm-%dOSDisk.vhd"
    caching       = "ReadOnly"
    create_option = "FromImage"
  }

  storage_data_disk {
    name              = "acctvm-%dDataDisk"
    create_option     = "Empty"
    disk_size_gb      = "1"
    lun               = 0
    managed_disk_type = "Standard_LRS"
  }

  storage_data_disk {
    name              = "acctvm-%dDataDisk"
    create_option     = "Empty"
    disk_size_gb      = "1"
    lun               = 1
    managed_disk_type = "Standard_LRS"
  }

  storage_data_disk {
    name              = "acctvm-%dDataDisk"
    create_option     = "Empty"
    disk_size_gb      = "1"
    lun               = 2
    managed_disk_type = "Standard_LRS"
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

`, rInt, location, rInt, rInt, rInt, rInt, rInt,rInt,rInt,rInt,rInt,rInt,rInt,rInt)
}

func testAccAzureRMMsSqlVirtualMachine_basic(rInt int, location string) string {
	vmconfig := testAccAzureRMVirtualMachineConfig(rInt, location)
	return fmt.Sprintf(`
%s
resource "azurerm_mssql_virtual_machine" "test" {
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  virtual_machine_resource_id = azurerm_virtual_machine.test.id
  sql_server_license_type     = "PAYG"
}
`, vmconfig)
}

func testAccAzureRMMsSqlVirtualMachine_complete(rInt int, location string) string {
	vmconfig := testAccAzureRMVirtualMachineConfig(rInt, location)
	return fmt.Sprintf(`
%s
resource "azurerm_mssql_virtual_machine" "test" {
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  virtual_machine_resource_id = azurerm_virtual_machine.test.id
  sql_server_license_type     = "PAYG"
  sql_image_sku               = "Developer"

  auto_patching_settings {
    day_of_week                      = "Sunday"
    enable                           = true
    maintenance_window_duration      = 60
    maintenance_window_starting_hour = 2
  }

  key_vault_credential_settings {
    enable = false
  }
  server_configurations_management_settings {
    is_r_services_enabled                  = false
    sql_connectivity_type                  = "PRIVATE"
    sql_connectivity_port                  = 1433
    sql_connectivity_auth_update_password  = "<password>"
    sql_connectivity_auth_update_user_name = "sqllogin"
  }
}
`, vmconfig)
}

func testAccAzureRMMsSqlVirtualMachine_withStorage(rInt int, location string) string {
	vmconfig := testAccAzureRMVirtualMachineConfig(rInt, location)
	return fmt.Sprintf(`
%s
resource "azurerm_mssql_virtual_machine" "test" {
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  virtual_machine_resource_id = azurerm_virtual_machine.test.id
  sql_server_license_type     = "PAYG"
  sql_image_sku               = "Developer"
  storage_configuration_settings {
    storage_workload_type         = "OLTP"
    sql_data_default_file_path    = "F:\\folderpath\\"
    sql_data_luns                 = [0]
    sql_log_default_file_path     = "E:\\folderpath\\"
    sql_log_luns                  = [1]
    sql_temp_db_default_file_path = "D:\\TEMP"
  }
}

`, vmconfig)
}
