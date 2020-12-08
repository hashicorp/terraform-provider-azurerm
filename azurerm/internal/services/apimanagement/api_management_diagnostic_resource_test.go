package apimanagement_test

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

func TestAccAzureRMApiManagementDiagnostic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_diagnostic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDiagnosticDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementDiagnostic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementDiagnosticExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementDiagnostic_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_diagnostic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDiagnosticDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementDiagnostic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementDiagnosticExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApiManagementDiagnostic_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementDiagnosticExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementDiagnostic_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_diagnostic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementDiagnostic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementDiagnosticExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementDiagnostic_requiresImport),
		},
	})
}

func testCheckAzureRMApiManagementDiagnosticDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.DiagnosticClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_diagnostic" {
			continue
		}

		diagnosticId, err := parse.DiagnosticID(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := client.Get(ctx, diagnosticId.ResourceGroup, diagnosticId.ServiceName, diagnosticId.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}
	return nil
}

func testCheckAzureRMApiManagementDiagnosticExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.DiagnosticClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		diagnosticId, err := parse.DiagnosticID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, diagnosticId.ResourceGroup, diagnosticId.ServiceName, diagnosticId.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: API Management Diagnostic %q (Resource Group %q / API Management Service %q) does not exist", diagnosticId.Name, diagnosticId.ResourceGroup, diagnosticId.ServiceName)
			}
			return fmt.Errorf("bad: Get on apiManagementDiagnosticClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementDiagnostic_template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMApiManagementDiagnostic_basic(data acceptance.TestData) string {
	config := testAccAzureRMApiManagementDiagnostic_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_diagnostic" "test" {
  identifier               = "applicationinsights"
  resource_group_name      = azurerm_resource_group.test.name
  api_management_name      = azurerm_api_management.test.name
  api_management_logger_id = azurerm_api_management_logger.test.id
}
`, config)
}

func testAccAzureRMApiManagementDiagnostic_update(data acceptance.TestData) string {
	config := testAccAzureRMApiManagementDiagnostic_template(data)
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

resource "azurerm_api_management_diagnostic" "test" {
  identifier               = "applicationinsights"
  resource_group_name      = azurerm_resource_group.test.name
  api_management_name      = azurerm_api_management.test.name
  api_management_logger_id = azurerm_api_management_logger.test2.id
}
`, config, data.RandomInteger)
}

func testAccAzureRMApiManagementDiagnostic_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementDiagnostic_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_diagnostic" "import" {
  identifier               = azurerm_api_management_diagnostic.test.identifier
  resource_group_name      = azurerm_api_management_diagnostic.test.resource_group_name
  api_management_name      = azurerm_api_management_diagnostic.test.api_management_name
  api_management_logger_id = azurerm_api_management_diagnostic.test.api_management_logger_id
}
`, template)
}
