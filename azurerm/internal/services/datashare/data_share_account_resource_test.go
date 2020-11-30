package datashare_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccDataShareAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_account", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataShareAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataShareAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataShareAccountExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.principal_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.tenant_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccDataShareAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_account", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataShareAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataShareAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataShareAccountExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccDataShareAccount_requiresImport),
		},
	})
}

func TestAccDataShareAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataShareAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataShareAccount_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataShareAccountExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.principal_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.tenant_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccDataShareAccount_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataShareAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataShareAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataShareAccountExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.principal_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.tenant_id"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccDataShareAccount_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataShareAccountExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.principal_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.tenant_id"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccDataShareAccount_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataShareAccountExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.principal_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.tenant_id"),
				),
			},
			{
				Config: testAccDataShareAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataShareAccountExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.principal_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.tenant_id"),
				),
			},
			data.ImportStep(),
			data.ImportStep(),
		},
	})
}

func testCheckDataShareAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataShare.AccountClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("dataShare Account not found: %s", resourceName)
		}
		id, err := parse.AccountID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: data_share account %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on DataShareAccountClient: %+v", err)
		}
		return nil
	}
}

func testCheckDataShareAccountDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DataShare.AccountClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_share_account" {
			continue
		}
		id, err := parse.AccountID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on data_share.accountClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccDataShareAccount_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-datashare-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccDataShareAccount_basic(data acceptance.TestData) string {
	template := testAccDataShareAccount_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_account" "test" {
  name                = "acctest-DSA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func testAccDataShareAccount_requiresImport(data acceptance.TestData) string {
	config := testAccDataShareAccount_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_account" "import" {
  name                = azurerm_data_share_account.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  identity {
    type = "SystemAssigned"
  }
}
`, config)
}

func testAccDataShareAccount_complete(data acceptance.TestData) string {
	template := testAccDataShareAccount_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_account" "test" {
  name                = "acctest-DSA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  identity {
    type = "SystemAssigned"
  }

  tags = {
    env = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccDataShareAccount_update(data acceptance.TestData) string {
	template := testAccDataShareAccount_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_account" "test" {
  name                = "acctest-DSA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  identity {
    type = "SystemAssigned"
  }

  tags = {
    env = "Stage"
  }
}
`, template, data.RandomInteger)
}
