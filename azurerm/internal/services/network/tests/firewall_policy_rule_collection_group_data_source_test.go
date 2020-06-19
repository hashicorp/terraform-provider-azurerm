package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceFirewallPolicyRuleCollectionGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_firewall_policy_rule_collection_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFirewallPolicyRuleCollectionGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourceFirewallPolicyRuleCollectionGroup_basic(data acceptance.TestData) string {
	config := testAccAzureRMFirewallPolicyRuleCollectionGroup_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_firewall_policy_rule_collection_group" "test" {
  name               = azurerm_firewall_policy_rule_collection_group.test.name
  firewall_policy_id = azurerm_firewall_policy_rule_collection_group.test.firewall_policy_id
}
`, config)
}
