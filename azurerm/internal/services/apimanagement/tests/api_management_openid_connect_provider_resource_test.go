package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementOpenIDConnectProvider_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_openid_connect_provider", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementOpenIDConnectProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementOpenIDConnectProvider_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementOpenIDConnectProviderExists(data.ResourceName),
				),
			},
			data.ImportStep("client_secret"),
		},
	})
}

func TestAccAzureRMApiManagementOpenIDConnectProvider_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_openid_connect_provider", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementOpenIDConnectProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementOpenIDConnectProvider_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementOpenIDConnectProviderExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementOpenIDConnectProvider_requiresImport),
		},
	})
}

func TestAccAzureRMApiManagementOpenIDConnectProvider_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_openid_connect_provider", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementOpenIDConnectProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementOpenIDConnectProvider_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementOpenIDConnectProviderExists(data.ResourceName),
				),
			},
			data.ImportStep("client_secret"),
			{
				Config: testAccAzureRMApiManagementOpenIDConnectProvider_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementOpenIDConnectProviderExists(data.ResourceName),
				),
			},
			data.ImportStep("client_secret"),
		},
	})
}

func testCheckAzureRMApiManagementOpenIDConnectProviderExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.OpenIdConnectClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("API Management OpenID Connect Provider not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		if resp, err := client.Get(ctx, resourceGroup, serviceName, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: OpenID Connect Provider %q (Resource Group %q / API Management Service %q) does not exist", name, resourceGroup, serviceName)
			}
			return fmt.Errorf("Bad: Get on apiManagement.OpenIdConnectClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMApiManagementOpenIDConnectProviderDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.OpenIdConnectClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_openid_connect_provider" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		if resp, err := client.Get(ctx, resourceGroup, serviceName, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on apiManagement.OpenIdConnectClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMApiManagementOpenIDConnectProvider_basic(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementOpenIDConnectProvider_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_openid_connect_provider" "test" {
  name                = "acctest-%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  client_id           = "00001111-2222-3333-%d"
  client_secret       = "%d-cwdavsxbacsaxZX-%d"
  display_name        = "Initial Name"
  metadata_endpoint   = "https://azacctest.hashicorptest.com/example/foo"
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementOpenIDConnectProvider_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementOpenIDConnectProvider_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_openid_connect_provider" "import" {
  name                = azurerm_api_management_openid_connect_provider.test.name
  api_management_name = azurerm_api_management_openid_connect_provider.test.api_management_name
  resource_group_name = azurerm_api_management_openid_connect_provider.test.resource_group_name
  client_id           = azurerm_api_management_openid_connect_provider.test.client_id
  client_secret       = azurerm_api_management_openid_connect_provider.test.client_secret
  display_name        = azurerm_api_management_openid_connect_provider.test.display_name
  metadata_endpoint   = azurerm_api_management_openid_connect_provider.test.metadata_endpoint
}
`, template)
}

func testAccAzureRMApiManagementOpenIDConnectProvider_complete(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementOpenIDConnectProvider_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_openid_connect_provider" "test" {
  name                = "acctest-%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  client_id           = "00001111-3333-2222-%d"
  client_secret       = "%d-423egvwdcsjx-%d"
  display_name        = "Updated Name"
  description         = "Example description"
  metadata_endpoint   = "https://azacctest.hashicorptest.com/example/updated"
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementOpenIDConnectProvider_template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
