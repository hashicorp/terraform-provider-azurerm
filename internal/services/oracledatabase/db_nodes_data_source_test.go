// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracledatabase_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"testing"
)

type DBNodesDataSource struct{}

func TestDBNodesDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracledatabase_db_nodes", "test")
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

data "azurerm_oracledatabase_db_nodes" "test" {
	resource_group_name = azurerm_resource_group.test.name
	cloud_vm_cluster_name = azurerm_oracledatabase_cloud_vm_cluster.test.name
}
`, d.template(data))
}

func (d DBNodesDataSource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

`, CloudVmClusterResource{}.basic(data))
}
