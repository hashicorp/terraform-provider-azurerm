package azurerm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMFunction_basic(t *testing.T) {
	dataSourceName := "data.azurerm_function.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()
	rs := strings.ToLower(acctest.RandString(11))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFunction_basic(rInt, location, rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "trigger_url"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMFunction_wait(t *testing.T) {
	dataSourceName := "data.azurerm_function.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()
	rs := strings.ToLower(acctest.RandString(11))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFunction_basic(rInt, location, rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "trigger_url"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key"),
				),
			},
		},
	})
}

func testAccDataSourceFunction_basic(rInt int, location, storage string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
 
resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "function-releases"

  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "container"
}

resource "azurerm_storage_blob" "javazip" {
  name = "testfunc.zip"

  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"

  type   = "block"
  source = "testdata/testfunc.zip"
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  app_settings = {
    WEBSITE_RUN_FROM_PACKAGE = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}/testfunc.zip"
    FUNCTIONS_WORKER_RUNTIME = "node"
  }
}

data "azurerm_function" "test" {
  function_app_name                = "${azurerm_function_app.test.name}"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  function_name       = "testfunc"
}
`, rInt, location, storage)
}
