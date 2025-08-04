// Copyright © 2025, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ExascaleDBNodesDataSource struct{}

func TestExascaleDBNodesDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_exascale_db_nodes", "test")
	r := ExascaleDBNodesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("exascale_db_nodes.0.cpu_core_count").Exists(),
			),
		},
	})
}

func (d ExascaleDBNodesDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_exascale_db_nodes" "test" {
  exa_db_vm_cluster_id = azurerm_oracle_exa_db_vm_cluster.test.id
}
`, ExadbVmClusterResource{}.basic(data))
}
