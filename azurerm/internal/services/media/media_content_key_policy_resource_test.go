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

type MediaContentKeyPolicyResource struct {
}

func TestAccMediaContentKeyPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_content_key_policy", "test")
	r := MediaContentKeyPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Policy-1"),
				check.That(data.ResourceName).Key("policy_option.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMediaContentKeyPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_content_key_policy", "test")
	r := MediaContentKeyPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Policy-1"),
				check.That(data.ResourceName).Key("policy_option.#").HasValue("1"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMediaContentKeyPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_content_key_policy", "test")
	r := MediaContentKeyPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Policy-1"),
				check.That(data.ResourceName).Key("policy_option.#").HasValue("4"),
			),
		},
		data.ImportStep(),
	})
}

func (r MediaContentKeyPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ContentKeyPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Media.ContentKeyPoliciesClient.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Content Key Policy %s (Media Account %s) (resource group: %s): %v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.ContentKeyPolicyProperties != nil), nil
}

func (r MediaContentKeyPolicyResource) basic(data acceptance.TestData) string {
	template := r.template(data)
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
`, template)
}

func (r MediaContentKeyPolicyResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_media_content_key_policy" "import" {
  name                        = azurerm_media_content_key_policy.test.name
  resource_group_name         = azurerm_media_content_key_policy.test.resource_group_name
  media_services_account_name = azurerm_media_content_key_policy.test.media_services_account_name
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
`, template)
}

func (r MediaContentKeyPolicyResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_content_key_policy" "test" {
  name                        = "Policy-1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  description                 = "My Policy Description"
  policy_option {
    name = "fairPlay"
    fairplay_configuration {
      ask                       = "myKeyForFairPlay"
      pfx                       = "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"
      pfx_password              = "password"
      rental_duration_seconds   = 2249
      rental_and_lease_key_type = "PersistentUnlimited"

    }
    open_restriction_enabled = true
  }
  policy_option {
    name = "playReady"
    playready_configuration_license {
      allow_test_devices = true
      begin_date         = "2017-10-16T18:22:53Z"
      play_right {
        scms_restriction                                         = 2
        digital_video_only_content_restriction                   = false
        image_constraint_for_analog_component_video_restriction  = false
        image_constraint_for_analog_computer_monitor_restriction = false
        allow_passing_video_content_to_unknown_output            = "NotAllowed"
        uncompressed_digital_video_opl                           = 100
        uncompressed_digital_audio_opl                           = 100
        analog_video_opl                                         = 150
        compressed_digital_audio_opl                             = 150
      }
      license_type                             = "Persistent"
      content_type                             = "UltraVioletDownload"
      content_key_location_from_header_enabled = true
    }
    open_restriction_enabled = true
  }
  policy_option {
    name                            = "clearKey"
    clear_key_configuration_enabled = true
    token_restriction {
      issuer                      = "urn:issuer"
      audience                    = "urn:audience"
      token_type                  = "Swt"
      primary_symmetric_token_key = "AAAAAAAAAAAAAAAAAAAAAA=="
    }
  }
  policy_option {
    name = "widevine"
    widevine_configuration_template = jsonencode({
      "allowed_track_types" : "SD_HD",
      "content_key_specs" : [{
        "track_type" : "SD",
        "security_level" : 1,
        "required_output_protection" : {
          "hdcp" : "HDCP_V2"
        },
      }],
      "policy_overrides" : {
        "can_play" : true,
        "can_persist" : true,
        "can_renew" : false,
      },
    })
    open_restriction_enabled = true
  }
}
`, r.template(data))
}

func (r MediaContentKeyPolicyResource) template(data acceptance.TestData) string {
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
