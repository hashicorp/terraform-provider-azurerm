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

type TriggerTumblingWindowResource struct{}

func TestAccDataFactoryTriggerTumblingWindow_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_tumbling_window", "test")
	r := TriggerTumblingWindowResource{}

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

func TestAccDataFactoryTriggerTumblingWindow_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_tumbling_window", "test")
	r := TriggerTumblingWindowResource{}

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

func TestAccDataFactoryTriggerTumblingWindow_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_tumbling_window", "test")
	r := TriggerTumblingWindowResource{}

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

func TestAccDataFactoryTriggerTumblingWindow_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_tumbling_window", "test")
	r := TriggerTumblingWindowResource{}

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

func (r TriggerTumblingWindowResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r TriggerTumblingWindowResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_trigger_tumbling_window" "test" {
  name            = "acctestdft%d"
  data_factory_id = azurerm_data_factory.test.id
  frequency       = "Minute"
  interval        = 15
  start_time      = "2022-09-21T00:00:00Z"

  pipeline {
    name = azurerm_data_factory_pipeline.test.name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r TriggerTumblingWindowResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_trigger_tumbling_window" "import" {
  name            = azurerm_data_factory_trigger_tumbling_window.test.name
  data_factory_id = azurerm_data_factory_trigger_tumbling_window.test.data_factory_id
  frequency       = azurerm_data_factory_trigger_tumbling_window.test.frequency
  interval        = azurerm_data_factory_trigger_tumbling_window.test.interval
  start_time      = azurerm_data_factory_trigger_tumbling_window.test.start_time

  dynamic "pipeline" {
    for_each = azurerm_data_factory_trigger_tumbling_window.test.pipeline
    content {
      name = pipeline.value.name
    }
  }
}
`, r.basic(data))
}

func (r TriggerTumblingWindowResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_trigger_tumbling_window" "test" {
  name            = "acctestdft%d"
  data_factory_id = azurerm_data_factory.test.id
  start_time      = "2022-09-21T00:00:00Z"
  end_time        = "2022-09-21T08:00:00Z"
  frequency       = "Minute"
  interval        = 15
  delay           = "16:00:00"

  activated   = false
  annotations = ["test1", "test2", "test3"]
  description = "test description"

  retry {
    count    = 1
    interval = 30
  }

  pipeline {
    name = azurerm_data_factory_pipeline.test.name
    parameters = {
      Env = "Test"
    }
  }

  // Self dependency
  trigger_dependency {
    size   = "24:00:00"
    offset = "-24:00:00"
  }

  trigger_dependency {
    size         = "06:00:00"
    offset       = "06:00:00"
    trigger_name = azurerm_data_factory_trigger_tumbling_window.dependency.name
  }

  additional_properties = {
    foo = "test1"
    bar = "test2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (TriggerTumblingWindowResource) template(data acceptance.TestData) string {
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
  name            = "acctestdfp%d"
  data_factory_id = azurerm_data_factory.test.id

  parameters = {
    test = "testparameter"
  }
}

resource "azurerm_data_factory_trigger_tumbling_window" "dependency" {
  name            = "acctestdft2%d"
  data_factory_id = azurerm_data_factory.test.id
  frequency       = "Minute"
  interval        = 15
  start_time      = "2022-09-21T00:00:00Z"

  pipeline {
    name = azurerm_data_factory_pipeline.test.name
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
