package applicationinsights_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type AppInsightsDataSource struct {
}

func TestAccDataSourceApplicationInsights_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_application_insights", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AppInsightsDataSource{}.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("instrumentation_key").Exists(),
				check.That(data.ResourceName).Key("app_id").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("application_type").HasValue("other"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.foo").HasValue("bar"),
			),
		},
	})
}

func (AppInsightsDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appinsights-%[1]d"
  location = "%[2]s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "other"

  tags = {
    "foo" = "bar"
  }
}

data "azurerm_application_insights" "test" {
  resource_group_name = azurerm_resource_group.test.name
  name                = azurerm_application_insights.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}
