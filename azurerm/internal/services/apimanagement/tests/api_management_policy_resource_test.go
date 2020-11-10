package tests

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementPolicyExists(data.ResourceName),
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

func TestAccAzureRMApiManagementPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementPolicyExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementPolicy_requiresImport),
		},
	})
}

func TestAccAzureRMApiManagementPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementPolicyExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMApiManagementPolicy_customPolicy(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementPolicyExists(data.ResourceName),
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

func TestAccAzureRMApiManagementPolicy_customPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementPolicy_customPolicy(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementPolicyExists(data.ResourceName),
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

func testCheckAzureRMApiManagementPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.PolicyClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		serviceID := rs.Primary.Attributes["api_management_id"]
		id, err := parse.ApiManagementID(serviceID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.ServiceName, apimanagement.PolicyExportFormatXML)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Policy (API Management Service %q / Resource Group %q) does not exist", id.ServiceName, id.ResourceGroup)
			}

			return fmt.Errorf("Bad: Get on apimanagement.PolicyClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMApiManagementPolicyDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.PolicyClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_policy" {
			continue
		}

		serviceID := rs.Primary.Attributes["api_management_id"]
		id, err := parse.ApiManagementID(serviceID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.ServiceName, apimanagement.PolicyExportFormatXML)
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

func testAccAzureRMApiManagementPolicy_basic(data acceptance.TestData) string {
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

resource "azurerm_api_management_policy" "test" {
  api_management_id = azurerm_api_management.test.id
  xml_link          = "https://gist.githubusercontent.com/riordanp/ca22f8113afae0eb38cc12d718fd048d/raw/d6ac89a2f35a6881a7729f8cb4883179dc88eea1/example.xml"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMApiManagementPolicy_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementPolicy_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_policy" "import" {
  api_management_id = azurerm_api_management_policy.test.api_management_id
  xml_link          = azurerm_api_management_policy.test.xml_link
}
`, template)
}

func testAccAzureRMApiManagementPolicy_customPolicy(data acceptance.TestData) string {
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

resource "azurerm_api_management_named_value" "test" {
  name                = "acctestnamedvalue-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "Test"
  value               = "test"
}

resource "azurerm_api_management_policy" "test" {
  api_management_id = azurerm_api_management.test.id

  xml_content = <<XML
<policies>
  <inbound>
    <set-variable name="abc" value="@(context.Request.Headers.GetValueOrDefault("X-Header-Name", ""))" />
    <find-and-replace from="xyz" to="{{${azurerm_api_management_named_value.test.name}}}" />
  </inbound>
</policies>
XML

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
