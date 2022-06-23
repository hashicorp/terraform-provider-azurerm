package dns_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DnsAAAARecordDataSource struct{}

func TestAccDataSourceDnsAAAARecord_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dns_aaaa_record", "test")
	r := DnsAAAARecordDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("zone_name").Exists(),
				check.That(data.ResourceName).Key("records").Exists(),
				check.That(data.ResourceName).Key("ttl").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("tags").Exists(),
				check.That(data.ResourceName).Key("target_resource_id").DoesNotExist(),
			),
		},
	})
}

func (DnsAAAARecordDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_dns_aaaa_record" "test" {
  name                = azurerm_dns_aaaa_record.test.name
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
}
`, DnsAAAARecordResource{}.basic(data))
}
