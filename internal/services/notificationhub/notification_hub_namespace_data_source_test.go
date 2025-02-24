// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package notificationhub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type NotificationHubNamespaceDataSource struct{}

func TestAccNotificationHubNamespaceDataSource_free(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_notification_hub_namespace", "test")
	d := NotificationHubNamespaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.free(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Free"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "zone_redundant", "false"),
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

func TestAccNotificationHubNamespaceDataSource_zoneRedundancy(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_notification_hub_namespace", "test")
	d := NotificationHubNamespaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.free(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Free"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "zone_redundant", "true"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
			),
		},
	})
}
