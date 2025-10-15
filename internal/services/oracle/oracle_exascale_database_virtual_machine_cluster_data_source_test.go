// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
)

type ExascaleDatabaseVirtualMachineClusterDataSource struct{}

func TestExascaleDatabaseVirtualMachineClusterDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, fmt.Sprintf("data.%[1]s", oracle.ExascaleDatabaseVirtualMachineClusterDataSource{}.ResourceType()), "test")
	r := ExascaleDatabaseVirtualMachineClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("exascale_database_storage_vault_id").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("enabled_ecpu_count").Exists(),
				check.That(data.ResourceName).Key("grid_image_ocid").Exists(),
				check.That(data.ResourceName).Key("total_ecpu_count").Exists(),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
			),
		},
	})
}

func (d ExascaleDatabaseVirtualMachineClusterDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_exascale_database_virtual_machine_cluster" "test" {
  name                = azurerm_oracle_exascale_database_virtual_machine_cluster.test.name
  resource_group_name = azurerm_oracle_exascale_database_virtual_machine_cluster.test.resource_group_name
}
`, ExadbVmClusterResource{}.basic(data))
}
