package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMSubnet_basic(t *testing.T) {

	ri := acctest.RandInt()
	config := testAccAzureRMSubnet_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists("azurerm_subnet.test"),
				),
			},
		},
	})
}

func TestAccAzureRMSubnet_routeTableUpdate(t *testing.T) {

	ri := acctest.RandInt()
	location := testLocation()
	initConfig := testAccAzureRMSubnet_routeTable(ri, location)
	updatedConfig := testAccAzureRMSubnet_updatedRouteTable(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: initConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists("azurerm_subnet.test"),
				),
			},

			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetRouteTableExists("azurerm_subnet.test", fmt.Sprintf("acctest-%d", ri)),
				),
			},
		},
	})
}

func TestAccAzureRMSubnet_bug7986(t *testing.T) {
	ri := acctest.RandInt()
	initConfig := testAccAzureRMSubnet_bug7986(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: initConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists("azurerm_subnet.first"),
					testCheckAzureRMSubnetExists("azurerm_subnet.second"),
				),
			},
		},
	})
}

func TestAccAzureRMSubnet_bug15204(t *testing.T) {
	ri := acctest.RandInt()
	initConfig := testAccAzureRMSubnet_bug15204(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: initConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists("azurerm_subnet.test"),
				),
			},
		},
	})
}

func TestAccAzureRMSubnet_disappears(t *testing.T) {

	ri := acctest.RandInt()
	config := testAccAzureRMSubnet_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists("azurerm_subnet.test"),
					testCheckAzureRMSubnetDisappears("azurerm_subnet.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMSubnetExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		log.Printf("[INFO] Checking Subnet addition.")

		name := rs.Primary.Attributes["name"]
		vnetName := rs.Primary.Attributes["virtual_network_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for subnet: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).subnetClient

		resp, err := conn.Get(resourceGroup, vnetName, name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on subnetClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Subnet %q (resource group: %q) does not exist", name, resourceGroup)
		}

		if resp.RouteTable == nil {
			return fmt.Errorf("Bad: Subnet %q (resource group: %q) does not contain route tables after add", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMSubnetRouteTableExists(subnetName string, routeTableId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[subnetName]
		if !ok {
			return fmt.Errorf("Not found: %s", subnetName)
		}

		log.Printf("[INFO] Checking Subnet update.")

		name := rs.Primary.Attributes["name"]
		vnetName := rs.Primary.Attributes["virtual_network_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for subnet: %s", name)
		}

		vnetConn := testAccProvider.Meta().(*ArmClient).vnetClient
		vnetResp, vnetErr := vnetConn.Get(resourceGroup, vnetName, "")
		if vnetErr != nil {
			return fmt.Errorf("Bad: Get on vnetClient: %+v", vnetErr)
		}

		if vnetResp.Subnets == nil {
			return fmt.Errorf("Bad: Vnet %q (resource group: %q) does not have subnets after update", vnetName, resourceGroup)
		}

		conn := testAccProvider.Meta().(*ArmClient).subnetClient

		resp, err := conn.Get(resourceGroup, vnetName, name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on subnetClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Subnet %q (resource group: %q) does not exist", subnetName, resourceGroup)
		}

		if resp.RouteTable == nil {
			return fmt.Errorf("Bad: Subnet %q (resource group: %q) does not contain route tables after update", subnetName, resourceGroup)
		}

		if !strings.Contains(*resp.RouteTable.ID, routeTableId) {
			return fmt.Errorf("Bad: Subnet %q (resource group: %q) does not have route table %q", subnetName, resourceGroup, routeTableId)
		}

		return nil
	}
}

func testCheckAzureRMSubnetDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		vnetName := rs.Primary.Attributes["virtual_network_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for subnet: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).subnetClient

		_, error := conn.Delete(resourceGroup, vnetName, name, make(chan struct{}))
		err := <-error
		if err != nil {
			return fmt.Errorf("Bad: Delete on subnetClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMSubnetDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).subnetClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_subnet" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		vnetName := rs.Primary.Attributes["virtual_network_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, vnetName, name, "")

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Subnet still exists:\n%#v", resp.SubnetPropertiesFormat)
		}
	}

	return nil
}

func testAccAzureRMSubnet_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
  route_table_id       = "${azurerm_route_table.test.id}"
}

resource "azurerm_route_table" "test" {
  name                = "acctestroutetable%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_route" "test" {
  name                = "acctestroute%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  route_table_name    = "${azurerm_route_table.test.name}"

  address_prefix         = "10.100.0.0/14"
  next_hop_type          = "VirtualAppliance"
  next_hop_in_ip_address = "10.10.1.1"
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMSubnet_routeTable(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
  route_table_id       = "${azurerm_route_table.test.id}"
}

resource "azurerm_route_table" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_route" "route_a" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  route_table_name    = "${azurerm_route_table.test.name}"

  address_prefix         = "10.100.0.0/14"
  next_hop_type          = "VirtualAppliance"
  next_hop_in_ip_address = "10.10.1.1"
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMSubnet_updatedRouteTable(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"

  tags {
    environment = "Testing"
  }
}

resource "azurerm_network_security_group" "test_secgroup" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  security_rule {
    name                       = "acctest-%d"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  tags {
    environment = "Testing"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags {
    environment = "Testing"
  }
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
  route_table_id       = "${azurerm_route_table.test.id}"
}

resource "azurerm_route_table" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags {
    environment = "Testing"
  }
}

resource "azurerm_route" "route_a" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  route_table_name    = "${azurerm_route_table.test.name}"

  address_prefix         = "10.100.0.0/14"
  next_hop_type          = "VirtualAppliance"
  next_hop_in_ip_address = "10.10.1.1"
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMSubnet_bug7986(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest%d-rg"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest%d-vn"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_route_table" "first" {
  name                = "acctest%d-private-1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_route" "first" {
  name                = "acctest%d-private-1"
  resource_group_name = "${azurerm_resource_group.test.name}"
  route_table_name    = "${azurerm_route_table.first.name}"
  address_prefix      = "0.0.0.0/0"
  next_hop_type       = "None"
}

resource "azurerm_subnet" "first" {
  name                 = "acctest%d-private-1"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.0.0/24"
  route_table_id       = "${azurerm_route_table.first.id}"
}

resource "azurerm_route_table" "second" {
  name                = "acctest%d-private-2"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_route" "second" {
  name                = "acctest%d-private-2"
  resource_group_name = "${azurerm_resource_group.test.name}"
  route_table_name    = "${azurerm_route_table.second.name}"
  address_prefix      = "0.0.0.0/0"
  next_hop_type       = "None"
}

resource "azurerm_subnet" "second" {
  name                 = "acctest%d-private-2"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
  route_table_id       = "${azurerm_route_table.second.id}"
}`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMSubnet_bug15204(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.85.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_network_security_group" "test" {
  name = "acctestnsg-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                      = "acctestsubnet-%d"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  virtual_network_name      = "${azurerm_virtual_network.test.name}"
  address_prefix            = "10.85.9.0/24"
  route_table_id            = "${azurerm_route_table.test.id}"
  network_security_group_id = "${azurerm_network_security_group.test.id}"
}
`, rInt, location, rInt, rInt, rInt, rInt)
}
