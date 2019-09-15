package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-06-01/network"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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
			"tierUpdate":                   testAccAzureRMExpressRouteCircuit_tierUpdate,
			"premiumMetered":               testAccAzureRMExpressRouteCircuit_premiumMetered,
			"premiumUnlimited":             testAccAzureRMExpressRouteCircuit_premiumUnlimited,
			"allowClassicOperationsUpdate": testAccAzureRMExpressRouteCircuit_allowClassicOperationsUpdate,
			"requiresImport":               testAccAzureRMExpressRouteCircuit_requiresImport,
			"data_basic":                   testAccDataSourceAzureRMExpressRoute_basicMetered,
		},
		"PrivatePeering": {
			"azurePrivatePeering":           testAccAzureRMExpressRouteCircuitPeering_azurePrivatePeering,
			"azurePrivatePeeringWithUpdate": testAccAzureRMExpressRouteCircuitPeering_azurePrivatePeeringWithCircuitUpdate,
			"requiresImport":                testAccAzureRMExpressRouteCircuitPeering_requiresImport,
		},
		"MicrosoftPeering": {
			"microsoftPeering": testAccAzureRMExpressRouteCircuitPeering_microsoftPeering,
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
	resourceName := "azurerm_express_route_circuit.test"
	var erc network.ExpressRouteCircuit
	ri := tf.AccRandTimeInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuit_basicMeteredConfig(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists(resourceName, &erc),
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

func testAccAzureRMExpressRouteCircuit_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_express_route_circuit.test"
	var erc network.ExpressRouteCircuit
	ri := tf.AccRandTimeInt()

	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuit_basicMeteredConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists(resourceName, &erc),
				),
			},
			{
				Config:      testAccAzureRMExpressRouteCircuit_requiresImportConfig(ri, location),
				ExpectError: testRequiresImportError("azurerm_express_route_circuit"),
			},
		},
	})
}

func testAccAzureRMExpressRouteCircuit_basicUnlimited(t *testing.T) {
	resourceName := "azurerm_express_route_circuit.test"
	var erc network.ExpressRouteCircuit
	ri := tf.AccRandTimeInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuit_basicUnlimitedConfig(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists(resourceName, &erc),
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

func testAccAzureRMExpressRouteCircuit_update(t *testing.T) {
	resourceName := "azurerm_express_route_circuit.test"
	var erc network.ExpressRouteCircuit
	ri := tf.AccRandTimeInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuit_basicMeteredConfig(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists(resourceName, &erc),
					resource.TestCheckResourceAttr(resourceName, "sku.0.family", "MeteredData"),
				),
			},
			{
				Config: testAccAzureRMExpressRouteCircuit_basicUnlimitedConfig(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists(resourceName, &erc),
					resource.TestCheckResourceAttr(resourceName, "sku.0.family", "UnlimitedData"),
				),
			},
		},
	})
}

func testAccAzureRMExpressRouteCircuit_tierUpdate(t *testing.T) {
	resourceName := "azurerm_express_route_circuit.test"
	var erc network.ExpressRouteCircuit
	ri := tf.AccRandTimeInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuit_sku(ri, testLocation(), "Standard", "MeteredData"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists("azurerm_express_route_circuit.test", &erc),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "Standard"),
				),
			},
			{
				Config: testAccAzureRMExpressRouteCircuit_sku(ri, testLocation(), "Premium", "MeteredData"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists("azurerm_express_route_circuit.test", &erc),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "Premium"),
				),
			},
		},
	})
}

func testAccAzureRMExpressRouteCircuit_premiumMetered(t *testing.T) {
	resourceName := "azurerm_express_route_circuit.test"
	var erc network.ExpressRouteCircuit
	ri := tf.AccRandTimeInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuit_sku(ri, testLocation(), "Premium", "MeteredData"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists("azurerm_express_route_circuit.test", &erc),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "Premium"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.family", "MeteredData"),
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

func testAccAzureRMExpressRouteCircuit_premiumUnlimited(t *testing.T) {
	resourceName := "azurerm_express_route_circuit.test"
	var erc network.ExpressRouteCircuit
	ri := tf.AccRandTimeInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuit_sku(ri, testLocation(), "Premium", "UnlimitedData"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists("azurerm_express_route_circuit.test", &erc),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "Premium"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.family", "UnlimitedData"),
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

func testAccAzureRMExpressRouteCircuit_allowClassicOperationsUpdate(t *testing.T) {
	resourceName := "azurerm_express_route_circuit.test"
	var erc network.ExpressRouteCircuit
	ri := tf.AccRandTimeInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuit_allowClassicOperations(ri, testLocation(), "false"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists("azurerm_express_route_circuit.test", &erc),
					resource.TestCheckResourceAttr(resourceName, "allow_classic_operations", "false"),
				),
			},
			{
				Config: testAccAzureRMExpressRouteCircuit_allowClassicOperations(ri, testLocation(), "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitExists("azurerm_express_route_circuit.test", &erc),
					resource.TestCheckResourceAttr(resourceName, "allow_classic_operations", "true"),
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

		client := testAccProvider.Meta().(*ArmClient).network.ExpressRouteCircuitsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
	client := testAccProvider.Meta().(*ArmClient).network.ExpressRouteCircuitsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAccAzureRMExpressRouteCircuit_basicMeteredConfig(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
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
`, rInt, location, rInt)
}

func testAccAzureRMExpressRouteCircuit_requiresImportConfig(rInt int, location string) string {
	template := testAccAzureRMExpressRouteCircuit_basicMeteredConfig(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_circuit" "import" {
  name                  = "${azurerm_express_route_circuit.test.name}"
  location              = "${azurerm_express_route_circuit.test.location}"
  resource_group_name   = "${azurerm_express_route_circuit.test.resource_group_name}"
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

func testAccAzureRMExpressRouteCircuit_basicUnlimitedConfig(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
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
`, rInt, location, rInt)
}

func testAccAzureRMExpressRouteCircuit_sku(rInt int, location string, tier string, family string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
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
`, rInt, location, rInt, tier, family)
}

func testAccAzureRMExpressRouteCircuit_allowClassicOperations(rInt int, location string, allowClassicOperations string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
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
`, rInt, location, rInt, allowClassicOperations)
}
