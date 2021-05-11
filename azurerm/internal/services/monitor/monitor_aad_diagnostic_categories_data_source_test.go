package monitor_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"testing"
)

type MonitorAADDiagnosticCategoriesDataSource struct{}

func TestAccDataSourceMonitorAADDiagnosticCategories_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_aad_diagnostic_categories", "test")
	r := MonitorAADDiagnosticCategoriesDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("logs.#").Exists(),
			),
		},
	})
}

func (MonitorAADDiagnosticCategoriesDataSource) basic(_ acceptance.TestData) string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_monitor_aad_diagnostic_categories" "test" {}
`
}
