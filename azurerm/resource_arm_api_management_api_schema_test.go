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

func TestAccAzureRMApiManagementApiSchema_basic(t *testing.T) {
	resourceName := "azurerm_api_management_api_schema.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiSchemaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiSchema_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiSchemaExists(resourceName),
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

func TestAccAzureRMApiManagementApiSchema_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_api_management_api_schema.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiSchemaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiSchema_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiSchemaExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMApiManagementApiSchema_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_api_management_api_schema"),
			},
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
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		schemaID := rs.Primary.Attributes["schema_id"]
		apiName := rs.Primary.Attributes["api_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ApiSchemasClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMApiManagementApiSchema_basic(rInt int, location string) string {
	template := testAccAzureRMApiManagementApiSchema_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_schema" "test" {
  api_name            = "${azurerm_api_management_api.test.name}"
  api_management_name = "${azurerm_api_management_api.test.api_management_name}"
  resource_group_name = "${azurerm_api_management_api.test.resource_group_name}"
  schema_id           = "acctestSchema%d"
  content_type        = "application/vnd.ms-azure-apim.xsd+xml"
  value               = "${file("testdata/api_management_api_schema.xml")}"
}
`, template, rInt)
}

func testAccAzureRMApiManagementApiSchema_requiresImport(rInt int, location string) string {
	template := testAccAzureRMApiManagementApiSchema_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_schema" "test" {
  api_name            = "${azurerm_api_management_api.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  schema_id           = "acctestSchema%d"
  content_type        = "application/vnd.ms-azure-apim.xsd+xml"
  value               = "${file("testdata/api_management_api_schema.xml")}"
}
`, template, rInt)
}

func testAccAzureRMApiManagementApiSchema_template(rInt int, location string) string {
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

  sku {
    name     = "Developer"
    capacity = 1
  }
}

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"
}
`, rInt, location, rInt, rInt)
}
