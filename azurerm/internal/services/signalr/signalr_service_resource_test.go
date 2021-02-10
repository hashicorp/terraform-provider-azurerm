package signalr_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/signalr/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SignalRServiceResource struct{}

func TestAccSignalRService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service", "test")
	r := SignalRServiceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Free_F1"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("1"),
				check.That(data.ResourceName).Key("hostname").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("public_port").Exists(),
				check.That(data.ResourceName).Key("server_port").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSignalRService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service", "test")
	r := SignalRServiceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Free_F1"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("1"),
				check.That(data.ResourceName).Key("hostname").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("public_port").Exists(),
				check.That(data.ResourceName).Key("server_port").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSignalRService_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service", "test")
	r := SignalRServiceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.standardWithCapacity(data, 1),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_S1"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("1"),
				check.That(data.ResourceName).Key("hostname").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("public_port").Exists(),
				check.That(data.ResourceName).Key("server_port").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSignalRService_standardWithCap2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service", "test")
	r := SignalRServiceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.standardWithCapacity(data, 2),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_S1"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("2"),
				check.That(data.ResourceName).Key("hostname").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("public_port").Exists(),
				check.That(data.ResourceName).Key("server_port").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSignalRService_skuUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service", "test")
	r := SignalRServiceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Free_F1"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("1"),
				check.That(data.ResourceName).Key("hostname").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("public_port").Exists(),
				check.That(data.ResourceName).Key("server_port").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
		{
			Config: r.standardWithCapacity(data, 1),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_S1"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("1"),
				check.That(data.ResourceName).Key("hostname").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("public_port").Exists(),
				check.That(data.ResourceName).Key("server_port").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Free_F1"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("1"),
				check.That(data.ResourceName).Key("hostname").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
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

func TestAccSignalRService_capacityUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service", "test")
	r := SignalRServiceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.standardWithCapacity(data, 1),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_S1"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("1"),
				check.That(data.ResourceName).Key("hostname").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("public_port").Exists(),
				check.That(data.ResourceName).Key("server_port").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
		{
			Config: r.standardWithCapacity(data, 5),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_S1"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("5"),
				check.That(data.ResourceName).Key("hostname").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("public_port").Exists(),
				check.That(data.ResourceName).Key("server_port").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
		{
			Config: r.standardWithCapacity(data, 1),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_S1"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("1"),
				check.That(data.ResourceName).Key("hostname").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
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

func TestAccSignalRService_skuAndCapacityUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service", "test")
	r := SignalRServiceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Free_F1"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("1"),
				check.That(data.ResourceName).Key("hostname").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("public_port").Exists(),
				check.That(data.ResourceName).Key("server_port").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
		{
			Config: r.standardWithCapacity(data, 2),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_S1"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("2"),
				check.That(data.ResourceName).Key("hostname").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("public_port").Exists(),
				check.That(data.ResourceName).Key("server_port").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Free_F1"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("1"),
				check.That(data.ResourceName).Key("hostname").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
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

func TestAccSignalRService_serviceMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service", "test")
	r := SignalRServiceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withServiceMode(data, "Serverless"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("hostname").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("public_port").Exists(),
				check.That(data.ResourceName).Key("server_port").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSignalRService_cors(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service", "test")
	r := SignalRServiceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withCors(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cors.#").HasValue("1"),
				check.That(data.ResourceName).Key("cors.0.allowed_origins.#").HasValue("2"),
				check.That(data.ResourceName).Key("hostname").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("public_port").Exists(),
				check.That(data.ResourceName).Key("server_port").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSignalRService_upstreamSetting(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service", "test")
	r := SignalRServiceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withUpstreamSettings(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("upstream_setting.#").HasValue("4"),
				check.That(data.ResourceName).Key("upstream_setting.0.category_pattern").DoesNotExist(),
				check.That(data.ResourceName).Key("upstream_setting.0.event_pattern").DoesNotExist(),
				check.That(data.ResourceName).Key("upstream_setting.0.hub_pattern").DoesNotExist(),
				check.That(data.ResourceName).Key("upstream_setting.0.url_template").HasValue("http://foo.com"),
				check.That(data.ResourceName).Key("upstream_setting.1.category_pattern").HasValue("connections,messages"),
				check.That(data.ResourceName).Key("upstream_setting.1.event_pattern").HasValue("*"),
				check.That(data.ResourceName).Key("upstream_setting.1.hub_pattern").HasValue("hub1"),
				check.That(data.ResourceName).Key("upstream_setting.1.url_template").HasValue("http://foo.com"),
				check.That(data.ResourceName).Key("upstream_setting.2.category_pattern").HasValue("*"),
				check.That(data.ResourceName).Key("upstream_setting.2.event_pattern").HasValue("connect,disconnect"),
				check.That(data.ResourceName).Key("upstream_setting.2.hub_pattern").HasValue("hub1,hub2"),
				check.That(data.ResourceName).Key("upstream_setting.2.url_template").HasValue("http://foo3.com"),
				check.That(data.ResourceName).Key("upstream_setting.3.category_pattern").HasValue("connections"),
				check.That(data.ResourceName).Key("upstream_setting.3.event_pattern").HasValue("disconnect"),
				check.That(data.ResourceName).Key("upstream_setting.3.hub_pattern").HasValue("*"),
				check.That(data.ResourceName).Key("upstream_setting.3.url_template").HasValue("http://foo4.com"),
			),
		},
	})
}

func (r SignalRServiceResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ServiceID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.SignalR.Client.Get(ctx, id.ResourceGroup, id.SignalRName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving SignalR Service %q (Resource Group %q): %+v", id.SignalRName, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r SignalRServiceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_signalr_service" "test" {
  name                = "acctestSignalR-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Free_F1"
    capacity = 1
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r SignalRServiceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_signalr_service" "import" {
  name                = azurerm_signalr_service.test.name
  location            = azurerm_signalr_service.test.location
  resource_group_name = azurerm_signalr_service.test.resource_group_name

  sku {
    name     = "Free_F1"
    capacity = 1
  }
}
`, r.basic(data))
}

func (r SignalRServiceResource) standardWithCapacity(data acceptance.TestData, capacity int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_signalr_service" "test" {
  name                = "acctestSignalR-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Standard_S1"
    capacity = %d
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, capacity)
}

func (r SignalRServiceResource) withCors(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_signalr_service" "test" {
  name                = "acctestSignalR-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Free_F1"
    capacity = 1
  }

  cors {
    allowed_origins = [
      "https://example.com",
      "https://contoso.com",
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r SignalRServiceResource) withServiceMode(data acceptance.TestData, serviceMode string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_signalr_service" "test" {
  name                = "acctestSignalR-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Free_F1"
    capacity = 1
  }

  features {
    flag  = "ServiceMode"
    value = "%s"
  }

  features {
    flag  = "EnableConnectivityLogs"
    value = "False"
  }

  features {
    flag  = "EnableMessagingLogs"
    value = "False"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, serviceMode)
}

func (r SignalRServiceResource) withUpstreamSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_signalr_service" "test" {
  name                = "acctestSignalR-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Free_F1"
    capacity = 1
  }

  features {
    flag  = "ServiceMode"
    value = "Serverless"
  }

  upstream_setting {
    url_template = "http://foo.com"
  }

  upstream_setting {
    category_pattern = "connections,messages"
    event_pattern    = "*"
    hub_pattern      = "hub1"
    url_template     = "http://foo.com"
  }

  upstream_setting {
    category_pattern = "*"
    event_pattern    = "connect,disconnect"
    hub_pattern      = "hub1,hub2"
    url_template     = "http://foo3.com"
  }

  upstream_setting {
    category_pattern = "connections"
    event_pattern    = "disconnect"
    hub_pattern      = "*"
    url_template     = "http://foo4.com"
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
