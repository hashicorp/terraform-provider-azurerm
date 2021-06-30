package notificationhub_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

type NotificationHubNamespaceDataSource struct{}

func TestAccNotificationHubNamespaceDataSource_free(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_notification_hub_namespace", "test")
	d := NotificationHubNamespaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.free(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestCheckResourceAttr(data.ResourceName, "namespace_type", "NotificationHub"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Free"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
			),
		},
	})
}

func (d NotificationHubNamespaceDataSource) free(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_notification_hub_namespace" "test" {
  name                = azurerm_notification_hub_namespace.test.name
  resource_group_name = azurerm_notification_hub_namespace.test.resource_group_name
}
`, NotificationHubNamespaceResource{}.free(data))
}
