package monitor_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type MonitorDiagnosticCategoriesDataSource struct {
}

func TestAccDataSourceMonitorDiagnosticCategories_appService(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_diagnostic_categories", "test")
	r := MonitorDiagnosticCategoriesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.appService(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("metrics.#").Exists(),
				check.That(data.ResourceName).Key("logs.#").Exists(),
			),
		},
	})
}

func TestAccDataSourceMonitorDiagnosticCategories_storageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_diagnostic_categories", "test")
	r := MonitorDiagnosticCategoriesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.storageAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("metrics.#").Exists(),
				check.That(data.ResourceName).Key("logs.#").Exists(),
			),
		},
	})
}

func (MonitorDiagnosticCategoriesDataSource) appService(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}

data "azurerm_monitor_diagnostic_categories" "test" {
  resource_id = azurerm_app_service.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (MonitorDiagnosticCategoriesDataSource) storageAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

data "azurerm_monitor_diagnostic_categories" "test" {
  resource_id = azurerm_storage_account.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
