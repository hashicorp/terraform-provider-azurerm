// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package media_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/liveevents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LiveEventResource struct{}

func TestAccLiveEvent_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_live_event", "test")
	r := LiveEventResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Event-1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLiveEvent_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_live_event", "test")
	r := LiveEventResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Event-1"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccLiveEvent_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_live_event", "test")
	r := LiveEventResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccLiveEvent_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_live_event", "test")
	r := LiveEventResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("encoding.#").HasValue("1"),
				check.That(data.ResourceName).Key("preview.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func (LiveEventResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := liveevents.ParseLiveEventID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Media.V20220801Client.LiveEvents.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
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

func (r LiveEventResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_live_event" "test" {
  name                        = "Event-1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  location                    = azurerm_resource_group.test.location
  description                 = "Updated Description"

  input {
    streaming_protocol = "RTMP"
    ip_access_control_allow {
      name                 = "Test"
      address              = "0.0.0.0"
      subnet_prefix_length = 4
    }
  }

  encoding {
    type               = "Standard"
    preset_name        = "Default720p"
    stretch_mode       = "AutoSize"
    key_frame_interval = "PT2S"
  }

  preview {
    ip_access_control_allow {
      name                 = "Allow"
      address              = "0.0.0.0"
      subnet_prefix_length = 4
    }
  }

  use_static_hostname     = true
  hostname_prefix         = "special-event-update"
  stream_options          = ["LowLatency"]
  transcription_languages = ["en-GB"]

  tags = {
    env = "test"
  }
}
`, r.template(data))
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
  stream_options          = ["LowLatency"]
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
