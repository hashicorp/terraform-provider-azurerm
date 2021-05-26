package network_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type WebApplicationFirewallDataSource struct {
}

func TestAccDataSourceAzureRMWebApplicationFirewallPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_web_application_firewall_policy", "test")
	r := WebApplicationFirewallDataSource{}
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestwafpolicy-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(resourceGroupName),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.env").HasValue("test"),
			),
		},
	})
}

func (WebApplicationFirewallDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_web_application_firewall_policy" "test" {
  resource_group_name = azurerm_resource_group.test.name
  name                = azurerm_web_application_firewall_policy.test.name

}
`, WebApplicationFirewallResource{}.complete(data))
}
