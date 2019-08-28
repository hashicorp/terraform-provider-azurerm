package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMVirtualNetworkPeering_basic(t *testing.T) {
	firstResourceName := "azurerm_virtual_network_peering.test1"
	secondResourceName := "azurerm_virtual_network_peering.test2"

	ri := tf.AccRandTimeInt()
	config := testAccAzureRMVirtualNetworkPeering_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkPeeringExists(firstResourceName),
					testCheckAzureRMVirtualNetworkPeeringExists(secondResourceName),
					resource.TestCheckResourceAttr(firstResourceName, "allow_virtual_network_access", "true"),
					resource.TestCheckResourceAttr(secondResourceName, "allow_virtual_network_access", "true"),
				),
			},
			{
				ResourceName:      firstResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkPeering_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	firstResourceName := "azurerm_virtual_network_peering.test1"
	secondResourceName := "azurerm_virtual_network_peering.test2"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkPeering_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkPeeringExists(firstResourceName),
					testCheckAzureRMVirtualNetworkPeeringExists(secondResourceName),
				),
			},
			{
				Config:      testAccAzureRMVirtualNetworkPeering_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_virtual_network_peering"),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkPeering_disappears(t *testing.T) {
	firstResourceName := "azurerm_virtual_network_peering.test1"
	secondResourceName := "azurerm_virtual_network_peering.test2"

	ri := tf.AccRandTimeInt()
	config := testAccAzureRMVirtualNetworkPeering_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkPeeringExists(firstResourceName),
					testCheckAzureRMVirtualNetworkPeeringExists(secondResourceName),
					resource.TestCheckResourceAttr(firstResourceName, "allow_virtual_network_access", "true"),
					resource.TestCheckResourceAttr(secondResourceName, "allow_virtual_network_access", "true"),
					testCheckAzureRMVirtualNetworkPeeringDisappears(firstResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkPeering_update(t *testing.T) {
	firstResourceName := "azurerm_virtual_network_peering.test1"
	secondResourceName := "azurerm_virtual_network_peering.test2"

	ri := tf.AccRandTimeInt()
	preConfig := testAccAzureRMVirtualNetworkPeering_basic(ri, testLocation())
	postConfig := testAccAzureRMVirtualNetworkPeering_basicUpdate(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkPeeringExists(firstResourceName),
					testCheckAzureRMVirtualNetworkPeeringExists(secondResourceName),
					resource.TestCheckResourceAttr(firstResourceName, "allow_virtual_network_access", "true"),
					resource.TestCheckResourceAttr(secondResourceName, "allow_virtual_network_access", "true"),
					resource.TestCheckResourceAttr(firstResourceName, "allow_forwarded_traffic", "false"),
					resource.TestCheckResourceAttr(secondResourceName, "allow_forwarded_traffic", "false"),
				),
			},

			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkPeeringExists(firstResourceName),
					testCheckAzureRMVirtualNetworkPeeringExists(secondResourceName),
					resource.TestCheckResourceAttr(firstResourceName, "allow_virtual_network_access", "true"),
					resource.TestCheckResourceAttr(secondResourceName, "allow_virtual_network_access", "true"),
					resource.TestCheckResourceAttr(firstResourceName, "allow_forwarded_traffic", "true"),
					resource.TestCheckResourceAttr(secondResourceName, "allow_forwarded_traffic", "true"),
				),
			},
		},
	})
}

func testCheckAzureRMVirtualNetworkPeeringExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		vnetName := rs.Primary.Attributes["virtual_network_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for virtual network peering: %s", name)
		}

		// Ensure resource group/virtual network peering combination exists in API
		client := testAccProvider.Meta().(*ArmClient).network.VnetPeeringsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, vnetName, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on vnetPeeringsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Virtual Network Peering %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMVirtualNetworkPeeringDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		vnetName := rs.Primary.Attributes["virtual_network_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for virtual network peering: %s", name)
		}

		// Ensure resource group/virtual network peering combination exists in API
		client := testAccProvider.Meta().(*ArmClient).network.VnetPeeringsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		future, err := client.Delete(ctx, resourceGroup, vnetName, name)
		if err != nil {
			return fmt.Errorf("Error deleting Peering %q (NW %q / RG %q): %+v", name, vnetName, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for deletion of Peering %q (NW %q / RG %q): %+v", name, vnetName, resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMVirtualNetworkPeeringDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).network.VnetPeeringsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_network_peering" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		vnetName := rs.Primary.Attributes["virtual_network_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, vnetName, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Virtual Network Peering sitll exists:\n%#v", resp.VirtualNetworkPeeringPropertiesFormat)
		}
	}

	return nil
}

func testAccAzureRMVirtualNetworkPeering_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test1" {
  name                = "acctestvirtnet-1-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.0.1.0/24"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctestvirtnet-2-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.0.2.0/24"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_virtual_network_peering" "test1" {
  name                         = "acctestpeer-1-%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  virtual_network_name         = "${azurerm_virtual_network.test1.name}"
  remote_virtual_network_id    = "${azurerm_virtual_network.test2.id}"
  allow_virtual_network_access = true
}

resource "azurerm_virtual_network_peering" "test2" {
  name                         = "acctestpeer-2-%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  virtual_network_name         = "${azurerm_virtual_network.test2.name}"
  remote_virtual_network_id    = "${azurerm_virtual_network.test1.id}"
  allow_virtual_network_access = true
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualNetworkPeering_requiresImport(rInt int, location string) string {
	template := testAccAzureRMVirtualNetworkPeering_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network_peering" "import" {
  name                         = "${azurerm_virtual_network_peering.test1.name}"
  resource_group_name          = "${azurerm_virtual_network_peering.test1.resource_group_name}"
  virtual_network_name         = "${azurerm_virtual_network_peering.test1.virtual_network_name}"
  remote_virtual_network_id    = "${azurerm_virtual_network_peering.test1.remote_virtual_network_id}"
  allow_virtual_network_access = "${azurerm_virtual_network_peering.test1.allow_virtual_network_access}"
}
`, template)
}

func testAccAzureRMVirtualNetworkPeering_basicUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test1" {
  name                = "acctestvirtnet-1-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.0.1.0/24"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctestvirtnet-2-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.0.2.0/24"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_virtual_network_peering" "test1" {
  name                         = "acctestpeer-1-%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  virtual_network_name         = "${azurerm_virtual_network.test1.name}"
  remote_virtual_network_id    = "${azurerm_virtual_network.test2.id}"
  allow_forwarded_traffic      = true
  allow_virtual_network_access = true
}

resource "azurerm_virtual_network_peering" "test2" {
  name                         = "acctestpeer-2-%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  virtual_network_name         = "${azurerm_virtual_network.test2.name}"
  remote_virtual_network_id    = "${azurerm_virtual_network.test1.id}"
  allow_forwarded_traffic      = true
  allow_virtual_network_access = true
}
`, rInt, location, rInt, rInt, rInt, rInt)
}
