package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func testAccAzureRMExpressRouteCircuitPeering_azurePrivatePeering(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_peering", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitPeering_privatePeering(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitPeeringExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "peering_type", "AzurePrivatePeering"),
					resource.TestCheckResourceAttr(data.ResourceName, "microsoft_peering_config.#", "0"),
				),
			},
			data.ImportStep("shared_key"), // is not returned by the API
		},
	})
}

func testAccAzureRMExpressRouteCircuitPeering_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_peering", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitPeering_privatePeering(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitPeeringExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMExpressRouteCircuitPeering_requiresImportConfig),
		},
	})
}

func testAccAzureRMExpressRouteCircuitPeering_microsoftPeering(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_peering", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitPeering_msPeering(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitPeeringExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "peering_type", "MicrosoftPeering"),
					resource.TestCheckResourceAttr(data.ResourceName, "microsoft_peering_config.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMExpressRouteCircuitPeering_microsoftPeeringIpv6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_peering", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitPeering_msPeeringIpv6(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitPeeringExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "peering_type", "MicrosoftPeering"),
					resource.TestCheckResourceAttr(data.ResourceName, "microsoft_peering_config.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "microsoft_peering_config_ipv6.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "microsoft_peering_config_ipv6.0.microsoft_peering_config.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "microsoft_peering_config_ipv6.0.primary_peer_address_prefix", "2002:db01::/126"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMExpressRouteCircuitPeering_microsoftPeeringIpv6CustomerRouting(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_peering", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitPeering_msPeeringIpv6CustomerRouting(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitPeeringExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "peering_type", "MicrosoftPeering"),
					resource.TestCheckResourceAttr(data.ResourceName, "microsoft_peering_config.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "microsoft_peering_config_ipv6.0.microsoft_peering_config.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "microsoft_peering_config_ipv6.0.microsoft_peering_config.0.customer_asn", "64511"),
					resource.TestCheckResourceAttr(data.ResourceName, "microsoft_peering_config_ipv6.0.microsoft_peering_config.0.routing_registry_name", "ARIN"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMExpressRouteCircuitPeering_microsoftPeeringIpv6WithRouteFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_peering", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitPeering_msPeeringIpv6WithRouteFilter(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitPeeringExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "peering_type", "MicrosoftPeering"),
					resource.TestCheckResourceAttr(data.ResourceName, "microsoft_peering_config.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "microsoft_peering_config_ipv6.0.microsoft_peering_config.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "microsoft_peering_config_ipv6.0.microsoft_peering_config.0.customer_asn", "64511"),
					resource.TestCheckResourceAttr(data.ResourceName, "microsoft_peering_config_ipv6.0.microsoft_peering_config.0.routing_registry_name", "ARIN"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMExpressRouteCircuitPeering_microsoftPeeringCustomerRouting(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_peering", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitPeering_msPeeringCustomerRouting(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitPeeringExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "peering_type", "MicrosoftPeering"),
					resource.TestCheckResourceAttr(data.ResourceName, "microsoft_peering_config.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "microsoft_peering_config.0.customer_asn", "64511"),
					resource.TestCheckResourceAttr(data.ResourceName, "microsoft_peering_config.0.routing_registry_name", "ARIN"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMExpressRouteCircuitPeering_azurePrivatePeeringWithCircuitUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_peering", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitPeering_privatePeering(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitPeeringExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "peering_type", "AzurePrivatePeering"),
					resource.TestCheckResourceAttr(data.ResourceName, "microsoft_peering_config.#", "0"),
				),
			},
			{
				Config: testAccAzureRMExpressRouteCircuitPeering_privatePeeringWithCircuitUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitPeeringExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "peering_type", "AzurePrivatePeering"),
					resource.TestCheckResourceAttr(data.ResourceName, "microsoft_peering_config.#", "0"),
				),
			},
			data.ImportStep("shared_key"), // is not returned by the API
		},
	})
}

