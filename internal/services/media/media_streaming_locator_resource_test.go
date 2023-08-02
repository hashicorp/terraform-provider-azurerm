// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package media_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/streamingpoliciesandstreaminglocators"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StreamingLocatorResource struct{}

func TestAccStreamingLocator_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_locator", "test")
	r := StreamingLocatorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Locator-1"),
				check.That(data.ResourceName).Key("asset_name").HasValue("test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStreamingLocator_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_locator", "test")
	r := StreamingLocatorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Locator-1"),
				check.That(data.ResourceName).Key("asset_name").HasValue("test"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccStreamingLocator_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_locator", "test")
	r := StreamingLocatorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("start_time").HasValue("2018-03-01T00:00:00Z"),
				check.That(data.ResourceName).Key("end_time").HasValue("2028-12-31T23:59:59Z"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStreamingLocator_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_locator", "test")
	r := StreamingLocatorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Locator-1"),
				check.That(data.ResourceName).Key("asset_name").HasValue("test"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("start_time").HasValue("2018-03-01T00:00:00Z"),
				check.That(data.ResourceName).Key("end_time").HasValue("2028-12-31T23:59:59Z"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Locator-1"),
				check.That(data.ResourceName).Key("asset_name").HasValue("test"),
			),
		},
		data.ImportStep(),
	})
}

func (StreamingLocatorResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := streamingpoliciesandstreaminglocators.ParseStreamingLocatorID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Media.V20220801Client.StreamingPoliciesAndStreamingLocators.StreamingLocatorsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r StreamingLocatorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_streaming_locator" "test" {
  name                        = "Locator-1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  streaming_policy_name       = "Predefined_ClearStreamingOnly"
  asset_name                  = azurerm_media_asset.test.name
}
`, r.template(data))
}

func (r StreamingLocatorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_streaming_locator" "import" {
  name                        = azurerm_media_streaming_locator.test.name
  resource_group_name         = azurerm_media_streaming_locator.test.resource_group_name
  media_services_account_name = azurerm_media_streaming_locator.test.media_services_account_name
  streaming_policy_name       = "Predefined_ClearStreamingOnly"
  asset_name                  = azurerm_media_asset.test.name
}
`, r.basic(data))
}

func (r StreamingLocatorResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_services_account_filter" "test" {
  name                        = "Filter-1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
}

resource "azurerm_media_streaming_locator" "test" {
  name                        = "Job-1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  streaming_policy_name       = "Predefined_DownloadOnly"
  asset_name                  = azurerm_media_asset.test.name
  start_time                  = "2018-03-01T00:00:00Z"
  end_time                    = "2028-12-31T23:59:59Z"
  streaming_locator_id        = "90000000-0000-0000-0000-000000000000"
  alternative_media_id        = "my-Alternate-MediaID"
  filter_names                = [azurerm_media_services_account_filter.test.name]
}
`, r.template(data))
}

func (StreamingLocatorResource) template(data acceptance.TestData) string {
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
