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

func TestAccDedicatedHostGroup_listBySubscriptionAndRG(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host_group", "testlist1")
	r := DedicatedHostGroupResource{}
	listResourceAddress := "azurerm_dedicated_host_group.list"

	acceptance.RunTest(t, resource.TestCase{
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
					querycheck.ExpectLengthAtLeast(listResourceAddress, 3),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryByResourceGroupName(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r DedicatedHostGroupResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-compute-%[1]d"
  location = "%[2]s"
}

resource "azurerm_dedicated_host_group" "test" {
  count                       = 3
  name                        = "acctestDHG${count.index}-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  platform_fault_domain_count = 2
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r DedicatedHostGroupResource) basicQuery() string {
	return `
list "azurerm_dedicated_host_group" "list" {
  provider = azurerm
  config {}
}
`
}

func (r DedicatedHostGroupResource) basicQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_dedicated_host_group" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestRG-compute-%d"
  }
}
`, data.RandomInteger)
}
