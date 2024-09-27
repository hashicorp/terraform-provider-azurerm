// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DesktopVirtualizationWorkspaceDataSource struct{}

func TestAccDesktopVirtualizationWorkspaceDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_desktop_workspace", "test")
	d := DesktopVirtualizationWorkspaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").IsNotEmpty(),
				check.That(data.ResourceName).Key("description").IsNotEmpty(),
				check.That(data.ResourceName).Key("resource_group_name").IsNotEmpty(),
				check.That(data.ResourceName).Key("friendly_name").HasValue("Acceptance Test!"),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Secondary)),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func (DesktopVirtualizationWorkspaceDataSource) complete(data acceptance.TestData) string {
	template := AzureRMDesktopVirtualizationWorkspaceResource{}.complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_virtual_desktop_workspace" "test" {
  name                = azurerm_virtual_desktop_workspace.test.name
  resource_group_name = azurerm_virtual_desktop_workspace.test.resource_group_name
}
`, template)
}
