// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containerapps_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ContainerAppJobDataSource struct{}

func TestAccContainerAppJobDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_app_job", "test")
	r := ContainerAppJobDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("container_app_environment_id").Exists(),
				check.That(data.ResourceName).Key("manual_trigger_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("template.#").HasValue("1"),
			),
		},
	})
}

func (d ContainerAppJobDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_app_job" "test" {
  name                = azurerm_container_app_job.test.name
  resource_group_name = azurerm_container_app_job.test.resource_group_name
}
`, ContainerAppJobResource{}.basic(data))
}
