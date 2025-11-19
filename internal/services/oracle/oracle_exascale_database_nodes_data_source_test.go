// Copyright © 2025, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ExascaleDatabaseNodesDataSource struct{}

func TestExascaleDatabaseNodesDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_exascale_database_nodes", "test")
	r := ExascaleDatabaseNodesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("exascale_database_nodes.0.cpu_core_count").Exists(),
			),
		},
	})
}

func (d ExascaleDatabaseNodesDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_exascale_database_nodes" "test" {
  exascale_database_virtual_machine_cluster_id = azurerm_oracle_exascale_database_virtual_machine_cluster.test.id
}
`, ExascaleDatabaseVirtualMachineClusterResource{}.basic(data))
}
