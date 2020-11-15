package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMServiceBusSubscription_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMServiceBusSubscription_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "auto_delete_on_idle"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "default_message_ttl"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "lock_duration"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "dead_lettering_on_message_expiration"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "dead_lettering_on_filter_evaluation_error"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "enable_batched_operations"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "max_delivery_count"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "requires_session"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "status"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMServiceBusSubscription_basic(data acceptance.TestData) string {
	template := testAccAzureRMServiceBusSubscription_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_subscription" "test" {
  name                = azurerm_servicebus_subscription.test.name
  namespace_name      = azurerm_servicebus_subscription.test.namespace_name
  topic_name          = azurerm_servicebus_subscription.test.topic_name
  resource_group_name = azurerm_servicebus_subscription.test.resource_group_name
}
`, template)
}
