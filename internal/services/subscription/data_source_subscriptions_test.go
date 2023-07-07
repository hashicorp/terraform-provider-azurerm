// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package subscription_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SubscriptionsDataSource struct{}

func TestAccDataSourceSubscriptions_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_subscriptions", "current")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SubscriptionsDataSource{}.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("subscriptions.0.id").Exists(),
				check.That(data.ResourceName).Key("subscriptions.0.subscription_id").Exists(),
				check.That(data.ResourceName).Key("subscriptions.0.display_name").Exists(),
				check.That(data.ResourceName).Key("subscriptions.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("subscriptions.0.state").Exists(),
			),
		},
	})
}

func (d SubscriptionsDataSource) basic() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_subscriptions" "current" {}
`
}
