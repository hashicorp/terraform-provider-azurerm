// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package vmware_test

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

func TestAccVmwarePrivateCloud_list_basic(t *testing.T) {
	r := VmwarePrivateCloudResource{}
	listResourceAddress := "azurerm_vmware_private_cloud.list"

	data := acceptance.BuildTestData(t, "azurerm_vmware_private_cloud", "test")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 2),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryByResourceGroupName(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 2),
				},
			},
		},
	})
}

func (r VmwarePrivateCloudResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  disable_correlation_request_id = true
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-Vmware-%[1]d"
  location = "%[2]s"
}

resource "azurerm_vmware_private_cloud" "test" {
  count = 2

  name                = "acctest-PC${count.index}-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "av36"

  management_cluster {
    size = 3
  }

  network_subnet_cidr = "192.168.${count.index * 4 + 48}.0/22"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r VmwarePrivateCloudResource) basicQuery() string {
	return `
list "azurerm_vmware_private_cloud" "list" {
  provider = azurerm
  config {}
}
`
}

func (r VmwarePrivateCloudResource) basicQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_vmware_private_cloud" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestRG-Vmware-%[1]d"
  }
}
`, data.RandomInteger)
}
