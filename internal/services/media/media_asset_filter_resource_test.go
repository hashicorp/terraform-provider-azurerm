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

type AssetFilterResource struct{}

func TestAccAssetFilter_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_asset_filter", "test")
	r := AssetFilterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Filter-1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAssetFilter_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_asset_filter", "test")
	r := AssetFilterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Filter-1"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAssetFilter_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_asset_filter", "test")
	r := AssetFilterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("first_quality_bitrate").HasValue("128000"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAssetFilter_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_asset_filter", "test")
	r := AssetFilterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Filter-1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("first_quality_bitrate").HasValue("128000"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Filter-1"),
			),
		},
		data.ImportStep(),
	})
}

func (AssetFilterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := assetsandassetfilters.ParseAssetFilterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Media.V20220801Client.AssetsAndAssetFilters.AssetFiltersGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r AssetFilterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_asset_filter" "test" {
  name     = "Filter-1"
  asset_id = azurerm_media_asset.test.id
}
`, r.template(data))
}

func (r AssetFilterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_asset_filter" "import" {
  name     = azurerm_media_asset_filter.test.name
  asset_id = azurerm_media_asset.test.id
}
`, r.basic(data))
}

func (r AssetFilterResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_asset_filter" "test" {
  name                  = "Filter-1"
  asset_id              = azurerm_media_asset.test.id
  first_quality_bitrate = 128000

  presentation_time_range {
    start_in_units                = 0
    end_in_units                  = 15
    presentation_window_in_units  = 90
    live_backoff_in_units         = 0
    unit_timescale_in_miliseconds = 1000
    force_end                     = false
  }

  track_selection {
    condition {
      property  = "Type"
      operation = "Equal"
      value     = "Audio"
    }

    condition {
      property  = "Language"
      operation = "NotEqual"
      value     = "en"
    }

    condition {
      property  = "FourCC"
      operation = "NotEqual"
      value     = "EC-3"
    }
  }


  track_selection {
    condition {
      property  = "Type"
      operation = "Equal"
      value     = "Video"
    }

    condition {
      property  = "Bitrate"
      operation = "Equal"
      value     = "3000000-5000000"
    }
  }
}
`, r.template(data))
}

func (AssetFilterResource) template(data acceptance.TestData) string {
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

resource "azurerm_media_asset" "test" {
  name                        = "test"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}
