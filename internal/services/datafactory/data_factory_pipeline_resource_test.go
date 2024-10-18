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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

func TestAccDataFactoryPipeline_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_pipeline", "test")
	r := PipelineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.update1(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update2(data),
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
	id, err := parse.PipelineID(state.ID)
	if err != nil {
		return nil, err
	}

	hackClient := azuresdkhacks.PipelinesClient{
		OriginalClient: clients.DataFactory.PipelinesClient,
	}
	resp, err := hackClient.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (t PipelineResource) appendVariableActivityNameIs(expected string) func(input []interface{}) (*bool, error) {
	return func(input []interface{}) (*bool, error) {
		if len(input) == 0 || input[0] == nil {
			return utils.Bool(false), nil
		}

		val, ok := input[0].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("nested item was not a dictionary")
		}

		actual, ok := val["name"].(string)
		if !ok {
			return nil, fmt.Errorf("name was not present in the json")
		}

		return utils.Bool(actual == expected), nil
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

func (PipelineResource) update1(data acceptance.TestData) string {
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

  parameters {
    name          = "teststring"
    type          = "String"
    default_value = "teststringvalue"
  }

  parameters {
    name          = "testint"
    type          = "Int"
    default_value = "123"
  }

  parameters {
    name          = "testfloat"
    type          = "Float"
    default_value = "123.45"
  }

  parameters {
    name          = "testbool"
    type          = "Bool"
    default_value = "true"
  }

  parameters {
    name          = "testarrayint"
    type          = "Array"
    default_value = "[1, 2, 3]"
  }

  parameters {
    name          = "testarraystring"
    type          = "Array"
    default_value = jsonencode(["a", "b", "c"])
  }

  parameters {
    name = "testobject"
    type = "Object"
    default_value = jsonencode({
      key1 = "value1"
      key2 = "value2"
    })
  }

  parameters {
    name          = "testsecurestring"
    type          = "SecureString"
    default_value = "securestringvalue"
  }

  variables {
    name          = "foo"
    type          = "String"
    default_value = "test1"
  }

  variables {
    name          = "qux"
    type          = "Array"
    default_value = jsonencode(["a", "b", "c"])
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (PipelineResource) update2(data acceptance.TestData) string {
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

  parameters {
    name          = "teststring"
    type          = "String"
    default_value = "teststringvalue"
  }

  parameters {
    name          = "teststring2"
    default_value = "teststringvalue"
  }

  parameters {
    name          = "testint"
    type          = "Int"
    default_value = "123"
  }

  parameters {
    name          = "testfloat"
    type          = "Float"
    default_value = "123.45"
  }

  parameters {
    name          = "testbool"
    type          = "Bool"
    default_value = "true"
  }

  parameters {
    name          = "testarrayint"
    type          = "Array"
    default_value = "[1, 2, 3]"
  }

  parameters {
    name          = "testarraystring"
    type          = "Array"
    default_value = jsonencode(["a", "b", "c"])
  }

  parameters {
    name = "testobject"
    type = "Object"
    default_value = jsonencode({
      key1 = "value1"
      key2 = "value2"
    })
  }

  parameters {
    name          = "testsecurestring"
    type          = "SecureString"
    default_value = "securestringvalue"
  }

  variables {
    name          = "foo"
    type          = "String"
    default_value = "test1"
  }

  variables {
    name          = "bar"
    default_value = "test2"
  }

  variables {
    name          = "qux"
    type          = "Array"
    default_value = jsonencode(["a", "b", "c"])
  }
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
  variables {
    name          = "bob"
    type          = "String"
    default_value = "item1"
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
      }
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
  variables {
    name          = "bob"
    type          = "String"
    default_value = "item1"
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
  variables {
    name          = "bob"
    type          = "String"
    default_value = "item1"
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
