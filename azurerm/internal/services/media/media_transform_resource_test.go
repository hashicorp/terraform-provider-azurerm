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

type MediaTransformResource struct {
}

func TestAccMediaTransform_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_transform", "test")
	r := MediaTransformResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Transform-1"),
				check.That(data.ResourceName).Key("output.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMediaTransform_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_transform", "test")
	r := MediaTransformResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Transform-1"),
				check.That(data.ResourceName).Key("output.#").HasValue("1"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMediaTransform_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_transform", "test")
	r := MediaTransformResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("description").HasValue("Transform description"),
				check.That(data.ResourceName).Key("output.#").HasValue("4"),
				check.That(data.ResourceName).Key("name").HasValue("Transform-1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMediaTransform_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_transform", "test")
	r := MediaTransformResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Transform-1"),
				check.That(data.ResourceName).Key("output.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("description").HasValue("Transform description"),
				check.That(data.ResourceName).Key("output.#").HasValue("4"),
				check.That(data.ResourceName).Key("name").HasValue("Transform-1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Transform-1"),
				check.That(data.ResourceName).Key("output.#").HasValue("1"),
				check.That(data.ResourceName).Key("description").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func (r MediaTransformResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.TransformID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Media.TransformsClient.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Transform %s (Media Account %s) (resource group: %s): %v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.TransformProperties != nil), nil
}

func (r MediaTransformResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_media_transform" "test" {
  name                        = "Transform-1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  output {
    relative_priority = "High"
    on_error_action   = "ContinueJob"
    builtin_preset {
      preset_name = "AACGoodQualityAudio"
    }
  }
}
`, template)
}

func (r MediaTransformResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_media_transform" "import" {
  name                        = azurerm_media_transform.test.name
  resource_group_name         = azurerm_media_transform.test.resource_group_name
  media_services_account_name = azurerm_media_transform.test.media_services_account_name

  output {
    relative_priority = "High"
    on_error_action   = "ContinueJob"
    builtin_preset {
      preset_name = "AACGoodQualityAudio"
    }
  }
}
`, template)
}

func (r MediaTransformResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_transform" "test" {
  name                        = "Transform-1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  description                 = "Transform description"
  output {
    relative_priority = "High"
    on_error_action   = "ContinueJob"
    builtin_preset {
      preset_name = "AACGoodQualityAudio"
    }
  }

  output {
    relative_priority = "High"
    on_error_action   = "StopProcessingJob"
    audio_analyzer_preset {
      audio_language      = "en-US"
      audio_analysis_mode = "Basic"
    }
  }

  output {
    relative_priority = "Low"
    on_error_action   = "StopProcessingJob"
    face_detector_preset {
      analysis_resolution = "StandardDefinition"
    }
  }

  output {
    relative_priority = "Normal"
    on_error_action   = "StopProcessingJob"
    video_analyzer_preset {
      audio_language      = "en-US"
      audio_analysis_mode = "Basic"
      insights_type       = "AllInsights"
    }
  }
}
`, r.template(data))
}

func (r MediaTransformResource) template(data acceptance.TestData) string {
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
