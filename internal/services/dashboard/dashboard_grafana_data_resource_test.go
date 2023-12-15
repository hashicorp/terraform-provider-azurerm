// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dashboard_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type DashboardGrafanaDataSource struct{}

func TestAccDashboardGrafanaDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dashboard_grafana", "test")
	r := DashboardGrafanaDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
	})
}

func (d DashboardGrafanaDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_dashboard_grafana" "test" {
  name                = azurerm_dashboard_grafana.test.name
  resource_group_name = azurerm_dashboard_grafana.test.resource_group_name
}
`, DashboardGrafanaResource{}.basic(data))
}
