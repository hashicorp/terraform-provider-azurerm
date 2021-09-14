package synapse_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CustomDatasetResource struct {
}

func TestAccSynapseDataset_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_dataset", "test")
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

func TestAccSynapseDataset_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_dataset", "test")
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

func TestAccSynapseDataset_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_dataset", "test")
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

func TestAccSynapseDataset_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_dataset", "test")
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

func TestAccSynapseDataset_delimitedText(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_dataset", "test")
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

func TestAccSynapseDataset_avro(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_dataset", "test")
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
	id, err := parse.DatasetID(state.ID)
	if err != nil {
		return nil, err
	}

	environment := clients.Account.Environment
	client, err := clients.Synapse.DatasetClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return nil, err
	}

	resp, err := client.GetDataset(ctx, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r CustomDatasetResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_dataset" "test" {
  name                 = "acctestds%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  type                 = "Json"

  linked_service {
    name = azurerm_synapse_linked_service.test.name
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

resource "azurerm_synapse_dataset" "import" {
  name                 = azurerm_synapse_dataset.test.name
  synapse_workspace_id = azurerm_synapse_dataset.test.synapse_workspace_id
  type                 = azurerm_synapse_dataset.test.type
  type_properties_json = azurerm_synapse_dataset.test.type_properties_json

  linked_service {
    name = azurerm_synapse_dataset.test.linked_service.0.name
  }
}
`, r.basic(data))
}

func (r CustomDatasetResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_dataset" "test" {
  name                 = "acctestds%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  type                 = "Json"

  linked_service {
    name = azurerm_synapse_linked_service.test.name
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

resource "azurerm_synapse_dataset" "test" {
  name                 = "acctestds%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  type                 = "DelimitedText"

  linked_service {
    name = azurerm_synapse_linked_service.test.name
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

resource "azurerm_synapse_dataset" "test" {
  name                 = "acctestds%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  type                 = "Avro"

  linked_service {
    name = azurerm_synapse_linked_service.test.name
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
  name                     = "acctestsa%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "content"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%d"
  storage_account_id = azurerm_storage_account.test.id
}

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestdf%d"
  location                             = azurerm_resource_group.test.location
  resource_group_name                  = azurerm_resource_group.test.name
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
  managed_virtual_network_enabled      = true
}

resource "azurerm_synapse_firewall_rule" "test" {
  name                 = "AllowAll"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  start_ip_address     = "0.0.0.0"
  end_ip_address       = "255.255.255.255"
}

resource "azurerm_synapse_linked_service" "test" {
  name                 = "acctestls%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  type                 = "AzureBlobStorage"
  type_properties_json = <<JSON
{
  "connectionString": "${azurerm_storage_account.test.primary_connection_string}"
}
JSON

  depends_on = [
    azurerm_synapse_firewall_rule.test,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
