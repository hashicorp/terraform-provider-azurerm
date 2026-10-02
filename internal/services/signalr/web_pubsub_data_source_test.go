// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package signalr_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type WebPubsubDataSource struct{}

func TestAccDataSourceWebPubsub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_web_pubsub", "test")
	r := WebPubsubDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku").HasValue("Standard_S1"),
				check.That(data.ResourceName).Key("capacity").HasValue("1"),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("local_auth_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("aad_auth_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("live_trace.#").HasValue("1"),
				check.That(data.ResourceName).Key("live_trace.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("live_trace.0.connectivity_logs_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("live_trace.0.messaging_logs_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("live_trace.0.http_request_logs_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("hostname").Exists(),
				check.That(data.ResourceName).Key("public_port").Exists(),
				check.That(data.ResourceName).Key("server_port").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
	})
}

func TestAccDataSourceWebPubsub_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_web_pubsub", "test")
	r := WebPubsubDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.identity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
			),
		},
	})
}

func (r WebPubsubDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-wps-%[1]d"
  location = "%[2]s"
}

resource "azurerm_web_pubsub" "test" {
  name                = "acctestWebPubsub-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku      = "Standard_S1"
  capacity = 1

  public_network_access_enabled = true

  live_trace {
    enabled                   = true
    connectivity_logs_enabled = false
    messaging_logs_enabled    = true
    http_request_logs_enabled = false
  }

  local_auth_enabled = true
  aad_auth_enabled   = true

}

data "azurerm_web_pubsub" "test" {
  name                = azurerm_web_pubsub.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r WebPubsubDataSource) identity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_web_pubsub" "test" {
  name                = azurerm_web_pubsub.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, WebPubsubResource{}.systemAssignedIdentity(data))
}
