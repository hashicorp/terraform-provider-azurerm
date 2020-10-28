package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMVPNGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVPNGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVPNGateway_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVPNGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVPNGateway_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVPNGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVPNGateway_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVPNGatewayExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMVPNGateway_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_vpn_gateway"),
			},
		},
	})
}

func TestAccAzureRMVPNGateway_bgpSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVPNGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVPNGateway_bgpSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVPNGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVPNGateway_scaleUnit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVPNGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVPNGateway_scaleUnit(data, 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVPNGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVPNGateway_scaleUnit(data, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVPNGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVPNGateway_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVPNGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVPNGateway_tags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVPNGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVPNGateway_tagsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVPNGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVPNGateway_tags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVPNGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMVPNGatewayExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VpnGatewaysClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		name := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on VpnGatewaysClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: VPN Gateway %q does not exist in Resource Group %q", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMVPNGatewayDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VpnGatewaysClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_vpn_gateway" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		name := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("VPN Gateway still exists:\n%#v", resp.VpnGatewayProperties)
		}
	}

	return nil
}

func testAccAzureRMVPNGateway_basic(data acceptance.TestData) string {
	template := testAccAzureRMVPNGateway_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "test" {
  name                = "acctestVPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id
}
`, template, data.RandomInteger)
}

func testAccAzureRMVPNGateway_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMVPNGateway_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "import" {
  name                = azurerm_vpn_gateway.test.name
  location            = azurerm_vpn_gateway.test.location
  resource_group_name = azurerm_vpn_gateway.test.resource_group_name
  virtual_hub_id      = azurerm_vpn_gateway.test.virtual_hub_id
}
`, template)
}

func testAccAzureRMVPNGateway_bgpSettings(data acceptance.TestData) string {
	template := testAccAzureRMVPNGateway_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "test" {
  name                = "acctestVPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id

  bgp_settings {
    asn         = 65515
    peer_weight = 0

    instance_bgp_peering_address {
      custom_ips = ["169.254.21.5"]
    }

    instance_bgp_peering_address {
      custom_ips = ["169.254.21.10"]
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMVPNGateway_scaleUnit(data acceptance.TestData, scaleUnit int) string {
	template := testAccAzureRMVPNGateway_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "test" {
  name                = "acctestVPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id
  scale_unit          = %d
}
`, template, data.RandomInteger, scaleUnit)
}

func testAccAzureRMVPNGateway_tags(data acceptance.TestData) string {
	template := testAccAzureRMVPNGateway_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "test" {
  name                = "acctestVPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id

  tags = {
    Hello = "World"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMVPNGateway_tagsUpdated(data acceptance.TestData) string {
	template := testAccAzureRMVPNGateway_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "test" {
  name                = "acctestVPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id

  tags = {
    Hello = "World"
    Rick  = "C-137"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMVPNGateway_template(data acceptance.TestData) string {
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

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctestvh-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_prefix      = "10.0.1.0/24"
  virtual_wan_id      = azurerm_virtual_wan.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
