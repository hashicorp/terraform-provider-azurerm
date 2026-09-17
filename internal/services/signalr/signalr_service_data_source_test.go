// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package signalr_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SignalRServiceDataSource struct{}

func TestAccDataSourceSignalRService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_signalr_service", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SignalRServiceDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("hostname").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("public_port").Exists(),
				check.That(data.ResourceName).Key("server_port").Exists(),
				check.That(data.ResourceName).Key("sku.#").HasValue("1"),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_S1"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("1"),
				check.That(data.ResourceName).Key("connectivity_logs_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("messaging_logs_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("http_request_logs_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("service_mode").HasValue("Default"),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
	})
}

func TestAccDataSourceSignalRService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_signalr_service", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SignalRServiceDataSource{}.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku.#").HasValue("1"),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_S1"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("1"),
				check.That(data.ResourceName).Key("service_mode").HasValue("Serverless"),
				check.That(data.ResourceName).Key("connectivity_logs_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("messaging_logs_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("http_request_logs_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("live_trace.#").HasValue("1"),
				check.That(data.ResourceName).Key("live_trace.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("live_trace.0.connectivity_logs_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("live_trace.0.messaging_logs_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("cors.#").HasValue("1"),
				check.That(data.ResourceName).Key("cors.0.allowed_origins.#").HasValue("2"),
				check.That(data.ResourceName).Key("upstream_endpoint.#").HasValue("4"),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
			),
		},
	})
}

func (r SignalRServiceDataSource) basic(data acceptance.TestData) string {
	template := SignalRServiceResource{}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_signalr_service" "test" {
  name                = azurerm_signalr_service.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func (r SignalRServiceDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_signalr_service" "test" {
  name                = "acctestSignalR-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Standard_S1"
    capacity = 1
  }

  service_mode              = "Serverless"
  connectivity_logs_enabled = true
  messaging_logs_enabled    = false
  http_request_logs_enabled = false

  live_trace {
    enabled                   = true
    connectivity_logs_enabled = false
    messaging_logs_enabled    = true
  }

  cors {
    allowed_origins = [
      "https://example.com",
      "https://contoso.com",
    ]
  }

  upstream_endpoint {
    category_pattern = ["*"]
    event_pattern    = ["*"]
    hub_pattern      = ["*"]
    url_template     = "http://foo.com/{hub}/api/{category}/{event}"
  }

  upstream_endpoint {
    category_pattern = ["connections", "messages"]
    event_pattern    = ["*"]
    hub_pattern      = ["hub1"]
    url_template     = "http://foo.com"
  }

  upstream_endpoint {
    category_pattern = ["*"]
    event_pattern    = ["connect", "disconnect"]
    hub_pattern      = ["hub1", "hub2"]
    url_template     = "http://foo3.com"
  }

  upstream_endpoint {
    category_pattern = ["connections"]
    event_pattern    = ["disconnect"]
    hub_pattern      = ["*"]
    url_template     = "http://foo4.com"
  }

  identity {
    type = "SystemAssigned"
  }
}

data "azurerm_signalr_service" "test" {
  name                = azurerm_signalr_service.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}
