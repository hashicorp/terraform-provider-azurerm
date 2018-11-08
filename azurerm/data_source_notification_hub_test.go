package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMNotificationHub_basic(t *testing.T) {
	dataSourceName := "data.azurerm_notification_hub.test"
	rInt := acctest.RandInt()
	location := testLocation()
	config := testAccDataSourceAzureRMNotificationHubBasic(rInt, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
