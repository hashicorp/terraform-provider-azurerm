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

func TestAccAzureRMApplicationSecurityGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_security_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationSecurityGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationSecurityGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationSecurityGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_security_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationSecurityGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationSecurityGroupExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMApplicationSecurityGroup_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_application_security_group"),
			},
		},
	})
}

func TestAccAzureRMApplicationSecurityGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_security_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationSecurityGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationSecurityGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Hello", "World"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationSecurityGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_security_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationSecurityGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationSecurityGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			{
				Config: testAccAzureRMApplicationSecurityGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationSecurityGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Hello", "World"),
				),
			},
		},
	})
}

func testCheckAzureRMApplicationSecurityGroupDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ApplicationSecurityGroupsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_application_security_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Application Security Group still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMApplicationSecurityGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ApplicationSecurityGroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Application Security Group: %q", name)
		}

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Application Security Group %q (resource group: %q) was not found: %+v", name, resourceGroup, err)
			}

			return fmt.Errorf("Bad: Get on applicationSecurityGroupsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApplicationSecurityGroup_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_security_group" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMApplicationSecurityGroup_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApplicationSecurityGroup_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_application_security_group" "import" {
  name                = azurerm_application_security_group.test.name
  location            = azurerm_application_security_group.test.location
  resource_group_name = azurerm_application_security_group.test.resource_group_name
}
`, template)
}

func testAccAzureRMApplicationSecurityGroup_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_security_group" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
