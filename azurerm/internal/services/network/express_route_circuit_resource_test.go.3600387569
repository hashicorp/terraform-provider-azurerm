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

type ExpressRouteCircuitResource struct {
}

func TestAccExpressRouteCircuit(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests due to
	// Azure only being happy about provisioning a couple at a time
	testCases := map[string]map[string]func(t *testing.T){
		"basic": {
			"metered":                      testAccExpressRouteCircuit_basicMetered,
			"unlimited":                    testAccExpressRouteCircuit_basicUnlimited,
			"update":                       testAccExpressRouteCircuit_update,
			"updateTags":                   testAccExpressRouteCircuit_updateTags,
			"tierUpdate":                   testAccExpressRouteCircuit_tierUpdate,
			"premiumMetered":               testAccExpressRouteCircuit_premiumMetered,
			"premiumUnlimited":             testAccExpressRouteCircuit_premiumUnlimited,
			"allowClassicOperationsUpdate": testAccExpressRouteCircuit_allowClassicOperationsUpdate,
			"requiresImport":               testAccExpressRouteCircuit_requiresImport,
			"data_basic":                   testAccDataSourceExpressRoute_basicMetered,
			"bandwidthReduction":           testAccExpressRouteCircuit_bandwidthReduction,
			"port":                         testAccExpressRouteCircuit_withExpressRoutePort,
			"updatePort":                   testAccExpressRouteCircuit_updateExpressRoutePort,
		},
		"PrivatePeering": {
			"azurePrivatePeering":           testAccExpressRouteCircuitPeering_azurePrivatePeering,
			"azurePrivatePeeringWithUpdate": testAccExpressRouteCircuitPeering_azurePrivatePeeringWithCircuitUpdate,
			"requiresImport":                testAccExpressRouteCircuitPeering_requiresImport,
		},
		"MicrosoftPeering": {
			"microsoftPeering":                    testAccExpressRouteCircuitPeering_microsoftPeering,
			"microsoftPeeringCustomerRouting":     testAccExpressRouteCircuitPeering_microsoftPeeringCustomerRouting,
			"microsoftPeeringWithRouteFilter":     testAccExpressRouteCircuitPeering_microsoftPeeringWithRouteFilter,
			"microsoftPeeringIpv6":                testAccExpressRouteCircuitPeering_microsoftPeeringIpv6,
			"microsoftPeeringIpv6CustomerRouting": testAccExpressRouteCircuitPeering_microsoftPeeringIpv6CustomerRouting,
			"microsoftPeeringIpv6WithRouteFilter": testAccExpressRouteCircuitPeering_microsoftPeeringIpv6WithRouteFilter,
		},
		"authorization": {
			"basic":          testAccExpressRouteCircuitAuthorization_basic,
			"multiple":       testAccExpressRouteCircuitAuthorization_multiple,
			"requiresImport": testAccExpressRouteCircuitAuthorization_requiresImport,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func testAccExpressRouteCircuit_basicMetered(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	r := ExpressRouteCircuitResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMeteredConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccExpressRouteCircuit_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	r := ExpressRouteCircuitResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMeteredConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImportConfig(data),
			ExpectError: acceptance.RequiresImportError("azurerm_express_route_circuit"),
		},
	})
}

func testAccExpressRouteCircuit_basicUnlimited(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	r := ExpressRouteCircuitResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicUnlimitedConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccExpressRouteCircuit_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	r := ExpressRouteCircuitResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMeteredConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.family").HasValue("MeteredData"),
			),
		},
		{
			Config: r.basicUnlimitedConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.family").HasValue("UnlimitedData"),
			),
		},
	})
}

func testAccExpressRouteCircuit_updateTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	r := ExpressRouteCircuitResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMeteredConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.Environment").HasValue("production"),
			),
		},
		{
			Config: r.basicMeteredConfigUpdatedTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.Environment").HasValue("test"),
			),
		},
	})
}

func testAccExpressRouteCircuit_tierUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	r := ExpressRouteCircuitResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.sku(data, "Standard", "MeteredData"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("Standard"),
			),
		},
		{
			Config: r.sku(data, "Premium", "MeteredData"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("Premium"),
			),
		},
	})
}

func testAccExpressRouteCircuit_premiumMetered(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	r := ExpressRouteCircuitResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.sku(data, "Premium", "MeteredData"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("Premium"),
				check.That(data.ResourceName).Key("sku.0.family").HasValue("MeteredData"),
			),
		},
		data.ImportStep(),
	})
}

