// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dns_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DnsCaaRecordDataSource struct{}

func TestAccDataSourceDnsCaaRecord_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dns_caa_record", "test")
	r := DnsCaaRecordDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("zone_name").Exists(),
				check.That(data.ResourceName).Key("record.#").HasValue("4"),
				check.That(data.ResourceName).Key("ttl").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func (DnsCaaRecordDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_dns_caa_record" "test" {
  name                = azurerm_dns_caa_record.test.name
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
}
`, DnsCaaRecordResource{}.basic(data))
}
