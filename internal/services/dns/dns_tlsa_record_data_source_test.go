// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dns_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DnsTLSARecordDataSource struct{}

func TestAccDataSourceDnsTLSARecord_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dns_tlsa_record", "test")
	r := DnsTLSARecordDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("dns_zone_id").Exists(),
				check.That(data.ResourceName).Key("record.#").HasValue("2"),
				check.That(data.ResourceName).Key("ttl").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func (DnsTLSARecordDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_dns_tlsa_record" "test" {
  name        = azurerm_dns_tlsa_record.test.name
  dns_zone_id = azurerm_dns_zone.test.id
}
`, DnsTLSARecordResource{}.basic(data))
}
