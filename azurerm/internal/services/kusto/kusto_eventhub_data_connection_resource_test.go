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

type KustoEventHubDataConnectionResource struct {
}

func TestAccKustoEventHubDataConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_eventhub_data_connection", "test")
	r := KustoEventHubDataConnectionResource{}

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

func TestAccKustoEventHubDataConnection_eventSystemProperties(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_eventhub_data_connection", "test")
	r := KustoEventHubDataConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.eventSystemProperties(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoEventHubDataConnection_unboundMapping1(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_eventhub_data_connection", "test")
	r := KustoEventHubDataConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.unboundMapping1(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoEventHubDataConnection_unboundMapping2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_eventhub_data_connection", "test")
	r := KustoEventHubDataConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.unboundMapping2(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (KustoEventHubDataConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.DataConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Kusto.DataConnectionsClient.Get(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	value, ok := resp.Value.AsEventHubDataConnection()
	if !ok {
		return nil, fmt.Errorf("%s is not an EventHubDataConnection", id.String())
	}

	return utils.Bool(value.EventHubConnectionProperties != nil), nil
}

func (r KustoEventHubDataConnectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_eventhub_data_connection" "test" {
  name                = "acctestkedc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_name        = azurerm_kusto_cluster.test.name
  database_name       = azurerm_kusto_database.test.name

  eventhub_id    = azurerm_eventhub.test.id
  consumer_group = azurerm_eventhub_consumer_group.test.name
}
`, r.template(data), data.RandomInteger)
}

func (r KustoEventHubDataConnectionResource) unboundMapping1(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_eventhub_data_connection" "test" {
  name                = "acctestkedc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_name        = azurerm_kusto_cluster.test.name
  database_name       = azurerm_kusto_database.test.name

  eventhub_id    = azurerm_eventhub.test.id
  consumer_group = azurerm_eventhub_consumer_group.test.name

  mapping_rule_name = "Json_Mapping"
}
`, r.template(data), data.RandomInteger)
}

func (r KustoEventHubDataConnectionResource) unboundMapping2(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_eventhub_data_connection" "test" {
  name                = "acctestkedc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_name        = azurerm_kusto_cluster.test.name
  database_name       = azurerm_kusto_database.test.name

  eventhub_id    = azurerm_eventhub.test.id
  consumer_group = azurerm_eventhub_consumer_group.test.name

  mapping_rule_name = "Json_Mapping"
  data_format       = "MULTIJSON"
}
`, r.template(data), data.RandomInteger)
}

func (r KustoEventHubDataConnectionResource) eventSystemProperties(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_eventhub_data_connection" "test" {
  name                = "acctestkedc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_name        = azurerm_kusto_cluster.test.name
  database_name       = azurerm_kusto_database.test.name

  eventhub_id    = azurerm_eventhub.test.id
  consumer_group = azurerm_eventhub_consumer_group.test.name

  mapping_rule_name = "Json_Mapping"
  data_format       = "MULTIJSON"
  event_system_properties = [
    "x-opt-publisher"
  ]
}
`, r.template(data), data.RandomInteger)
}

func (KustoEventHubDataConnectionResource) template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
