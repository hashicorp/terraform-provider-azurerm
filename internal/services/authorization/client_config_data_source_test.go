// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ClientConfigDataSource struct{}

func TestAccClientConfigDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_client_config", "current")
	clientId := os.Getenv("ARM_CLIENT_ID")
	tenantId := os.Getenv("ARM_TENANT_ID")
	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")
	objectIdRegex := regexp.MustCompile("^[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}$")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: ClientConfigDataSource{}.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("client_id").HasValue(clientId),
				check.That(data.ResourceName).Key("tenant_id").HasValue(tenantId),
				check.That(data.ResourceName).Key("subscription_id").HasValue(subscriptionId),
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
