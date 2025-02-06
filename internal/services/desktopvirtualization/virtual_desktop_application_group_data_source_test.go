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

type DesktopVirtualizationApplicationGroupDataSource struct{}

func TestAccDesktopVirtualizationApplicationGroupDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_desktop_application_group", "test")
	d := DesktopVirtualizationApplicationGroupDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").IsNotEmpty(),
				check.That(data.ResourceName).Key("resource_group_name").IsNotEmpty(),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Secondary)),
				check.That(data.ResourceName).Key("type").HasValue("RemoteApp"),
				check.That(data.ResourceName).Key("host_pool_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("friendly_name").HasValue("TestAppGroup"),
				check.That(data.ResourceName).Key("description").HasValue("Acceptance Test: An application group"),
				check.That(data.ResourceName).Key("tags.Purpose").HasValue("Acceptance-Testing"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
	})
}

func (DesktopVirtualizationApplicationGroupDataSource) complete(data acceptance.TestData) string {
	template := VirtualDesktopApplicationResource{}.complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_virtual_desktop_application_group" "test" {
  name                = azurerm_virtual_desktop_application_group.test.name
  resource_group_name = azurerm_virtual_desktop_application_group.test.resource_group_name
}
`, template)
}
