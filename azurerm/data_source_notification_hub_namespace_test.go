package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMNotificationHubNamespace_free(t *testing.T) {
	dataSourceName := "data.azurerm_notification_hub_namespace.test"
	rInt := acctest.RandInt()
	location := testLocation()
	config := testAccDataSourceAzureRMNotificationHubNamespaceFree(rInt, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNotificationHubNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "namespace_type", "NotificationHub"),
					resource.TestCheckResourceAttr(dataSourceName, "sku.0.name", "Free"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMNotificationHubNamespaceFree(rInt int, location string) string {
	r := testAccAzureRMNotificationHubNamespace_free(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_notification_hub_namespace" "test" {
  name                = "${azurerm_notification_hub_namespace.test.name}"
  resource_group_name = "${azurerm_notification_hub_namespace.test.resource_group_name}"
}
`, r)
}
