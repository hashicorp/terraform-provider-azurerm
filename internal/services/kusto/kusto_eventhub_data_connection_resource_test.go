// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/dataconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KustoEventHubDataConnectionResource struct{}

func TestAccKustoEventHubDataConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_eventhub_data_connection", "test")
	r := KustoEventHubDataConnectionResource{}

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

func TestAccKustoEventHubDataConnection_eventSystemProperties(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_eventhub_data_connection", "test")
	r := KustoEventHubDataConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.eventSystemProperties(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoEventHubDataConnection_unboundMapping1(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_eventhub_data_connection", "test")
	r := KustoEventHubDataConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.unboundMapping1(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoEventHubDataConnection_unboundMapping2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_eventhub_data_connection", "test")
	r := KustoEventHubDataConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.unboundMapping2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoEventHubDataConnection_systemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_eventhub_data_connection", "test")
	r := KustoEventHubDataConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoEventHubDataConnection_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_eventhub_data_connection", "test")
	r := KustoEventHubDataConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoEventHubDataConnection_databaseRoutingType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_eventhub_data_connection", "test")
	r := KustoEventHubDataConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.databaseRoutingType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("database_routing_type").HasValue("Multi"),
			),
		},
		data.ImportStep(),
	})
}

func (KustoEventHubDataConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dataconnections.ParseDataConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Kusto.DataConnectionsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	if resp.Model != nil {
		value, ok := resp.Model.(dataconnections.EventHubDataConnection)
		if !ok {
			return nil, fmt.Errorf("%s is not an EventHubDataConnection", id.String())
		}
		exists := value.Properties != nil
		return &exists, nil
	} else {
		return nil, fmt.Errorf("response model is empty")
	}
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

func (r KustoEventHubDataConnectionResource) systemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_eventhub.test.id
  role_definition_name = "Azure Event Hubs Data Receiver"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_kusto_eventhub_data_connection" "test" {
  name                = "acctestkedc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_name        = azurerm_kusto_cluster.test.name
  database_name       = azurerm_kusto_database.test.name

  eventhub_id    = azurerm_eventhub.test.id
  consumer_group = azurerm_eventhub_consumer_group.test.name

  identity_id = azurerm_user_assigned_identity.test.id

  depends_on = [azurerm_role_assignment.test]
}
`, r.template(data), data.RandomInteger)
}

func (r KustoEventHubDataConnectionResource) userAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_eventhub.test.id
  role_definition_name = "Azure Event Hubs Data Receiver"
  principal_id         = azurerm_kusto_cluster.test.identity.0.principal_id
}

resource "azurerm_kusto_eventhub_data_connection" "test" {
  name                = "acctestkedc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_name        = azurerm_kusto_cluster.test.name
  database_name       = azurerm_kusto_database.test.name

  eventhub_id    = azurerm_eventhub.test.id
  consumer_group = azurerm_eventhub_consumer_group.test.name

  identity_id = azurerm_kusto_cluster.test.id

  depends_on = [azurerm_role_assignment.test]
}
`, r.template(data), data.RandomInteger)
}

func (r KustoEventHubDataConnectionResource) databaseRoutingType(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_kusto_eventhub_data_connection" "test" {
  name                  = "acctestkedc-%d"
  resource_group_name   = azurerm_resource_group.test.name
  location              = azurerm_resource_group.test.location
  cluster_name          = azurerm_kusto_cluster.test.name
  database_name         = azurerm_kusto_database.test.name
  eventhub_id           = azurerm_eventhub.test.id
  consumer_group        = azurerm_eventhub_consumer_group.test.name
  mapping_rule_name     = "Json_Mapping"
  data_format           = "MULTIJSON"
  database_routing_type = "Multi"
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

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
