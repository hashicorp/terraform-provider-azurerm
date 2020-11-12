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

func TestAccAzureRMApiManagementIdentityProviderAAD_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_aad", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementIdentityProviderAADDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementIdentityProviderAAD_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementIdentityProviderAADExists(data.ResourceName),
				),
			},
			data.ImportStep("client_secret"),
		},
	})
}

func TestAccAzureRMApiManagementIdentityProviderAAD_basicDeprecated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_aad", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementIdentityProviderAADDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementIdentityProviderAAD_basicDeprecated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementIdentityProviderAADExists(data.ResourceName),
				),
			},
			data.ImportStep("client_secret", "api_management_name", "resource_group_name"),
		},
	})
}

func TestAccAzureRMApiManagementIdentityProviderAAD_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_aad", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementIdentityProviderAADDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementIdentityProviderAAD_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementIdentityProviderAADExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "client_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_secret", "00000000000000000000000000000000"),
					resource.TestCheckResourceAttr(data.ResourceName, "allowed_tenants.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "allowed_tenants.0", data.Client().TenantID),
				),
			},
			{
				Config: testAccAzureRMApiManagementIdentityProviderAAD_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementIdentityProviderAADExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "client_id", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_secret", "11111111111111111111111111111111"),
					resource.TestCheckResourceAttr(data.ResourceName, "allowed_tenants.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "allowed_tenants.0", data.Client().TenantID),
					resource.TestCheckResourceAttr(data.ResourceName, "allowed_tenants.1", data.Client().TenantID),
				),
			},
			data.ImportStep("client_secret"),
		},
	})
}

func TestAccAzureRMApiManagementIdentityProviderAAD_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_aad", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementIdentityProviderAADDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementIdentityProviderAAD_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementIdentityProviderAADExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementIdentityProviderAAD_requiresImport),
		},
	})
}

func testCheckAzureRMApiManagementIdentityProviderAADDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.IdentityProviderClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_identity_provider_aad" {
			continue
		}

		apiManagementId := rs.Primary.Attributes["api_management_id"]
		id, err := parse.ApiManagementID(apiManagementId)
		if err != nil {
			return fmt.Errorf("Error parsing API Management ID %q: %+v", apiManagementId, err)
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, apimanagement.Aad)

		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}
	return nil
}

func testCheckAzureRMApiManagementIdentityProviderAADExists(resourceName string) resource.TestCheckFunc {
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

		resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, apimanagement.Aad)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: API Management Identity Provider %q (Resource Group %q / API Management Service %q) does not exist", apimanagement.Aad, id.ResourceGroup, id.ServiceName)
			}
			return fmt.Errorf("Bad: Get on apiManagementIdentityProviderClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementIdentityProviderAAD_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-api-%[1]d"
  location = "%[2]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_identity_provider_aad" "test" {
  api_management_id = azurerm_api_management.test.id
  client_id         = "00000000-0000-0000-0000-000000000000"
  client_secret     = "00000000000000000000000000000000"
  signin_tenant     = "00000000-0000-0000-0000-000000000000"
  allowed_tenants   = ["%[3]s"]
}
`, data.RandomInteger, data.Locations.Primary, data.Client().TenantID)
}

func testAccAzureRMApiManagementIdentityProviderAAD_basicDeprecated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-api-%[1]d"
  location = "%[2]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_identity_provider_aad" "test" {
  resource_group_name = azurerm_api_management.test.resource_group_name
  api_management_name = azurerm_api_management.test.name
  client_id           = "00000000-0000-0000-0000-000000000000"
  client_secret       = "00000000000000000000000000000000"
  signin_tenant       = "00000000-0000-0000-0000-000000000000"
  allowed_tenants     = ["%[3]s"]
}
`, data.RandomInteger, data.Locations.Primary, data.Client().TenantID)
}

func testAccAzureRMApiManagementIdentityProviderAAD_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-api-%[1]d"
  location = "%[2]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_identity_provider_aad" "test" {
  api_management_id = azurerm_api_management.test.id
  client_id         = "11111111-1111-1111-1111-111111111111"
  client_secret     = "11111111111111111111111111111111"
  allowed_tenants   = ["%[3]s", "%[3]s"]
}
`, data.RandomInteger, data.Locations.Primary, data.Client().TenantID)
}

func testAccAzureRMApiManagementIdentityProviderAAD_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementIdentityProviderAAD_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_identity_provider_aad" "import" {
  api_management_id = azurerm_api_management_identity_provider_aad.test.api_management_id
  client_id         = azurerm_api_management_identity_provider_aad.test.client_id
  client_secret     = azurerm_api_management_identity_provider_aad.test.client_secret
  allowed_tenants   = azurerm_api_management_identity_provider_aad.test.allowed_tenants
}
`, template)
}
