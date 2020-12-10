package servicebus_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMServiceBusTopicAuthorizationRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_topic_authorization_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMServiceBusTopicAuthorizationRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicAuthorizationRuleExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "namespace_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMServiceBusTopicAuthorizationRule_basic(data acceptance.TestData) string {
	template := testAccAzureRMServiceBusTopicAuthorizationRule_base(data, true, true, true)
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_topic_authorization_rule" "test" {
  name                = azurerm_servicebus_topic_authorization_rule.test.name
  namespace_name      = azurerm_servicebus_topic_authorization_rule.test.namespace_name
  resource_group_name = azurerm_servicebus_topic_authorization_rule.test.resource_group_name
  topic_name          = azurerm_servicebus_topic_authorization_rule.test.topic_name
}
`, template)
}
