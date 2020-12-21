package firewall_test

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccFirewallPolicyDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_firewall_policy", "test")
	dataParent := acceptance.BuildTestData(t, "data.azurerm_firewall_policy", "test-parent")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallDataSourcePolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "location", location.Normalize(data.Locations.Primary)),
					resource.TestCheckResourceAttrSet(data.ResourceName, "base_policy_id"),
					resource.TestCheckResourceAttr(dataParent.ResourceName, "child_policies.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "dns.0.proxy_enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "dns.0.servers.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "threat_intelligence_mode", string(network.AzureFirewallThreatIntelModeAlert)),
					resource.TestCheckResourceAttr(data.ResourceName, "threat_intelligence_allowlist.0.ip_addresses.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "threat_intelligence_allowlist.0.fqdns.#", "2"),
				),
			},
		},
	})
}

func testAccFirewallDataSourcePolicy_basic(data acceptance.TestData) string {
	config := testAccFirewallPolicy_inherit(data)

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
`, config)
}
