// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package authorization_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ClientConfigDataSource struct{}

func TestAccClientConfigDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_client_config", "current")
	clientData := data.Client()
	objectIdRegex := regexp.MustCompile("^[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}$")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: ClientConfigDataSource{}.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("client_id").HasValue(clientData.Default.ClientID),
				check.That(data.ResourceName).Key("tenant_id").HasValue(clientData.TenantID),
				check.That(data.ResourceName).Key("subscription_id").HasValue(clientData.SubscriptionID),
				check.That(data.ResourceName).Key("object_id").MatchesRegex(objectIdRegex),
			),
		},
	})
}

func (d ClientConfigDataSource) basic() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}
`
}
