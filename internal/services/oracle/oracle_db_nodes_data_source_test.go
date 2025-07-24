// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DBNodesDataSource struct{}

func TestDBNodesDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_db_nodes", "test")
	r := DBNodesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("db_nodes.0.cpu_core_count").Exists(),
			),
		},
	})
}

func (d DBNodesDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_db_nodes" "test" {
  cloud_vm_cluster_id = azurerm_oracle_cloud_vm_cluster.test.id
}
`, CloudVmClusterResource{}.basic(data))
}
