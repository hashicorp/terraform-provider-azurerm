package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMMsSqlVirtualMachine_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSqlVirtualMachine_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_machine_resource_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "sql_license_type"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "auto_patching.0.enable"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "key_vault_credential.0.enable"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_configuration.0.is_r_services_enabled"),
				),
			},
		},
	})
}

func testAccDataSourceSqlVirtualMachine_basic(data acceptance.TestData) string {
	config := testAccAzureRMMsSqlVirtualMachine_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_mssql_virtual_machine" "test" {
  resource_group_name      = "${azurerm_mssql_virtual_machine.test.resource_group_name}"
  name                     = "${azurerm_mssql_virtual_machine.test.name}"
}
`, config)
}
