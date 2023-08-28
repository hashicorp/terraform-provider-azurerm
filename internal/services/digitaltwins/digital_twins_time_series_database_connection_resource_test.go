// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package digitaltwins_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/timeseriesdatabaseconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type TimeSeriesDatabaseConnectionResource struct{}

func TestAccTimeSeriesDatabaseConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_time_series_database_connection", "test")
	r := TimeSeriesDatabaseConnectionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("eventhub_consumer_group_name").HasValue("$Default"),
				check.That(data.ResourceName).Key("kusto_table_name").IsSet(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccTimeSeriesDatabaseConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_time_series_database_connection", "test")
	r := TimeSeriesDatabaseConnectionResource{}
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

func (r TimeSeriesDatabaseConnectionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := timeseriesdatabaseconnections.ParseTimeSeriesDatabaseConnectionID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DigitalTwins.TimeSeriesDatabaseConnectionsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving TimeSeriesDatabaseConnection %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r TimeSeriesDatabaseConnectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_digital_twins_time_series_database_connection" "test" {
  name                            = "connection-%[2]d"
  digital_twins_id                = azurerm_digital_twins_instance.test.id
  eventhub_name                   = azurerm_eventhub.test.name
  eventhub_namespace_id           = azurerm_eventhub_namespace.test.id
  eventhub_namespace_endpoint_uri = "sb://${azurerm_eventhub_namespace.test.name}.servicebus.windows.net"
  kusto_cluster_id                = azurerm_kusto_cluster.test.id
  kusto_cluster_uri               = azurerm_kusto_cluster.test.uri
  kusto_database_name             = azurerm_kusto_database.test.name

  depends_on = [
    azurerm_role_assignment.database_contributor,
    azurerm_role_assignment.eventhub_data_owner,
    azurerm_kusto_database_principal_assignment.test
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r TimeSeriesDatabaseConnectionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_eventhub_consumer_group" "test" {
  name                = "acctesteventhubcg-%[2]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_digital_twins_time_series_database_connection" "test" {
  name                            = "connection-%[2]d"
  digital_twins_id                = azurerm_digital_twins_instance.test.id
  eventhub_name                   = azurerm_eventhub.test.name
  eventhub_namespace_id           = azurerm_eventhub_namespace.test.id
  eventhub_namespace_endpoint_uri = "sb://${azurerm_eventhub_namespace.test.name}.servicebus.windows.net"
  kusto_cluster_id                = azurerm_kusto_cluster.test.id
  kusto_cluster_uri               = azurerm_kusto_cluster.test.uri
  kusto_database_name             = azurerm_kusto_database.test.name

  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.test.name
  kusto_table_name             = "mytable"

  depends_on = [
    azurerm_role_assignment.database_contributor,
    azurerm_role_assignment.eventhub_data_owner,
    azurerm_kusto_database_principal_assignment.test
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r TimeSeriesDatabaseConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-digitaltwin-%[2]d"
  location = "%[1]s"
}

resource "azurerm_digital_twins_instance" "test" {
  name                = "acctest-DT-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%[2]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 7
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "test" {
  name                = "acctestkd-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_name        = azurerm_kusto_cluster.test.name
}

resource "azurerm_role_assignment" "database_contributor" {
  scope                = azurerm_kusto_database.test.id
  principal_id         = azurerm_digital_twins_instance.test.identity.0.principal_id
  role_definition_name = "Contributor"
}

resource "azurerm_role_assignment" "eventhub_data_owner" {
  scope                = azurerm_eventhub.test.id
  principal_id         = azurerm_digital_twins_instance.test.identity.0.principal_id
  role_definition_name = "Azure Event Hubs Data Owner"
}

resource "azurerm_kusto_database_principal_assignment" "test" {
  name                = "acctestkdpa%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  cluster_name        = azurerm_kusto_cluster.test.name
  database_name       = azurerm_kusto_database.test.name

  tenant_id      = azurerm_digital_twins_instance.test.identity.0.tenant_id
  principal_id   = azurerm_digital_twins_instance.test.identity.0.principal_id
  principal_type = "App"
  role           = "Admin"
}

`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
