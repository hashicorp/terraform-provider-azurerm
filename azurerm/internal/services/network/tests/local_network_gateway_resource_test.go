package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMLocalNetworkGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLocalNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "gateway_address", "127.0.0.1"),
					resource.TestCheckResourceAttr(data.ResourceName, "address_space.0", "127.0.0.0/8"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLocalNetworkGateway_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLocalNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMLocalNetworkGatewayConfig_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_local_network_gateway"),
			},
		},
	})
}

func TestAccAzureRMLocalNetworkGateway_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLocalNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "gateway_address", "127.0.0.1"),
					resource.TestCheckResourceAttr(data.ResourceName, "address_space.0", "127.0.0.0/8"),
					testCheckAzureRMLocalNetworkGatewayDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMLocalNetworkGateway_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLocalNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_tags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "acctest"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLocalNetworkGateway_bgpSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLocalNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_bgpSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "gateway_address", "127.0.0.1"),
					resource.TestCheckResourceAttr(data.ResourceName, "address_space.0", "127.0.0.0/8"),
					resource.TestCheckResourceAttr(data.ResourceName, "bgp_settings.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLocalNetworkGateway_bgpSettingsDisable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLocalNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_bgpSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "gateway_address", "127.0.0.1"),
					resource.TestCheckResourceAttr(data.ResourceName, "address_space.0", "127.0.0.0/8"),
					resource.TestCheckResourceAttr(data.ResourceName, "bgp_settings.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "bgp_settings.0.asn", "2468"),
					resource.TestCheckResourceAttr(data.ResourceName, "bgp_settings.0.bgp_peering_address", "10.104.1.1"),
				),
			},
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "gateway_address", "127.0.0.1"),
					resource.TestCheckResourceAttr(data.ResourceName, "address_space.0", "127.0.0.0/8"),
					resource.TestCheckResourceAttr(data.ResourceName, "bgp_settings.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMLocalNetworkGateway_bgpSettingsEnable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLocalNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "gateway_address", "127.0.0.1"),
					resource.TestCheckResourceAttr(data.ResourceName, "address_space.0", "127.0.0.0/8"),
					resource.TestCheckResourceAttr(data.ResourceName, "bgp_settings.#", "0"),
				),
			},
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_bgpSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "gateway_address", "127.0.0.1"),
					resource.TestCheckResourceAttr(data.ResourceName, "address_space.0", "127.0.0.0/8"),
					resource.TestCheckResourceAttr(data.ResourceName, "bgp_settings.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "bgp_settings.0.asn", "2468"),
					resource.TestCheckResourceAttr(data.ResourceName, "bgp_settings.0.bgp_peering_address", "10.104.1.1"),
				),
			},
		},
	})
}

func TestAccAzureRMLocalNetworkGateway_bgpSettingsComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLocalNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_bgpSettingsComplete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "gateway_address", "127.0.0.1"),
					resource.TestCheckResourceAttr(data.ResourceName, "address_space.0", "127.0.0.0/8"),
					resource.TestCheckResourceAttr(data.ResourceName, "bgp_settings.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "bgp_settings.0.asn", "2468"),
					resource.TestCheckResourceAttr(data.ResourceName, "bgp_settings.0.bgp_peering_address", "10.104.1.1"),
					resource.TestCheckResourceAttr(data.ResourceName, "bgp_settings.0.peer_weight", "15"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLocalNetworkGateway_updateAddressSpace(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLocalNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_multipleAddressSpace(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_multipleAddressSpaceUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLocalNetworkGateway_fqdn(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLocalNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_fqdn(data, "www.foo.com"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_fqdn(data, "www.bar.com"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

// testCheckAzureRMLocalNetworkGatewayExists returns the resource.TestCheckFunc
// which checks whether or not the expected local network gateway exists both
// in the schema, and on Azure.
func testCheckAzureRMLocalNetworkGatewayExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// first, check that it exists on Azure:
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.LocalNetworkGatewaysClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// then, check within the schema for the local network gateway:
		res, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Local network gateway '%s' not found.", resourceName)
		}

		// and finally, extract the name and the resource group:
		id, err := azure.ParseAzureResourceID(res.Primary.ID)
		if err != nil {
			return err
		}
		localNetName := id.Path["localNetworkGateways"]
		resGrp := id.ResourceGroup

		resp, err := client.Get(ctx, resGrp, localNetName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Local network gateway %q (resource group %q) does not exist on Azure.", localNetName, resGrp)
			}

			return fmt.Errorf("Error reading the state of local network gateway %q: %+v", localNetName, err)
		}

		return nil
	}
}

func testCheckAzureRMLocalNetworkGatewayDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// first, check that it exists on Azure:
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.LocalNetworkGatewaysClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// then check within the schema for the local network gateway:
		res, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Local network gateway '%s' not found.", resourceName)
		}

		// and finally, extract the name and the resource group:
		id, err := azure.ParseAzureResourceID(res.Primary.ID)
		if err != nil {
			return err
		}
		localNetName := id.Path["localNetworkGateways"]
		resourceGroup := id.ResourceGroup

		future, err := client.Delete(ctx, resourceGroup, localNetName)
		if err != nil {
			if response.WasNotFound(future.Response()) {
				return fmt.Errorf("Local network gateway %q (resource group %q) does not exist on Azure.", localNetName, resourceGroup)
			}
			return fmt.Errorf("Error deleting the state of local network gateway %q: %+v", localNetName, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for deletion of the local network gateway %q to complete: %+v", localNetName, err)
		}

		return nil
	}
}

func testCheckAzureRMLocalNetworkGatewayDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.LocalNetworkGatewaysClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, res := range s.RootModule().Resources {
		if res.Type != "azurerm_local_network_gateway" {
			continue
		}

		id, err := azure.ParseAzureResourceID(res.Primary.ID)
		if err != nil {
			return err
		}
		localNetName := id.Path["localNetworkGateways"]
		resourceGroup := id.ResourceGroup

		resp, err := client.Get(ctx, resourceGroup, localNetName)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Local network gateway still exists:\n%#v", resp.LocalNetworkGatewayPropertiesFormat)
	}

	return nil
}

func testAccAzureRMLocalNetworkGatewayConfig_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/8"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMLocalNetworkGatewayConfig_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLocalNetworkGatewayConfig_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_local_network_gateway" "import" {
  name                = azurerm_local_network_gateway.test.name
  location            = azurerm_local_network_gateway.test.location
  resource_group_name = azurerm_local_network_gateway.test.resource_group_name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/8"]
}
`, template)
}

func testAccAzureRMLocalNetworkGatewayConfig_tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/8"]

  tags = {
    environment = "acctest"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMLocalNetworkGatewayConfig_bgpSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/8"]

  bgp_settings {
    asn                 = 2468
    bgp_peering_address = "10.104.1.1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMLocalNetworkGatewayConfig_bgpSettingsComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/8"]

  bgp_settings {
    asn                 = 2468
    bgp_peering_address = "10.104.1.1"
    peer_weight         = 15
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMLocalNetworkGatewayConfig_multipleAddressSpace(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/24", "127.0.1.0/24"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMLocalNetworkGatewayConfig_multipleAddressSpaceUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.1.0/24", "127.0.0.0/24"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMLocalNetworkGatewayConfig_fqdn(data acceptance.TestData, fqdn string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_fqdn        = %q
  address_space       = ["127.0.0.0/8"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, fqdn)
}
