package applicationinsights_test

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMApplicationInsightsAPIKey_no_permission(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_api_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMApplicationInsightsAPIKey_basic(data, "[]", "[]"),
				ExpectError: regexp.MustCompile("The API Key needs to have a Role"),
			},
		},
	})
}

func TestAccAzureRMApplicationInsightsAPIKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_api_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsightsAPIKey_basic(data, "[]", `["annotations"]`),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsAPIKeyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "read_permissions.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "write_permissions.#", "1"),
				),
			},
			{
				Config:      testAccAzureRMApplicationInsightsAPIKey_requiresImport(data, "[]", `["annotations"]`),
				ExpectError: acceptance.RequiresImportError("azurerm_application_insights_api_key"),
			},
		},
	})
}

func TestAccAzureRMApplicationInsightsAPIKey_read_telemetry_permissions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_api_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsightsAPIKey_basic(data, `["aggregate", "api", "draft", "extendqueries", "search"]`, "[]"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsAPIKeyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "read_permissions.#", "5"),
					resource.TestCheckResourceAttr(data.ResourceName, "write_permissions.#", "0"),
				),
			},
			data.ImportStep("api_key"),
		},
	})
}

func TestAccAzureRMApplicationInsightsAPIKey_write_annotations_permission(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_api_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsightsAPIKey_basic(data, "[]", `["annotations"]`),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsAPIKeyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "read_permissions.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "write_permissions.#", "1"),
				),
			},
			data.ImportStep("api_key"),
		},
	})
}

func TestAccAzureRMApplicationInsightsAPIKey_authenticate_permission(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_api_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsightsAPIKey_basic(data, `["agentconfig"]`, "[]"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsAPIKeyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "read_permissions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "write_permissions.#", "0"),
				),
			},
			data.ImportStep("api_key"),
		},
	})
}

func TestAccAzureRMApplicationInsightsAPIKey_full_permissions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_api_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsightsAPIKey_basic(data, `["agentconfig", "aggregate", "api", "draft", "extendqueries", "search"]`, `["annotations"]`),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsAPIKeyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "read_permissions.#", "6"),
					resource.TestCheckResourceAttr(data.ResourceName, "write_permissions.#", "1"),
				),
			},
			data.ImportStep("api_key"),
		},
	})
}

func testCheckAzureRMApplicationInsightsAPIKeyDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).AppInsights.APIKeysClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_application_insights_api_key" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		id, err := azure.ParseAzureResourceID(rs.Primary.Attributes["id"])
		if err != nil {
			return err
		}
		resGroup := id.ResourceGroup
		appInsightsName := id.Path["components"]

		resp, err := conn.Get(ctx, resGroup, appInsightsName, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Application Insights API Key still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMApplicationInsightsAPIKeyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).AppInsights.APIKeysClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := azure.ParseAzureResourceID(rs.Primary.Attributes["id"])
		if err != nil {
			return err
		}
		keyID := id.Path["APIKeys"]
		resGroup := id.ResourceGroup
		appInsightsName := id.Path["components"]

		resp, err := conn.Get(ctx, resGroup, appInsightsName, keyID)
		if err != nil {
			return fmt.Errorf("Bad: Get on appInsightsAPIKeyClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Application Insights API Key '%q' (resource group: '%q') does not exist", keyID, resGroup)
		}

		return nil
	}
}

func testAccAzureRMApplicationInsightsAPIKey_basic(data acceptance.TestData, readPerms, writePerms string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights_api_key" "test" {
  name                    = "acctestappinsightsapikey-%d"
  application_insights_id = azurerm_application_insights.test.id
  read_permissions        = %s
  write_permissions       = %s
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, readPerms, writePerms)
}

func testAccAzureRMApplicationInsightsAPIKey_requiresImport(data acceptance.TestData, readPerms, writePerms string) string {
	template := testAccAzureRMApplicationInsightsAPIKey_basic(data, readPerms, writePerms)
	return fmt.Sprintf(`
%s

resource "azurerm_application_insights_api_key" "import" {
  name                    = azurerm_application_insights_api_key.test.name
  application_insights_id = azurerm_application_insights_api_key.test.application_insights_id
  read_permissions        = azurerm_application_insights_api_key.test.read_permissions
  write_permissions       = azurerm_application_insights_api_key.test.write_permissions
}
`, template)
}
