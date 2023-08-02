// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatedns_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type PrivateDnsCNameRecordDataSource struct{}

func TestAccDataSourcePrivateDnsCNameRecord_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_cname_record", "test")
	r := PrivateDnsCNameRecordDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("zone_name").Exists(),
				check.That(data.ResourceName).Key("record").HasValue("contoso.com"),
				check.That(data.ResourceName).Key("ttl").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func (PrivateDnsCNameRecordDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_private_dns_cname_record" "test" {
  name                = azurerm_private_dns_cname_record.test.name
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
}
`, PrivateDnsCNameRecordResource{}.basic(data))
}
