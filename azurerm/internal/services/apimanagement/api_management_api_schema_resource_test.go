package apimanagement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementApiSchema_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_schema", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiSchemaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiSchema_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiSchemaExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementApiSchema_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_schema", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiSchemaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiSchema_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiSchemaExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementApiSchema_requiresImport),
		},
	})
}

func testCheckAzureRMApiManagementApiSchemaDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ApiSchemasClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_api_schema" {
			continue
		}

		schemaID := rs.Primary.Attributes["schema_id"]
		apiName := rs.Primary.Attributes["api_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, serviceName, apiName, schemaID)
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

func testCheckAzureRMApiManagementApiSchemaExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ApiSchemasClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		schemaID := rs.Primary.Attributes["schema_id"]
		apiName := rs.Primary.Attributes["api_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, serviceName, apiName, schemaID)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: API Schema %q (API %q / API Management Service %q / Resource Group: %q) does not exist", schemaID, apiName, serviceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on apiManagementApiSchemasClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementApiSchema_basic(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApiSchema_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_schema" "test" {
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management_api.test.api_management_name
  resource_group_name = azurerm_api_management_api.test.resource_group_name
  schema_id           = "acctestSchema%d"
  content_type        = "application/vnd.ms-azure-apim.xsd+xml"
  value               = file("testdata/api_management_api_schema.xml")
}
`, template, data.RandomInteger)
}

func testAccAzureRMApiManagementApiSchema_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApiSchema_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_schema" "import" {
  api_name            = azurerm_api_management_api_schema.test.api_name
  api_management_name = azurerm_api_management_api_schema.test.api_management_name
  resource_group_name = azurerm_api_management_api_schema.test.resource_group_name
  schema_id           = azurerm_api_management_api_schema.test.schema_id
  content_type        = azurerm_api_management_api_schema.test.content_type
  value               = azurerm_api_management_api_schema.test.value
}
`, template)
}

func testAccAzureRMApiManagementApiSchema_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
