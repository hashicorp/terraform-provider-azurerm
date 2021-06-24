package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ExpressRouteCircuitPeeringResource struct {
}

func testAccExpressRouteCircuitPeering_azurePrivatePeering(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_peering", "test")
	r := ExpressRouteCircuitPeeringResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.privatePeering(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("peering_type").HasValue("AzurePrivatePeering"),
				check.That(data.ResourceName).Key("microsoft_peering_config.#").HasValue("0"),
			),
		},
		data.ImportStep("shared_key"), // is not returned by the API
	})
}

func testAccExpressRouteCircuitPeering_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_peering", "test")
	r := ExpressRouteCircuitPeeringResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.privatePeering(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImportConfig),
	})
}

func testAccExpressRouteCircuitPeering_microsoftPeering(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_peering", "test")
	r := ExpressRouteCircuitPeeringResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.msPeering(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("peering_type").HasValue("MicrosoftPeering"),
				check.That(data.ResourceName).Key("microsoft_peering_config.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func testAccExpressRouteCircuitPeering_microsoftPeeringIpv6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_peering", "test")
	r := ExpressRouteCircuitPeeringResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.msPeeringIpv6(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccExpressRouteCircuitPeering_microsoftPeeringIpv6CustomerRouting(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_peering", "test")
	r := ExpressRouteCircuitPeeringResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.msPeeringIpv6CustomerRouting(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccExpressRouteCircuitPeering_microsoftPeeringIpv6WithRouteFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_peering", "test")
	r := ExpressRouteCircuitPeeringResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.msPeeringIpv6WithRouteFilter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccExpressRouteCircuitPeering_microsoftPeeringCustomerRouting(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_peering", "test")
	r := ExpressRouteCircuitPeeringResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.msPeeringCustomerRouting(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("peering_type").HasValue("MicrosoftPeering"),
				check.That(data.ResourceName).Key("microsoft_peering_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("microsoft_peering_config.0.customer_asn").HasValue("64511"),
				check.That(data.ResourceName).Key("microsoft_peering_config.0.routing_registry_name").HasValue("ARIN"),
			),
		},
		data.ImportStep(),
	})
}

func testAccExpressRouteCircuitPeering_azurePrivatePeeringWithCircuitUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_peering", "test")
	r := ExpressRouteCircuitPeeringResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.privatePeering(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("peering_type").HasValue("AzurePrivatePeering"),
				check.That(data.ResourceName).Key("microsoft_peering_config.#").HasValue("0"),
			),
		},
		{
			Config: r.privatePeeringWithCircuitUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("peering_type").HasValue("AzurePrivatePeering"),
				check.That(data.ResourceName).Key("microsoft_peering_config.#").HasValue("0"),
			),
		},
		data.ImportStep("shared_key"), // is not returned by the API
	})
}

