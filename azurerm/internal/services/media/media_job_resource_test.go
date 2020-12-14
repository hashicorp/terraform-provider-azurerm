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

type MediaJobResource struct {
}

func TestAccMediaJob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_job", "test")
	r := MediaJobResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Job-1"),
				check.That(data.ResourceName).Key("description").HasValue("Job description"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMediaJob_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_job", "test")
	r := MediaJobResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Job-1"),
				check.That(data.ResourceName).Key("description").HasValue("Job description"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMediaJob_label(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_job", "test")
	r := MediaJobResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.label(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("input_asset.0.label").HasValue("Input"),
				check.That(data.ResourceName).Key("output_asset.0.label").HasValue("Output"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMediaJob_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_job", "test")
	r := MediaJobResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Job-1"),
				check.That(data.ResourceName).Key("description").HasValue("Job description"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("description").HasValue("Updated description"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("description").HasValue("Job description"),
			),
		},
		data.ImportStep(),
	})
}

func (MediaJobResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.JobID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Media.JobsClient.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.TransformName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Job %s (Media Services Account %s) (resource group: %s): %v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.JobProperties != nil), nil
}

func (r MediaJobResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_job" "test" {
  name                        = "Job-1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  transform_name              = azurerm_media_transform.test.name
  description                 = "Job description"
  priority                    = "Normal"
  input_asset {
    asset_name = azurerm_media_asset.input.name
  }
  output_asset {
    asset_name = azurerm_media_asset.output.name
  }
}

`, r.template(data))
}

func (r MediaJobResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_job" "import" {
  name                        = azurerm_media_job.test.name
  resource_group_name         = azurerm_media_job.test.resource_group_name
  media_services_account_name = azurerm_media_job.test.media_services_account_name
}

`, r.basic(data))
}

func (r MediaJobResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_job" "test" {
  name                        = "Job-1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  transform_name              = azurerm_media_transform.test.name
  description                 = "Updated description"
  priority                    = "Low"
  input_asset {
    asset_name = azurerm_media_asset.input.name
  }
  output_asset {
    asset_name = azurerm_media_asset.output.name
  }
}

`, r.template(data))
}

func (r MediaJobResource) label(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_job" "test" {
  name                        = "Job-1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  transform_name              = azurerm_media_transform.test.name
  description                 = "Job description"
  priority                    = "Normal"
  input_asset {
    asset_name = azurerm_media_asset.input.name
    label      = "Input"
  }
  output_asset {
    asset_name = azurerm_media_asset.output.name
    label      = "Output"
  }
}

`, r.template(data))
}

func (MediaJobResource) template(data acceptance.TestData) string {
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

resource "azurerm_media_transform" "test" {
  name                        = "transform1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  output {
    relative_priority = "Normal"
    on_error_action   = "ContinueJob"
    builtin_preset {
      preset_name = "AACGoodQualityAudio"
    }
  }
}

resource "azurerm_media_asset" "input" {
  name                        = "input"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  description                 = "Input Asset description"
}

resource "azurerm_media_asset" "output" {
  name                        = "output"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  description                 = "Output Asset description"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}
