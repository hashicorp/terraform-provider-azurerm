package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceArmMonitorDiagnosticCategories_appService(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_diagnostic_categories", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArmMonitorDiagnosticCategories_appService(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "metrics.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "logs.#", "6"),
				),
			},
		},
	})
}

func TestAccDataSourceArmMonitorDiagnosticCategories_storageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_diagnostic_categories", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArmMonitorDiagnosticCategories_storageAccount(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "metrics.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "logs.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceArmMonitorDiagnosticCategories_appService(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

data "azurerm_monitor_diagnostic_categories" "test" {
  resource_id = "${azurerm_app_service.test.id}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccDataSourceArmMonitorDiagnosticCategories_storageAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

data "azurerm_monitor_diagnostic_categories" "test" {
  resource_id = "${azurerm_storage_account.test.id}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