func testAccExpressRouteCircuit_premiumUnlimited(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	r := ExpressRouteCircuitResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.sku(data, "Premium", "UnlimitedData"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("Premium"),
				check.That(data.ResourceName).Key("sku.0.family").HasValue("UnlimitedData"),
			),
		},
		data.ImportStep(),
	})
}

func testAccExpressRouteCircuit_allowClassicOperationsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	r := ExpressRouteCircuitResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.allowClassicOperations(data, "false"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allow_classic_operations").HasValue("false"),
			),
		},
		{
			Config: r.allowClassicOperations(data, "true"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allow_classic_operations").HasValue("true"),
			),
		},
	})
}

func testAccExpressRouteCircuit_bandwidthReduction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	r := ExpressRouteCircuitResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.bandwidthReductionConfig(data, "100"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("bandwidth_in_mbps").HasValue("100"),
			),
		},
		{
			Config: r.bandwidthReductionConfig(data, "50"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("bandwidth_in_mbps").HasValue("50"),
			),
		},
	})
}

func testAccExpressRouteCircuit_withExpressRoutePort(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	r := ExpressRouteCircuitResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.withExpressRoutePort(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccExpressRouteCircuit_updateExpressRoutePort(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	r := ExpressRouteCircuitResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.withExpressRoutePort(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateExpressRoutePort(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t ExpressRouteCircuitResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resGroup := id.ResourceGroup
	name := id.Path["expressRouteCircuits"]

	resp, err := clients.Network.ExpressRouteCircuitsClient.Get(ctx, resGroup, name)
	if err != nil {
		return nil, fmt.Errorf("reading Express Route Circuit (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ExpressRouteCircuitResource) basicMeteredConfig(data acceptance.TestData) string {
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

  allow_classic_operations = false

  tags = {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ExpressRouteCircuitResource) basicMeteredConfigUpdatedTags(data acceptance.TestData) string {
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

  allow_classic_operations = false

  tags = {
    Environment = "test"
    Purpose     = "UpdatedAcceptanceTests"
    Caller      = "Additional Value"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ExpressRouteCircuitResource) requiresImportConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_circuit" "import" {
  name                  = azurerm_express_route_circuit.test.name
  location              = azurerm_express_route_circuit.test.location
  resource_group_name   = azurerm_express_route_circuit.test.resource_group_name
  service_provider_name = "Equinix"
  peering_location      = "Silicon Valley"
  bandwidth_in_mbps     = 50

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }

  allow_classic_operations = false

  tags = {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}
`, r.basicMeteredConfig(data))
}

func (ExpressRouteCircuitResource) basicUnlimitedConfig(data acceptance.TestData) string {
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
    family = "UnlimitedData"
  }

  allow_classic_operations = false

  tags = {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ExpressRouteCircuitResource) sku(data acceptance.TestData, tier string, family string) string {
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
    tier   = "%s"
    family = "%s"
  }

  allow_classic_operations = false

  tags = {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, tier, family)
}

func (ExpressRouteCircuitResource) allowClassicOperations(data acceptance.TestData, allowClassicOperations string) string {
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

  allow_classic_operations = %s

  tags = {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, allowClassicOperations)
}

func (ExpressRouteCircuitResource) bandwidthReductionConfig(data acceptance.TestData, bandwidth string) string {
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
  bandwidth_in_mbps     = %s

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }

  allow_classic_operations = false

  tags = {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, bandwidth)
}

func (ExpressRouteCircuitResource) withExpressRoutePort(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-expressroutecircuit-%d"
  location = "%s"
}

resource "azurerm_express_route_port" "test" {
  name                = "acctest-ERP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  peering_location    = "Airtel-Chennai2-CLS"
  bandwidth_in_gbps   = 10
  encapsulation       = "Dot1Q"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-ExpressRouteCircuit-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  express_route_port_id = azurerm_express_route_port.test.id
  bandwidth_in_gbps     = 5

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ExpressRouteCircuitResource) updateExpressRoutePort(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-expressroutecircuit-%d"
  location = "%s"
}

resource "azurerm_express_route_port" "test" {
  name                = "acctest-ERP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  peering_location    = "Airtel-Chennai2-CLS"
  bandwidth_in_gbps   = 10
  encapsulation       = "Dot1Q"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-ExpressRouteCircuit-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  express_route_port_id = azurerm_express_route_port.test.id
  bandwidth_in_gbps     = 10

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
