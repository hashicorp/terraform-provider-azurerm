// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ContainerAppDataSource struct{}

func TestAccContainerAppDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_app", "test")
	r := ContainerAppDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("secret.#").Exists(),
			),
		},
	})
}

func TestAccContainerAppDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_app", "test")
	r := ContainerAppDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data, "rev1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("secret.#").Exists(),
			),
		},
	})
}

func TestAccContainerAppDataSource_readSecretsDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_app", "test")
	r := ContainerAppDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.readSecretsDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("read_secrets").HasValue("false"),
				check.That(data.ResourceName).Key("secret.#").HasValue("0"),
			),
		},
	})
}

func (d ContainerAppDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_app" "test" {
  name                = azurerm_container_app.test.name
  resource_group_name = azurerm_container_app.test.resource_group_name
}
`, ContainerAppResource{}.basic(data))
}

func (d ContainerAppDataSource) complete(data acceptance.TestData, revisionSuffix string) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_app" "test" {
  name                = azurerm_container_app.test.name
  resource_group_name = azurerm_container_app.test.resource_group_name
}
`, ContainerAppResource{}.complete(data, revisionSuffix))
}

func (d ContainerAppDataSource) readSecretsDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_app" "test" {
  name                = azurerm_container_app.test.name
  resource_group_name = azurerm_container_app.test.resource_group_name
  read_secrets        = false
}
`, ContainerAppResource{}.basic(data))
}
