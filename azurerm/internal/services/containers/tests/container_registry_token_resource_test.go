package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMContainerRegistryToken_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_token", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryToken_basic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "container_registry_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "status", "enabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "scope_map_id", "_repositories_pull"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMContainerRegistryToken_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_token", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryToken_basic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "container_registry_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "status", "enabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "scope_map_id", "_repositories_pull"),
				),
			},
			{
				Config:      testAccAzureRMContainerRegistryToken_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_container_registry_token"),
			},
		},
	})
}

func TestAccAzureRMContainerRegistryToken_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_token", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryToken_complete(data, "enabled", "_repositories_pull"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "container_registry_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "status", "enabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "scope_map_id", "_repositories_pull"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMContainerRegistryToken_completeUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_token", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryToken_complete(data, "enabled", "_repositories_pull"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "container_registry_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "status", "enabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "scope_map_id", "_repositories_pull"),
				),
			},
			{
				Config: testAccAzureRMContainerRegistryToken_complete(data, "disabled", "_repositories_push"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "container_registry_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "status", "disabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "scope_map_id", "_repositories_push"),
				),
			},
		},
	})
}

func testAccAzureRMContainerRegistryToken_basic_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Basic"

  # make sure network_rule_set is empty for basic SKU
  # premium SKU will automatically populate network_rule_set.default_action to allow
  network_rule_set = []
}

resource "azurerm_container_registry_token" "test" {
	name = "testtoken%d"
	resource_group_name = azurerm_resource_group.test.name
	container_registry_name = azurerm_container_registry.test.name
	scope_map_id = "_repositories_pull"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMContainerRegistryToken_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMContainerRegistryToken_basic_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_token" "import" {
  name                = azurerm_container_registry.test.name
  resource_group_name = azurerm_container_registry.test.resource_group_name
  container_registry_name = azurerm_container_registry.test.name
  scope_map_id = "_repositories_pull"
}
`, template)
}

func testAccAzureRMContainerRegistryToken_complete(data acceptance.TestData, status string, scopeMapId string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  admin_enabled       = false
  sku                 = "Basic"

  tags = {
    environment = "production"
  }
}

resource "azurerm_container_registry_token" "test" {
	name = "testtoken%d"
	resource_group_name = azurerm_resource_group.test.name
	container_registry_name = azurerm_container_registry.test.name
	scope_map_id = "%s"
	status = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, scopeMapId, status)
}

func testCheckAzureRMContainerRegistryTokenDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Containers.TokensClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_container_registry_token" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		containerRegistryName := rs.Primary.Attributes["container_registry_name"]

		token, err := client.Get(ctx, resourceGroup, containerRegistryName, name)
		if err != nil {
			if utils.ResponseWasNotFound(token.Response) {
				return nil
			}
		}

		return fmt.Errorf("Bad: Container registry token %q (Storage Container: %q) still exists", name, containerRegistryName)
	}

	return nil
}
