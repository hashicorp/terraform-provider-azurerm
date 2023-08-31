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

type StreamingPolicyResource struct{}

func TestAccStreamingPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_policy", "test")
	r := StreamingPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Policy-1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStreamingPolicy_clearKeyEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_policy", "test")
	r := StreamingPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.clearKeyEncryption(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Policy-1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStreamingPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_policy", "test")
	r := StreamingPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Policy-1"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccStreamingPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_policy", "test")
	r := StreamingPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Policy-1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("common_encryption_cenc.#").HasValue("1"),
				check.That(data.ResourceName).Key("common_encryption_cbcs.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Policy-1"),
			),
		},
		data.ImportStep(),
	})
}

func (StreamingPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := streamingpoliciesandstreaminglocators.ParseStreamingPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Media.V20220801Client.StreamingPoliciesAndStreamingLocators.StreamingPoliciesGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
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

func (r StreamingPolicyResource) clearKeyEncryption(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_content_key_policy" "test" {
  name                        = "Policy-1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  description                 = "My Policy Description"
  policy_option {
    name                            = "ClearKeyOption"
    clear_key_configuration_enabled = true
    token_restriction {
      issuer                      = "urn:issuer"
      audience                    = "urn:audience"
      token_type                  = "Swt"
      primary_symmetric_token_key = "AAAAAAAAAAAAAAAAAAAAAA=="
    }
  }
}

resource "azurerm_media_streaming_policy" "test" {
  name                            = "Policy-1"
  resource_group_name             = azurerm_resource_group.test.name
  media_services_account_name     = azurerm_media_services_account.test.name
  default_content_key_policy_name = azurerm_media_content_key_policy.test.name
  common_encryption_cenc {
    default_content_key {
      label = "aesDefaultKey"
    }

    clear_track {
      condition {
        property  = "FourCC"
        operation = "Equal"
        value     = "hev1"
      }
    }

    enabled_protocols {
      download         = false
      dash             = true
      hls              = false
      smooth_streaming = true
    }

    clear_key_encryption {
      custom_keys_acquisition_url_template = "https://contoso.com/{AlternativeMediaId}/clearkey/"
    }
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

resource "azurerm_media_content_key_policy" "test" {
  name                        = "Policy-1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  description                 = "My Policy Description"
  policy_option {
    name                            = "ClearKeyOption"
    clear_key_configuration_enabled = true
    token_restriction {
      issuer                      = "urn:issuer"
      audience                    = "urn:audience"
      token_type                  = "Swt"
      primary_symmetric_token_key = "AAAAAAAAAAAAAAAAAAAAAA=="
    }
  }
}

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

    clear_track {
      condition {
        property  = "FourCC"
        operation = "Equal"
        value     = "hev2"
      }
    }

    clear_track {
      condition {
        property  = "FourCC"
        operation = "Equal"
        value     = "hev1"
      }
    }

    default_content_key {
      label       = "aesDefaultKey"
      policy_name = azurerm_media_content_key_policy.test.name
    }

    content_key_to_track_mapping {
      label       = "aesKey"
      policy_name = azurerm_media_content_key_policy.test.name
      track {
        condition {
          property  = "FourCC"
          operation = "Equal"
          value     = "hev1"
        }
      }
    }

    drm_playready {
      custom_license_acquisition_url_template = "https://contoso.com/{AssetAlternativeId}/playready/{ContentKeyId}"
      custom_attributes                       = "PlayReady CustomAttributes"
    }
    drm_widevine_custom_license_acquisition_url_template = "https://contoso.com/{AssetAlternativeId}/widevine/{ContentKeyId}"
  }

  common_encryption_cbcs {
    default_content_key {
      label       = "aesDefaultKey"
      policy_name = azurerm_media_content_key_policy.test.name
    }
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

  envelope_encryption {
    default_content_key {
      label       = "aesDefaultKey"
      policy_name = azurerm_media_content_key_policy.test.name
    }
    custom_keys_acquisition_url_template = "https://contoso.com/{AssetAlternativeId}/envelope/{ContentKeyId}"
    enabled_protocols {
      dash             = true
      download         = false
      hls              = true
      smooth_streaming = true
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
