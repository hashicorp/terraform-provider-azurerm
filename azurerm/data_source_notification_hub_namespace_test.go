package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMNotificationHubNamespace_free(t *testing.T) {
	dataSourceName := "data.azurerm_notification_hub_namespace.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccDataSourceAzureRMNotificationHubNamespaceFree(rInt, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
