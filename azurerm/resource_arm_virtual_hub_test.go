package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVirtualHub_basic(t *testing.T) {
	resourceName := "azurerm_virtual_hub.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHub_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(resourceName),
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

func TestAccAzureRMVirtualHub_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_virtual_hub.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHub_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMVirtualHub_requiresImport(ri, acceptance.Location()),
				ExpectError: acceptance.RequiresImportError("azurerm_virtual_hub"),
			},
		},
	})
}

func TestAccAzureRMVirtualHub_update(t *testing.T) {
	resourceName := "azurerm_virtual_hub.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHub_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMVirtualHub_updated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(resourceName),
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

func TestAccAzureRMVirtualHub_routes(t *testing.T) {
	resourceName := "azurerm_virtual_hub.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHub_route(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMVirtualHub_routeUpdated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(resourceName),
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

func TestAccAzureRMVirtualHub_tags(t *testing.T) {
	resourceName := "azurerm_virtual_hub.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHub_tags(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(resourceName),
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

func testCheckAzureRMVirtualHubExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Virtual Hub not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VirtualHubClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Virtual Hub %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on network.VirtualHubClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMVirtualHubDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VirtualHubClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_hub" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on network.VirtualHubClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMVirtualHub_basic(rInt int, location string) string {
	template := testAccAzureRMVirtualHub_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub" "test" {
  name                = "acctestVHUB-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"
}
`, template, rInt)
}

func testAccAzureRMVirtualHub_requiresImport(rInt int, location string) string {
	template := testAccAzureRMVirtualHub_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub" "import" {
  name                = azurerm_virtual_hub.test.name
  location            = azurerm_virtual_hub.test.location
  resource_group_name = azurerm_virtual_hub.test.resource_group_name
  address_prefix      = "10.0.1.0/24"
}
`, template)
}

func testAccAzureRMVirtualHub_updated(rInt int, location string) string {
	template := testAccAzureRMVirtualHub_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub" "test" {
  name                = "acctestVHUB-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.2.0/24"
}
`, template, rInt)
}

func testAccAzureRMVirtualHub_route(rInt int, location string) string {
	template := testAccAzureRMVirtualHub_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub" "test" {
  name                = "acctestVHUB-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"

  route {
    address_prefixes    = [ "172.0.1.0/24" ]
    next_hop_ip_address = "12.34.56.78"
  }
}
`, template, rInt)
}

func testAccAzureRMVirtualHub_routeUpdated(rInt int, location string) string {
	template := testAccAzureRMVirtualHub_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub" "test" {
  name                = "acctestVHUB-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"

  route {
    address_prefixes    = [ "172.0.1.0/24" ]
    next_hop_ip_address = "87.65.43.21"
  }
}
`, template, rInt)
}

func testAccAzureRMVirtualHub_tags(rInt int, location string) string {
	template := testAccAzureRMVirtualHub_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub" "test" {
  name                = "acctestVHUB-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"

  tags = {
    Hello = "World"
  }
}
`, template, rInt)
}

func testAccAzureRMVirtualHub_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, rInt, location, rInt)
}
