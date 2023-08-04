// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatedns_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type PrivateDnsSoaRecordDataSource struct{}

func TestAccDataSourcePrivateDnsSoaRecord_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_soa_record", "test")
	r := PrivateDnsSoaRecordDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicWithDataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("zone_name").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("name").HasValue("@"),
				check.That(data.ResourceName).Key("email").HasValue("testemail.com"),
				check.That(data.ResourceName).Key("host_name").HasValue("azureprivatedns.net"),
				check.That(data.ResourceName).Key("expire_time").HasValue("2419200"),
				check.That(data.ResourceName).Key("minimum_ttl").HasValue("10"),
				check.That(data.ResourceName).Key("refresh_time").HasValue("3600"),
				check.That(data.ResourceName).Key("retry_time").HasValue("300"),
				check.That(data.ResourceName).Key("serial_number").HasValue("1"),
				check.That(data.ResourceName).Key("ttl").HasValue("3600"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func (PrivateDnsSoaRecordDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name

  soa_record {
    email = "testemail.com"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (d PrivateDnsSoaRecordDataSource) basicWithDataSource(data acceptance.TestData) string {
	config := d.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_private_dns_soa_record" "test" {
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
}
`, config)
}
