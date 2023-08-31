// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type TriggerBlobEventResource struct{}

func TestAccDataFactoryTriggerBlobEvent_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_blob_event", "test")
	r := TriggerBlobEventResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryTriggerBlobEvent_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_blob_event", "test")
	r := TriggerBlobEventResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDataFactoryTriggerBlobEvent_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_blob_event", "test")
	r := TriggerBlobEventResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryTriggerBlobEvent_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_blob_event", "test")
	r := TriggerBlobEventResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryTriggerBlobEvent_startStop(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_blob_event", "test")
	r := TriggerBlobEventResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.activated(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("activated").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.activated(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("activated").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.activated(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("activated").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func (t TriggerBlobEventResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.TriggerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataFactory.TriggersClient.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r TriggerBlobEventResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_trigger_blob_event" "test" {
  name                  = "acctestdf%d"
  data_factory_id       = azurerm_data_factory.test.id
  storage_account_id    = azurerm_storage_account.test.id
  events                = ["Microsoft.Storage.BlobCreated"]
  blob_path_begins_with = "/abc/blobs"
  activated             = false

  pipeline {
    name = azurerm_data_factory_pipeline.test.name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r TriggerBlobEventResource) activated(data acceptance.TestData, activated bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_trigger_blob_event" "test" {
  name                  = "acctestdf%d"
  data_factory_id       = azurerm_data_factory.test.id
  storage_account_id    = azurerm_storage_account.test.id
  events                = ["Microsoft.Storage.BlobCreated"]
  blob_path_begins_with = "/${azurerm_storage_container.test.name}/blobs/"
  ignore_empty_blobs    = true
  activated             = %t

  pipeline {
    name = azurerm_data_factory_pipeline.test.name
  }
}
`, r.template(data), data.RandomInteger, activated)
}

func (r TriggerBlobEventResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_trigger_blob_event" "import" {
  name                  = azurerm_data_factory_trigger_blob_event.test.name
  data_factory_id       = azurerm_data_factory_trigger_blob_event.test.data_factory_id
  storage_account_id    = azurerm_data_factory_trigger_blob_event.test.storage_account_id
  events                = azurerm_data_factory_trigger_blob_event.test.events
  blob_path_begins_with = azurerm_data_factory_trigger_blob_event.test.blob_path_begins_with

  dynamic "pipeline" {
    for_each = azurerm_data_factory_trigger_blob_event.test.pipeline
    content {
      name = pipeline.value.name
    }
  }
}
`, r.basic(data))
}

func (r TriggerBlobEventResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_trigger_blob_event" "test" {
  name                  = "acctestdf%d"
  data_factory_id       = azurerm_data_factory.test.id
  storage_account_id    = azurerm_storage_account.test.id
  events                = ["Microsoft.Storage.BlobCreated", "Microsoft.Storage.BlobDeleted"]
  blob_path_begins_with = "/${azurerm_storage_container.test.name}/blobs/"
  blob_path_ends_with   = ".txt"
  ignore_empty_blobs    = true
  activated             = true

  annotations = ["test1", "test2", "test3"]
  description = "test description"

  pipeline {
    name = azurerm_data_factory_pipeline.test.name
    parameters = {
      Env = "Test"
    }
  }

  additional_properties = {
    foo = "test1"
    bar = "test2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (TriggerBlobEventResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_pipeline" "test" {
  name            = "acctest%d"
  data_factory_id = azurerm_data_factory.test.id

  parameters = {
    test = "testparameter"
  }

  variables = {
    test = "testvariable"
  }

  activities_json = <<JSON
[
    {
        "name": "Append variable",
        "type": "AppendVariable",
        "dependsOn": [],
        "userProperties": [],
        "typeProperties": {
            "variableName": "test",
            "value": "something"
        }
    }
]
  JSON
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test-sc"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_data_factory_linked_service_azure_blob_storage" "blob_link" {
  name                 = "acctestsalink%s"
  data_factory_id      = azurerm_data_factory.test.id
  use_managed_identity = true

  service_endpoint = azurerm_storage_account.test.primary_blob_endpoint
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomString)
}
