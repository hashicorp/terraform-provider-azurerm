package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementApiDiagnostic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_diagnostic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDiagnosticDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiDiagnostic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiDiagnosticExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementApiDiagnostic_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_diagnostic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDiagnosticDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiDiagnostic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiDiagnosticExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApiManagementApiDiagnostic_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiDiagnosticExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementApiDiagnostic_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_diagnostic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDiagnosticDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiDiagnostic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiDiagnosticExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementApiDiagnostic_requiresImport),
		},
	})
}

func testCheckAzureRMApiManagementApiDiagnosticDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ApiDiagnosticClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_api_diagnostic" {
			continue
		}

		diagnosticId, err := parse.ApiManagementApiDiagnosticID(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := client.Get(ctx, diagnosticId.ResourceGroup, diagnosticId.ServiceName, diagnosticId.ApiName, diagnosticId.Name)

		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}
	return nil
}

func testCheckAzureRMApiManagementApiDiagnosticExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ApiDiagnosticClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		diagnosticId, err := parse.ApiManagementApiDiagnosticID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, diagnosticId.ResourceGroup, diagnosticId.ServiceName, diagnosticId.ApiName, diagnosticId.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: API Management Diagnostic %q (Resource Group %q / API Management Service %q / API %q) does not exist", diagnosticId.Name, diagnosticId.ResourceGroup, diagnosticId.ServiceName, diagnosticId.ApiName)
			}
			return fmt.Errorf("bad: Get on apiManagementApiDiagnosticClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementApiDiagnostic_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_logger" "test" {
  name                = "acctestapimnglogger-%[1]d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name

  application_insights {
    instrumentation_key = azurerm_application_insights.test.instrumentation_key
  }
}

resource "azurerm_api_management_api" "test" {
  name                = "acctestAMA-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  revision            = "1"
  display_name        = "Test API"
  path                = "test"
  protocols           = ["https"]

  import {
    content_format = "swagger-link-json"
    content_value  = "http://conferenceapi.azurewebsites.net/?format=json"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMApiManagementApiDiagnostic_basic(data acceptance.TestData) string {
	config := testAccAzureRMApiManagementApiDiagnostic_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_diagnostic" "test" {
  identifier               = "applicationinsights"
  resource_group_name      = azurerm_resource_group.test.name
  api_management_name      = azurerm_api_management.test.name
  api_name                 = azurerm_api_management_api.test.name
  api_management_logger_id = azurerm_api_management_logger.test.id
}
`, config)
}

func testAccAzureRMApiManagementApiDiagnostic_update(data acceptance.TestData) string {
	config := testAccAzureRMApiManagementApiDiagnostic_template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_application_insights" "test2" {
  name                = "acctestappinsightsUpdate-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_api_management_logger" "test2" {
  name                = "acctestapimngloggerUpdate-%[2]d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name

  application_insights {
    instrumentation_key = azurerm_application_insights.test2.instrumentation_key
  }
}

resource "azurerm_api_management_api_diagnostic" "test" {
  identifier               = "applicationinsights"
  resource_group_name      = azurerm_resource_group.test.name
  api_management_name      = azurerm_api_management.test.name
  api_name                 = azurerm_api_management_api.test.name
  api_management_logger_id = azurerm_api_management_logger.test2.id
}
`, config, data.RandomInteger)
}

func testAccAzureRMApiManagementApiDiagnostic_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApiDiagnostic_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_diagnostic" "import" {
  identifier               = azurerm_api_management_api_diagnostic.test.identifier
  resource_group_name      = azurerm_api_management_api_diagnostic.test.resource_group_name
  api_management_name      = azurerm_api_management_api_diagnostic.test.api_management_name
  api_name                 = azurerm_api_management_api.test.name
  api_management_logger_id = azurerm_api_management_api_diagnostic.test.api_management_logger_id
}
`, template)
}
