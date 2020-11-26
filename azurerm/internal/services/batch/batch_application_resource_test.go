package batch_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccBatchApplication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_application", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchApplication_template(data, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchApplicationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccBatchApplication_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_application", "test")
	displayName := fmt.Sprintf("TestAccDisplayName-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchApplication_template(data, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchApplicationExists(data.ResourceName),
				),
			},
			{
				Config: testAccBatchApplication_template(data, fmt.Sprintf(`display_name = "%s"`, displayName)),
				Check: resource.ComposeTestCheckFunc(
					testCheckBatchApplicationExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", displayName),
				),
			},
		},
	})
}

func testCheckBatchApplicationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Batch.ApplicationClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Batch Application not found: %s", resourceName)
		}

		id, err := parse.ApplicationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.BatchAccountName, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Batch Application %q (Account Name %q / Resource Group %q) does not exist", id.Name, id.BatchAccountName, id.ResourceGroup)
			}
			return fmt.Errorf("Bad: Get on batchApplicationClient: %+v", err)
		}

		return nil
	}
}

func testCheckBatchApplicationDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Batch.ApplicationClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_batch_application" {
			continue
		}

		id, err := parse.ApplicationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.BatchAccountName, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on batchApplicationClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccBatchApplication_template(data acceptance.TestData, displayName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "test" {
  name                 = "acctestba%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  pool_allocation_mode = "BatchService"
  storage_account_id   = azurerm_storage_account.test.id
}

resource "azurerm_batch_application" "test" {
  name                = "acctestbatchapp-%d"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  %s
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger, displayName)
}
