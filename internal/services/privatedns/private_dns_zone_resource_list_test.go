// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package privatedns_test

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

func TestAccNetworkProfile_list_basic(t *testing.T) {
	r := PrivateDnsZoneResource{}

	data := acceptance.BuildTestData(t, "azurerm_private_dns_zone", "test1")

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
				Query:             true,
				Config:            r.basicQuery(),
				ConfigQueryChecks: []querycheck.QueryCheck{}, // TODO
			},
			{
				Query:             true,
				Config:            r.basicQueryByResourceGroupName(data),
				ConfigQueryChecks: []querycheck.QueryCheck{}, // TODO
			},
		},
	})
}

func (r PrivateDnsZoneResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test1" {
  name                = "acctestzone1%[1]d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_zone" "test2" {
  name                = "acctestzone2%[1]d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_zone" "test3" {
  name                = "acctestzone3%[1]d.com"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r PrivateDnsZoneResource) basicQuery() string {
	return `
list "azurerm_private_dns_zone" "list" {
  provider = azurerm
  config {}
}
`
}

func (r PrivateDnsZoneResource) basicQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_private_dns_zone" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestRG-%[1]d"
  }
}
`, data.RandomInteger)
}
