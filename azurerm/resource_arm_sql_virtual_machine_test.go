package azurerm

import (
    "fmt"
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
                    resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
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
                    resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
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

        resourceGroupName := rs.Primary.Attributes["resource_group"]
        name := rs.Primary.Attributes["sql_virtual_machine_name"]

        client := testAccProvider.Meta().(*ArmClient).MSSQLVM.SQLVirtualMachinesClient
        ctx := testAccProvider.Meta().(*ArmClient).StopContext

        if resp, err := client.Get(ctx, resourceGroupName, name,""); err != nil {
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

        resourceGroupName := rs.Primary.Attributes["resource_group"]
        name := rs.Primary.Attributes["sql_virtual_machine_name"]

        if resp, err := client.Get(ctx, resourceGroupName, name,""); err != nil {
            if !utils.ResponseWasNotFound(resp.Response) {
                return fmt.Errorf("Bad: Get on sqlVirtualMachinesClient: %+v", err)
            }
        }

        return nil
    }

    return nil
}

func testAccAzureRMSqlVirtualMachine_basic(rInt int, location string) string {
    vmconfig :=  testAccAzureRMVirtualMachine_winTimeZone(rInt,location)
    return fmt.Sprintf(`
%s

resource "azurerm_mssql_virtual_machines" "test" {
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sql_virtual_machine_name    = "accsqlvm%d"
  virtual_machine_resource_id = azurerm_virtual_machine.test.id
  sql_server_license_type     = "PAYG"
}
`, vmconfig, rInt)
}

func testAccAzureRMSqlVirtualMachine_complete(rInt int, location string) string {
    vmconfig :=  testAccAzureRMVirtualMachine_winTimeZone(rInt,location)
    storageconfig := testAccAzureRMStorageAccount_basic(rInt,string(rInt),location)
    return fmt.Sprintf(`
%s

%s

resource "azurerm_mssql_virtual_machines" "test" {
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sql_virtual_machine_name    = "accsqlvm%d"
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
    storage_access_key       = azurerm_storage_account.testsa.primary_access_key
    storage_account_url      = "azurerm_storage_account.testsa.primary_blob_endpoint
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
`, vmconfig, storageconfig,rInt)
}
