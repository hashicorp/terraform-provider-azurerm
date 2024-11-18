package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DatasetCosmosDbMongoDbResource struct{}

func TestAccDataFactoryDatasetCosmosDbMongoDb_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_dataset_cosmosdb_mongoapi", "test")
	r := DatasetCosmosDbMongoDbResource{}

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

func TestAccDataFactoryDatasetCosmosDbMongoDb_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_dataset_cosmosdb_mongoapi", "test")
	r := DatasetCosmosDbMongoDbResource{}

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

func TestAccDataFactoryDatasetCosmosDbMongoDb_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_dataset_cosmosdb_mongoapi", "test")
	r := DatasetCosmosDbMongoDbResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryDatasetCosmosDbMongoDb_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_dataset_cosmosdb_mongoapi", "test")
	r := DatasetCosmosDbMongoDbResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("azurerm_data_factory_dataset_cosmosdb_mongoapi.test", plancheck.ResourceActionUpdate),
				},
			},
		},
		data.ImportStep(),
	})
}

func (t DatasetCosmosDbMongoDbResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DataSetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataFactory.DatasetClient.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func commonConfig(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_cosmosdb_mongoapi" "test" {
  connection_string = "mongodb://acc:pass@foobar.documents.azure.com:10255"
  name              = "ls-cosmosdb-mongoapi-%d"
  data_factory_id   = azurerm_data_factory.test.id
  database          = "mydbname"
}
	`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (DatasetCosmosDbMongoDbResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_dataset_cosmosdb_mongoapi" "test" {
  collection_name     = "collection-1"
  data_factory_id     = azurerm_data_factory.test.id
  linked_service_name = azurerm_data_factory_linked_service_cosmosdb_mongoapi.test.name
  name                = "name-1"
}
	`, commonConfig(data))
}

func (r DatasetCosmosDbMongoDbResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_dataset_cosmosdb_mongoapi" "import" {
  collection_name     = "collection-1"
  data_factory_id     = azurerm_data_factory.test.id
  linked_service_name = azurerm_data_factory_linked_service_cosmosdb_mongoapi.test.name
  name                = "name-1"
}
	`, r.basic(data))
}

func (DatasetCosmosDbMongoDbResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_dataset_cosmosdb_mongoapi" "test" {
  additional_properties = {
	"additionalProp1" = "value1"
  }
  annotations = [ "annotation1" ]
  collection_name = "collection-1"
  data_factory_id = azurerm_data_factory.test.id
  description = "some-description"
  folder = "folder-1"
  linked_service_name = azurerm_data_factory_linked_service_cosmosdb_mongoapi.test.name
  name = "name-1"
  parameters = {
	"param1" = "value1"
  }
}
  `, commonConfig(data))
}
