// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package subscription_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SubscriptionDataSource struct{}

func TestAccDataSourceSubscription_current(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_subscription", "current")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SubscriptionDataSource{}.currentConfig(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("subscription_id").HasValue(data.Client().SubscriptionID),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("tenant_id").Exists(),
				check.That(data.ResourceName).Key("state").HasValue("Enabled"),
			),
		},
	})
}

func TestAccDataSourceSubscription_specific(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_subscription", "specific")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SubscriptionDataSource{}.specificConfig(data.Client().SubscriptionID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("subscription_id").HasValue(data.Client().SubscriptionID),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("tenant_id").Exists(),
				check.That(data.ResourceName).Key("location_placement_id").Exists(),
				check.That(data.ResourceName).Key("quota_id").Exists(),
				check.That(data.ResourceName).Key("spending_limit").Exists(),
			),
		},
	})
}

func (d SubscriptionDataSource) currentConfig() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}
`
}

func (d SubscriptionDataSource) specificConfig(subscriptionId string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  subscription_id = "%s"
}

data "azurerm_subscription" "specific" {
  subscription_id = "%s"
}
`, subscriptionId, subscriptionId)
}
