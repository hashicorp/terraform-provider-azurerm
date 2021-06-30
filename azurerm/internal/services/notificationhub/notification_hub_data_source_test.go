package notificationhub_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

type NotificationHubDataSource struct{}

func TestAccNotificationHubDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_notification_hub", "test")
	d := NotificationHubDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestCheckResourceAttr(data.ResourceName, "apns_credential.#", "0"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "gcm_credential.#", "0"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
			),
		},
	})
}

func (d NotificationHubDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_notification_hub" "test" {
  name                = azurerm_notification_hub.test.name
  namespace_name      = azurerm_notification_hub_namespace.test.name
  resource_group_name = azurerm_notification_hub_namespace.test.resource_group_name
}
`, NotificationHubResource{}.basic(data))
}
