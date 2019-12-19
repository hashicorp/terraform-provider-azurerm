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

func TestAccAzureRMSqlVirtualMachine_basic(t *testing.T) {
	resourceName := "azurerm_sql_virtual_machine.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSqlVirtualMachine_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlVirtualMachineExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSqlVirtualMachine_complete(t *testing.T) {
	resourceName := "azurerm_sql_virtual_machine.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSqlVirtualMachine_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sql_server_license_type", "PAYG"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"wsfc_domain_credentials.0.cluster_bootstrap_account_password,wsfc_domain_credentials.0.cluster_operator_account_password,wsfc_domain_credentials.0.sql_service_account_password,auto_backup_settings.0.password,server_configurations_management_settings.0.sql_connectivity_auth_update_password"},
			},
		},
	})
}

func TestAccAzureRMSqlVirtualMachine_update(t *testing.T) {
	resourceName := "azurerm_sql_virtual_machine.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSqlVirtualMachine_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
				),
			},
			{
				Config: testAccAzureRMSqlVirtualMachine_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(resourceName, "account_tier", "Standard"),
				),
			},
		},
	})
}

func testCheckAzureRMSqlVirtualMachineExists(resourceName string) resource.TestCheckFunc {
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

func testCheckAzureRMSqlVirtualMachineDestroy(s *terraform.State) error {
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

func testAccAzureRMVirtualMachine(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.test.id}"]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  storage_os_disk {
    name          = "myosdisk1"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
  }

  os_profile {
    computer_name  = "winhost01"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_windows_config {
    timezone = "Pacific Standard Time"
    provision_vm_agent = true
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMStorageAccount(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%d"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "RA-GRS"
  is_hns_enabled           = true
}
`, rInt)
}

func testAccAzureRMSqlVirtualMachine_basic(rInt int, location string) string {
	vmconfig := testAccAzureRMVirtualMachine(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_virtual_machines" "test" {
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  virtual_machine_resource_id = azurerm_virtual_machine.test.id
  sql_server_license_type     = "PAYG"
}
`, vmconfig)
}

func testAccAzureRMSqlVirtualMachine_complete(rInt int, location string) string {
	vmconfig := testAccAzureRMVirtualMachine(rInt, location)
	//storageconfig := testAccAzureRMStorageAccount(rInt)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_virtual_machines" "test" {
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  virtual_machine_resource_id = azurerm_virtual_machine.test.id
  sql_server_license_type     = "PAYG"
  sql_management              = "Full"
  sql_image_sku               = "Enterprise"
  wsfc_domain_credentials {
    cluster_bootstrap_account_password = "<password>"
    cluster_operator_account_password  = "<password>"
    sql_service_account_password       = "<password>"
  }
  auto_patching_settings {
    day_of_week                      = "Sunday"
    enable                           = true
    maintenance_window_duration      = 60
    maintenance_window_starting_hour = 2
  }
  auto_backup_settings {
    backup_schedule_type     = "Manual"
    backup_system_dbs        = true
    enable                   = true
    enable_encryption        = true
    full_backup_frequency    = "Daily"
    full_backup_start_time   = 6
    full_backup_window_hours = 11
    log_backup_frequency     = 10
    password                 = "<Password>"
    retention_period         = 17
    storage_access_key       = azurerm_storage_account.test.primary_access_key
    storage_account_url      = azurerm_storage_account.test.primary_blob_endpoint
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
    sql_storage_disk_configuration_type    = "NEW"
    sql_storage_disk_count                 = 1
    sql_storage_starting_device_id         = 2
    sql_workload_type                      = "OLTP"
  }
}
`, vmconfig)
}
