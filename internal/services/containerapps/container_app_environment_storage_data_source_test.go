// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containerapps_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ContainerAppEnvironmentStorageDataSource struct{}

func TestAccContainerAppEnvironmentStorageDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_app_environment_storage", "test")
	r := ContainerAppEnvironmentStorageDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("account_name").IsSet(),
				check.That(data.ResourceName).Key("share_name").IsSet(),
				check.That(data.ResourceName).Key("access_mode").IsSet(),
			),
		},
	})
}

func (d ContainerAppEnvironmentStorageDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_app_environment_storage" "test" {
  name                         = azurerm_container_app_environment_storage.test.name
  container_app_environment_id = azurerm_container_app_environment_storage.test.container_app_environment_id
}
`, ContainerAppEnvironmentStorageResource{}.basic(data))
}
