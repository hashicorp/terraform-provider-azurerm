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

func TestAccAzureRMApiManagementAPIOperationPolicy_basic(t *testing.T) {
	resourceName := "azurerm_api_management_api_operation_policy.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementAPIOperationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementAPIOperationPolicy_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementAPIOperationPolicyExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"xml_link"},
			},
		},
	})
}

func TestAccAzureRMApiManagementAPIOperationPolicy_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_api_management_api_operation_policy.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementAPIOperationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementAPIOperationPolicy_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementAPIOperationPolicyExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMApiManagementAPIOperationPolicy_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_api_management_api_policy"),
			},
		},
	})
}

func TestAccAzureRMApiManagementAPIOperationPolicy_update(t *testing.T) {
	resourceName := "azurerm_api_management_api_operation_policy.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementAPIOperationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementAPIOperationPolicy_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementAPIOperationPolicyExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMApiManagementAPIOperationPolicy_updated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementAPIOperationPolicyExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"xml_link"},
			},
		},
	})
}

func testCheckAzureRMApiManagementAPIOperationPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		apiName := rs.Primary.Attributes["api_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		operationID := rs.Primary.Attributes["operation_id"]

		conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ApiOperationPoliciesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := conn.Get(ctx, resourceGroup, serviceName, apiName, operationID)
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

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_api_operation_policy" {
			continue
		}

		apiName := rs.Primary.Attributes["api_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		operationID := rs.Primary.Attributes["operation_id"]

		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := conn.Get(ctx, resourceGroup, serviceName, apiName, operationID)
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

func testAccAzureRMApiManagementAPIOperationPolicy_basic(rInt int, location string) string {
	template := testAccAzureRMApiManagementApiOperation_basic(rInt, location)

	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation_policy" "test" {
  api_name            = "${azurerm_api_management_api.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  operation_id        = "${azurerm_api_management_api_operation.test.operation_id}"
  xml_link            = "https://gist.githubusercontent.com/riordanp/ca22f8113afae0eb38cc12d718fd048d/raw/d6ac89a2f35a6881a7729f8cb4883179dc88eea1/example.xml"
}
`, template)
}

func testAccAzureRMApiManagementAPIOperationPolicy_requiresImport(rInt int, location string) string {
	template := testAccAzureRMApiManagementAPIOperationPolicy_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation_policy" "import" {
  api_name            = "${azurerm_api_management_api_policy.test.api_name}"
  api_management_name = "${azurerm_api_management_api_policy.test.api_management_name}"
  resource_group_name = "${azurerm_api_management_api_policy.test.resource_group_name}"
  operation_id        = "${azurerm_api_management_api_operation.test.operation_id}"
  xml_link            = "${azurerm_api_management_api_policy.test.xml_link}"
}
`, template)
}

func testAccAzureRMApiManagementAPIOperationPolicy_updated(rInt int, location string) string {
	template := testAccAzureRMApiManagementApiOperation_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation_policy" "test" {
  api_name            = "${azurerm_api_management_api.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  operation_id        = "${azurerm_api_management_api_operation.test.operation_id}"

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
