package tests

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementIdentityProviderAADB2C_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_aadb2c", "test")
	b2cConfig := testAccAzureRMApiManagementIdentityProviderAADB2C_getB2CConfig(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementIdentityProviderAADB2CDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementIdentityProviderAADB2C_basic(data, b2cConfig),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementIdentityProviderAADB2CExists(data.ResourceName),
				),
			},
			data.ImportStep("client_secret"),
		},
	})
}

func TestAccAzureRMApiManagementIdentityProviderAADB2C_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_aadb2c", "test")
	b2cConfig := testAccAzureRMApiManagementIdentityProviderAADB2C_getB2CConfig(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementIdentityProviderAADB2CDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementIdentityProviderAADB2C_basic(data, b2cConfig),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementIdentityProviderAADB2CExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMApiManagementIdentityProviderAADB2C_requiresImport(data, b2cConfig),
				ExpectError: acceptance.RequiresImportError(data.ResourceType),
			},
		},
	})
}

func testAccAzureRMApiManagementIdentityProviderAADB2C_getB2CConfig(t *testing.T) map[string]string {
	config := map[string]string{
		"tenant_id":     "",
		"tenant_slug":   "",
		"client_id":     "",
		"client_secret": "",
	}

	for k, _ := range config {
		e := fmt.Sprintf("ARM_TEST_B2C_%s", strings.ToUpper(k))
		if v := os.Getenv(e); v != "" {
			config[k] = v
			continue
		}
		t.Fatalf("`%s` must be set for acceptance tests for resource `azurerm_api_management_identity_provider_aadb2c`!", e)
	}

	return config
}

func testCheckAzureRMApiManagementIdentityProviderAADB2CDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.IdentityProviderClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_identity_provider_aadb2c" {
			continue
		}

		apiManagementId := rs.Primary.Attributes["api_management_id"]
		id, err := parse.ApiManagementID(apiManagementId)
		if err != nil {
			return fmt.Errorf("Error parsing API Management ID %q: %+v", apiManagementId, err)
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, apimanagement.AadB2C)

		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}
	return nil
}

func testCheckAzureRMApiManagementIdentityProviderAADB2CExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.IdentityProviderClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		apiManagementId := rs.Primary.Attributes["api_management_id"]
		id, err := parse.ApiManagementID(apiManagementId)
		if err != nil {
			return fmt.Errorf("Error parsing API Management ID %q: %+v", apiManagementId, err)
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, apimanagement.AadB2C)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: API Management Identity Provider %q (Resource Group %q / API Management Service %q) does not exist", apimanagement.AadB2C, id.ResourceGroup, id.ServiceName)
			}
			return fmt.Errorf("Bad: Get on apiManagementIdentityProviderClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementIdentityProviderAADB2C_basic(data acceptance.TestData, b2cConfig map[string]string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {
  tenant_id     = "%[1]s"
  client_id     = "%[2]s"
  client_secret = "%[3]s"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-api-%[5]d"
  location = "%[6]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[5]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}

resource "azuread_application" "test" {
  name                       = "acctestAM-%[5]d"
  oauth2_allow_implicit_flow = true
  reply_urls                 = ["https://${azurerm_api_management.test.name}.developer.azure-api.net/signin"]
}

resource "azuread_application_password" "test" {
  application_object_id = azuread_application.test.object_id
  end_date_relative     = "36h"
  value                 = "P@55w0rD!%[7]s"
}

resource "azurerm_api_management_identity_provider_aadb2c" "test" {
  api_management_id      = azurerm_api_management.test.id
  client_id              = azuread_application.test.application_id
  client_secret          = "P@55w0rD!%[7]s"
  allowed_tenant         = "%[4]s.onmicrosoft.com"
  signin_tenant          = "%[4]s.onmicrosoft.com"
  authority              = "%[4]s.b2clogin.com"
  signin_policy          = "B2C_1_Login"
  signup_policy          = "B2C_1_Signup"
  profile_editing_policy = "B2C_1_EditProfile"
  password_reset_policy  = "B2C_1_ResetPassword"

  depends_on = [azuread_application_password.test]
}
`, b2cConfig["tenant_id"], b2cConfig["client_id"], b2cConfig["client_secret"], b2cConfig["tenant_slug"], data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMApiManagementIdentityProviderAADB2C_requiresImport(data acceptance.TestData, b2cConfig map[string]string) string {
	template := testAccAzureRMApiManagementIdentityProviderAADB2C_basic(data, b2cConfig)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_identity_provider_aadb2c" "import" {
  api_management_id = azurerm_api_management_identity_provider_aadb2c.test.api_management_id
  client_id         = azurerm_api_management_identity_provider_aadb2c.test.client_id
  client_secret     = azurerm_api_management_identity_provider_aadb2c.test.client_secret
  allowed_tenant    = azurerm_api_management_identity_provider_aadb2c.test.allowed_tenant
  signin_tenant     = azurerm_api_management_identity_provider_aadb2c.test.signin_tenant
  authority         = azurerm_api_management_identity_provider_aadb2c.test.authority
  signup_policy     = azurerm_api_management_identity_provider_aadb2c.test.signup_policy
  signin_policy     = azurerm_api_management_identity_provider_aadb2c.test.signin_policy
}
`, template)
}
