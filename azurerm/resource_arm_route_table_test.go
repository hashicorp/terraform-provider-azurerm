package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMRouteTable_basic(t *testing.T) {
	resourceName := "azurerm_route_table.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRouteTable_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "disable_bgp_route_propagation", "false"),
					resource.TestCheckResourceAttr(resourceName, "route.#", "0"),
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

func TestAccAzureRMRouteTable_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_route_table.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRouteTable_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "disable_bgp_route_propagation", "false"),
					resource.TestCheckResourceAttr(resourceName, "route.#", "0"),
				),
			},
			{
				Config:      testAccAzureRMRouteTable_requiresImport(ri, acceptance.Location()),
				ExpectError: acceptance.RequiresImportError("azurerm_route_table"),
			},
		},
	})
}

func TestAccAzureRMRouteTable_complete(t *testing.T) {
	resourceName := "azurerm_route_table.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRouteTable_complete(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists("azurerm_route_table.test"),
					resource.TestCheckResourceAttr(resourceName, "disable_bgp_route_propagation", "true"),
					resource.TestCheckResourceAttr(resourceName, "route.#", "1"),
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

func TestAccAzureRMRouteTable_update(t *testing.T) {
	resourceName := "azurerm_route_table.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRouteTable_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists("azurerm_route_table.test"),
					resource.TestCheckResourceAttr(resourceName, "disable_bgp_route_propagation", "false"),
					resource.TestCheckResourceAttr(resourceName, "route.#", "0"),
				),
			},
			{
				Config: testAccAzureRMRouteTable_basicAppliance(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists("azurerm_route_table.test"),
					resource.TestCheckResourceAttr(resourceName, "disable_bgp_route_propagation", "false"),
					resource.TestCheckResourceAttr(resourceName, "route.#", "1"),
				),
			},
			{
				Config: testAccAzureRMRouteTable_complete(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists("azurerm_route_table.test"),
					resource.TestCheckResourceAttr(resourceName, "disable_bgp_route_propagation", "true"),
					resource.TestCheckResourceAttr(resourceName, "route.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMRouteTable_singleRoute(t *testing.T) {
	resourceName := "azurerm_route_table.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMRouteTable_singleRoute(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(resourceName),
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

func TestAccAzureRMRouteTable_removeRoute(t *testing.T) {
	resourceName := "azurerm_route_table.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMRouteTable_singleRoute(ri, acceptance.Location())
	noBlocksConfig := testAccAzureRMRouteTable_noRouteBlocks(ri, acceptance.Location())
	blocksEmptyConfig := testAccAzureRMRouteTable_singleRouteRemoved(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				// This configuration includes a single explicit route block
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "route.#", "1"),
				),
			},
			{
				// This configuration has no route blocks at all.
				Config: noBlocksConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(resourceName),
					// The route from the first step is preserved because no
					// blocks at all means "ignore existing blocks".
					resource.TestCheckResourceAttr(resourceName, "route.#", "1"),
				),
			},
			{
				// This configuration sets route to [] explicitly using the
				// attribute syntax.
				Config: blocksEmptyConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(resourceName),
					// The route from the first step is now removed, leaving us
					// with no routes at all.
					resource.TestCheckResourceAttr(resourceName, "route.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMRouteTable_disappears(t *testing.T) {
	resourceName := "azurerm_route_table.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMRouteTable_basic(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(resourceName),
					testCheckAzureRMRouteTableDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMRouteTable_withTags(t *testing.T) {
	resourceName := "azurerm_route_table.test"
	ri := tf.AccRandTimeInt()
	preConfig := testAccAzureRMRouteTable_withTags(ri, acceptance.Location())
	postConfig := testAccAzureRMRouteTable_withTagsUpdate(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(resourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "staging"),
				),
			},
		},
	})
}

