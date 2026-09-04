// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containerapps_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ContainerAppEnvironmentHttpRouteConfigDataSource struct{}

func TestAccContainerAppEnvironmentHttpRouteConfigDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_app_environment_http_route_config", "test")
	r := ContainerAppEnvironmentHttpRouteConfigDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("rules.#").HasValue("1"),
				check.That(data.ResourceName).Key("rules.0.targets.#").HasValue("1"),
				check.That(data.ResourceName).Key("rules.0.targets.0.container_app").IsNotEmpty(),
			),
		},
	})
}

func (r ContainerAppEnvironmentHttpRouteConfigDataSource) basic(data acceptance.TestData) string {
	resource := ContainerAppEnvironmentHttpRouteConfigResource{}
	return fmt.Sprintf(`
%s

data "azurerm_container_app_environment_http_route_config" "test" {
  name                         = azurerm_container_app_environment_http_route_config.test.name
  container_app_environment_id = azurerm_container_app_environment_http_route_config.test.container_app_environment_id
}
`, resource.basic(data))
}
