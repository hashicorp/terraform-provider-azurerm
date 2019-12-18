package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementDiagnostic_basic(t *testing.T) {
	resourceName := "azurerm_api_management_diagnostic.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMApiManagementDiagnostic_basic(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDiagnosticDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementDiagnosticExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApiManagementDiagnostic_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_api_management_diagnostic.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementDiagnostic_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementDiagnosticExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMApiManagementDiagnostic_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_api_management_diagnostic"),
			},
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

func testAccAzureRMApiManagementDiagnostic_basic(rInt int, location string) string {
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
`, rInt, location, rInt)
}

func testAccAzureRMApiManagementDiagnostic_requiresImport(rInt int, location string) string {
	template := testAccAzureRMApiManagementDiagnostic_basic(rInt, location)
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
