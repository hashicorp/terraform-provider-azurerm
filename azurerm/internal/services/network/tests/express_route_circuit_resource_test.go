package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMExpressRouteCircuit(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests due to
	// Azure only being happy about provisioning a couple at a time
	testCases := map[string]map[string]func(t *testing.T){
		"basic": {
			"metered":                      testAccAzureRMExpressRouteCircuit_basicMetered,
			"unlimited":                    testAccAzureRMExpressRouteCircuit_basicUnlimited,
			"update":                       testAccAzureRMExpressRouteCircuit_update,
			"updateTags":                   testAccAzureRMExpressRouteCircuit_updateTags,
			"tierUpdate":                   testAccAzureRMExpressRouteCircuit_tierUpdate,
			"premiumMetered":               testAccAzureRMExpressRouteCircuit_premiumMetered,
			"premiumUnlimited":             testAccAzureRMExpressRouteCircuit_premiumUnlimited,
			"allowClassicOperationsUpdate": testAccAzureRMExpressRouteCircuit_allowClassicOperationsUpdate,
			"requiresImport":               testAccAzureRMExpressRouteCircuit_requiresImport,
			"data_basic":                   testAccDataSourceAzureRMExpressRoute_basicMetered,
			"bandwidthReduction":           testAccAzureRMExpressRouteCircuit_bandwidthReduction,
		},
		"PrivatePeering": {
			"azurePrivatePeering":           testAccAzureRMExpressRouteCircuitPeering_azurePrivatePeering,
			"azurePrivatePeeringWithUpdate": testAccAzureRMExpressRouteCircuitPeering_azurePrivatePeeringWithCircuitUpdate,
			"requiresImport":                testAccAzureRMExpressRouteCircuitPeering_requiresImport,
		},
		"MicrosoftPeering": {
			"microsoftPeering":                    testAccAzureRMExpressRouteCircuitPeering_microsoftPeering,
			"microsoftPeeringCustomerRouting":     testAccAzureRMExpressRouteCircuitPeering_microsoftPeeringCustomerRouting,
			"microsoftPeeringWithRouteFilter":     testAccAzureRMExpressRouteCircuitPeering_microsoftPeeringWithRouteFilter,
			"microsoftPeeringIpv6":                testAccAzureRMExpressRouteCircuitPeering_microsoftPeeringIpv6,
			"microsoftPeeringIpv6CustomerRouting": testAccAzureRMExpressRouteCircuitPeering_microsoftPeeringIpv6CustomerRouting,
			"microsoftPeeringIpv6WithRouteFilter": testAccAzureRMExpressRouteCircuitPeering_microsoftPeeringIpv6WithRouteFilter,
		},
		"authorization": {
			"basic":          testAccAzureRMExpressRouteCircuitAuthorization_basic,
			"multiple":       testAccAzureRMExpressRouteCircuitAuthorization_multiple,
			"requiresImport": testAccAzureRMExpressRouteCircuitAuthorization_requiresImport,
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

func testAccAzureRMExpressRouteCircuit_basicMetered(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	var erc network.ExpressRouteCircuit

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuit_basicMeteredConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists(data.ResourceName, &erc),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMExpressRouteCircuit_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	var erc network.ExpressRouteCircuit

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuit_basicMeteredConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists(data.ResourceName, &erc),
				),
			},
			{
				Config:      testAccAzureRMExpressRouteCircuit_requiresImportConfig(data),
				ExpectError: acceptance.RequiresImportError("azurerm_express_route_circuit"),
			},
		},
	})
}

func testAccAzureRMExpressRouteCircuit_basicUnlimited(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	var erc network.ExpressRouteCircuit

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuit_basicUnlimitedConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists(data.ResourceName, &erc),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMExpressRouteCircuit_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	var erc network.ExpressRouteCircuit

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuit_basicMeteredConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists(data.ResourceName, &erc),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.family", "MeteredData"),
				),
			},
			{
				Config: testAccAzureRMExpressRouteCircuit_basicUnlimitedConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists(data.ResourceName, &erc),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.family", "UnlimitedData"),
				),
			},
		},
	})
}

func testAccAzureRMExpressRouteCircuit_updateTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	var erc network.ExpressRouteCircuit

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuit_basicMeteredConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists(data.ResourceName, &erc),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Environment", "production"),
				),
			},
			{
				Config: testAccAzureRMExpressRouteCircuit_basicMeteredConfigUpdatedTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists(data.ResourceName, &erc),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Environment", "test"),
				),
			},
		},
	})
}

