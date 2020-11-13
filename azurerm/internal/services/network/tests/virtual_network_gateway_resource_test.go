package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVirtualNetworkGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGateway_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "Basic"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualNetworkGateway_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGateway_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMVirtualNetworkGateway_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_virtual_network_gateway"),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGateway_lowerCaseSubnetName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGateway_lowerCaseSubnetName(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "Basic"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualNetworkGateway_vpnGw1(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGateway_vpnGw1(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists("azurerm_virtual_network_gateway.test"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGateway_activeActive(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGateway_activeActive(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists("azurerm_virtual_network_gateway.test"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGateway_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGateway_sku(data, "Standard"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "Standard"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGateway_vpnGw2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGateway_sku(data, "VpnGw2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "VpnGw2"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGateway_vpnGw3(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGateway_sku(data, "VpnGw3"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "VpnGw3"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGateway_generation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGateway_generation(data, "Generation2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "generation", "Generation2"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGateway_vpnClientConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGateway_vpnClientConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "vpn_client_configuration.0.radius_server_address", "1.2.3.4"),
					resource.TestCheckResourceAttr(data.ResourceName, "vpn_client_configuration.0.vpn_client_protocols.#", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGateway_vpnClientConfigAzureAdAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGateway_vpnClientConfigAzureAdAuth(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "vpn_client_configuration.0.aad_tenant", fmt.Sprintf("https://login.microsoftonline.com/%s/", data.Client().TenantID)),
					resource.TestCheckResourceAttr(data.ResourceName, "vpn_client_configuration.0.aad_audience", "41b23e61-6c1e-4545-b367-cd054e0ed4b4"),
					resource.TestCheckResourceAttr(data.ResourceName, "vpn_client_configuration.0.aad_issuer", fmt.Sprintf("https://sts.windows.net/%s/", data.Client().TenantID)),
					resource.TestCheckResourceAttr(data.ResourceName, "vpn_client_configuration.0.vpn_client_protocols.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGateway_vpnClientConfigOpenVPN(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGateway_vpnClientConfigOpenVPN(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "vpn_client_configuration.0.vpn_client_protocols.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGateway_enableBgp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGateway_enableBgp(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_bgp", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "bgp_settings.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGateway_expressRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGateway_expressRoute(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "type", "ExpressRoute"),
					resource.TestCheckResourceAttr(data.ResourceName, "bgp_settings.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGateway_privateIpAddressEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGateway_privateIpAddressEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists("azurerm_virtual_network_gateway.test"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVirtualNetworkGateway_privateIpAddressEnabledUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists("azurerm_virtual_network_gateway.test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMVirtualNetworkGatewayExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VnetGatewayClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		gatewayName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, gatewayName)
		if err != nil {
			return fmt.Errorf("Bad: Get on vnetGatewayClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Virtual Network Gateway %q (resource group: %q) does not exist", gatewayName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMVirtualNetworkGatewayDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VnetGatewayClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_network_gateway" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Virtual Network Gateway still exists:\n%#v", resp.VirtualNetworkGatewayPropertiesFormat)
		}
	}

	return nil
}

func testAccAzureRMVirtualNetworkGateway_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualNetworkGateway_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMVirtualNetworkGateway_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network_gateway" "import" {
  name                = azurerm_virtual_network_gateway.test.name
  location            = azurerm_virtual_network_gateway.test.location
  resource_group_name = azurerm_virtual_network_gateway.test.resource_group_name
  type                = azurerm_virtual_network_gateway.test.type
  vpn_type            = azurerm_virtual_network_gateway.test.vpn_type
  sku                 = azurerm_virtual_network_gateway.test.sku

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
`, template)
}

func testAccAzureRMVirtualNetworkGateway_lowerCaseSubnetName(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "gatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualNetworkGateway_vpnGw1(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "VpnGw1"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualNetworkGateway_activeActive(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "first" {
  name                = "acctestpip1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_public_ip" "second" {
  name = "acctestpip2-%d"

  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  depends_on = [
    azurerm_public_ip.first,
    azurerm_public_ip.second,
  ]
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "VpnGw1"

  active_active = true
  enable_bgp    = true

  ip_configuration {
    name                 = "gw-ip1"
    public_ip_address_id = azurerm_public_ip.first.id

    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  ip_configuration {
    name                          = "gw-ip2"
    public_ip_address_id          = azurerm_public_ip.second.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  bgp_settings {
    asn = "65010"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualNetworkGateway_vpnClientConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  depends_on          = [azurerm_public_ip.test]
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "VpnGw1"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  vpn_client_configuration {
    address_space        = ["10.2.0.0/24"]
    vpn_client_protocols = ["SSTP", "IkeV2"]

    radius_server_address = "1.2.3.4"
    radius_server_secret  = "1234"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualNetworkGateway_vpnClientConfigAzureAdAuth(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "VpnGw1"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  vpn_client_configuration {
    address_space        = ["10.2.0.0/24"]
    vpn_client_protocols = ["OpenVPN"]

    aad_tenant   = "https://login.microsoftonline.com/%s/"
    aad_audience = "41b23e61-6c1e-4545-b367-cd054e0ed4b4"
    aad_issuer   = "https://sts.windows.net/%s/"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.Client().TenantID, data.Client().TenantID)
}

func testAccAzureRMVirtualNetworkGateway_vpnClientConfigOpenVPN(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  depends_on          = [azurerm_public_ip.test]
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "VpnGw1"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  vpn_client_configuration {
    address_space        = ["10.2.0.0/24"]
    vpn_client_protocols = ["OpenVPN"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualNetworkGateway_sku(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "%s"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, sku)
}

func testAccAzureRMVirtualNetworkGateway_enableBgp(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type       = "Vpn"
  vpn_type   = "RouteBased"
  sku        = "VpnGw1"
  enable_bgp = true

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualNetworkGateway_expressRoute(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "ExpressRoute"
  vpn_type = "PolicyBased"
  sku      = "Standard"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualNetworkGateway_generation(data acceptance.TestData, generation string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type       = "Vpn"
  vpn_type   = "RouteBased"
  sku        = "VpnGw2"
  generation = "%s"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, generation)
}

func testAccAzureRMVirtualNetworkGateway_privateIpAddressEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                       = "Vpn"
  vpn_type                   = "RouteBased"
  sku                        = "VpnGw1"
  private_ip_address_enabled = true

  custom_route {
    address_prefixes = [
      "101.168.0.6/32"
    ]
  }

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualNetworkGateway_privateIpAddressEnabledUpdate(data acceptance.TestData) string {
  return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                       = "Vpn"
  vpn_type                   = "RouteBased"
  sku                        = "VpnGw1"
  private_ip_address_enabled = false

  custom_route {
    address_prefixes = [
      "101.168.0.6/32"
    ]
  }

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
