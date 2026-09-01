// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package monitor_test

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

func TestAccMonitorActionGroup_listBySubscriptionAndRG(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")
	r := MonitorActionGroupResource{}

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
				Config: r.basicListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_monitor_action_group.list", 3),
				},
			},
			{
				Query:  true,
				Config: r.basicListQueryByResourceGroupName(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_monitor_action_group.list", 3),
				},
			},
		},
	})
}

func (MonitorActionGroupResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_monitor_action_group" "test" {
  count               = 3
  name                = "acctestActionGroup${count.index}-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r MonitorActionGroupResource) basicListQuery() string {
	return `
list "azurerm_monitor_action_group" "list" {
  provider = azurerm
  config {}
}
`
}

func (r MonitorActionGroupResource) basicListQueryByResourceGroupName() string {
	return `
list "azurerm_monitor_action_group" "list" {
  provider = azurerm
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}
