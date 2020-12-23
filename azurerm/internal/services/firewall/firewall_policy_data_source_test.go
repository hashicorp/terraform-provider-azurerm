package firewall_test

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type FirewallPolicyDataSource struct {
}

func TestAccFirewallPolicyDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_firewall_policy", "test")
	r := FirewallPolicyDataSource{}
	dataParent := acceptance.BuildTestData(t, "data.azurerm_firewall_policy", "test-parent")

	data.DataSourceTest(t, []resource.TestStep{
			{
				Config: r.basic(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Key( "name").Exists(),
					check.That(data.ResourceName).Key( "resource_group_name").Exists(),
					check.That(data.ResourceName).Key( "location").HasValue( location.Normalize(data.Locations.Primary)),
					check.That(data.ResourceName).Key( "base_policy_id").Exists(),
					resource.TestCheckResourceAttr(dataParent.ResourceName, "child_policies.#", "1"),
					check.That(data.ResourceName).Key( "dns.0.proxy_enabled").HasValue( "true"),
					check.That(data.ResourceName).Key( "dns.0.servers.#").HasValue( "2"),
					check.That(data.ResourceName).Key( "threat_intelligence_mode").HasValue( string(network.AzureFirewallThreatIntelModeAlert)),
					check.That(data.ResourceName).Key( "threat_intelligence_allowlist.0.ip_addresses.#").HasValue( "2"),
					check.That(data.ResourceName).Key( "threat_intelligence_allowlist.0.fqdns.#").HasValue( "2"),
				),
			},
	})
}

func (FirewallPolicyDataSource) basic(data acceptance.TestData) string {
	// We deliberately set add a dependency between "data.azurerm_firewall_policy.test-parent"
	// and "azurerm_firewall_policy.test" so that we can test "data.azurerm_firewall_policy.test-parent.child_policies"
	return fmt.Sprintf(`
%s

data "azurerm_firewall_policy" "test-parent" {
  name                = azurerm_firewall_policy.test-parent.name
  resource_group_name = azurerm_firewall_policy.test.resource_group_name
}

data "azurerm_firewall_policy" "test" {
  name                = azurerm_firewall_policy.test.name
  resource_group_name = azurerm_firewall_policy.test.resource_group_name
}
`, FirewallPolicyResource{}.inherit(data))
}
