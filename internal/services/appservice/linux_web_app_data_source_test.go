// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type LinuxWebAppDataSource struct{}

func TestAccLinuxWebAppDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_linux_web_app", "test")
	d := LinuxWebAppDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
			),
		},
	})
}

func TestAccLinuxWebAppDataSource_completeAuthV2(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_linux_web_app", "test")
	d := LinuxWebAppDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.completeAuthV2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
			),
		},
	})
}

func (LinuxWebAppDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data azurerm_linux_web_app test {
  name                = azurerm_linux_web_app.test.name
  resource_group_name = azurerm_linux_web_app.test.resource_group_name
}
`, LinuxWebAppResource{}.complete(data))
}

func (LinuxWebAppDataSource) completeAuthV2(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data azurerm_linux_web_app test {
  name                = azurerm_linux_web_app.test.name
  resource_group_name = azurerm_linux_web_app.test.resource_group_name
}
`, LinuxWebAppResource{}.completeAuthV2(data))
}
