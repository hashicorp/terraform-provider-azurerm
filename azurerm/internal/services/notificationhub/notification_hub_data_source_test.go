package notificationhub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMNotificationHub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_notification_hub", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNotificationHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMNotificationHubBasic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "apns_credential.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "gcm_credential.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMNotificationHubBasic(data acceptance.TestData) string {
	template := testAccAzureRMNotificationHub_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_notification_hub" "test" {
  name                = azurerm_notification_hub.test.name
  namespace_name      = azurerm_notification_hub_namespace.test.name
  resource_group_name = azurerm_notification_hub_namespace.test.resource_group_name
}
`, template)
}
