package batch_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type BatchApplicationResource struct {
}

func TestAccBatchApplication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_application", "test")
	r := BatchApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.template(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBatchApplication_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_application", "test")
	r := BatchApplicationResource{}
	displayName := fmt.Sprintf("TestAccDisplayName-%d", data.RandomInteger)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.template(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.template(data, fmt.Sprintf(`display_name = "%s"`, displayName)),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue(displayName),
			),
		},
	})
}

func (t BatchApplicationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ApplicationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Batch.ApplicationClient.Get(ctx, id.ResourceGroup, id.BatchAccountName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Batch Application %q (Account Name %q / Resource Group %q) does not exist", id.Name, id.BatchAccountName, id.ResourceGroup)
	}

	return utils.Bool(resp.ApplicationProperties != nil), nil
}

func (BatchApplicationResource) template(data acceptance.TestData, displayName string) string {
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
