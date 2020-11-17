package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDataShare_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShare_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataShare_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShare_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMDataShare_requiresImport),
		},
	})
}

func TestAccAzureRMDataShare_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShare_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataShare_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShare_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDataShare_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDataShare_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataShare_snapshotSchedule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share", "test")
	startTime := time.Now().Add(time.Hour * 7).Format(time.RFC3339)
	startTime2 := time.Now().Add(time.Hour * 8).Format(time.RFC3339)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShare_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDataShare_snapshotSchedule(data, startTime),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDataShare_snapshotScheduleUpdated(data, startTime2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDataShare_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDataShareExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataShare.SharesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("dataShare not found: %s", resourceName)
		}
		id, err := parse.DataShareID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.AccountName, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: data_share share %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on DataShareClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMDataShareDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DataShare.SharesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_share" {
			continue
		}
		id, err := parse.DataShareID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.AccountName, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on data_share.shareClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMDataShare_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-datashare-%d"
  location = "%s"
}

resource "azurerm_data_share_account" "test" {
  name                = "acctest-dsa-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  identity {
    type = "SystemAssigned"
  }

  tags = {
    env = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMDataShare_basic(data acceptance.TestData) string {
	template := testAccAzureRMDataShare_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_share" "test" {
  name       = "acctest_ds_%d"
  account_id = azurerm_data_share_account.test.id
  kind       = "CopyBased"
}
`, template, data.RandomInteger)
}

func testAccAzureRMDataShare_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMDataShare_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_share" "import" {
  name       = azurerm_data_share.test.name
  account_id = azurerm_data_share_account.test.id
  kind       = azurerm_data_share.test.kind
}
`, config)
}

func testAccAzureRMDataShare_complete(data acceptance.TestData) string {
	template := testAccAzureRMDataShare_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_share" "test" {
  name        = "acctest_ds_%d"
  account_id  = azurerm_data_share_account.test.id
  kind        = "CopyBased"
  description = "share desc"
  terms       = "share terms"
}
`, template, data.RandomInteger)
}

func testAccAzureRMDataShare_update(data acceptance.TestData) string {
	template := testAccAzureRMDataShare_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_share" "test" {
  name        = "acctest_ds_%d"
  account_id  = azurerm_data_share_account.test.id
  kind        = "CopyBased"
  description = "share desc 2"
  terms       = "share terms 2"
}
`, template, data.RandomInteger)
}

func testAccAzureRMDataShare_snapshotSchedule(data acceptance.TestData, startTime string) string {
	template := testAccAzureRMDataShare_template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_share" "test" {
  name       = "acctest_ds_%[2]d"
  account_id = azurerm_data_share_account.test.id
  kind       = "CopyBased"

  snapshot_schedule {
    name       = "acctest-ss-%[2]d"
    recurrence = "Day"
    start_time = "%[3]s"
  }
}
`, template, data.RandomInteger, startTime)
}

func testAccAzureRMDataShare_snapshotScheduleUpdated(data acceptance.TestData, startTime string) string {
	template := testAccAzureRMDataShare_template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_share" "test" {
  name       = "acctest_ds_%[2]d"
  account_id = azurerm_data_share_account.test.id
  kind       = "CopyBased"

  snapshot_schedule {
    name       = "acctest-ss2-%[2]d"
    recurrence = "Hour"
    start_time = "%[3]s"
  }
}
`, template, data.RandomInteger, startTime)
}
