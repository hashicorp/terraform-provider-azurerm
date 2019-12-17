package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMServiceBusNamespaceRule_basic(t *testing.T) {
	dataSourceName := "data.azurerm_servicebus_namespace_authorization_rule.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMServiceBusNamespaceAuthorizationRule_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(dataSourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secondary_connection_string"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secondary_key"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMServiceBusNamespaceAuthorizationRule_basic(rInt int, location string) string {
	template := testAccAzureRMServiceBusNamespaceAuthorizationRule_base(rInt, location, true, true, true)
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_namespace_authorization_rule" "test" {
  name                = "${azurerm_servicebus_namespace_authorization_rule.test.name}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template)
}