func testAccAzureRMExpressRouteCircuitPeering_microsoftPeeringWithRouteFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_peering", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitPeering_msPeeringWithRouteFilter(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitPeeringExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "peering_type", "MicrosoftPeering"),
					resource.TestCheckResourceAttr(data.ResourceName, "microsoft_peering_config.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "route_filter_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMExpressRouteCircuitPeeringExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ExpressRoutePeeringsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		peeringType := rs.Primary.Attributes["peering_type"]
		circuitName := rs.Primary.Attributes["express_route_circuit_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Express Route Circuit Peering: %s", peeringType)
		}

		resp, err := client.Get(ctx, resourceGroup, circuitName, peeringType)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Express Route Circuit Peering %q (Circuit %q / Resource Group %q) does not exist", peeringType, circuitName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on expressRoutePeeringsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMExpressRouteCircuitPeeringDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ExpressRoutePeeringsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_express_route_circuit_peering" {
			continue
		}

		peeringType := rs.Primary.Attributes["peering_type"]
		circuitName := rs.Primary.Attributes["express_route_circuit_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, circuitName, peeringType)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Express Route Circuit Peering still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMExpressRouteCircuitPeering_privatePeering(data acceptance.TestData) string {
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

func testAccAzureRMExpressRouteCircuitPeering_requiresImportConfig(data acceptance.TestData) string {
	template := testAccAzureRMExpressRouteCircuitPeering_privatePeering(data)
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
`, template)
}

func testAccAzureRMExpressRouteCircuitPeering_msPeering(data acceptance.TestData) string {
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

func testAccAzureRMExpressRouteCircuitPeering_msPeeringIpv6(data acceptance.TestData) string {
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
  service_provider_name = "EquinixTest" //Equinix
  peering_location      = "Area51" //Silicon Valley
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
  }

  microsoft_peering_config_ipv6 {
    primary_peer_address_prefix   = "2002:db01::/126"
    secondary_peer_address_prefix = "2003:db01::/126"
  
    microsoft_peering_config {
      advertised_public_prefixes = ["2002:db01::/126"]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMExpressRouteCircuitPeering_msPeeringIpv6CustomerRouting(data acceptance.TestData) string {
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
  }
  microsoft_peering_config_ipv6 {
    primary_peer_address_prefix   = "2002:db01::/126"
    secondary_peer_address_prefix = "2003:db01::/126"
  
    microsoft_peering_config {
      advertised_public_prefixes = ["2002:db01::/126"]
	  customer_asn          = 64511
      routing_registry_name = "ARIN"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMExpressRouteCircuitPeering_msPeeringIpv6WithRouteFilter(data acceptance.TestData) string {
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
  service_provider_name = "EquinixTest" //Equinix
  peering_location      = "Area51" //Silicon Valley
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
  route_filter_id               = azurerm_route_filter.test.id

  microsoft_peering_config {
    advertised_public_prefixes = ["123.2.0.0/24"]
  }

  microsoft_peering_config_ipv6 {
    primary_peer_address_prefix   = "2002:db01::/126"
    secondary_peer_address_prefix = "2003:db01::/126"
  
    microsoft_peering_config {
      advertised_public_prefixes = ["2002:db01::/126"]
	  customer_asn          = 64511
      routing_registry_name = "ARIN"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMExpressRouteCircuitPeering_msPeeringCustomerRouting(data acceptance.TestData) string {
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
    // https://tools.ietf.org/html/rfc5398
    customer_asn          = 64511
    routing_registry_name = "ARIN"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMExpressRouteCircuitPeering_privatePeeringWithCircuitUpdate(data acceptance.TestData) string {
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

func testAccAzureRMExpressRouteCircuitPeering_msPeeringWithRouteFilter(data acceptance.TestData) string {
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
  primary_peer_address_prefix   = "192.168.1.0/30"
  secondary_peer_address_prefix = "192.168.2.0/30"
  vlan_id                       = 300
  route_filter_id               = azurerm_route_filter.test.id

  microsoft_peering_config {
    advertised_public_prefixes = ["123.1.0.0/24"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
