// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_gremlin_graph_test

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type CosmosdbGremlinDatabaseResource struct{}

func (r CosmosdbGremlinDatabaseResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_gremlin_database" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}
`, r.template(data, []string{"EnableGremlin"}), data.RandomInteger)
}

func (r CosmosdbGremlinDatabaseResource) serverless(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_gremlin_database" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}
`, r.template(data, []string{"EnableGremlin", "EnableServerless"}), data.RandomInteger)
}

func (r CosmosdbGremlinDatabaseResource) template(data acceptance.TestData, capabilities []string) string {
	var capeTfBuilder strings.Builder
	for _, c := range capabilities {
		fmt.Fprintf(&capeTfBuilder, "capabilities {name = \"%s\"}\n", c)
	}
	capeTf := capeTfBuilder.String()

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level = "Strong"
  }

  %s

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, capeTf)
}
