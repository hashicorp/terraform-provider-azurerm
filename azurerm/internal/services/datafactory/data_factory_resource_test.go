package datafactory_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDataFactory_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactory_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataFactory_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactory_tags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "production"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataFactory_tagsUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactory_tags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "production"),
				),
			},
			{
				Config: testAccAzureRMDataFactory_tagsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.updated", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataFactory_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactory_identity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.type"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.principal_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.tenant_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataFactory_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactory_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryExists(data.ResourceName),
					testCheckAzureRMDataFactoryDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMDataFactory_github(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactory_github(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "github_configuration.0.account_name", fmt.Sprintf("acctestGH-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "github_configuration.0.git_url", "https://github.com/terraform-providers/"),
					resource.TestCheckResourceAttr(data.ResourceName, "github_configuration.0.repository_name", "terraform-provider-azurerm"),
					resource.TestCheckResourceAttr(data.ResourceName, "github_configuration.0.branch_name", "master"),
					resource.TestCheckResourceAttr(data.ResourceName, "github_configuration.0.root_folder", "/"),
				),
			},
			{
				Config: testAccAzureRMDataFactory_githubUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "github_configuration.0.account_name", fmt.Sprintf("acctestGitHub-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "github_configuration.0.git_url", "https://github.com/terraform-providers/"),
					resource.TestCheckResourceAttr(data.ResourceName, "github_configuration.0.repository_name", "terraform-provider-azuread"),
					resource.TestCheckResourceAttr(data.ResourceName, "github_configuration.0.branch_name", "stable-website"),
					resource.TestCheckResourceAttr(data.ResourceName, "github_configuration.0.root_folder", "/azuread"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDataFactoryExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.FactoriesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Data Factory: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on dataFactoryClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Data Factory %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDataFactoryDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.FactoriesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Data Factory: %s", name)
		}

		resp, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp) {
				return fmt.Errorf("Bad: Delete on dataFactoryClient: %+v", err)
			}
		}

		return nil
	}
}

func testCheckAzureRMDataFactoryDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.FactoriesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_factory" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Data Factory still exists:\n%#v", resp.FactoryProperties)
		}
	}

	return nil
}

func testAccAzureRMDataFactory_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMDataFactory_tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMDataFactory_tagsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "production"
    updated     = "true"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMDataFactory_identity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMDataFactory_github(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  github_configuration {
    git_url         = "https://github.com/terraform-providers/"
    repository_name = "terraform-provider-azurerm"
    branch_name     = "master"
    root_folder     = "/"
    account_name    = "acctestGH-%d"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDataFactory_githubUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  github_configuration {
    git_url         = "https://github.com/terraform-providers/"
    repository_name = "terraform-provider-azuread"
    branch_name     = "stable-website"
    root_folder     = "/azuread"
    account_name    = "acctestGitHub-%d"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
