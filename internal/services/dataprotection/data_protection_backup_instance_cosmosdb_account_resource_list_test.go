// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccDataProtectionBackupInstanceCosmosdbAccount_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_cosmosdb_account", "test")
	r := DataProtectionBackupInstanceCosmosdbAccountResource{}
	listResourceAddress := "azurerm_data_protection_backup_instance_cosmosdb_account.list"

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 2),
					querycheck.ExpectIdentity(
						listResourceAddress,
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"backup_vault_name":   knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r DataProtectionBackupInstanceCosmosdbAccountResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_account" "test2" {
  name                = "acctestcosmos2%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level = "Session"
  }

  backup {
    type = "Continuous"
    tier = "Continuous7Days"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

resource "azurerm_role_assignment" "cosmos_operator2" {
  scope                = azurerm_cosmosdb_account.test2.id
  role_definition_name = "Cosmos DB Operator"
  principal_id         = azurerm_data_protection_backup_vault.test.identity[0].principal_id
}

resource "azurerm_data_protection_backup_instance_cosmosdb_account" "test2" {
  name                              = "acctest-dbi-cosmos2-%d"
  location                          = azurerm_resource_group.test.location
  data_protection_backup_vault_id   = azurerm_data_protection_backup_vault.test.id
  backup_policy_cosmosdb_account_id = azurerm_data_protection_backup_policy_cosmosdb_account.test.id
  cosmosdb_account_id               = azurerm_cosmosdb_account.test2.id

  depends_on = [
    azurerm_role_assignment.reader,
    azurerm_role_assignment.cosmos_operator2,
  ]
}
`, r.basic(data), data.RandomInteger, data.RandomInteger)
}

func (DataProtectionBackupInstanceCosmosdbAccountResource) basicListQuery() string {
	return `
list "azurerm_data_protection_backup_instance_cosmosdb_account" "list" {
  provider = azurerm
  config {
    data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  }
}
`
}
