// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccVirtualMachineScaleSetManagedDisk_list(t *testing.T) {
	r := VirtualMachineScaleSetManagedDiskResource{}
	listResourceAddress := "azurerm_virtual_machine_scale_set_managed_disk.list"

	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_managed_disk", "test")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.listConfig(data),
			},
			{
				Query:  true,
				Config: r.queryConfig(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 3),
				},
			},
			{
				Query:  true,
				Config: r.queryByResourceGroupConfig(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r VirtualMachineScaleSetManagedDiskResource) listConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_managed_disk" "test" {
  count = 3

  name                 = "acctestd-${count.index}-%d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  storage_account_type = "Standard_LRS"

  creation {
    option = "Empty"
  }

  disk_size_gb = 1
}
`, r.template(data), data.RandomInteger)
}

func (VirtualMachineScaleSetManagedDiskResource) queryConfig() string {
	return `
list "azurerm_virtual_machine_scale_set_managed_disk" "list" {
  provider = azurerm
  config {}
}
`
}

func (VirtualMachineScaleSetManagedDiskResource) queryByResourceGroupConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_virtual_machine_scale_set_managed_disk" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestRG-%d"
  }
}
`, data.RandomInteger)
}
