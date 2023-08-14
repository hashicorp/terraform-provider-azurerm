// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SynapseWorkspaceDataSource struct{}

func TestAccDataSourceSynapseWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_synapse_workspace", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SynapseWorkspaceDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("connectivity_endpoints.%").Exists(),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
			),
		},
	})
}

func (d SynapseWorkspaceDataSource) basic(data acceptance.TestData) string {
	config := SynapseWorkspaceResource{}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_synapse_workspace" "test" {
  name                = azurerm_synapse_workspace.test.name
  resource_group_name = azurerm_synapse_workspace.test.resource_group_name
}
`, config)
}
