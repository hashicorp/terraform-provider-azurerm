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
					resource.TestCheckResourceAttr(data.ResourceName, "status", "enabled"),
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
		CheckDestroy: testCheckAzureRMContainerRegistryTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryToken_basic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "status", "enabled"),
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
		CheckDestroy: testCheckAzureRMContainerRegistryTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryToken_complete(data, "enabled"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "status", "enabled"),
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
		CheckDestroy: testCheckAzureRMContainerRegistryTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryToken_complete(data, "enabled"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "status", "enabled"),
				),
			},
			{
				Config: testAccAzureRMContainerRegistryToken_complete(data, "disabled"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "status", "disabled"),
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
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
}

# use system wide scope map for tests
data "azurerm_container_registry_scope_map" "pull_repos" {
	name                    = "_repositories_pull"
	container_registry_name = azurerm_container_registry.test.name
	resource_group_name     = azurerm_container_registry.test.resource_group_name
}

resource "azurerm_container_registry_token" "test" {
	name                    = "testtoken%d"
	resource_group_name     = azurerm_resource_group.test.name
	container_registry_name = azurerm_container_registry.test.name
	scope_map_id            = data.azurerm_container_registry_scope_map.pull_repos.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMContainerRegistryToken_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMContainerRegistryToken_basic_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_token" "import" {
  name                    = azurerm_container_registry_token.test.name
  resource_group_name     = azurerm_container_registry_token.test.resource_group_name
  container_registry_name = azurerm_container_registry_token.test.container_registry_name
  scope_map_id            = azurerm_container_registry_token.test.scope_map_id
}
`, template)
}

func testAccAzureRMContainerRegistryToken_complete(data acceptance.TestData, status string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  admin_enabled       = false
  sku                 = "Premium"

  tags = {
    environment = "production"
  }
}

# use system wide scope map for tests
data "azurerm_container_registry_scope_map" "pull_repos" {
	name                    = "_repositories_pull"
	container_registry_name = azurerm_container_registry.test.name
	resource_group_name     = azurerm_container_registry.test.resource_group_name
}

resource "azurerm_container_registry_token" "test" {
	name                    = "testtoken%d"
	resource_group_name     = azurerm_resource_group.test.name
	container_registry_name = azurerm_container_registry.test.name
	scope_map_id            = data.azurerm_container_registry_scope_map.pull_repos.id
	status                  = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, status)
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
