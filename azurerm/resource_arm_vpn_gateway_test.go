package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMVPNGateway_basic(t *testing.T) {
	resourceName := "azurerm_vpn_gateway.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVPNGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVPNGateway_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVPNGatewayExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMVPNGateway_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_vpn_gateway.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVPNGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVPNGateway_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVPNGatewayExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMVPNGateway_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_vpn_gateway"),
			},
		},
	})
}

func TestAccAzureRMVPNGateway_bgpSettings(t *testing.T) {
	resourceName := "azurerm_vpn_gateway.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVPNGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVPNGateway_bgpSettings(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVPNGatewayExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMVPNGateway_scaleUnit(t *testing.T) {
	resourceName := "azurerm_vpn_gateway.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVPNGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVPNGateway_scaleUnit(ri, location, 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVPNGatewayExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMVPNGateway_scaleUnit(ri, location, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVPNGatewayExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMVPNGateway_tags(t *testing.T) {
	resourceName := "azurerm_vpn_gateway.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVPNGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVPNGateway_tags(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVPNGatewayExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMVPNGateway_tagsUpdated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVPNGatewayExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMVPNGateway_tags(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVPNGatewayExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMVPNGatewayExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		name := rs.Primary.Attributes["name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VpnGatewaysClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
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
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_vpn_gateway" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		name := rs.Primary.Attributes["name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VpnGatewaysClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
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

func testAccAzureRMVPNGateway_basic(rInt int, location string) string {
	template := testAccAzureRMVPNGateway_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "test" {
  name                = "acctestVPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id
}
`, template, rInt)
}

func testAccAzureRMVPNGateway_requiresImport(rInt int, location string) string {
	template := testAccAzureRMVPNGateway_basic(rInt, location)
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

func testAccAzureRMVPNGateway_bgpSettings(rInt int, location string) string {
	template := testAccAzureRMVPNGateway_template(rInt, location)
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
  }
}
`, template, rInt)
}

func testAccAzureRMVPNGateway_scaleUnit(rInt int, location string, scaleUnit int) string {
	template := testAccAzureRMVPNGateway_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "test" {
  name                = "acctestVPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id
  scale_unit          = %d
}
`, template, rInt, scaleUnit)
}

func testAccAzureRMVPNGateway_tags(rInt int, location string) string {
	template := testAccAzureRMVPNGateway_template(rInt, location)
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
`, template, rInt)
}

func testAccAzureRMVPNGateway_tagsUpdated(rInt int, location string) string {
	template := testAccAzureRMVPNGateway_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "test" {
  name                = "acctestVPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id

  tags = {
    Hello = "World"
    Rick = "C-137"
  }
}
`, template, rInt)
}

func testAccAzureRMVPNGateway_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
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
`, rInt, location, rInt, rInt, rInt)
}
