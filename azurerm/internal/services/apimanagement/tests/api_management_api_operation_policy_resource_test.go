package tests

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementAPIOperationPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementAPIOperationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementAPIOperationPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementAPIOperationPolicyExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"xml_link"},
			},
		},
	})
}

func TestAccAzureRMApiManagementAPIOperationPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementAPIOperationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementAPIOperationPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementAPIOperationPolicyExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementAPIOperationPolicy_requiresImport),
		},
	})
}

func TestAccAzureRMApiManagementAPIOperationPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementAPIOperationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementAPIOperationPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementAPIOperationPolicyExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMApiManagementAPIOperationPolicy_updated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementAPIOperationPolicyExists(data.ResourceName),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"xml_link"},
			},
		},
	})
}

func TestAccAzureRMApiManagementAPIOperationPolicy_rawXml(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementAPIOperationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementAPIOperationPolicy_rawXml(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementAPIOperationPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMApiManagementAPIOperationPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ApiOperationPoliciesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		apiName := rs.Primary.Attributes["api_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		operationID := rs.Primary.Attributes["operation_id"]

		resp, err := conn.Get(ctx, resourceGroup, serviceName, apiName, operationID, apimanagement.PolicyExportFormatXML)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: API Policy (API Management Service %q / API %q / Operation %q / Resource Group %q) does not exist", serviceName, apiName, operationID, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on apiManagementAPIOperationPoliciesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMApiManagementAPIOperationPolicyDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ApiOperationPoliciesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_api_operation_policy" {
			continue
		}

		apiName := rs.Primary.Attributes["api_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		operationID := rs.Primary.Attributes["operation_id"]

		resp, err := conn.Get(ctx, resourceGroup, serviceName, apiName, operationID, apimanagement.PolicyExportFormatXML)
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

func testAccAzureRMApiManagementAPIOperationPolicy_basic(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApiOperation_basic(data)

	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation_policy" "test" {
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  operation_id        = azurerm_api_management_api_operation.test.operation_id
  xml_link            = "https://gist.githubusercontent.com/riordanp/ca22f8113afae0eb38cc12d718fd048d/raw/d6ac89a2f35a6881a7729f8cb4883179dc88eea1/example.xml"
}
`, template)
}

func testAccAzureRMApiManagementAPIOperationPolicy_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementAPIOperationPolicy_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation_policy" "import" {
  api_name            = azurerm_api_management_api_operation_policy.test.api_name
  api_management_name = azurerm_api_management_api_operation_policy.test.api_management_name
  resource_group_name = azurerm_api_management_api_operation_policy.test.resource_group_name
  operation_id        = azurerm_api_management_api_operation_policy.test.operation_id
  xml_link            = azurerm_api_management_api_operation_policy.test.xml_link
}
`, template)
}

func testAccAzureRMApiManagementAPIOperationPolicy_updated(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApiOperation_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation_policy" "test" {
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  operation_id        = azurerm_api_management_api_operation.test.operation_id

  xml_content = <<XML
<policies>
  <inbound>
    <set-variable name="abc" value="@(context.Request.Headers.GetValueOrDefault("X-Header-Name", ""))" />
    <find-and-replace from="xyz" to="abc" />
  </inbound>
</policies>
XML

}
`, template)
}

func testAccAzureRMApiManagementAPIOperationPolicy_rawXml(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApiOperation_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation_policy" "test" {
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  operation_id        = azurerm_api_management_api_operation.test.operation_id

  xml_content = file("testdata/api_management_api_operation_policy.xml")
}
`, template)
}
