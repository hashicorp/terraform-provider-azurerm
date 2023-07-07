// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package healthcare_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type HealthCareWorkspaceDataSource struct{}

func TestAccHealthCareWorkspaceDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_healthcare_workspace", "test")
	r := HealthCareWorkspaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists()),
		},
	})
}

func (HealthCareWorkspaceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_healthcare_workspace" "test" {
  name                = azurerm_healthcare_workspace.test.name
  resource_group_name = azurerm_healthcare_workspace.test.resource_group_name
}
`, HealthCareWorkspaceResource{}.complete(data))
}
