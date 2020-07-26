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

func TestAccAzureRMApiManagementIdentityProviderAADB2C_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_aadb2c", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementIdentityProviderAADB2CDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementIdentityProviderAADB2C_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementIdentityProviderAADB2CExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementIdentityProviderAADB2C_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_aadb2c", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementIdentityProviderAADB2CDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementIdentityProviderAADB2C_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementIdentityProviderAADB2CExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "client_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_secret", "00000000000000000000000000000000"),
					resource.TestCheckResourceAttr(data.ResourceName, "allowed_tenants.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "allowed_tenants.0", data.Client().TenantID),

					resource.TestCheckResourceAttr(data.ResourceName, "signin_tenant", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckResourceAttr(data.ResourceName, "authority", "ExampleAuthority"),
					resource.TestCheckResourceAttr(data.ResourceName, "signup_policy", "ExampleSignupPolicy"),
					resource.TestCheckResourceAttr(data.ResourceName, "signin_policy", "ExampleSigninPolicy"),
				),
			},
			{
				Config: testAccAzureRMApiManagementIdentityProviderAADB2C_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementIdentityProviderAADB2CExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "client_id", "22222222-2222-2222-2222-222222222222"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_secret", "22222222222222222222222222222222"),
					resource.TestCheckResourceAttr(data.ResourceName, "allowed_tenants.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "allowed_tenants.0", data.Client().TenantID),
					resource.TestCheckResourceAttr(data.ResourceName, "signin_tenant", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckResourceAttr(data.ResourceName, "authority", "ExampleAuthority"),
					resource.TestCheckResourceAttr(data.ResourceName, "signup_policy", "ExampleSignupPolicy"),
					resource.TestCheckResourceAttr(data.ResourceName, "signin_policy", "ExampleSigninPolicy")),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementIdentityProviderAADB2C_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_aadb2c", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementIdentityProviderAADB2CDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementIdentityProviderAADB2C_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementIdentityProviderAADB2CExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementIdentityProviderAADB2C_requiresImport),
		},
	})
}

func testCheckAzureRMApiManagementIdentityProviderAADB2CDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.IdentityProviderClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_identity_provider_aadb2c" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		resp, err := client.Get(ctx, resourceGroup, serviceName, apimanagement.AadB2C)

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

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		resp, err := client.Get(ctx, resourceGroup, serviceName, apimanagement.AadB2C)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: API Management Identity Provider %q (Resource Group %q / API Management Service %q) does not exist", apimanagement.AadB2C, resourceGroup, serviceName)
			}
			return fmt.Errorf("Bad: Get on apiManagementIdentityProviderClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementIdentityProviderAADB2C_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-api-%d"
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

resource "azurerm_api_management_identity_provider_aadb2c" "test" {
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  client_id           = "00000000-0000-0000-0000-000000000000"
  client_secret       = "00000000000000000000000000000000"
  allowed_tenants     = ["%s"]
  signin_tenant       = "%s"
  authority           = "${azurerm_api_management.test.name}.b2clogin.com"
  signup_policy       = "ExampleSignupPolicy"
  signin_policy       = "ExampleSigninPolicy"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Client().TenantID, data.Client().TenantID)
}

func testAccAzureRMApiManagementIdentityProviderAADB2C_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-api-%d"
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

resource "azurerm_api_management_identity_provider_aadb2c" "test" {
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  client_id           = "22222222-2222-2222-2222-222222222222"
  client_secret       = "22222222222222222222222222222222"
  allowed_tenants     = ["%s"]
  signin_tenant       = "%s"
  authority           = "ExampleAuthority"
  signup_policy       = "ExampleSignupPolicy"
  signin_policy       = "ExampleSigninPolicy"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Client().TenantID, data.Client().TenantID)
}

func testAccAzureRMApiManagementIdentityProviderAADB2C_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementIdentityProviderAADB2C_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_identity_provider_aadb2c" "import" {
  resource_group_name = azurerm_api_management_identity_provider_aadb2c.test.resource_group_name
  api_management_name = azurerm_api_management_identity_provider_aadb2c.test.api_management_name
  client_id           = azurerm_api_management_identity_provider_aadb2c.test.client_id
  client_secret       = azurerm_api_management_identity_provider_aadb2c.test.client_secret
  allowed_tenants     = azurerm_api_management_identity_provider_aadb2c.test.allowed_tenants
  signin_tenant       = azurerm_api_management_identity_provider_aadb2c.test.signin_tenant
  authority           = azurerm_api_management_identity_provider_aadb2c.test.authority
  signup_policy       = azurerm_api_management_identity_provider_aadb2c.test.signup_policy
  signin_policy       = azurerm_api_management_identity_provider_aadb2c.test.signin_policy
}
`, template)
}
