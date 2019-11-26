package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	networkSvc "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVirtualHubConnection_basic(t *testing.T) {
	resourceName := "azurerm_virtual_hub_connection.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualHubConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubConnection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubConnectionExists(resourceName),
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

func TestAccAzureRMVirtualHubConnection_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_virtual_hub_connection.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualHubConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubConnection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubConnectionExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMVirtualHubConnection_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_virtual_hub_connection"),
			},
		},
	})
}

func TestAccAzureRMVirtualHubConnection_complete(t *testing.T) {
	resourceName := "azurerm_virtual_hub_connection.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualHubConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubConnection_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubConnectionExists(resourceName),
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

func TestAccAzureRMVirtualHubConnection_allowHubToRemoteVnetTransit(t *testing.T) {
	resourceName := "azurerm_virtual_hub_connection.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualHubConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubConnection_allowHubToRemoteVnetTransit(ri, location, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubConnectionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMVirtualHubConnection_allowHubToRemoteVnetTransit(ri, location, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubConnectionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMVirtualHubConnection_allowHubToRemoteVnetTransit(ri, location, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubConnectionExists(resourceName),
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

func TestAccAzureRMVirtualHubConnection_allowRemoteVnetToUseHubVnetGateways(t *testing.T) {
	resourceName := "azurerm_virtual_hub_connection.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualHubConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubConnection_allowRemoteVnetToUseHubVnetGateways(ri, location, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubConnectionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMVirtualHubConnection_allowRemoteVnetToUseHubVnetGateways(ri, location, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubConnectionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMVirtualHubConnection_allowRemoteVnetToUseHubVnetGateways(ri, location, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubConnectionExists(resourceName),
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

func TestAccAzureRMVirtualHubConnection_enableInternetSecurity(t *testing.T) {
	resourceName := "azurerm_virtual_hub_connection.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualHubConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubConnection_enableInternetSecurity(ri, location, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubConnectionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMVirtualHubConnection_enableInternetSecurity(ri, location, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubConnectionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMVirtualHubConnection_enableInternetSecurity(ri, location, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubConnectionExists(resourceName),
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

func testCheckAzureRMVirtualHubConnectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Virtual Hub Connection not found: %s", resourceName)
		}

		virtualHubId := rs.Primary.Attributes["virtual_hub_id"]
		id, err := networkSvc.ParseVirtualHubID(virtualHubId)
		if err != nil {
			return err
		}

		resourceGroup := id.Base.ResourceGroup
		hubName := id.Name
		name := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).Network.VirtualHubClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, hubName)
		if err != nil {
			return fmt.Errorf("Bad: Get on network.VirtualHubClient: %+v", err)
		}

		if resp.VirtualHubProperties == nil {
			return fmt.Errorf("VirtualHubProperties was nil!")
		}

		props := *resp.VirtualHubProperties
		if props.VirtualNetworkConnections == nil {
			return fmt.Errorf("props.VirtualNetworkConnections was nil")
		}

		conns := *props.VirtualNetworkConnections

		found := false
		for _, conn := range conns {
			if conn.Name != nil && *conn.Name == name {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Connection %q was not found", name)
		}

		return nil
	}
}

func testCheckAzureRMVirtualHubConnectionDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_hub_connection" {
			continue
		}

		virtualHubId := rs.Primary.Attributes["virtual_hub_id"]
		id, err := networkSvc.ParseVirtualHubID(virtualHubId)
		if err != nil {
			return err
		}

		resourceGroup := id.Base.ResourceGroup
		hubName := id.Name
		name := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).Network.VirtualHubClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, hubName)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on network.VirtualHubClient: %+v", err)
			}
		}

		if resp.VirtualHubProperties == nil {
			return fmt.Errorf("VirtualHubProperties was nil!")
		}

		props := *resp.VirtualHubProperties
		if props.VirtualNetworkConnections == nil {
			return fmt.Errorf("props.VirtualNetworkConnections was nil")
		}

		conns := *props.VirtualNetworkConnections

		for _, conn := range conns {
			if conn.Name != nil && *conn.Name == name {
				return fmt.Errorf("Connection %q still exists", name)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMVirtualHubConnection_basic(rInt int, location string) string {
	template := testAccAzureRMVirtualHubConnection_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctestvhub-%d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id
}
`, template, rInt)
}

func testAccAzureRMVirtualHubConnection_requiresImport(rInt int, location string) string {
	template := testAccAzureRMVirtualHubConnection_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_connection" "import" {
  name                      = azurerm_virtual_hub_connection.test.name
  virtual_hub_id            = azurerm_virtual_hub_connection.test.virtual_hub_id
  remote_virtual_network_id = azurerm_virtual_hub_connection.test.remote_virtual_network_id
}
`, template)
}

func testAccAzureRMVirtualHubConnection_allowHubToRemoteVnetTransit(rInt int, location string, enabled bool) string {
	template := testAccAzureRMVirtualHubConnection_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_connection" "test" {
  name                             = "acctestvhub-%d"
  virtual_hub_id                   = azurerm_virtual_hub.test.id
  remote_virtual_network_id        = azurerm_virtual_network.test.id
  allow_hub_to_remote_vnet_transit = %t
}
`, template, rInt, enabled)
}

func testAccAzureRMVirtualHubConnection_allowRemoteVnetToUseHubVnetGateways(rInt int, location string, enabled bool) string {
	template := testAccAzureRMVirtualHubConnection_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_connection" "test" {
  name                                       = "acctestvhub-%d"
  virtual_hub_id                             = azurerm_virtual_hub.test.id
  remote_virtual_network_id                  = azurerm_virtual_network.test.id
  allow_remote_vnet_to_use_hub_vnet_gateways = %t
}
`, template, rInt, enabled)
}

func testAccAzureRMVirtualHubConnection_enableInternetSecurity(rInt int, location string, enabled bool) string {
	template := testAccAzureRMVirtualHubConnection_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctestvhub-%d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id
  enable_internet_security  = %t
}
`, template, rInt, enabled)
}

func testAccAzureRMVirtualHubConnection_complete(rInt int, location string) string {
	template := testAccAzureRMVirtualHubConnection_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_connection" "test" {
  name                                       = "acctestvhub-%d"
  virtual_hub_id                             = azurerm_virtual_hub.test.id
  remote_virtual_network_id                  = azurerm_virtual_network.test.id
  allow_hub_to_remote_vnet_transit           = true
  allow_remote_vnet_to_use_hub_vnet_gateways = true
  enable_internet_security                   = true
}
`, template, rInt)
}

func testAccAzureRMVirtualHubConnection_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["172.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctestvhub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"
}
`, rInt, location, rInt, rInt, rInt)
}
