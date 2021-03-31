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

type LiveEventResource struct{}

func TestAccLiveEvent_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_live_event", "test")
	r := LiveEventResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Event-1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLiveEvent_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_live_event", "test")
	r := LiveEventResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Event-1"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccLiveEvent_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_live_event", "test")
	r := LiveEventResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("encoding.#").HasValue("1"),
				check.That(data.ResourceName).Key("preview.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func (LiveEventResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.LiveEventID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Media.LiveEventsClient.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Live Event %s (Media Services Account %s) (resource group: %s): %v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.LiveEventProperties != nil), nil
}

func (r LiveEventResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_live_event" "test" {
  name                        = "Event-1"
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

`, r.template(data))
}

func (r LiveEventResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_live_event" "import" {
  name                        = azurerm_media_live_event.test.name
  resource_group_name         = azurerm_media_live_event.test.resource_group_name
  media_services_account_name = azurerm_media_live_event.test.media_services_account_name
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

`, r.basic(data))
}

func (r LiveEventResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_live_event" "test" {
  name                        = "Event-1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  location                    = azurerm_resource_group.test.location
  description                 = "My Event Description"

  input {
    streaming_protocol = "RTMP"
    ip_access_control_allow {
      name                 = "AllowAll"
      address              = "0.0.0.0"
      subnet_prefix_length = 0
    }
  }

  encoding {
    type               = "Standard"
    preset_name        = "Default720p"
    stretch_mode       = "AutoFit"
    key_frame_interval = "PT2S"
  }

  preview {
    ip_access_control_allow {
      name                 = "AllowAll"
      address              = "0.0.0.0"
      subnet_prefix_length = 0
    }
  }

  use_static_hostname     = true
  hostname_prefix         = "special-event"
  transcription_languages = ["en-US"]
}

`, r.template(data))
}

func (LiveEventResource) template(data acceptance.TestData) string {
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
