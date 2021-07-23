package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type CustomDatasetResource struct {
}

func TestAccDataFactoryCustomDataset_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_custom_dataset", "test")
	r := CustomDatasetResource{}

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

func TestAccDataFactoryCustomDataset_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_custom_dataset", "test")
	r := CustomDatasetResource{}

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

func TestAccDataFactoryCustomDataset_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_custom_dataset", "test")
	r := CustomDatasetResource{}

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

func TestAccDataFactoryCustomDataset_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_custom_dataset", "test")
	r := CustomDatasetResource{}

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

func TestAccDataFactoryCustomDataset_delimitedText(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_custom_dataset", "test")
	r := CustomDatasetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.delimitedText(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryCustomDataset_avro(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_custom_dataset", "test")
	r := CustomDatasetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.avro(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t CustomDatasetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DataSetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataFactory.DatasetClient.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r CustomDatasetResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_custom_dataset" "test" {
  name            = "acctestds%d"
  data_factory_id = azurerm_data_factory.test.id
  type            = "Json"

  linked_service {
    name = azurerm_data_factory_linked_custom_service.test.name
  }

  type_properties_json = <<JSON
{
  "location": {
    "container": "${azurerm_storage_container.test.name}",
    "type": "AzureBlobStorageLocation"
  }
}
JSON
}
`, r.template(data), data.RandomInteger)
}

func (r CustomDatasetResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_custom_dataset" "import" {
  name                 = azurerm_data_factory_custom_dataset.test.name
  data_factory_id      = azurerm_data_factory_custom_dataset.test.data_factory_id
  type                 = azurerm_data_factory_custom_dataset.test.type
  type_properties_json = azurerm_data_factory_custom_dataset.test.type_properties_json

  linked_service {
    name = azurerm_data_factory_custom_dataset.test.linked_service.0.name
  }
}
`, r.basic(data))
}

func (r CustomDatasetResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_custom_dataset" "test" {
  name            = "acctestds%d"
  data_factory_id = azurerm_data_factory.test.id
  type            = "Json"

  linked_service {
    name = azurerm_data_factory_linked_custom_service.test.name
    parameters = {
      key1 = "value1"
      key2 = "value2"
    }
  }

  type_properties_json = <<JSON
{
  "location": {
    "container":"${azurerm_storage_container.test.name}",
    "fileName":"foo.txt",
    "folderPath": "foo/bar/",
    "type":"AzureBlobStorageLocation"
  },
  "encodingName":"UTF-8"
}
JSON

  description = "test description"
  annotations = ["test1", "test2", "test3"]
  folder      = "testFolder"

  parameters = {
    foo = "test1"
    Bar = "Test2"
  }

  additional_properties = {
    foo = "test1"
    bar = "test2"
  }

  schema_json = <<JSON
{
  "type": "object",
  "properties": {
    "name": {
      "type": "object",
      "properties": {
        "firstName": {
          "type": "string"
        },
        "lastName": {
          "type": "string"
        }
      }
    },
    "age": {
      "type": "integer"
    }
  }
}
JSON
}
`, r.template(data), data.RandomInteger)
}

func (r CustomDatasetResource) delimitedText(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_custom_dataset" "test" {
  name            = "acctestds%d"
  data_factory_id = azurerm_data_factory.test.id
  type            = "DelimitedText"

  linked_service {
    name = azurerm_data_factory_linked_custom_service.test.name
  }

  type_properties_json = <<JSON
{
  "location": {
    "container":"test",
    "fileName":"foo.txt",
    "folderPath": "foo/bar/",
    "type":"AzureBlobStorageLocation"
  },
  "columnDelimiter": "\n",
  "rowDelimiter": "\t",
  "encodingName": "UTF-8",
  "compressionCodec": "bzip2",
  "compressionLevel": "Farest",
  "quoteChar": "",
  "escapeChar": "",
  "firstRowAsHeader": false,
  "nullValue": ""
}
JSON

  schema_json = <<JSON
[
  {
    "name": "col1",
    "type": "INT_32"
  },
  {
    "name": "col2",
    "type": "Decimal",
    "precision": "38",
    "scale": "2"
  }
]
JSON
}
`, r.template(data), data.RandomInteger)
}

func (r CustomDatasetResource) avro(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_custom_dataset" "test" {
  name            = "acctestds%d"
  data_factory_id = azurerm_data_factory.test.id
  type            = "Avro"

  linked_service {
    name = azurerm_data_factory_linked_custom_service.test.name
  }

  type_properties_json = <<JSON
{
  "location": {
    "fileName":".avro",
    "folderPath": "foo",
    "type":"AzureBlobStorageLocation"
  },
  "avroCompressionCodec": "deflate",
  "avroCompressionLevel": 4
}
JSON

  schema_json = <<JSON
{
  "type": "record",
  "namespace": "com.example",
  "name": "test",
  "fields": [
    {
      "name": "first",
      "type": "string"
    },
    {
      "name": "last",
      "type": "int"
    },
    {
      "name": "Hobby",
      "type": {
        "type": "array",
        "items": "string"
      }
    }
  ]
}
JSON
}
`, r.template(data), data.RandomInteger)
}

func (CustomDatasetResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestdf%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "content"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}
