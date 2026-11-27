// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_container_test

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type CosmosDBAccountResource struct{}

func (CosmosDBAccountResource) analyticalStorage(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel, enableAnalyticalStorage bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%[3]s"

  analytical_storage_enabled = %[4]t

  consistency_policy {
    consistency_level = "%[5]s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, string(kind), enableAnalyticalStorage, string(consistency))
}

func (CosmosDBAccountResource) basic(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%[3]s"

  consistency_policy {
    consistency_level = "%[4]s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, string(kind), string(consistency))
}

func (CosmosDBAccountResource) capabilities(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, capabilities []string) string {
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
  name     = "acctestRG-cosmos-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%[3]s"

  consistency_policy {
    consistency_level = "Strong"
  }

  %[4]s

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, string(kind), capeTf)
}

type CosmosSqlDatabaseResource struct{}

func (CosmosSqlDatabaseResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}
`, CosmosDBAccountResource{}.basic(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelStrong), data.RandomInteger)
}

func (CosmosSqlDatabaseResource) serverless(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}
`, CosmosDBAccountResource{}.capabilities(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableServerless"}), data.RandomInteger)
}
