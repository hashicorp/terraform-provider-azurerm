// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dashboard_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DashboardGrafanaDataSource struct{}

func TestAccDashboardGrafanaDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dashboard_grafana", "test")
	r := DashboardGrafanaDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("grafana_major_version").HasValue("11"),
			),
		},
	})
}

func (DashboardGrafanaDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_dashboard_grafana" "test" {
  name                = azurerm_dashboard_grafana.test.name
  resource_group_name = azurerm_dashboard_grafana.test.resource_group_name
}
`, DashboardGrafanaResource{}.basic(data))
}
