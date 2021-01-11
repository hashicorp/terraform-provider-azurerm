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

type StreamingPolicyResource struct{}

func TestAccStreamingPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_policy", "test")
	r := StreamingPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Policy-1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStreamingPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_policy", "test")
	r := StreamingPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Policy-1"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccStreamingPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_policy", "test")
	r := StreamingPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("common_encryption_cenc.#").HasValue("1"),
				check.That(data.ResourceName).Key("common_encryption_cbcs.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStreamingPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_policy", "test")
	r := StreamingPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Policy-1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("common_encryption_cenc.#").HasValue("1"),
				check.That(data.ResourceName).Key("common_encryption_cbcs.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Policy-1"),
			),
		},
		data.ImportStep(),
	})
}

func (StreamingPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.StreamingPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Media.StreamingPoliciesClient.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Streaming Policy %s (Media Services Account %s) (resource group: %s): %v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.StreamingPolicyProperties != nil), nil
}

func (r StreamingPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_streaming_policy" "test" {
  name                        = "Policy-1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  no_encryption_enabled_protocols {
    download         = true
    dash             = true
    hls              = true
    smooth_streaming = true
  }
}

`, r.template(data))
}

func (r StreamingPolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_streaming_policy" "import" {
  name                        = azurerm_media_streaming_policy.test.name
  resource_group_name         = azurerm_media_streaming_policy.test.resource_group_name
  media_services_account_name = azurerm_media_streaming_policy.test.media_services_account_name
  no_encryption_enabled_protocols {
    download         = true
    dash             = true
    hls              = true
    smooth_streaming = true
  }
}

`, r.basic(data))
}

func (r StreamingPolicyResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_streaming_policy" "test" {
  name                        = "Policy-1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  common_encryption_cenc {
    enabled_protocols {
      download         = false
      dash             = true
      hls              = false
      smooth_streaming = false
    }
    drm_playready {
      custom_license_acquisition_url_template = "https://contoso.com/{AssetAlternativeId}/playready/{ContentKeyId}"
      custom_attributes                       = "PlayReady CustomAttributes"
    }
    drm_widevine_custom_license_acquisition_url_template = "https://contoso.com/{AssetAlternativeId}/widevine/{ContentKeyId"
  }

  common_encryption_cbcs {
    enabled_protocols {
      download         = false
      dash             = true
      hls              = false
      smooth_streaming = false
    }
    drm_fairplay {
      custom_license_acquisition_url_template = "https://contoso.com/{AssetAlternativeId}/fairplay/{ContentKeyId}"
      allow_persistent_license                = true
    }
  }
}

`, r.template(data))
}

func (StreamingPolicyResource) template(data acceptance.TestData) string {
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
