// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ApiManagementSubscriptionDataSource struct{}

func TestAccDataSourceApiManagementSubscription_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management_subscription", "test")
	r := ApiManagementSubscriptionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("subscription_id").HasValue("test-subscription"),
				check.That(data.ResourceName).Key("display_name").HasValue("Test Subscription"),
				check.That(data.ResourceName).Key("allow_tracing").HasValue("true"),
				check.That(data.ResourceName).Key("state").HasValue("active"),
			),
		},
	})
}

func (ApiManagementSubscriptionDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "accTestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_api_management_subscription" "test" {
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  subscription_id     = "test-subscription"
  allow_tracing       = true
  display_name        = "Test Subscription"
  state               = "active"
}

data "azurerm_api_management_subscription" "test" {
  resource_group_name = azurerm_api_management_subscription.test.resource_group_name
  api_management_name = azurerm_api_management_subscription.test.api_management_name
  subscription_id     = azurerm_api_management_subscription.test.subscription_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
