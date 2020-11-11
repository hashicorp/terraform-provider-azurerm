package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMUserAssignedIdentity_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_user_assigned_identity", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMUserAssignedIdentityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMUserAssignedIdentity_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMUserAssignedIdentityExists(data.ResourceName),
					resource.TestMatchResourceAttr(data.ResourceName, "principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "client_id", validate.UUIDRegExp),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMUserAssignedIdentity_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_user_assigned_identity", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMUserAssignedIdentityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMUserAssignedIdentity_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMUserAssignedIdentityExists(data.ResourceName),
					resource.TestMatchResourceAttr(data.ResourceName, "principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "client_id", validate.UUIDRegExp),
				),
			},
			{
				Config:      testAccAzureRMUserAssignedIdentity_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_user_assigned_identity"),
			},
		},
	})
}

func testCheckAzureRMUserAssignedIdentityExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MSI.UserAssignedIdentitiesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for virtual machine: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on userAssignedIdentitiesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: User Assigned Identity %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMUserAssignedIdentityDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MSI.UserAssignedIdentitiesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_machine" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("User Assigned Identity still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMUserAssignedIdentity_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMUserAssignedIdentity_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "import" {
  name                = azurerm_user_assigned_identity.test.name
  resource_group_name = azurerm_user_assigned_identity.test.resource_group_name
  location            = azurerm_user_assigned_identity.test.location
}
`, testAccAzureRMUserAssignedIdentity_basic(data))
}
