// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package trafficmanager_test

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

func TestAccTrafficManagerProfile_list_basic(t *testing.T) {
	r := TrafficManagerProfileResource{}
	listResourceAddress := "azurerm_traffic_manager_profile.list"

	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test1")

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
					querycheck.ExpectLengthAtLeast(listResourceAddress, 3),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryByResourceGroupName(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r TrafficManagerProfileResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-list-%[1]d"
  location = "%[2]s"
}

resource "azurerm_traffic_manager_profile" "test" {
  count = 3

  name                   = "acctest-tmprofile-${count.index}-%[1]d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Weighted"

  dns_config {
    relative_name = "acctest-tmprofile-${count.index}-%[1]d"
    ttl           = 30
  }

  monitor_config {
    protocol = "HTTPS"
    port     = 443
    path     = "/"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r TrafficManagerProfileResource) basicQuery() string {
	return `
list "azurerm_traffic_manager_profile" "list" {
  provider = azurerm
  config {}
}
`
}

func (r TrafficManagerProfileResource) basicQueryByResourceGroupName() string {
	return `
list "azurerm_traffic_manager_profile" "list" {
  provider = azurerm
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}
