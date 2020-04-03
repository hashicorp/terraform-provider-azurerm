package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDataShareAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShareAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.ENV", "Test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataShareAccount_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_data_share_account", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShareAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.ENV", "Test"),
				),
			},
			{
				Config:      testAccAzureRMDataShareAccount_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_data_share_account"),
			},
		},
	})
}

func TestAccAzureRMDataShareAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShareAccount_complete(data),
				Check:  resource.ComposeTestCheckFunc(),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataShareAccount_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShareAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.ENV", "Test"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDataShareAccount_complete(data),
				Check:  resource.ComposeTestCheckFunc(),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDataShareAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataShare.AccountClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("DataShare Account not found: %s", resourceName)

		}
		id, err := parse.DataShareAccountID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: data_share account %q does not exist", id.Name)
			}
			return fmt.Errorf("Bad: Get on DataShareAccountClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMDataShareAccountDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DataShare.AccountClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_share_account" {
			continue
		}
		id, err := parse.DataShareAccountID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on data_share.accountClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMDataShareAccount_basic(data acceptance.TestData) string {
	template := testAccAzureRMDataShareAccount_template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_data_share_account" "test" {
  name = "acctest-dsa-%d"
  location = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMDataShareAccount_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMDataShareAccount_basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_data_share_account" "import" {
  name = azurerm_data_share_account.test.name
  location = azurerm_data_share_account.test.location
  resource_group_name = azurerm_data_share_account.test.resource_group_name
  tags = azurerm_data_share_account.test.tags
}
`, config)
}

func testAccAzureRMDataShareAccount_complete(data acceptance.TestData) string {
	template := testAccAzureRMDataShareAccount_template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_data_share_account" "test" {
  name = "acctest-dsa-%d"
  location = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMDataShareAccount_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
