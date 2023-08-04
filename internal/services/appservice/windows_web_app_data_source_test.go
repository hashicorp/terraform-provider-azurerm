// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type WindowsWebAppDataSource struct{}

func TestAccWindowsWebAppDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_windows_web_app", "test")
	d := WindowsWebAppDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
			),
		},
	})
}

func TestAccWindowsWebAppDataSource_completeAuthV2(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_windows_web_app", "test")
	d := WindowsWebAppDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.completeAuthV2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
			),
		},
	})
}

func (WindowsWebAppDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data azurerm_windows_web_app test {
  name                = azurerm_windows_web_app.test.name
  resource_group_name = azurerm_windows_web_app.test.resource_group_name
}
`, WindowsWebAppResource{}.complete(data))
}

func (WindowsWebAppDataSource) completeAuthV2(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data azurerm_windows_web_app test {
  name                = azurerm_windows_web_app.test.name
  resource_group_name = azurerm_windows_web_app.test.resource_group_name
}
`, WindowsWebAppResource{}.completeAuthV2(data))
}
