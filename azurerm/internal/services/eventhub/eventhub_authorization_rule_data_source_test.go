package eventhub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccEventHubAuthorizationRuleDataSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventhub_authorization_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckEventHubAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEventHubAuthorizationRuleDataSource_base(data, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckEventHubAuthorizationRuleExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "namespace_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "eventhub_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
		},
	})
}

func TestAccEventHubAuthorizationRuleDataSource_withAliasConnectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventhub_authorization_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckEventHubAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEventHubAuthorizationRuleDataSource_withAliasConnectionString(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckEventHubAuthorizationRuleExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string_alias"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string_alias"),
				),
			},
		},
	})
}

func testAccEventHubAuthorizationRuleDataSource_base(data acceptance.TestData, listen, send, manage bool) string {
	template := testAccEventHubAuthorizationRule_base(data, listen, send, manage)
	return fmt.Sprintf(`
%s

data "azurerm_eventhub_authorization_rule" "test" {
  name                = azurerm_eventhub_authorization_rule.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func testAccEventHubAuthorizationRuleDataSource_withAliasConnectionString(data acceptance.TestData) string {
	template := testAccEventHubAuthorizationRule_withAliasConnectionString(data)
	return fmt.Sprintf(`
%s

data "azurerm_eventhub_authorization_rule" "test" {
  name                = azurerm_eventhub_authorization_rule.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
