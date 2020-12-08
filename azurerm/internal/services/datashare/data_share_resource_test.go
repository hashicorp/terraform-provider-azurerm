package datashare_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DataShareResource struct {
}

func TestAccDataShare_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share", "test")
	r := DataShareResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataShare_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share", "test")
	r := DataShareResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDataShare_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share", "test")
	r := DataShareResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataShare_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share", "test")
	r := DataShareResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataShare_snapshotSchedule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share", "test")
	r := DataShareResource{}
	startTime := time.Now().Add(time.Hour * 7).Format(time.RFC3339)
	startTime2 := time.Now().Add(time.Hour * 8).Format(time.RFC3339)

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.snapshotSchedule(data, startTime),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.snapshotScheduleUpdated(data, startTime2),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t DataShareResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ShareID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataShare.SharesClient.Get(ctx, id.ResourceGroup, id.AccountName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Data Share %q (resource group: %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.ShareProperties != nil), nil
}

func (DataShareResource) template(data acceptance.TestData) string {
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

func (r DataShareResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_share" "test" {
  name       = "acctest_ds_%d"
  account_id = azurerm_data_share_account.test.id
  kind       = "CopyBased"
}
`, r.template(data), data.RandomInteger)
}

func (r DataShareResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_share" "import" {
  name       = azurerm_data_share.test.name
  account_id = azurerm_data_share_account.test.id
  kind       = azurerm_data_share.test.kind
}
`, r.basic(data))
}

func (r DataShareResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_share" "test" {
  name        = "acctest_ds_%d"
  account_id  = azurerm_data_share_account.test.id
  kind        = "CopyBased"
  description = "share desc"
  terms       = "share terms"
}
`, r.template(data), data.RandomInteger)
}

func (r DataShareResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_share" "test" {
  name        = "acctest_ds_%d"
  account_id  = azurerm_data_share_account.test.id
  kind        = "CopyBased"
  description = "share desc 2"
  terms       = "share terms 2"
}
`, r.template(data), data.RandomInteger)
}

func (r DataShareResource) snapshotSchedule(data acceptance.TestData, startTime string) string {
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
`, r.template(data), data.RandomInteger, startTime)
}

func (r DataShareResource) snapshotScheduleUpdated(data acceptance.TestData, startTime string) string {
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
`, r.template(data), data.RandomInteger, startTime)
}