func testAccExpressRouteCircuitPeering_microsoftPeeringWithRouteFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_peering", "test")
	r := ExpressRouteCircuitPeeringResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.msPeeringWithRouteFilter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("peering_type").HasValue("MicrosoftPeering"),
				check.That(data.ResourceName).Key("microsoft_peering_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("route_filter_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (t ExpressRouteCircuitPeeringResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	circuitName := id.Path["expressRouteCircuits"]
	peeringType := id.Path["peerings"]

	resp, err := clients.Network.ExpressRoutePeeringsClient.Get(ctx, resourceGroup, circuitName, peeringType)
	if err != nil {
		return nil, fmt.Errorf("reading Express Route Circuit Peering (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ExpressRouteCircuitPeeringResource) privatePeering(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  service_provider_name = "Equinix"
  peering_location      = "Silicon Valley"
  bandwidth_in_mbps     = 50

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }

  tags = {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}

resource "azurerm_express_route_circuit_peering" "test" {
  peering_type                  = "AzurePrivatePeering"
  express_route_circuit_name    = azurerm_express_route_circuit.test.name
  resource_group_name           = azurerm_resource_group.test.name
  shared_key                    = "SSSSsssssshhhhhItsASecret"
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.1.0/30"
  secondary_peer_address_prefix = "192.168.2.0/30"
  vlan_id                       = 100
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ExpressRouteCircuitPeeringResource) requiresImportConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_circuit_peering" "import" {
  peering_type                  = azurerm_express_route_circuit_peering.test.peering_type
  express_route_circuit_name    = azurerm_express_route_circuit_peering.test.express_route_circuit_name
  resource_group_name           = azurerm_express_route_circuit_peering.test.resource_group_name
  shared_key                    = azurerm_express_route_circuit_peering.test.shared_key
  peer_asn                      = azurerm_express_route_circuit_peering.test.peer_asn
  primary_peer_address_prefix   = azurerm_express_route_circuit_peering.test.primary_peer_address_prefix
  secondary_peer_address_prefix = azurerm_express_route_circuit_peering.test.secondary_peer_address_prefix
  vlan_id                       = azurerm_express_route_circuit_peering.test.vlan_id
}
`, r.privatePeering(data))
}

func (ExpressRouteCircuitPeeringResource) msPeering(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  service_provider_name = "Equinix"
  peering_location      = "Silicon Valley"
  bandwidth_in_mbps     = 50

  sku {
    tier   = "Premium"
    family = "MeteredData"
  }

  tags = {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}

resource "azurerm_express_route_circuit_peering" "test" {
  peering_type                  = "MicrosoftPeering"
  express_route_circuit_name    = azurerm_express_route_circuit.test.name
  resource_group_name           = azurerm_resource_group.test.name
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.1.0/30"
  secondary_peer_address_prefix = "192.168.2.0/30"
  vlan_id                       = 300

  microsoft_peering_config {
    advertised_public_prefixes = ["123.1.0.0/24"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ExpressRouteCircuitPeeringResource) msPeeringIpv6(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-expressroute-%d"
  location = "%s"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  service_provider_name = "Equinix"
  peering_location      = "Silicon Valley"
  bandwidth_in_mbps     = 50

  sku {
    tier   = "Premium"
    family = "MeteredData"
  }

  tags = {
    Env     = "Test"
    Purpose = "AcceptanceTests"
  }
}

resource "azurerm_express_route_circuit_peering" "test" {
  peering_type                  = "MicrosoftPeering"
  express_route_circuit_name    = azurerm_express_route_circuit.test.name
  resource_group_name           = azurerm_resource_group.test.name
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.7.0/30"
  secondary_peer_address_prefix = "192.168.8.0/30"
  vlan_id                       = 300

  microsoft_peering_config {
    advertised_public_prefixes = ["123.4.0.0/24"]
  }

  ipv6 {
    primary_peer_address_prefix   = "2002:db03::/126"
    secondary_peer_address_prefix = "2003:db03::/126"

    microsoft_peering {
      advertised_public_prefixes = ["2002:db01::/126"]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ExpressRouteCircuitPeeringResource) msPeeringIpv6CustomerRouting(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-expressroute-%d"
  location = "%s"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  service_provider_name = "Equinix"
  peering_location      = "Silicon Valley"
  bandwidth_in_mbps     = 50

  sku {
    tier   = "Premium"
    family = "MeteredData"
  }

  tags = {
    Env     = "Test"
    Purpose = "AcceptanceTests"
  }
}

resource "azurerm_express_route_circuit_peering" "test" {
  peering_type                  = "MicrosoftPeering"
  express_route_circuit_name    = azurerm_express_route_circuit.test.name
  resource_group_name           = azurerm_resource_group.test.name
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.9.0/30"
  secondary_peer_address_prefix = "192.168.10.0/30"
  vlan_id                       = 300

  microsoft_peering_config {
    advertised_public_prefixes = ["123.5.0.0/24"]
  }
  ipv6 {
    primary_peer_address_prefix   = "2002:db05::/126"
    secondary_peer_address_prefix = "2003:db05::/126"

    microsoft_peering {
      advertised_public_prefixes = ["2002:db05::/126"]
      customer_asn               = 64511
      routing_registry_name      = "ARIN"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ExpressRouteCircuitPeeringResource) msPeeringIpv6WithRouteFilter(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-expressroute-%d"
  location = "%s"
}

resource "azurerm_route_filter" "test" {
  name                = "acctestrf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  rule {
    name        = "acctestrule%d"
    access      = "Allow"
    rule_type   = "Community"
    communities = ["12076:52005", "12076:52006"]
  }
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  service_provider_name = "Equinix"
  peering_location      = "Silicon Valley"
  bandwidth_in_mbps     = 50

  sku {
    tier   = "Premium"
    family = "MeteredData"
  }

  tags = {
    Env     = "Test"
    Purpose = "AcceptanceTests"
  }
}

resource "azurerm_express_route_circuit_peering" "test" {
  peering_type                  = "MicrosoftPeering"
  express_route_circuit_name    = azurerm_express_route_circuit.test.name
  resource_group_name           = azurerm_resource_group.test.name
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.11.0/30"
  secondary_peer_address_prefix = "192.168.12.0/30"
  vlan_id                       = 300
  route_filter_id               = azurerm_route_filter.test.id

  microsoft_peering_config {
    advertised_public_prefixes = ["123.3.0.0/24"]
  }

  ipv6 {
    primary_peer_address_prefix   = "2002:db02::/126"
    secondary_peer_address_prefix = "2003:db02::/126"
    route_filter_id               = azurerm_route_filter.test.id

    microsoft_peering {
      advertised_public_prefixes = ["2002:db01::/126"]
      customer_asn               = 64511
      routing_registry_name      = "ARIN"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (ExpressRouteCircuitPeeringResource) msPeeringCustomerRouting(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  service_provider_name = "Equinix"
  peering_location      = "Silicon Valley"
  bandwidth_in_mbps     = 50

  sku {
    tier   = "Premium"
    family = "MeteredData"
  }

  tags = {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}

resource "azurerm_express_route_circuit_peering" "test" {
  peering_type                  = "MicrosoftPeering"
  express_route_circuit_name    = azurerm_express_route_circuit.test.name
  resource_group_name           = azurerm_resource_group.test.name
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.3.0/30"
  secondary_peer_address_prefix = "192.168.4.0/30"
  vlan_id                       = 300

  microsoft_peering_config {
    advertised_public_prefixes = ["123.2.0.0/24"]
    // https://tools.ietf.org/html/rfc5398
    customer_asn          = 64511
    routing_registry_name = "ARIN"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ExpressRouteCircuitPeeringResource) privatePeeringWithCircuitUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  service_provider_name = "Equinix"
  peering_location      = "Silicon Valley"
  bandwidth_in_mbps     = 50

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }

  tags = {
    Environment = "prod"
    Purpose     = "AcceptanceTests"
  }
}

resource "azurerm_express_route_circuit_peering" "test" {
  peering_type                  = "AzurePrivatePeering"
  express_route_circuit_name    = azurerm_express_route_circuit.test.name
  resource_group_name           = azurerm_resource_group.test.name
  shared_key                    = "SSSSsssssshhhhhItsASecret"
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.1.0/30"
  secondary_peer_address_prefix = "192.168.2.0/30"
  vlan_id                       = 100
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ExpressRouteCircuitPeeringResource) msPeeringWithRouteFilter(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_filter" "test" {
  name                = "acctestrf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  rule {
    name        = "acctestrule%d"
    access      = "Allow"
    rule_type   = "Community"
    communities = ["12076:52005", "12076:52006"]
  }
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  service_provider_name = "Equinix"
  peering_location      = "Silicon Valley"
  bandwidth_in_mbps     = 50

  sku {
    tier   = "Premium"
    family = "MeteredData"
  }

  tags = {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}

resource "azurerm_express_route_circuit_peering" "test" {
  peering_type                  = "MicrosoftPeering"
  express_route_circuit_name    = azurerm_express_route_circuit.test.name
  resource_group_name           = azurerm_resource_group.test.name
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.5.0/30"
  secondary_peer_address_prefix = "192.168.6.0/30"
  vlan_id                       = 300
  route_filter_id               = azurerm_route_filter.test.id

  microsoft_peering_config {
    advertised_public_prefixes = ["123.1.0.0/24"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
