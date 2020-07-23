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

func TestAccAzureRMResourceProviderRegistration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_provider_registration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMResourceProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMResourceProviderRegistration_basic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceProviderRegistrationExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMResourceProviderRegistration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_provider_registration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMResourceProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMResourceProviderRegistration_basic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceProviderRegistrationExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMResourceProviderRegistration_requiresImport(),
				ExpectError: acceptance.RequiresImportError("azurerm_resource_provider_registration"),
			},
		},
	})
}

func testCheckAzureRMResourceProviderRegistrationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Resource.ProvidersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		resourceProviderNamespace := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceProviderNamespace, "")

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Resource Provider Namespace %q is not found", resourceProviderNamespace)
			}
			return fmt.Errorf("Bad: Get on ProvidersClient: %+v", err)
		}

		if resp.RegistrationState != nil && *resp.RegistrationState != "Registered" {
			return fmt.Errorf("Bad: Resource Provider Namespace %q is not registered", resourceProviderNamespace)
		}

		return nil
	}
}

func testCheckAzureRMResourceProviderDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Resource.ProvidersClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_resource_provider_registration" {
			continue
		}

		resourceProviderNamespace := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceProviderNamespace, "")

		if err != nil {
			return err
		}

		if resp.RegistrationState != nil && *resp.RegistrationState != "Unregistered" {
			return fmt.Errorf("Bad: Resource Provider Namespace %q is not unregistered", resourceProviderNamespace)
		}

		return nil
	}

	return nil
}

func testAccAzureRMResourceProviderRegistration_basic() string {
	return (`
provider "azurerm" {
  features {}
  skip_provider_registration = true
}

resource "azurerm_resource_provider_registration" "test" {
  name = "Microsoft.PolicyInsights"
}
`)
}

func testAccAzureRMResourceProviderRegistration_requiresImport() string {
	template := testAccAzureRMResourceProviderRegistration_basic()
	return fmt.Sprintf(`
%s
resource "azurerm_resource_provider_registration" "import" {
  name = azurerm_resource_provider_registration.test.name
}
`, template)
}
