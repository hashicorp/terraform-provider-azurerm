package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMStaticSite_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_site", "test")

	if ok := skipStaticSite(); ok {
		t.Skip("Skipping as `ARM_TEST_GITHUB_TOKEN` was not specified")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStaticSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStaticSite_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStaticSiteExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"github_configuration.0.api_location",
				"github_configuration.0.app_location",
				"github_configuration.0.artifact_location",
				"github_configuration.0.repo_token"),
		},
	})
}

func TestAccAzureRMStaticSite_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_site", "test")

	if ok := skipStaticSite(); ok {
		t.Skip("Skipping as `ARM_TEST_GITHUB_TOKEN` was not specified")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStaticSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStaticSite_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStaticSiteExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMStaticSite_requiresImport),
		},
	})
}

func testCheckAzureRMStaticSiteDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Web.StaticSitesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_static_site" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetStaticSite(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMStaticSiteExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Web.StaticSitesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		staticSiteName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Static Site: %s", staticSiteName)
		}

		resp, err := client.GetStaticSite(ctx, resourceGroup, staticSiteName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Static Site %q (resource group: %q) does not exist", staticSiteName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on StaticSitesClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMStaticSite_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_static_site" "test" {
  name                = "acctestSS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  github_configuration {
    repo_url   = "https://github.com/aristosvo/azure-static-web-app"
    branch     = "master"
    repo_token = "%s"

    app_location      = "/"
    api_location      = ""
    artifact_location = "dist/angular-basic"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, os.Getenv("ARM_TEST_GITHUB_TOKEN"))
}

func testAccAzureRMStaticSite_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMStaticSite_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_static_site" "import" {
  name                = azurerm_static_site.test.name
  location            = azurerm_static_site.test.location
  resource_group_name = azurerm_static_site.test.resource_group_name

  github_configuration {
    repo_url   = azurerm_static_site.test.github_configuration.0.repo_url
    branch     = azurerm_static_site.test.github_configuration.0.branch
    repo_token = azurerm_static_site.test.github_configuration.0.repo_token

    app_location      = azurerm_static_site.test.github_configuration.0.app_location
    api_location      = azurerm_static_site.test.github_configuration.0.api_location
    artifact_location = azurerm_static_site.test.github_configuration.0.artifact_location
  }
}
`, template)
}
