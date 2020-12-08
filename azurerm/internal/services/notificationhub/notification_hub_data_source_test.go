package notificationhub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

type NotificationHubDataSource struct{}

func TestAccNotificationHubDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_notification_hub", "test")
	d := NotificationHubDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: d.basic(data),
			Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "apns_credential.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "gcm_credential.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
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
