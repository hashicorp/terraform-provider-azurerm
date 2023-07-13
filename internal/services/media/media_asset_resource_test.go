// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package media_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/assetsandassetfilters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MediaAssetResource struct{}

func TestAccMediaAsset_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_asset", "test")
	r := MediaAssetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Asset-Content1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMediaAsset_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_asset", "test")
	r := MediaAssetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Asset-Content1"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMediaAsset_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_asset", "test")
	r := MediaAssetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("alternate_id").HasValue("Asset-alternateid"),
				check.That(data.ResourceName).Key("storage_account_name").HasValue(fmt.Sprintf("acctestsa1%s", data.RandomString)),
				check.That(data.ResourceName).Key("container").HasValue("asset-container"),
				check.That(data.ResourceName).Key("description").HasValue("Asset description"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMediaAsset_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_asset", "test")
	r := MediaAssetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Asset-Content1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("alternate_id").HasValue("Asset-alternateid"),
				check.That(data.ResourceName).Key("storage_account_name").HasValue(fmt.Sprintf("acctestsa1%s", data.RandomString)),
				check.That(data.ResourceName).Key("container").HasValue("asset-container"),
				check.That(data.ResourceName).Key("description").HasValue("Asset description"),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Asset-Content1"),
				check.That(data.ResourceName).Key("description").HasValue(""),
				check.That(data.ResourceName).Key("alternate_id").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func (MediaAssetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := assetsandassetfilters.ParseAssetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Media.V20220801Client.AssetsAndAssetFilters.AssetsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r MediaAssetResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_media_asset" "test" {
  name                        = "Asset-Content1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
}
`, template)
}

func (r MediaAssetResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_media_asset" "import" {
  name                        = azurerm_media_asset.test.name
  resource_group_name         = azurerm_media_asset.test.resource_group_name
  media_services_account_name = azurerm_media_asset.test.media_services_account_name
}
`, template)
}

func (MediaAssetResource) complete(data acceptance.TestData) string {
	template := MediaAssetResource{}.template(data)
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

func (MediaAssetResource) template(data acceptance.TestData) string {
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
