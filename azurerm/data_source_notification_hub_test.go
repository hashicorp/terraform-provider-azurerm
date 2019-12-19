package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMNotificationHub_basic(t *testing.T) {
	dataSourceName := "data.azurerm_notification_hub.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccDataSourceAzureRMNotificationHubBasic(rInt, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNotificationHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "apns_credential.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "gcm_credential.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMNotificationHubBasic(rInt int, location string) string {
	r := testAccAzureRMNotificationHub_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_notification_hub" "test" {
  name                = "${azurerm_notification_hub.test.name}"
  namespace_name      = "${azurerm_notification_hub_namespace.test.name}"
  resource_group_name = "${azurerm_notification_hub_namespace.test.resource_group_name}"
}
`, r)
}
