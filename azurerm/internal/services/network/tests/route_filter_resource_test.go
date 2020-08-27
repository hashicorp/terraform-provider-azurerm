package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMRouteFilter_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_filter", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRouteFilter_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteFilterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRouteFilter_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_filter", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRouteFilter_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteFilterExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMRouteFilter_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_route_filter"),
			},
		},
	})
}

func TestAccAzureRMRouteFilter_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_filter", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRouteFilter_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteFilterExists("azurerm_route_filter.test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRouteFilter_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_filter", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRouteFilter_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteFilterExists(data.ResourceName),
					testCheckAzureRMRouteFilterDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMRouteFilter_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_filter", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRouteFilter_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteFilterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: testAccAzureRMRouteFilter_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteFilterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "staging"),
				),
			},
		},
	})
}

func TestAccAzureRMRouteFilter_withRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_filter", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRouteFilter_withRules(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteFilterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.access", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.rule_type", "Community"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.communities.0", "12076:53005"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.communities.1", "12076:53006"),
				),
			},
			{
				Config: testAccAzureRMRouteFilter_withRulesUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteFilterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.access", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.rule_type", "Community"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.communities.0", "12076:52005"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.communities.1", "12076:52006"),
				),
			},
		},
	})
}

func testCheckAzureRMRouteFilterExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for route filter: %q", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.RouteFiltersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Route Filter %q (resource group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on routeFiltersClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMRouteFilterDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for route filter: %q", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.RouteFiltersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		future, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			if !response.WasNotFound(future.Response()) {
				return fmt.Errorf("Error deleting Route Filter %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for deletion of Route Filter %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMRouteFilterDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.RouteFiltersClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_route_filter" {
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

		return fmt.Errorf("Route Filter still exists:\n%#v", resp.RouteFilterPropertiesFormat)
	}

	return nil
}

func testAccAzureRMRouteFilter_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_filter" "test" {
  name                = "acctestrf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMRouteFilter_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_route_filter" "import" {
  name                = azurerm_route_filter.test.name
  location            = azurerm_route_filter.test.location
  resource_group_name = azurerm_route_filter.test.resource_group_name
}
`, testAccAzureRMRouteFilter_basic(data))
}

func testAccAzureRMRouteFilter_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_filter" "test" {
  name                = "acctestrf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMRouteFilter_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_filter" "test" {
  name                = "acctestrf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMRouteFilter_withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_filter" "test" {
  name                = "acctestrf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMRouteFilter_withRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_filter" "test" {
  name                = "acctestrf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  rule {
    name        = "acctestrule%d"
    access      = "Allow"
    rule_type   = "Community"
    communities = ["12076:53005", "12076:53006"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMRouteFilter_withRulesUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_filter" "test" {
  name                = "acctestrf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  rule {
    name        = "acctestrule%d"
    access      = "Allow"
    rule_type   = "Community"
    communities = ["12076:52005", "12076:52006"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
