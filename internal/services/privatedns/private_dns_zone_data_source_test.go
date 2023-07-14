// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatedns_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type PrivateDnsZoneDatasource struct{}

func TestAccDataSourcePrivateDNSZone_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_zone", "test")
	r := PrivateDnsZoneDatasource{}

	resourceName := "azurerm_private_dns_zone.test"

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").MatchesOtherKey(check.That(resourceName).Key("id")),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourcePrivateDNSZone_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_zone", "test")
	r := PrivateDnsZoneDatasource{}

	resourceName := "azurerm_private_dns_zone.test"

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.tags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").MatchesOtherKey(check.That(resourceName).Key("id")),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.hello").HasValue("world"),
			),
		},
	})
}

func TestAccDataSourcePrivateDNSZone_withoutResourceGroupName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_zone", "test")
	r := PrivateDnsZoneDatasource{}

	resourceName := "azurerm_private_dns_zone.test"

	// This test is split across multiple test steps to avoid an API race
	// condition that occures when running multiple test cases in parallel
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.template(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(resourceName).Key("id").Exists(),
			),
		},
		{
			Config: r.withoutResourceGroupName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").MatchesOtherKey(check.That(resourceName).Key("id")),
				check.That(data.ResourceName).Key("resource_group_name").MatchesOtherKey(check.That(resourceName).Key("resource_group_name")),
			),
		},
	})
}

func (r PrivateDnsZoneDatasource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_private_dns_zone" "test" {
  name                = azurerm_private_dns_zone.test.name
  resource_group_name = azurerm_private_dns_zone.test.resource_group_name
}
`, r.template(data))
}

func (PrivateDnsZoneDatasource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%[1]d.internal"
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    hello = "world"
  }
}

data "azurerm_private_dns_zone" "test" {
  name                = azurerm_private_dns_zone.test.name
  resource_group_name = azurerm_private_dns_zone.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (PrivateDnsZoneDatasource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%[1]d.internal"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r PrivateDnsZoneDatasource) withoutResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_private_dns_zone" "test" {
  name = azurerm_private_dns_zone.test.name
}
`, r.template(data))
}
