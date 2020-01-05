package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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

func TestAccAzureRMApiManagementDiagnostic_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
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
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_diagnostic" {
			continue
		}

		identifier := rs.Primary.Attributes["identifier"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, resourceGroup, serviceName, identifier)

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
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		identifier := rs.Primary.Attributes["identifier"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.DiagnosticClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, resourceGroup, serviceName, identifier)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: API Management Diagnostic %q (Resource Group %q / API Management Service %q) does not exist", identifier, resourceGroup, serviceName)
			}
			return fmt.Errorf("Bad: Get on apiManagementDiagnosticClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementDiagnostic_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name = "Developer_1"
}

resource "azurerm_api_management_diagnostic" "test" {
  identifier          = "applicationinsights"
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  enabled             = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMApiManagementDiagnostic_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementDiagnostic_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_diagnostic" "import" {
  identifier          = "${azurerm_api_management_diagnostic.test.identifier}"
  resource_group_name = "${azurerm_api_management_diagnostic.test.resource_group_name}"
  api_management_name = "${azurerm_api_management_diagnostic.test.api_management_name}"
  enabled             = "${azurerm_api_management_diagnostic.test.enabled}"
}
`, template)
}
