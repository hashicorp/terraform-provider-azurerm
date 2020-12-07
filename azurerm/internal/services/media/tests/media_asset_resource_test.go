package media_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AssetResource struct {
}

func TestAccAsset_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_asset", "test")
	r := AssetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Asset-Content1"),
			),
		},
		data.ImportStep(),
	})
}

func TestMediaAccAsset_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_asset", "test")
	r := AssetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("alternate_id").HasValue("Asset-alternateid"),
				check.That(data.ResourceName).Key("storage_account_name").HasValue(fmt.Sprintf("acctestsa1%s", data.RandomString)),
				check.That(data.ResourceName).Key("container").HasValue("asset-container"),
				check.That(data.ResourceName).Key("description").HasValue("Asset description"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAsset_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_asset", "test")
	r := AssetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Asset-Content1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("alternate_id").HasValue("Asset-alternateid"),
				check.That(data.ResourceName).Key("storage_account_name").HasValue(fmt.Sprintf("acctestsa1%s", data.RandomString)),
				check.That(data.ResourceName).Key("container").HasValue("asset-container"),
				check.That(data.ResourceName).Key("description").HasValue("Asset description"),
			),
		},
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Asset-Content1"),
				check.That(data.ResourceName).Key("description").HasValue(""),
				check.That(data.ResourceName).Key("alternate_id").HasValue(""),
			),
		},
		data.ImportStep(),
		data.ImportStep(),
	})
}

func (AssetResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.AssetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Media.AssetsClient.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Asset %s (Media Services Account %s) (resource group: %s): %v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.AssetProperties != nil), nil
}

func (AssetResource) basic(data acceptance.TestData) string {
	template := AssetResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_media_asset" "test" {
  name                        = "Asset-Content1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
}

`, template)
}

func (AssetResource) complete(data acceptance.TestData) string {
	template := AssetResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_media_asset" "test" {
  name                        = "Asset-Content1"
  description                 = "Asset description"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  storage_account_name        = azurerm_storage_account.test.name
  alternate_id                = "Asset-alternateid"
  container                   = "asset-container"
}

`, template)
}

func (AssetResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-media-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa1%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  storage_account {
    id         = azurerm_storage_account.test.id
    is_primary = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}
