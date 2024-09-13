// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

type AppServiceEnvironmentDataSource struct{}

func TestAccDataSourceAppServiceEnvironment_basic(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skip("skipping as removed in 4.0")
	}
	data := acceptance.BuildTestData(t, "data.azurerm_app_service_environment", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: AppServiceEnvironmentDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("front_end_scale_factor").Exists(),
				check.That(data.ResourceName).Key("pricing_tier").Exists(),
				check.That(data.ResourceName).Key("service_ip_address").Exists(),
				check.That(data.ResourceName).Key("cluster_setting.#").HasValue("2"),
			),
		},
	})
}

func (d AppServiceEnvironmentDataSource) basic(data acceptance.TestData) string {
	config := AppServiceEnvironmentResource{}.clusterSettings(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service_environment" "test" {
  name                = azurerm_app_service_environment.test.name
  resource_group_name = azurerm_app_service_environment.test.resource_group_name
}
`, config)
}
