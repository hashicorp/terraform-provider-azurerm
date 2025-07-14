// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/pipelines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PipelineResource struct{}

func TestAccDataFactoryPipeline_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_pipeline", "test")
	r := PipelineResource{}

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

func TestAccDataFactoryPipeline_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_pipeline", "test")
	r := PipelineResource{}

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

func TestAccDataFactoryPipeline_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_pipeline", "test")
	r := PipelineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.update(data),
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

func TestAccDataFactoryPipeline_activities(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_pipeline", "test")
	r := PipelineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.webActivityHeaders(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.webActivityHeaders(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.activities(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("activities_json").Exists(),
				check.That(data.ResourceName).Key("activities_json").ContainsJsonValue(r.appendVariableActivityNameIs("Append variable1")),
			),
		},
		data.ImportStep(),
		{
			Config: r.activitiesUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("activities_json").Exists(),
				check.That(data.ResourceName).Key("activities_json").ContainsJsonValue(r.appendVariableActivityNameIs("Append variable1")),
			),
		},
		data.ImportStep(),
		{
			Config: r.activities(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("activities_json").Exists(),
				check.That(data.ResourceName).Key("activities_json").ContainsJsonValue(r.appendVariableActivityNameIs("Append variable1")),
			),
		},
		data.ImportStep(),
	})
}

func (t PipelineResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := pipelines.ParsePipelineID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataFactory.PipelinesClient.Get(ctx, *id, pipelines.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (t PipelineResource) appendVariableActivityNameIs(expected string) func(input []interface{}) (*bool, error) {
	return func(input []interface{}) (*bool, error) {
		if len(input) == 0 || input[0] == nil {
			return pointer.To(false), nil
		}

		val, ok := input[0].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("nested item was not a dictionary")
		}

		actual, ok := val["name"].(string)
		if !ok {
			return nil, fmt.Errorf("name was not present in the json")
		}

		return pointer.To(actual == expected), nil
	}
}

func (PipelineResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfv2%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_pipeline" "test" {
  name            = "acctest%d"
  data_factory_id = azurerm_data_factory.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (PipelineResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfv2%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_pipeline" "test" {
  name            = "acctest%d"
  data_factory_id = azurerm_data_factory.test.id
  annotations     = ["test1", "test2", "test3"]
  description     = "test description"

  parameters = {
    test = "testparameter"
  }

  variables = {
    foo = "test1"
    bar = "test2"
  }

  activities_json = <<JSON
[
  {
    "name": "test append variable",
    "type": "AppendVariable",
    "dependsOn": [],
    "userProperties": [],
    "typeProperties": {
      "variableName": "bob",
      "value": "something"
    }
  },
  {
    "name": "test web activity",
    "type": "WebActivity",
    "dependsOn": [],
    "userProperties": [],
    "typeProperties": {
	  "url": "https://test.com",
	  "method": "POST",
      "headers": {
        "authorization": {
          "value": "foo",
          "type": "Expression"
        },
        "content_type": "application/x-www-form-urlencoded"
      }
    }
  },
  {
    "name": "test filter",
    "type": "Filter",
    "dependsOn": [
      {
        "activity": "Filter something",
        "dependencyConditions": ["Succeeded"]
      }
    ],
    "userProperties": [],
    "typeProperties": {
      "items": {
        "value": "@json(activity('Filter Something').output.response)",
        "type": "Expression"
      },
      "condition": {
        "value": "@equals(coalesce(item().Authorised, 0), 1)",
        "type": "Expression"
      }
    }
  }
]
JSON
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (PipelineResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfv2%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_pipeline" "test" {
  name            = "acctest%d"
  data_factory_id = azurerm_data_factory.test.id
  annotations     = ["test1", "test2"]
  description     = "updated description"

  parameters = {
    test  = "testparameter"
    test2 = "testparameter2"
  }

  variables = {
    foo = "test1"
    bar = "test2"
    baz = "test3"
  }

  activities_json = <<JSON
[
  {
    "name": "test append variable",
    "type": "AppendVariable",
    "dependsOn": [],
    "userProperties": [],
    "typeProperties": {
      "variableName": "bob",
      "value": "something"
    }
  },
  {
    "name": "test web activity",
    "type": "WebActivity",
    "dependsOn": [],
    "userProperties": [],
    "typeProperties": {
	  "url": "https://test.com",
	  "method": "POST"
    }
  }
]
JSON
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (PipelineResource) activities(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfv2%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_pipeline" "test" {
  name            = "acctest%d"
  data_factory_id = azurerm_data_factory.test.id
  variables = {
    "bob" = "item1"
  }
  activities_json = <<JSON
[
  {
    "name": "Append variable1",
    "type": "AppendVariable",
    "dependsOn": [],
    "userProperties": [],
    "typeProperties": {
      "variableName": "bob",
      "value": "something"
    }
  }
]
JSON
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (PipelineResource) webActivityHeaders(data acceptance.TestData, withHeader bool) string {
	headerBlock := `
      "headers": {
        "authorization": {
          "value": "foo",
          "type": "Expression"
        },
        "content_type": "application/x-www-form-urlencoded"
      },
  `
	if !withHeader {
		headerBlock = ``
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfv2%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_pipeline" "test" {
  name            = "acctest%d"
  data_factory_id = azurerm_data_factory.test.id
  variables = {
    "bob" = "item1"
  }
  activities_json = <<JSON
[
  {
    "name": "test webactivity",
    "type": "WebActivity",
    "dependsOn": [],
    "userProperties": [],
    "typeProperties": {
    %s
	  "url": "https://test.com",
	  "method": "POST"
    }
  }
]
JSON
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, headerBlock)
}

func (PipelineResource) activitiesUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfv2%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_pipeline" "test" {
  name            = "acctest%d"
  data_factory_id = azurerm_data_factory.test.id
  variables = {
    "bob" = "item1"
  }
  activities_json = <<JSON
[
  {
    "name": "Append variable1",
    "type": "AppendVariable",
    "dependsOn": [],
    "userProperties": [],
    "typeProperties": {
      "variableName": "bob",
      "value": "something"
    }
  }
]
JSON
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
