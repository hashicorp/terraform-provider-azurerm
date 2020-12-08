package notificationhub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceNotificationHubNamespace_free(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_notification_hub_namespace", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckNotificationHubNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNotificationHubNamespaceFree(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "namespace_type", "NotificationHub"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Free"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
		},
	})
}

func testAccDataSourceNotificationHubNamespaceFree(data acceptance.TestData) string {
	template := testAccNotificationHubNamespace_free(data)
	return fmt.Sprintf(`
%s

data "azurerm_notification_hub_namespace" "test" {
  name                = azurerm_notification_hub_namespace.test.name
  resource_group_name = azurerm_notification_hub_namespace.test.resource_group_name
}
`, template)
}
