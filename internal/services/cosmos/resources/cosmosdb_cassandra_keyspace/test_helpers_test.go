// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_cassandra_keyspace_test

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type CosmosDBAccountResource struct{}

func (CosmosDBAccountResource) capabilities(data acceptance.TestData, capabilities []string) string {
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
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level = "Strong"
  }

  %[3]s

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, capeTf)
}
