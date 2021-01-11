package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMWebApplicationFirewallPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_web_application_firewall_policy", "test")
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMWebApplicationFirewallPolicyBasic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctestwafpolicy-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.env", "test"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMWebApplicationFirewallPolicyBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_web_application_firewall_policy" "test" {
  resource_group_name = azurerm_resource_group.test.name
  name                = azurerm_web_application_firewall_policy.test.name

}
`, testAccAzureRMWebApplicationFirewallPolicy_complete(data))
}
