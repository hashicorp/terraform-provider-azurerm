package media_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type LiveOutputResource struct{}

func TestAccLiveOutput_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_live_event_output", "test")
	r := LiveOutputResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Output-1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLiveOutput_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_live_event_output", "test")
	r := LiveOutputResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Output-1"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccLiveOutput_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_live_event_output", "test")
	r := LiveOutputResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("manifest_name").HasValue("testmanifest"),
				check.That(data.ResourceName).Key("hls_fragments_per_ts_segment").HasValue("5"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLiveOutput_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_live_event_output", "test")
	r := LiveOutputResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Output-1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("manifest_name").HasValue("testmanifest"),
				check.That(data.ResourceName).Key("hls_fragments_per_ts_segment").HasValue("5"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Output-1"),
			),
		},
		data.ImportStep(),
	})
}

func (LiveOutputResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.LiveOutputID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Media.LiveOutputsClient.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.LiveeventName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Live Event Output %s (Media Services Account %s) (resource group: %s): %v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.LiveOutputProperties != nil), nil
}

func (r LiveOutputResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_live_event_output" "test" {
  name                    = "Output-1"
  live_event_id           = azurerm_media_live_event.test.id
  archive_window_duration = "PT5M"
  asset_name              = azurerm_media_asset.test.name
}

`, r.template(data))
}

func (r LiveOutputResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_live_event_output" "import" {
  name                    = azurerm_media_live_event_output.test.name
  live_event_id           = azurerm_media_live_event.test.id
  archive_window_duration = "PT5M"
  asset_name              = azurerm_media_asset.test.name
}

`, r.basic(data))
}

func (r LiveOutputResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_live_event_output" "test" {
  name                         = "Output-2"
  live_event_id                = azurerm_media_live_event.test.id
  archive_window_duration      = "PT5M"
  asset_name                   = azurerm_media_asset.test.name
  description                  = "Test live output 1"
  manifest_name                = "testmanifest"
  output_snap_time_in_seconds  = 0
  hls_fragments_per_ts_segment = 5
}

`, r.template(data))
}

func (LiveOutputResource) template(data acceptance.TestData) string {
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
  name                        = "inputAsset"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
}

resource "azurerm_media_live_event" "test" {
  name                        = "event"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  location                    = azurerm_resource_group.test.location

  input {
    streaming_protocol          = "RTMP"
    key_frame_interval_duration = "PT6S"
    ip_access_control_allow {
      name                 = "AllowAll"
      address              = "0.0.0.0"
      subnet_prefix_length = 0
    }
  }
}

`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}