func TestAccAzureRMRouteTable_multipleRoutes(t *testing.T) {
	resourceName := "azurerm_route_table.test"
	ri := tf.AccRandTimeInt()
	preConfig := testAccAzureRMRouteTable_singleRoute(ri, acceptance.Location())
	postConfig := testAccAzureRMRouteTable_multipleRoutes(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "route.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "route.0.name", "route1"),
					resource.TestCheckResourceAttr(resourceName, "route.0.address_prefix", "10.1.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "route.0.next_hop_type", "VnetLocal"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "route.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "route.0.name", "route1"),
					resource.TestCheckResourceAttr(resourceName, "route.0.address_prefix", "10.1.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "route.0.next_hop_type", "VnetLocal"),
					resource.TestCheckResourceAttr(resourceName, "route.1.name", "route2"),
					resource.TestCheckResourceAttr(resourceName, "route.1.address_prefix", "10.2.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "route.1.next_hop_type", "VnetLocal"),
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

func TestAccAzureRMRouteTable_withTagsSubnet(t *testing.T) {
	ri := tf.AccRandTimeInt()
	configSetup := testAccAzureRMRouteTable_withTagsSubnet(ri, acceptance.Location())
	configTest := testAccAzureRMRouteTable_withAddTagsSubnet(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: configSetup,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists("azurerm_route_table.test"),
					testCheckAzureRMSubnetExists("azurerm_subnet.subnet1"),
					resource.TestCheckResourceAttrSet("azurerm_subnet.subnet1", "route_table_id"),
				),
			},
			{
				Config: configTest,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists("azurerm_route_table.test"),
					testCheckAzureRMSubnetExists("azurerm_subnet.subnet1"),
					resource.TestCheckResourceAttrSet("azurerm_subnet.subnet1", "route_table_id"),
				),
			},
		},
	})
}

func testCheckAzureRMRouteTableExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for route table: %q", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.RouteTablesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Route Table %q (resource group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on routeTablesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMRouteTableDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for route table: %q", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.RouteTablesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		future, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			if !response.WasNotFound(future.Response()) {
				return fmt.Errorf("Error deleting Route Table %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for deletion of Route Table %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMRouteTableDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.RouteTablesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_route_table" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Route Table still exists:\n%#v", resp.RouteTablePropertiesFormat)
	}

	return nil
}

func testAccAzureRMRouteTable_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}

func testAccAzureRMRouteTable_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_route_table" "import" {
  name                = "${azurerm_route_table.test.name}"
  location            = "${azurerm_route_table.test.location}"
  resource_group_name = "${azurerm_route_table.test.resource_group_name}"
}
`, testAccAzureRMRouteTable_basic(rInt, location))
}

func testAccAzureRMRouteTable_basicAppliance(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  route {
    name                   = "route1"
    address_prefix         = "10.1.0.0/16"
    next_hop_type          = "VirtualAppliance"
    next_hop_in_ip_address = "192.168.0.1"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMRouteTable_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  route {
    name           = "acctestRoute"
    address_prefix = "10.1.0.0/16"
    next_hop_type  = "vnetlocal"
  }

  disable_bgp_route_propagation = true
}
`, rInt, location, rInt)
}

func testAccAzureRMRouteTable_singleRoute(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  route {
    name           = "route1"
    address_prefix = "10.1.0.0/16"
    next_hop_type  = "vnetlocal"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMRouteTable_noRouteBlocks(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}

func testAccAzureRMRouteTable_singleRouteRemoved(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  route = []
}
`, rInt, location, rInt)
}

func testAccAzureRMRouteTable_multipleRoutes(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  route {
    name           = "route1"
    address_prefix = "10.1.0.0/16"
    next_hop_type  = "vnetlocal"
  }

  route {
    name           = "route2"
    address_prefix = "10.2.0.0/16"
    next_hop_type  = "vnetlocal"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMRouteTable_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  route {
    name           = "route1"
    address_prefix = "10.1.0.0/16"
    next_hop_type  = "vnetlocal"
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMRouteTable_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  route {
    name           = "route1"
    address_prefix = "10.1.0.0/16"
    next_hop_type  = "vnetlocal"
  }

  tags = {
    environment = "staging"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMRouteTable_withTagsSubnet(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.0.0.0/16"]

  tags = {
    environment = "staging"
  }
}

resource "azurerm_subnet" "subnet1" {
  name                 = "subnet1"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
  route_table_id       = "${azurerm_route_table.test.id}"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  route {
    name           = "route1"
    address_prefix = "10.1.0.0/16"
    next_hop_type  = "vnetlocal"
  }

  tags = {
    environment = "staging"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMRouteTable_withAddTagsSubnet(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"

  tags = {
    environment = "staging"
    cloud       = "Azure"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.0.0.0/16"]

  tags = {
    environment = "staging"
    cloud       = "Azure"
  }
}

resource "azurerm_subnet" "subnet1" {
  name                 = "subnet1"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
  route_table_id       = "${azurerm_route_table.test.id}"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  route {
    name           = "route1"
    address_prefix = "10.1.0.0/16"
    next_hop_type  = "vnetlocal"
  }

  tags = {
    environment = "staging"
    cloud       = "Azure"
  }
}
`, rInt, location, rInt, rInt)
}
