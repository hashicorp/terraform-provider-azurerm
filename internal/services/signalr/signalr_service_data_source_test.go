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

func TestAccDataSourceSignalRService_cors(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_signalr_service", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SignalRServiceDataSource{}.cors(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("cors.#").HasValue("1"),
				check.That(data.ResourceName).Key("cors.0.allowed_origins.#").HasValue("2"),
			),
		},
	})
}

func TestAccDataSourceSignalRService_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_signalr_service", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SignalRServiceDataSource{}.identity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
			),
		},
	})
}

func TestAccDataSourceSignalRService_liveTrace(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_signalr_service", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SignalRServiceDataSource{}.liveTrace(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("live_trace.#").HasValue("1"),
				check.That(data.ResourceName).Key("live_trace.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("live_trace.0.connectivity_logs_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("live_trace.0.messaging_logs_enabled").HasValue("false"),
			),
		},
	})
}

func TestAccDataSourceSignalRService_resourceLogs(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_signalr_service", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SignalRServiceDataSource{}.resourceLogs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("connectivity_logs_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("messaging_logs_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("http_request_logs_enabled").HasValue("false"),
			),
		},
	})
}

func TestAccDataSourceSignalRService_serviceMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_signalr_service", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SignalRServiceDataSource{}.serviceMode(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("service_mode").HasValue("Serverless"),
			),
		},
	})
}

func TestAccDataSourceSignalRService_upstreamEndpoint(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_signalr_service", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SignalRServiceDataSource{}.upstreamEndpoint(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("upstream_endpoint.#").HasValue("4"),
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

func (r SignalRServiceDataSource) cors(data acceptance.TestData) string {
	template := SignalRServiceResource{}.withCors(data)
	return fmt.Sprintf(`
%s

data "azurerm_signalr_service" "test" {
  name                = azurerm_signalr_service.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func (r SignalRServiceDataSource) identity(data acceptance.TestData) string {
	template := SignalRServiceResource{}.systemAssignedIdentity(data)
	return fmt.Sprintf(`
%s

data "azurerm_signalr_service" "test" {
  name                = azurerm_signalr_service.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func (r SignalRServiceDataSource) liveTrace(data acceptance.TestData) string {
	template := SignalRServiceResource{}.liveTrace(data)
	return fmt.Sprintf(`
%s

data "azurerm_signalr_service" "test" {
  name                = azurerm_signalr_service.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func (r SignalRServiceDataSource) resourceLogs(data acceptance.TestData) string {
	template := SignalRServiceResource{}.resourceLogs(data)
	return fmt.Sprintf(`
%s

data "azurerm_signalr_service" "test" {
  name                = azurerm_signalr_service.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func (r SignalRServiceDataSource) serviceMode(data acceptance.TestData) string {
	template := SignalRServiceResource{}.withServiceMode(data, "Serverless")
	return fmt.Sprintf(`
%s

data "azurerm_signalr_service" "test" {
  name                = azurerm_signalr_service.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func (r SignalRServiceDataSource) upstreamEndpoint(data acceptance.TestData) string {
	template := SignalRServiceResource{}.withUpstreamEndpoints(data)
	return fmt.Sprintf(`
%s

data "azurerm_signalr_service" "test" {
  name                = azurerm_signalr_service.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
