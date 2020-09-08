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

func TestAccAzureRMApiManagementGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", "Test Group"),
					resource.TestCheckResourceAttr(data.ResourceName, "type", "custom"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", "Test Group"),
					resource.TestCheckResourceAttr(data.ResourceName, "type", "custom"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementGroup_requiresImport),
		},
	})
}

func TestAccAzureRMApiManagementGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementGroup_complete(data, "Test Group", "A test description."),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", "Test Group"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "A test description."),
					resource.TestCheckResourceAttr(data.ResourceName, "type", "external"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementGroup_descriptionDisplayNameUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementGroup_complete(data, "Original Group", "The original description."),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", "Original Group"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "The original description."),
					resource.TestCheckResourceAttr(data.ResourceName, "type", "external"),
				),
			},
			{
				Config: testAccAzureRMApiManagementGroup_complete(data, "Modified Group", "A modified description."),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", "Modified Group"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "A modified description."),
					resource.TestCheckResourceAttr(data.ResourceName, "type", "external"),
				),
			},
			{
				Config: testAccAzureRMApiManagementGroup_complete(data, "Original Group", "The original description."),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", "Original Group"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "The original description."),
					resource.TestCheckResourceAttr(data.ResourceName, "type", "external"),
				),
			},
		},
	})
}

func testCheckAzureRMAPIManagementGroupDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.GroupClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		resp, err := client.Get(ctx, resourceGroup, serviceName, name)

		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}
	return nil
}

func testCheckAzureRMAPIManagementGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.GroupClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		resp, err := client.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: API Management Group %q (Resource Group %q / API Management Service %q) does not exist", name, resourceGroup, serviceName)
			}
			return fmt.Errorf("Bad: Get on apiManagement.GroupClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementGroup_basic(data acceptance.TestData) string {
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

  sku_name = "Developer_1"
}

resource "azurerm_api_management_group" "test" {
  name                = "acctestAMGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "Test Group"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementGroup_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementGroup_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_group" "import" {
  name                = azurerm_api_management_group.test.name
  resource_group_name = azurerm_api_management_group.test.resource_group_name
  api_management_name = azurerm_api_management_group.test.api_management_name
  display_name        = azurerm_api_management_group.test.display_name
}
`, template)
}

func testAccAzureRMApiManagementGroup_complete(data acceptance.TestData, displayName, description string) string {
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

  sku_name = "Developer_1"
}

resource "azurerm_api_management_group" "test" {
  name                = "acctestAMGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "%s"
  description         = "%s"
  type                = "external"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, displayName, description)
}
