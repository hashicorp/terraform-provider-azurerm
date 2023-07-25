// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type ContainerAppDataSource struct{}

func TestAccContainerAppDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_app", "test")
	r := ContainerAppDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  acceptance.ComposeTestCheckFunc(),
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
