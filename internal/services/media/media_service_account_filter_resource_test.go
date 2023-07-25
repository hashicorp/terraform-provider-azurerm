// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package media_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/accountfilters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AccountFilterResource struct{}

func TestAccMediaServicesAccountFilter_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account_filter", "test")
	r := AccountFilterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMediaServicesAccountFilter_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account_filter", "test")
	r := AccountFilterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMediaServicesAccountFilter_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account_filter", "test")
	r := AccountFilterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMediaServicesAccountFilter_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account_filter", "test")
	r := AccountFilterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (AccountFilterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := accountfilters.ParseAccountFilterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Media.V20220801Client.AccountFilters.AccountFiltersGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r AccountFilterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_services_account_filter" "test" {
  name                        = "Filter-1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
}
`, r.template(data))
}

func (r AccountFilterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_services_account_filter" "import" {
  name                        = azurerm_media_services_account_filter.test.name
  resource_group_name         = azurerm_media_services_account_filter.test.resource_group_name
  media_services_account_name = azurerm_media_services_account_filter.test.media_services_account_name
}
`, r.basic(data))
}

func (r AccountFilterResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_services_account_filter" "test" {
  name                        = "Filter-1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  first_quality_bitrate       = 128000

  presentation_time_range {
    start_in_units                 = 0
    end_in_units                   = 15
    presentation_window_in_units   = 90
    live_backoff_in_units          = 0
    unit_timescale_in_milliseconds = 1000
    force_end                      = false
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

func (AccountFilterResource) template(data acceptance.TestData) string {
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
