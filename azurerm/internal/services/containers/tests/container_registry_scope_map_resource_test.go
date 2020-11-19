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

func TestAccAzureRMContainerRegistryScopeMap_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_scope_map", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryScopeMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryScopeMap_basic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "actions.0", "repositories/testrepo/content/read"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMContainerRegistryScopeMap_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_scope_map", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryScopeMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryScopeMap_basic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "actions.0", "repositories/testrepo/content/read"),
				),
			},
			{
				Config:      testAccAzureRMContainerRegistryScopeMap_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_container_registry_scope_map"),
			},
		},
	})
}

func TestAccAzureRMContainerRegistryScopeMap_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_scope_map", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryScopeMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryScopeMap_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "actions.0", "repositories/testrepo/content/read"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMContainerRegistryScopeMap_completeUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_scope_map", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryScopeMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryScopeMap_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "container_registry_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "description"),
					resource.TestCheckResourceAttr(data.ResourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "actions.0", "repositories/testrepo/content/read"),
				),
			},
			{
				Config: testAccAzureRMContainerRegistryScopeMap_completeUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "container_registry_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "description"),
					resource.TestCheckResourceAttr(data.ResourceName, "actions.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "actions.0", "repositories/testrepo/content/read"),
					resource.TestCheckResourceAttr(data.ResourceName, "actions.1", "repositories/testrepo/content/delete"),
				),
			},
		},
	})
}

func testAccAzureRMContainerRegistryScopeMap_basic_basic(data acceptance.TestData) string {
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

resource "azurerm_container_registry_scope_map" "test" {
  name                    = "testscopemap%d"
  resource_group_name     = azurerm_resource_group.test.name
  container_registry_name = azurerm_container_registry.test.name
  actions                 = ["repositories/testrepo/content/read"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMContainerRegistryScopeMap_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMContainerRegistryScopeMap_basic_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_scope_map" "import" {
  name                    = azurerm_container_registry_scope_map.test.name
  resource_group_name     = azurerm_container_registry_scope_map.test.resource_group_name
  container_registry_name = azurerm_container_registry_scope_map.test.container_registry_name
  actions                 = azurerm_container_registry_scope_map.test.actions
}
`, template)
}

func testAccAzureRMContainerRegistryScopeMap_complete(data acceptance.TestData) string {
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

resource "azurerm_container_registry_scope_map" "test" {
  name                    = "testscopemap%d"
  description             = "An example scope map"
  resource_group_name     = azurerm_resource_group.test.name
  container_registry_name = azurerm_container_registry.test.name
  actions                 = ["repositories/testrepo/content/read"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMContainerRegistryScopeMap_completeUpdated(data acceptance.TestData) string {
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

resource "azurerm_container_registry_scope_map" "test" {
  name                    = "testscopemap%d"
  description             = "An example scope map"
  resource_group_name     = azurerm_resource_group.test.name
  container_registry_name = azurerm_container_registry.test.name
  actions                 = ["repositories/testrepo/content/read", "repositories/testrepo/content/delete"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testCheckAzureRMContainerRegistryScopeMapDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Containers.ScopeMapsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_container_registry_scope_map" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		containerRegistryName := rs.Primary.Attributes["container_registry_name"]

		scopeMap, err := client.Get(ctx, resourceGroup, containerRegistryName, name)
		if err != nil {
			if utils.ResponseWasNotFound(scopeMap.Response) {
				return nil
			}
		}

		return fmt.Errorf("Bad: Container registry scope map %q (Storage Container: %q) still exists", name, containerRegistryName)
	}

	return nil
}
