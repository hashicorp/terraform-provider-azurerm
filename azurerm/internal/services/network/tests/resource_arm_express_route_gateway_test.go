package tests

import (
	"fmt"
	//"net/http"
	"testing"

	//"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMExpressRouteGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteGateway_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMExpressRouteGateway_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_express_route_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteGateway_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteGatewayExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMExpressRouteGateway_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_express_route_gateway"),
			},
		},
	})
}

func TestAccAzureRMExpressRouteGateway_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteGateway_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMExpressRouteGateway_updateScaleUnits(data, 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "scale_units", "2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMExpressRouteGateway_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteGateway_tags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMExpressRouteGatewayExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ExpressRouteGatewaysClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("ExpressRoute Gateway not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: ExpressRoute Gateway %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on network.ExpressRouteGatewaysClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMExpressRouteGatewayDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ExpressRouteGatewaysClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_express_route_gateway" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on network.ExpressRouteGatewaysClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMExpressRouteGateway_basic(data acceptance.TestData) string {
	template := testAccAzureRMExpressRouteGateway_template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_express_route_gateway" "test" {
  name                = "acctestER-gateway-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
	virtual_hub_id      = azurerm_virtual_hub.test.id
	scale_units         = 1
}
`, template, data.RandomInteger)
}

func testAccAzureRMExpressRouteGateway_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMExpressRouteGateway_basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_express_route_gateway" "test" {
  name                = azurerm_express_route_gateway.test.name
  resource_group_name = azurerm_express_route_gateway.test.name
  location            = azurerm_express_route_gateway.test.location
	scale_units         = 1
}
`, template)
}

func testAccAzureRMExpressRouteGateway_updateScaleUnits(data acceptance.TestData, scaleUnits int) string {
	template := testAccAzureRMExpressRouteGateway_template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_express_route_gateway" "test" {
  name                = "acctestER-gateway-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
	virtual_hub_id      = azurerm_virtual_hub.test.id
	scale_units         = %d
}
`, template, data.RandomInteger, scaleUnits)
}

func testAccAzureRMExpressRouteGateway_tags(data acceptance.TestData) string {
	template := testAccAzureRMExpressRouteGateway_template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_express_route_gateway" "test" {
  name                = "acctestER-gateway-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
	virtual_hub_id      = azurerm_virtual_hub.test.id
	scale_units         = 1
  tags = {
    Hello = "World"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMExpressRouteGateway_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}