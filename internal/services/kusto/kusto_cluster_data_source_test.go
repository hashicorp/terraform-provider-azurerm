// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccKustoClusterDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kusto_cluster", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: testAccDataSourceKustoCluster_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(KustoClusterResource{}),
				check.That(data.ResourceName).Key("uri").IsSet(),
				check.That(data.ResourceName).Key("data_ingestion_uri").IsSet(),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsSet(),
			),
		},
	})
}

func testAccDataSourceKustoCluster_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kusto_cluster" "test" {
  name                = azurerm_kusto_cluster.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, KustoClusterResource{}.identitySystemAssigned(data))
}
