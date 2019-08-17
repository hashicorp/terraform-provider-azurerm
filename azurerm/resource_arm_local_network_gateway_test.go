package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMLocalNetworkGateway_basic(t *testing.T) {
	resourceName := "azurerm_local_network_gateway.test"

	rInt := tf.AccRandTimeInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLocalNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_basic(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gateway_address", "127.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "address_space.0", "127.0.0.0/8"),
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

func TestAccAzureRMLocalNetworkGateway_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_local_network_gateway.test"

	rInt := tf.AccRandTimeInt()
	location := testLocation()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLocalNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMLocalNetworkGatewayConfig_requiresImport(rInt, location),
				ExpectError: testRequiresImportError("azurerm_local_network_gateway"),
			},
		},
	})
}

func TestAccAzureRMLocalNetworkGateway_disappears(t *testing.T) {
	resourceName := "azurerm_local_network_gateway.test"
	rInt := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLocalNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_basic(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gateway_address", "127.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "address_space.0", "127.0.0.0/8"),
					testCheckAzureRMLocalNetworkGatewayDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMLocalNetworkGateway_tags(t *testing.T) {
	resourceName := "azurerm_local_network_gateway.test"

	rInt := tf.AccRandTimeInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLocalNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_tags(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "acctest"),
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

func TestAccAzureRMLocalNetworkGateway_bgpSettings(t *testing.T) {
	resourceName := "azurerm_local_network_gateway.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLocalNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_bgpSettings(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gateway_address", "127.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "address_space.0", "127.0.0.0/8"),
					resource.TestCheckResourceAttr(resourceName, "bgp_settings.#", "1"),
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

func TestAccAzureRMLocalNetworkGateway_bgpSettingsDisable(t *testing.T) {
	resourceName := "azurerm_local_network_gateway.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLocalNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_bgpSettings(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gateway_address", "127.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "address_space.0", "127.0.0.0/8"),
					resource.TestCheckResourceAttr(resourceName, "bgp_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bgp_settings.0.asn", "2468"),
					resource.TestCheckResourceAttr(resourceName, "bgp_settings.0.bgp_peering_address", "10.104.1.1"),
				),
			},
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gateway_address", "127.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "address_space.0", "127.0.0.0/8"),
					resource.TestCheckResourceAttr(resourceName, "bgp_settings.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMLocalNetworkGateway_bgpSettingsEnable(t *testing.T) {
	resourceName := "azurerm_local_network_gateway.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLocalNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gateway_address", "127.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "address_space.0", "127.0.0.0/8"),
					resource.TestCheckResourceAttr(resourceName, "bgp_settings.#", "0"),
				),
			},
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_bgpSettings(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gateway_address", "127.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "address_space.0", "127.0.0.0/8"),
					resource.TestCheckResourceAttr(resourceName, "bgp_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bgp_settings.0.asn", "2468"),
					resource.TestCheckResourceAttr(resourceName, "bgp_settings.0.bgp_peering_address", "10.104.1.1"),
				),
			},
		},
	})
}

func TestAccAzureRMLocalNetworkGateway_bgpSettingsComplete(t *testing.T) {
	resourceName := "azurerm_local_network_gateway.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLocalNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_bgpSettingsComplete(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLocalNetworkGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gateway_address", "127.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "address_space.0", "127.0.0.0/8"),
					resource.TestCheckResourceAttr(resourceName, "bgp_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bgp_settings.0.asn", "2468"),
					resource.TestCheckResourceAttr(resourceName, "bgp_settings.0.bgp_peering_address", "10.104.1.1"),
					resource.TestCheckResourceAttr(resourceName, "bgp_settings.0.peer_weight", "15"),
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

// testCheckAzureRMLocalNetworkGatewayExists returns the resource.TestCheckFunc
// which checks whether or not the expected local network gateway exists both
// in the schema, and on Azure.
func testCheckAzureRMLocalNetworkGatewayExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// first check within the schema for the local network gateway:
		res, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Local network gateway '%s' not found.", resourceName)
		}

		// then, extract the name and the resource group:
		id, err := azure.ParseAzureResourceID(res.Primary.ID)
		if err != nil {
			return err
		}
		localNetName := id.Path["localNetworkGateways"]
		resGrp := id.ResourceGroup

		// and finally, check that it exists on Azure:
		client := testAccProvider.Meta().(*ArmClient).network.LocalNetworkGatewaysClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
		// first check within the schema for the local network gateway:
		res, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Local network gateway '%s' not found.", resourceName)
		}

		// then, extract the name and the resource group:
		id, err := azure.ParseAzureResourceID(res.Primary.ID)
		if err != nil {
			return err
		}
		localNetName := id.Path["localNetworkGateways"]
		resourceGroup := id.ResourceGroup

		// and finally, check that it exists on Azure:
		client := testAccProvider.Meta().(*ArmClient).network.LocalNetworkGatewaysClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

		client := testAccProvider.Meta().(*ArmClient).network.LocalNetworkGatewaysClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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

func testAccAzureRMLocalNetworkGatewayConfig_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/8"]
}
`, rInt, location, rInt)
}

func testAccAzureRMLocalNetworkGatewayConfig_requiresImport(rInt int, location string) string {
	template := testAccAzureRMLocalNetworkGatewayConfig_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_local_network_gateway" "import" {
  name                = "${azurerm_local_network_gateway.test.name}"
  location            = "${azurerm_local_network_gateway.test.location}"
  resource_group_name = "${azurerm_local_network_gateway.test.resource_group_name}"
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/8"]
}
`, template)
}

func testAccAzureRMLocalNetworkGatewayConfig_tags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/8"]

  tags = {
    environment = "acctest"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMLocalNetworkGatewayConfig_bgpSettings(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/8"]

  bgp_settings {
    asn                 = 2468
    bgp_peering_address = "10.104.1.1"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMLocalNetworkGatewayConfig_bgpSettingsComplete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/8"]

  bgp_settings {
    asn                 = 2468
    bgp_peering_address = "10.104.1.1"
    peer_weight         = 15
  }
}
`, rInt, location, rInt)
}
