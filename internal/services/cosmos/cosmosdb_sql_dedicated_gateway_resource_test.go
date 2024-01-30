// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-05-15/sqldedicatedgateway"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CosmosDbSqlDedicatedGatewayResource struct{}

func TestAccCosmosDbSqlDedicatedGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_dedicated_gateway", "test")
	r := CosmosDbSqlDedicatedGatewayResource{}

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

func TestAccCosmosDbSqlDedicatedGateway_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_dedicated_gateway", "test")
	r := CosmosDbSqlDedicatedGatewayResource{}

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

func TestAccCosmosDbSqlDedicatedGateway_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_dedicated_gateway", "test")
	r := CosmosDbSqlDedicatedGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CosmosDbSqlDedicatedGatewayResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := sqldedicatedgateway.ParseServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cosmos.SqlDedicatedGatewayClient
	resp, err := client.ServiceGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r CosmosDbSqlDedicatedGatewayResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level = "BoundedStaleness"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CosmosDbSqlDedicatedGatewayResource) basic(data acceptance.TestData) string {
	template := r.template(data)

	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_sql_dedicated_gateway" "test" {
  cosmosdb_account_id = azurerm_cosmosdb_account.test.id
  instance_size       = "Cosmos.D4s"
  instance_count      = 1
}
`, template)
}

func (r CosmosDbSqlDedicatedGatewayResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)

	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_sql_dedicated_gateway" "import" {
  cosmosdb_account_id = azurerm_cosmosdb_sql_dedicated_gateway.test.cosmosdb_account_id
  instance_count      = azurerm_cosmosdb_sql_dedicated_gateway.test.instance_count
  instance_size       = azurerm_cosmosdb_sql_dedicated_gateway.test.instance_size
}
`, config)
}

func (r CosmosDbSqlDedicatedGatewayResource) update(data acceptance.TestData) string {
	template := r.template(data)

	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_sql_dedicated_gateway" "test" {
  cosmosdb_account_id = azurerm_cosmosdb_account.test.id
  instance_size       = "Cosmos.D4s"
  instance_count      = 2
}
`, template)
}
