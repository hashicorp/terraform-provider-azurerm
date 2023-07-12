// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package healthcare_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type HealthCareWorkspaceIotConnectorDataSource struct{}

func TestAccHealthCareWorkspaceIotConnectorDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_healthcare_medtech_service", "test")
	r := HealthCareWorkspaceIotConnectorDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists()),
		},
	})
}

func (HealthCareWorkspaceIotConnectorDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_healthcare_medtech_service" "test" {
  name         = azurerm_healthcare_medtech_service.test.name
  workspace_id = azurerm_healthcare_workspace.test.id
}
`, HealthCareWorkspaceMedTechServiceResource{}.basic(data))
}
