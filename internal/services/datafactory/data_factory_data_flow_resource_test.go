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

type DataFlowResource struct{}

func TestAccDataFactoryDataFlow_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_data_flow", "test")
	r := DataFlowResource{}

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

func TestAccDataFactoryDataFlow_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_data_flow", "test")
	r := DataFlowResource{}

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

func TestAccDataFactoryDataFlow_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_data_flow", "test")
	r := DataFlowResource{}

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

func TestAccDataFactoryDataFlow_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_data_flow", "test")
	r := DataFlowResource{}

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

func (t DataFlowResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DataFlowID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataFactory.DataFlowClient.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r DataFlowResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_data_flow" "test" {
  name            = "acctestdf%d"
  data_factory_id = azurerm_data_factory.test.id

  source {
    name = "source1"

    linked_service {
      name = azurerm_data_factory_linked_custom_service.test.name
    }
  }

  sink {
    name = "sink1"

    linked_service {
      name = azurerm_data_factory_linked_custom_service.test.name
    }
  }

  script = <<EOT
source(
  allowSchemaDrift: true, 
  validateSchema: false, 
  limit: 100, 
  ignoreNoFilesFound: false, 
  documentForm: 'documentPerLine') ~> source1 
source1 sink(
  allowSchemaDrift: true, 
  validateSchema: false, 
  skipDuplicateMapInputs: true, 
  skipDuplicateMapOutputs: true) ~> sink1
EOT
}
`, r.template(data), data.RandomInteger)
}

func (r DataFlowResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_data_flow" "import" {
  name            = azurerm_data_factory_data_flow.test.name
  data_factory_id = azurerm_data_factory_data_flow.test.data_factory_id
  script          = azurerm_data_factory_data_flow.test.script
  source {
    name = azurerm_data_factory_data_flow.test.source.0.name
    linked_service {
      name = azurerm_data_factory_data_flow.test.source.0.linked_service.0.name
    }
  }

  sink {
    name = azurerm_data_factory_data_flow.test.sink.0.name
    linked_service {
      name = azurerm_data_factory_data_flow.test.sink.0.linked_service.0.name
    }
  }
}
`, r.basic(data))
}

func (r DataFlowResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_data_flow" "test" {
  name            = "acctestdf%d"
  data_factory_id = azurerm_data_factory.test.id
  description     = "description for data flow"
  annotations     = ["anno1", "anno2"]
  folder          = "folder1"

  source {
    name        = "source1"
    description = "description for source1"

    linked_service {
      name = azurerm_data_factory_linked_custom_service.test.name
      parameters = {
        "Key1" = "value1"
      }
    }

    schema_linked_service {
      name = azurerm_data_factory_linked_custom_service.test.name
      parameters = {
        "Key1" = "value1"
      }
    }
  }

  sink {
    name        = "sink1"
    description = "description for sink1"

    linked_service {
      name = azurerm_data_factory_linked_custom_service.test.name
      parameters = {
        "Key1" = "value1"
      }
    }

    schema_linked_service {
      name = azurerm_data_factory_linked_custom_service.test.name
      parameters = {
        "Key1" = "value1"
      }
    }
  }

  transformation {
    name        = "filter1"
    description = "description for filter1"

    dataset {
      name = azurerm_data_factory_dataset_json.test1.name
      parameters = {
        "Key1" = "value1"
      }
    }

    linked_service {
      name = azurerm_data_factory_linked_custom_service.test.name
      parameters = {
        "Key1" = "value1"
      }
    }
  }

  script_lines = [<<EOT
source(output(
		movie as string,
		title as string,
		genres as string,
		year as string,
		Rating as string,
		{Rotton Tomato} as string
	),
	allowSchemaDrift: true,
	validateSchema: false,
	limit: 100,
	ignoreNoFilesFound: false) ~> source1
source1 filter(toInteger(year) >= 1910 && toInteger(year) <= 2000) ~> Filter1
Filter1 sink(allowSchemaDrift: true,
	validateSchema: false,
	skipDuplicateMapInputs: true,
	skipDuplicateMapOutputs: true,
	saveOrder: 0,
	partitionBy('roundRobin', 3)) ~> sink1
EOT,
<<EOT
source(output(
		movie as string,
		title as string,
		genres as string,
		year as string,
		Rating as string,
		{Rotton Tomato} as string
	),
	allowSchemaDrift: true,
	validateSchema: false,
	limit: 100,
	ignoreNoFilesFound: false) ~> source1
source1 filter(toInteger(year) >= 1910 && toInteger(year) <= 2000) ~> Filter1
Filter1 sink(allowSchemaDrift: true,
	validateSchema: false,
	skipDuplicateMapInputs: true,
	skipDuplicateMapOutputs: true,
	saveOrder: 0,
	partitionBy('roundRobin', 3)) ~> sink1
EOT
  ]

  script = <<EOT
source(output(
		movie as string,
		title as string,
		genres as string,
		year as string,
		Rating as string,
		{Rotton Tomato} as string
	),
	allowSchemaDrift: true,
	validateSchema: false,
	limit: 100,
	ignoreNoFilesFound: false) ~> source1
source1 filter(toInteger(year) >= 1910 && toInteger(year) <= 2000) ~> Filter1
Filter1 sink(allowSchemaDrift: true,
	validateSchema: false,
	skipDuplicateMapInputs: true,
	skipDuplicateMapOutputs: true,
	saveOrder: 0,
	partitionBy('roundRobin', 3)) ~> sink1
EOT
}
`, r.template(data), data.RandomInteger)
}

func (DataFlowResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_linked_custom_service" "test" {
  name                 = "acctestls%d"
  data_factory_id      = azurerm_data_factory.test.id
  type                 = "AzureBlobStorage"
  type_properties_json = <<JSON
{
  "connectionString": "${azurerm_storage_account.test.primary_connection_string}"
}
JSON
}

resource "azurerm_data_factory_dataset_json" "test1" {
  name                = "acctestds1%d"
  data_factory_id     = azurerm_data_factory.test.id
  linked_service_name = azurerm_data_factory_linked_custom_service.test.name

  azure_blob_storage_location {
    container = "container"
    path      = "foo/bar/"
    filename  = "foo.txt"
  }

  encoding = "UTF-8"
}

resource "azurerm_data_factory_dataset_json" "test2" {
  name                = "acctestds2%d"
  data_factory_id     = azurerm_data_factory.test.id
  linked_service_name = azurerm_data_factory_linked_custom_service.test.name

  azure_blob_storage_location {
    container = "container"
    path      = "foo/bar/"
    filename  = "bar.txt"
  }

  encoding = "UTF-8"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
