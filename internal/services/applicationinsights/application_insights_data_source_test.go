// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AppInsightsDataSource struct{}

func TestAccApplicationInsightsDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_application_insights", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: AppInsightsDataSource{}.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("instrumentation_key").Exists(),
				check.That(data.ResourceName).Key("app_id").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("workspace_id").Exists(),
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

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "other"
  workspace_id        = azurerm_log_analytics_workspace.test.id
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
