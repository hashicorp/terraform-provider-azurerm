package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMRouteTable_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMRouteTable_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists("azurerm_route_table.test"),
				),
			},
		},
	})
}

func TestAccAzureRMRouteTable_singleRoute(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMRouteTable_singleRoute(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists("azurerm_route_table.test"),
				),
			},
		},
	})
}

func TestAccAzureRMRouteTable_removeRoute(t *testing.T) {
	resourceName := "azurerm_route_table.test"
	ri := acctest.RandInt()
	config := testAccAzureRMRouteTable_singleRoute(ri, testLocation())
	updatedConfig := testAccAzureRMRouteTable_singleRouteRemoved(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "route.#", "1"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "route.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMRouteTable_disappears(t *testing.T) {
	resourceName := "azurerm_route_table.test"
	ri := acctest.RandInt()
	config := testAccAzureRMRouteTable_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
	ri := acctest.RandInt()
	preConfig := testAccAzureRMRouteTable_withTags(ri, testLocation())
	postConfig := testAccAzureRMRouteTable_withTagsUpdate(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
	ri := acctest.RandInt()
	preConfig := testAccAzureRMRouteTable_singleRoute(ri, testLocation())
	postConfig := testAccAzureRMRouteTable_multipleRoutes(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
		},
	})
}

func TestAccAzureRMRouteTable_withTagsSubnet(t *testing.T) {
	ri := acctest.RandInt()
	configSetup := testAccAzureRMRouteTable_withTagsSubnet(ri, testLocation())
	configTest := testAccAzureRMRouteTable_withAddTagsSubnet(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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

func testCheckAzureRMRouteTableExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %q", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for route table: %q", name)
		}

		client := testAccProvider.Meta().(*ArmClient).routeTablesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testCheckAzureRMRouteTableDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %q", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for route table: %q", name)
		}

		client := testAccProvider.Meta().(*ArmClient).routeTablesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		future, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			if !response.WasNotFound(future.Response()) {
				return fmt.Errorf("Error deleting Route Table %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		err = future.WaitForCompletion(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("Error waiting for deletion of Route Table %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMRouteTableDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).routeTablesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_route_table" "test" {
    name = "acctestrt%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}

func testAccAzureRMRouteTable_singleRoute(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_route_table" "test" {
    name = "acctestrt%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    route {
    	name = "route1"
		address_prefix = "10.1.0.0/16"
		next_hop_type = "vnetlocal"
    }
}
`, rInt, location, rInt)
}

func testAccAzureRMRouteTable_singleRouteRemoved(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_route_table" "test" {
    name = "acctestrt%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
	route = []
}
`, rInt, location, rInt)
}

func testAccAzureRMRouteTable_multipleRoutes(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_route_table" "test" {
    name = "acctestrt%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    route {
    	name = "route1"
		address_prefix = "10.1.0.0/16"
		next_hop_type = "vnetlocal"
    }

    route {
    	name = "route2"
		address_prefix = "10.2.0.0/16"
		next_hop_type = "vnetlocal"
    }
}
`, rInt, location, rInt)
}

func testAccAzureRMRouteTable_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_route_table" "test" {
    name = "acctestrt%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    route {
    	name = "route1"
    	address_prefix = "10.1.0.0/16"
    	next_hop_type = "vnetlocal"
    }

    tags {
		environment = "Production"
		cost_center = "MSFT"
    }
}
`, rInt, location, rInt)
}

func testAccAzureRMRouteTable_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_route_table" "test" {
    name = "acctestrt%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    route {
    	name = "route1"
    	address_prefix = "10.1.0.0/16"
    	next_hop_type = "vnetlocal"
    }

    tags {
		environment = "staging"
    }
}
`, rInt, location, rInt)
}

func testAccAzureRMRouteTable_withTagsSubnet(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name 		= "acctestRG-%d"
	location 	= "%s"
    tags {
		environment = "staging"
    }
}

resource "azurerm_virtual_network" "test" {
    name 				= "acctestvirtnet%d"
    location 			= "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    address_space 		= ["10.0.0.0/16"]
 
	tags {
		environment = "staging"
    }
}
  
resource "azurerm_subnet" "subnet1" {
	name                 = "subnet1"
    resource_group_name = "${azurerm_resource_group.test.name}"
	virtual_network_name = "${azurerm_virtual_network.test.name}"
	address_prefix       = "10.0.1.0/24"
	route_table_id       = "${azurerm_route_table.test.id}"
}

resource "azurerm_route_table" "test" {
    name 				= "acctestrt%d"
    location 			= "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    route {
    	name = "route1"
    	address_prefix = "10.1.0.0/16"
    	next_hop_type = "vnetlocal"
    }

    tags {
		environment = "staging"
    }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMRouteTable_withAddTagsSubnet(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name 		= "acctestRG-%d"
	location 	= "%s"
    tags {
		environment = "staging"
		cloud = "Azure"
    }
}

resource "azurerm_virtual_network" "test" {
    name 				= "acctestvirtnet%d"
    location 			= "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    address_space 		= ["10.0.0.0/16"]
 
	tags {
		environment = "staging"
		cloud = "Azure"
    }
}
  
resource "azurerm_subnet" "subnet1" {
	name                 = "subnet1"
    resource_group_name = "${azurerm_resource_group.test.name}"
	virtual_network_name = "${azurerm_virtual_network.test.name}"
	address_prefix       = "10.0.1.0/24"
	route_table_id       = "${azurerm_route_table.test.id}"
}

resource "azurerm_route_table" "test" {
    name 				= "acctestrt%d"
    location 			= "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    route {
    	name = "route1"
    	address_prefix = "10.1.0.0/16"
    	next_hop_type = "vnetlocal"
    }

    tags {
		environment = "staging"
		cloud = "Azure"
    }
}
`, rInt, location, rInt, rInt)
}
