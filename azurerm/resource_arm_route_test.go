package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMRoute_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoute_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteExists("azurerm_route.test"),
				),
			},
			{
				ResourceName:      "azurerm_route.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRoute_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoute_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteExists("azurerm_route.test"),
				),
			},
			{
				Config:      testAccAzureRMRoute_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_route"),
			},
		},
	})
}

func TestAccAzureRMRoute_update(t *testing.T) {
	resourceName := "azurerm_route.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoute_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "next_hop_type", "VnetLocal"),
					resource.TestCheckResourceAttr(resourceName, "next_hop_in_ip_address", ""),
				),
			},
			{
				Config: testAccAzureRMRoute_basicAppliance(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "next_hop_type", "VirtualAppliance"),
					resource.TestCheckResourceAttr(resourceName, "next_hop_in_ip_address", "192.168.0.1"),
				),
			},
			{
				Config: testAccAzureRMRoute_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "next_hop_type", "VnetLocal"),
					resource.TestCheckResourceAttr(resourceName, "next_hop_in_ip_address", ""),
				),
			},
		},
	})
}

func TestAccAzureRMRoute_disappears(t *testing.T) {
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMRoute_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteExists("azurerm_route.test"),
					testCheckAzureRMRouteDisappears("azurerm_route.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMRoute_multipleRoutes(t *testing.T) {
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMRoute_basic(ri, location)
	postConfig := testAccAzureRMRoute_multipleRoutes(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteExists("azurerm_route.test"),
				),
			},

			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteExists("azurerm_route.test1"),
				),
			},
		},
	})
}

func testCheckAzureRMRouteExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		rtName := rs.Primary.Attributes["route_table_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for route: %q", name)
		}

		client := testAccProvider.Meta().(*ArmClient).network.RoutesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, rtName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Route %q (resource group: %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on routesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMRouteDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		rtName := rs.Primary.Attributes["route_table_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for route: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).network.RoutesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		future, err := client.Delete(ctx, resourceGroup, rtName, name)
		if err != nil {
			return fmt.Errorf("Error deleting Route %q (Route Table %q / Resource Group %q): %+v", name, rtName, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for deletion of Route %q (Route Table %q / Resource Group %q): %+v", name, rtName, resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMRouteDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).network.RoutesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_route" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		rtName := rs.Primary.Attributes["route_table_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, rtName, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Route still exists:\n%#v", resp.RoutePropertiesFormat)
		}
	}

	return nil
}

func testAccAzureRMRoute_basic(rInt int, location string) string {
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

resource "azurerm_route" "test" {
  name                = "acctestroute%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  route_table_name    = "${azurerm_route_table.test.name}"

  address_prefix = "10.1.0.0/16"
  next_hop_type  = "vnetlocal"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMRoute_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s
resource "azurerm_route" "import" {
  name                = "${azurerm_route.test.name}"
  resource_group_name = "${azurerm_route.test.resource_group_name}"
  route_table_name    = "${azurerm_route.test.route_table_name}"

  address_prefix = "${azurerm_route.test.address_prefix}"
  next_hop_type  = "${azurerm_route.test.next_hop_type}"
}
`, testAccAzureRMRoute_basic(rInt, location))
}

func testAccAzureRMRoute_basicAppliance(rInt int, location string) string {
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

resource "azurerm_route" "test" {
  name                = "acctestroute%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  route_table_name    = "${azurerm_route_table.test.name}"

  address_prefix         = "10.1.0.0/16"
  next_hop_type          = "VirtualAppliance"
  next_hop_in_ip_address = "192.168.0.1"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMRoute_multipleRoutes(rInt int, location string) string {
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

resource "azurerm_route" "test1" {
  name                = "acctestroute%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  route_table_name    = "${azurerm_route_table.test.name}"

  address_prefix = "10.2.0.0/16"
  next_hop_type  = "none"
}
`, rInt, location, rInt, rInt)
}
