package kusto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type KustoEventGridDataConnectionResource struct {
}

func TestAccKustoEventGridDataConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_eventgrid_data_connection", "test")
	r := KustoEventGridDataConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoEventGridDataConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_eventgrid_data_connection", "test")
	r := KustoEventGridDataConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccKustoEventGridDataConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_eventgrid_data_connection", "test")
	r := KustoEventGridDataConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoEventGridDataConnection_mappingRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_eventgrid_data_connection", "test")
	r := KustoEventGridDataConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.mappingRule(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoEventGridDataConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_eventgrid_data_connection", "test")
	r := KustoEventGridDataConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (KustoEventGridDataConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.DataConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Kusto.DataConnectionsClient.Get(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	value, ok := resp.Value.AsEventGridDataConnection()
	if !ok {
		return nil, fmt.Errorf("%s is not an Event Grid Data Connection", id.String())
	}

	return utils.Bool(value.EventGridConnectionProperties != nil), nil
}

func (r KustoEventGridDataConnectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_eventgrid_data_connection" "test" {
  name                         = "acctestkrgdc-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  cluster_name                 = azurerm_kusto_cluster.test.name
  database_name                = azurerm_kusto_database.test.name
  storage_account_id           = azurerm_storage_account.test.id
  eventhub_id                  = azurerm_eventhub.test.id
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.test.name

  depends_on = [azurerm_eventgrid_event_subscription.test]
}
`, r.template(data), data.RandomInteger)
}

func (r KustoEventGridDataConnectionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_eventgrid_data_connection" "import" {
  name                         = azurerm_kusto_eventgrid_data_connection.test.name
  resource_group_name          = azurerm_kusto_eventgrid_data_connection.test.resource_group_name
  location                     = azurerm_kusto_eventgrid_data_connection.test.location
  cluster_name                 = azurerm_kusto_eventgrid_data_connection.test.cluster_name
  database_name                = azurerm_kusto_eventgrid_data_connection.test.database_name
  storage_account_id           = azurerm_kusto_eventgrid_data_connection.test.storage_account_id
  eventhub_id                  = azurerm_kusto_eventgrid_data_connection.test.eventhub_id
  eventhub_consumer_group_name = azurerm_kusto_eventgrid_data_connection.test.eventhub_consumer_group_name
}
`, r.basic(data))
}

func (r KustoEventGridDataConnectionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_eventgrid_data_connection" "test" {
  name                         = "acctestkrgdc-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  cluster_name                 = azurerm_kusto_cluster.test.name
  database_name                = azurerm_kusto_database.test.name
  storage_account_id           = azurerm_storage_account.test.id
  eventhub_id                  = azurerm_eventhub.test.id
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.test.name

  blob_storage_event_type = "Microsoft.Storage.BlobRenamed"
  skip_first_record       = true

  depends_on = [azurerm_eventgrid_event_subscription.test]
}
`, r.template(data), data.RandomInteger)
}

func (r KustoEventGridDataConnectionResource) mappingRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_eventgrid_data_connection" "test" {
  name                         = "acctestkrgdc-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  cluster_name                 = azurerm_kusto_cluster.test.name
  database_name                = azurerm_kusto_database.test.name
  storage_account_id           = azurerm_storage_account.test.id
  eventhub_id                  = azurerm_eventhub.test.id
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.test.name

  blob_storage_event_type = "Microsoft.Storage.BlobRenamed"
  skip_first_record       = true

  mapping_rule_name = "Json_Mapping"
  data_format       = "MULTIJSON"

  depends_on = [azurerm_eventgrid_event_subscription.test]
}
`, r.template(data), data.RandomInteger)
}

func (KustoEventGridDataConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "test" {
  name                = "acctestkd-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_name        = azurerm_kusto_cluster.test.name
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 1
  message_retention   = 1
}

resource "azurerm_eventhub_consumer_group" "test" {
  name                = "acctesteventhubcg-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_eventgrid_event_subscription" "test" {
  name                  = "acctest-eg-%d"
  scope                 = azurerm_storage_account.test.id
  eventhub_endpoint_id  = azurerm_eventhub.test.id
  event_delivery_schema = "EventGridSchema"
  included_event_types  = ["Microsoft.Storage.BlobCreated", "Microsoft.Storage.BlobRenamed"]

  retry_policy {
    event_time_to_live    = 144
    max_delivery_attempts = 10
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
