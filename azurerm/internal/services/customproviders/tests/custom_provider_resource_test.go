package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/customproviders/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMCustomProvider_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_provider", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCustomProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCustomProvider_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCustomProviderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCustomProvider_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_provider", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCustomProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCustomProvider_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCustomProviderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCustomProvider_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCustomProviderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCustomProvider_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCustomProviderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCustomProvider_action(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_provider", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCustomProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCustomProvider_action(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCustomProviderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCustomProvider_actionUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCustomProviderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCustomProvider_action(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCustomProviderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMCustomProviderExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).CustomProviders.CustomProviderClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id, err := parse.CustomProviderID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on CustomProviderClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Custom Provider %q (resource group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMCustomProviderDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).CustomProviders.CustomProviderClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_custom_provider" {
			continue
		}

		id, err := parse.CustomProviderID(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Custom Provider still exists: %q", id.Name)
		}
	}

	return nil
}

func testAccAzureRMCustomProvider_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cp-%d"
  location = "%s"
}
resource "azurerm_custom_provider" "test" {
  name                = "accTEst_saa%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  resource_type {
    name     = "dEf1"
    endpoint = "https://testendpoint.com/"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMCustomProvider_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cp-%d"
  location = "%s"
}
resource "azurerm_custom_provider" "test" {
  name                = "accTEst_saa%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  resource_type {
    name     = "dEf1"
    endpoint = "https://testendpoint.com/"
  }

  action {
    name     = "dEf2"
    endpoint = "https://example.com/"
  }

  validation {
    specification = "https://raw.githubusercontent.com/Azure/azure-custom-providers/master/CustomRPWithSwagger/Artifacts/Swagger/pingaction.json"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMCustomProvider_action(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cp-%d"
  location = "%s"
}
resource "azurerm_custom_provider" "test" {
  name                = "accTEst_saa%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  action {
    name     = "dEf1"
    endpoint = "https://testendpoint.com/"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMCustomProvider_actionUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cp-%d"
  location = "%s"
}
resource "azurerm_custom_provider" "test" {
  name                = "accTEst_saa%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  action {
    name     = "dEf2"
    endpoint = "https://example.com/"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