func testAccAzureRMExpressRouteCircuit_tierUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	var erc network.ExpressRouteCircuit

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuit_sku(data, "Standard", "MeteredData"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists("azurerm_express_route_circuit.test", &erc),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "Standard"),
				),
			},
			{
				Config: testAccAzureRMExpressRouteCircuit_sku(data, "Premium", "MeteredData"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists("azurerm_express_route_circuit.test", &erc),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "Premium"),
				),
			},
		},
	})
}

func testAccAzureRMExpressRouteCircuit_premiumMetered(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	var erc network.ExpressRouteCircuit

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuit_sku(data, "Premium", "MeteredData"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists("azurerm_express_route_circuit.test", &erc),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "Premium"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.family", "MeteredData"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMExpressRouteCircuit_premiumUnlimited(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	var erc network.ExpressRouteCircuit

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuit_sku(data, "Premium", "UnlimitedData"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists("azurerm_express_route_circuit.test", &erc),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "Premium"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.family", "UnlimitedData"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMExpressRouteCircuit_allowClassicOperationsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	var erc network.ExpressRouteCircuit

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuit_allowClassicOperations(data, "false"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists("azurerm_express_route_circuit.test", &erc),
					resource.TestCheckResourceAttr(data.ResourceName, "allow_classic_operations", "false"),
				),
			},
			{
				Config: testAccAzureRMExpressRouteCircuit_allowClassicOperations(data, "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists("azurerm_express_route_circuit.test", &erc),
					resource.TestCheckResourceAttr(data.ResourceName, "allow_classic_operations", "true"),
				),
			},
		},
	})
}

func testAccAzureRMExpressRouteCircuit_bandwidthReduction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit", "test")
	var erc network.ExpressRouteCircuit

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuit_bandwidthReductionConfig(data, "100"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists("azurerm_express_route_circuit.test", &erc),
					resource.TestCheckResourceAttr(data.ResourceName, "bandwidth_in_mbps", "100"),
				),
			},
			{
				Config: testAccAzureRMExpressRouteCircuit_bandwidthReductionConfig(data, "50"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists("azurerm_express_route_circuit.test", &erc),
					resource.TestCheckResourceAttr(data.ResourceName, "bandwidth_in_mbps", "50"),
				),
			},
		},
	})
}

func testCheckAzureRMExpressRouteCircuitExists(resourceName string, erc *network.ExpressRouteCircuit) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		expressRouteCircuitName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Express Route Circuit: %s", expressRouteCircuitName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ExpressRouteCircuitsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, resourceGroup, expressRouteCircuitName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Express Route Circuit %q (resource group: %q) does not exist", expressRouteCircuitName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on expressRouteCircuitClient: %+v", err)
		}

		*erc = resp

		return nil
	}
}

func testCheckAzureRMExpressRouteCircuitDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ExpressRouteCircuitsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_express_route_circuit" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Express Route Circuit still exists:\n%#v", resp.ExpressRouteCircuitPropertiesFormat)
		}
	}

	return nil
}

func testAccAzureRMExpressRouteCircuit_basicMeteredConfig(data acceptance.TestData) string {
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
  service_provider_name = "EquinixTest"
  peering_location      = "Area51"
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

func testAccAzureRMExpressRouteCircuit_basicMeteredConfigUpdatedTags(data acceptance.TestData) string {
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
  service_provider_name = "EquinixTest"
  peering_location      = "Area51"
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

func testAccAzureRMExpressRouteCircuit_requiresImportConfig(data acceptance.TestData) string {
	template := testAccAzureRMExpressRouteCircuit_basicMeteredConfig(data)
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
`, template)
}

func testAccAzureRMExpressRouteCircuit_basicUnlimitedConfig(data acceptance.TestData) string {
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

func testAccAzureRMExpressRouteCircuit_sku(data acceptance.TestData, tier string, family string) string {
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

func testAccAzureRMExpressRouteCircuit_allowClassicOperations(data acceptance.TestData, allowClassicOperations string) string {
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

func testAccAzureRMExpressRouteCircuit_bandwidthReductionConfig(data acceptance.TestData, bandwidth string) string {
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
