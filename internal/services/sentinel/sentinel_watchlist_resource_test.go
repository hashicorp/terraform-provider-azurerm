// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/watchlists"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WatchlistResource struct{}

func TestAccWatchlist_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_watchlist", "test")
	r := WatchlistResource{}

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

func TestAccWatchlist_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_watchlist", "test")
	r := WatchlistResource{}

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

func TestAccWatchlist_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_watchlist", "test")
	r := WatchlistResource{}

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

func (r WatchlistResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.Sentinel.WatchlistsClient

	id, err := watchlists.ParseWatchlistID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (r WatchlistResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_watchlist" "test" {
  name                       = "accTestWL-%d"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  display_name               = "test"
  item_search_key            = "Key"
}
`, template, data.RandomInteger)
}

func (r WatchlistResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_watchlist" "test" {
  name                       = "accTestWL-%d"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  display_name               = "test"
  description                = "description"
  labels                     = ["label1", "laebl2"]
  default_duration           = "P2DT3H"
  item_search_key            = "Key"
}
`, template, data.RandomInteger)
}

func (r WatchlistResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_watchlist" "import" {
  name                       = azurerm_sentinel_watchlist.test.name
  log_analytics_workspace_id = azurerm_sentinel_watchlist.test.log_analytics_workspace_id
  display_name               = azurerm_sentinel_watchlist.test.display_name
  item_search_key            = azurerm_sentinel_watchlist.test.item_search_key
}
`, template)
}

func (r WatchlistResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sentinel-%d"
  location = %q
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-workspace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "test" {
  workspace_id = azurerm_log_analytics_workspace.test.id
}


`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
