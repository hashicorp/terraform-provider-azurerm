// Copyright © 2025, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
)

type ExadbVmClusterDataSource struct{}

func TestExadbVmClusterDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ExadbVmClusterDataSource{}.ResourceType(), "test")
	r := ExadbVmClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("exascale_db_storage_vault_id").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("enabled_ecpu_count").Exists(),
				check.That(data.ResourceName).Key("grid_image_ocid").Exists(),
				check.That(data.ResourceName).Key("total_ecpu_count").Exists(),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
			),
		},
	})
}

func (d ExadbVmClusterDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_exa_db_vm_cluster" "test" {
  name                = azurerm_oracle_exa_db_vm_cluster.test.name
  resource_group_name = azurerm_oracle_exa_db_vm_cluster.test.resource_group_name
}
`, ExadbVmClusterResource{}.basic(data))
}
