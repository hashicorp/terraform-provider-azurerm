// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ContainerAppEnvironmentDataSource struct{}

func TestAccContainerAppEnvironmentDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_app_environment", "test")
	r := ContainerAppEnvironmentDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("log_analytics_workspace_name").IsSet(),
				check.That(data.ResourceName).Key("location").IsSet(),
				check.That(data.ResourceName).Key("internal_load_balancer_enabled").HasValue("true"),
			),
		},
	})
}

func (d ContainerAppEnvironmentDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_app_environment" "test" {
  name                = azurerm_container_app_environment.test.name
  resource_group_name = azurerm_container_app_environment.test.resource_group_name
}

`, ContainerAppEnvironmentResource{}.complete(data))
}
