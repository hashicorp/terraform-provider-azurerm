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

type TriggerCustomEventResource struct{}

func TestAccDataFactoryTriggerCustomEvent_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_custom_event", "test")
	r := TriggerCustomEventResource{}

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

func TestAccDataFactoryTriggerCustomEvent_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_custom_event", "test")
	r := TriggerCustomEventResource{}

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

func TestAccDataFactoryTriggerCustomEvent_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_custom_event", "test")
	r := TriggerCustomEventResource{}

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

func TestAccDataFactoryTriggerCustomEvent_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_custom_event", "test")
	r := TriggerCustomEventResource{}

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

func (t TriggerCustomEventResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.TriggerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataFactory.TriggersClient.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r TriggerCustomEventResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_trigger_custom_event" "test" {
  name                = "acctestdf%d"
  data_factory_id     = azurerm_data_factory.test.id
  eventgrid_topic_id  = azurerm_eventgrid_topic.test.id
  events              = ["event1"]
  subject_begins_with = "abc"

  pipeline {
    name = azurerm_data_factory_pipeline.test.name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r TriggerCustomEventResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_trigger_custom_event" "import" {
  name                = azurerm_data_factory_trigger_custom_event.test.name
  data_factory_id     = azurerm_data_factory_trigger_custom_event.test.data_factory_id
  eventgrid_topic_id  = azurerm_data_factory_trigger_custom_event.test.eventgrid_topic_id
  events              = azurerm_data_factory_trigger_custom_event.test.events
  subject_begins_with = azurerm_data_factory_trigger_custom_event.test.subject_begins_with

  dynamic "pipeline" {
    for_each = azurerm_data_factory_trigger_custom_event.test.pipeline
    content {
      name = pipeline.value.name
    }
  }
}
`, r.basic(data))
}

func (r TriggerCustomEventResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_trigger_custom_event" "test" {
  name                = "acctestdf%d"
  data_factory_id     = azurerm_data_factory.test.id
  eventgrid_topic_id  = azurerm_eventgrid_topic.test.id
  events              = ["event1", "event2"]
  subject_begins_with = "abc"
  subject_ends_with   = "xyz"

  activated   = false
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

func (TriggerCustomEventResource) template(data acceptance.TestData) string {
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
    foo = "bar"
  }
}

resource "azurerm_eventgrid_topic" "test" {
  name                = "acctesteg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
