// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package media_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-07-01/encodings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MediaJobResource struct{}

func TestAccMediaJob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_job", "test")
	r := MediaJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Job-1"),
				check.That(data.ResourceName).Key("input_asset.0.name").HasValue("inputAsset"),
				check.That(data.ResourceName).Key("output_asset.0.name").HasValue("outputAsset"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMediaJob_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_job", "test")
	r := MediaJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Job-1"),
				check.That(data.ResourceName).Key("input_asset.0.name").HasValue("inputAsset"),
				check.That(data.ResourceName).Key("output_asset.0.name").HasValue("outputAsset"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMediaJob_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_job", "test")
	r := MediaJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("description").HasValue("Job description"),
				check.That(data.ResourceName).Key("priority").HasValue("Normal"),
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Job-1"),
				check.That(data.ResourceName).Key("input_asset.0.name").HasValue("inputAsset"),
				check.That(data.ResourceName).Key("output_asset.0.name").HasValue("outputAsset"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("description").HasValue("Job description"),
				check.That(data.ResourceName).Key("priority").HasValue("Normal"),
				check.That(data.ResourceName).Key("input_asset.0.label").HasValue("Input"),
				check.That(data.ResourceName).Key("output_asset.0.label").HasValue("Output"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Job-1"),
				check.That(data.ResourceName).Key("input_asset.0.name").HasValue("inputAsset"),
				check.That(data.ResourceName).Key("output_asset.0.name").HasValue("outputAsset"),
			),
		},
		data.ImportStep(),
	})
}

func (MediaJobResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := encodings.ParseJobID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Media.V20220701Client.Encodings.JobsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r MediaJobResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_job" "test" {
  name                        = "Job-1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  transform_name              = azurerm_media_transform.test.name
  input_asset {
    name = azurerm_media_asset.input.name
  }
  output_asset {
    name = azurerm_media_asset.output.name
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
  transform_name              = azurerm_media_job.test.transform_name
  input_asset {
    name = azurerm_media_job.test.input_asset[0].name
  }
  output_asset {
    name = azurerm_media_job.test.output_asset[0].name
  }
}
`, r.basic(data))
}

func (r MediaJobResource) complete(data acceptance.TestData) string {
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
    name  = azurerm_media_asset.input.name
    label = "Input"
  }
  output_asset {
    name  = azurerm_media_asset.output.name
    label = "Output"
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
  name                        = "inputAsset"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  description                 = "Input Asset description"
}

resource "azurerm_media_asset" "output" {
  name                        = "outputAsset"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  description                 = "Output Asset description"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}
