package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMMsSqlVirtualMachine_basic(t *testing.T) {
	dataSourceName := "data.azurerm_mssql_virtual_machine.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSqlVirtualMachine_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "location"),
					resource.TestCheckResourceAttrSet(dataSourceName, "virtual_machine_resource_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "sql_server_license_type"),
				),
			},
		},
	})
}
func TestAccDataSourceAzureRMMsSqlVirtualMachine_complete(t *testing.T) {
	dataSourceName := "data.azurerm_mssql_virtual_machine.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSqlVirtualMachine_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "location"),
					resource.TestCheckResourceAttrSet(dataSourceName, "virtual_machine_resource_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "sql_server_license_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "sql_image_sku"),
					resource.TestCheckResourceAttrSet(dataSourceName, "auto_patching_settings.0.day_of_week"),
					resource.TestCheckResourceAttrSet(dataSourceName, "auto_patching_settings.0.enable"),
					resource.TestCheckResourceAttrSet(dataSourceName, "auto_patching_settings.0.maintenance_window_duration"),
					resource.TestCheckResourceAttrSet(dataSourceName, "auto_patching_settings.0.maintenance_window_starting_hour"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_vault_credential_settings.0.enable"),
					resource.TestCheckResourceAttrSet(dataSourceName, "server_configurations_management_settings.0.is_r_services_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "server_configurations_management_settings.0.sql_connectivity_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "server_configurations_management_settings.0.sql_connectivity_port"),
					resource.TestCheckResourceAttrSet(dataSourceName, "server_configurations_management_settings.0.sql_connectivity_auth_update_password"),
					resource.TestCheckResourceAttrSet(dataSourceName, "server_configurations_management_settings.0.sql_connectivity_auth_update_user_name"),
				),
			},
		},
	})
}

func testAccDataSourceSqlVirtualMachine_basic(rInt int, location string) string {
	config := testAccAzureRMMsSqlVirtualMachine_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_sql_virtual_machine" "test" {
  resource_group           = "${azurerm_sql_virtual_machine.test.resource_group}"
  name                     = "${azurerm_sql_virtual_machine.test.sql_virtual_machine_name}"
}
`, config)
}

func testAccDataSourceSqlVirtualMachine_complete(rInt int, location string) string {
	config := testAccAzureRMMsSqlVirtualMachine_complete(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_sql_virtual_machine" "test" {
  resource_group           = "${azurerm_sql_virtual_machine.test.resource_group}"
  name                     = "${azurerm_sql_virtual_machine.test.sql_virtual_machine_name}"
}
`, config)
}
