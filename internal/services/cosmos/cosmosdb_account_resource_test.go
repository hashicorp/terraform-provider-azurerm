// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CosmosDBAccountResource struct{}

func TestAccCosmosDBAccount_basic_global_boundedStaleness(t *testing.T) {
	testAccCosmosDBAccount_basicDocumentDbWith(t, documentdb.DefaultConsistencyLevelBoundedStaleness)
}

func TestAccCosmosDBAccount_basic_global_consistentPrefix(t *testing.T) {
	testAccCosmosDBAccount_basicDocumentDbWith(t, documentdb.DefaultConsistencyLevelConsistentPrefix)
}

func TestAccCosmosDBAccount_basic_global_eventual(t *testing.T) {
	testAccCosmosDBAccount_basicDocumentDbWith(t, documentdb.DefaultConsistencyLevelEventual)
}

func TestAccCosmosDBAccount_basic_global_session(t *testing.T) {
	testAccCosmosDBAccount_basicDocumentDbWith(t, documentdb.DefaultConsistencyLevelSession)
}

func TestAccCosmosDBAccount_basic_global_strong(t *testing.T) {
	testAccCosmosDBAccount_basicDocumentDbWith(t, documentdb.DefaultConsistencyLevelStrong)
}

func TestAccCosmosDBAccount_basic_mongo_boundedStaleness(t *testing.T) {
	testAccCosmosDBAccount_basicMongoDBWith(t, documentdb.DefaultConsistencyLevelBoundedStaleness)
}

func TestAccCosmosDBAccount_basic_mongo_consistentPrefix(t *testing.T) {
	testAccCosmosDBAccount_basicMongoDBWith(t, documentdb.DefaultConsistencyLevelConsistentPrefix)
}

func TestAccCosmosDBAccount_basic_mongo_eventual(t *testing.T) {
	testAccCosmosDBAccount_basicMongoDBWith(t, documentdb.DefaultConsistencyLevelEventual)
}

func TestAccCosmosDBAccount_basic_mongo_session(t *testing.T) {
	testAccCosmosDBAccount_basicMongoDBWith(t, documentdb.DefaultConsistencyLevelSession)
}

func TestAccCosmosDBAccount_basic_mongo_strong(t *testing.T) {
	testAccCosmosDBAccount_basicMongoDBWith(t, documentdb.DefaultConsistencyLevelStrong)
}

func TestAccCosmosDBAccount_basic_mongo_strong_without_capability(t *testing.T) {
	testAccCosmosDBAccount_basicMongoDBWith(t, documentdb.DefaultConsistencyLevelStrong)
}

func TestAccCosmosDBAccount_basic_parse_boundedStaleness(t *testing.T) {
	testAccCosmosDBAccount_basicWith(t, documentdb.DatabaseAccountKindParse, documentdb.DefaultConsistencyLevelBoundedStaleness)
}

func TestAccCosmosDBAccount_basic_parse_consistentPrefix(t *testing.T) {
	testAccCosmosDBAccount_basicWith(t, documentdb.DatabaseAccountKindParse, documentdb.DefaultConsistencyLevelConsistentPrefix)
}

func TestAccCosmosDBAccount_basic_parse_eventual(t *testing.T) {
	testAccCosmosDBAccount_basicWith(t, documentdb.DatabaseAccountKindParse, documentdb.DefaultConsistencyLevelEventual)
}

func TestAccCosmosDBAccount_basic_parse_session(t *testing.T) {
	testAccCosmosDBAccount_basicWith(t, documentdb.DatabaseAccountKindParse, documentdb.DefaultConsistencyLevelSession)
}

func TestAccCosmosDBAccount_basic_parse_strong(t *testing.T) {
	testAccCosmosDBAccount_basicWith(t, documentdb.DatabaseAccountKindParse, documentdb.DefaultConsistencyLevelStrong)
}

func TestAccCosmosDBAccount_public_network_access_enabled(t *testing.T) {
	testAccCosmosDBAccount_public_network_access_enabled(t, documentdb.DatabaseAccountKindMongoDB, documentdb.DefaultConsistencyLevelStrong)
}

func testAccCosmosDBAccount_public_network_access_enabled(t *testing.T, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.network_access_enabled(data, kind, consistency),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, consistency, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_keyVaultUri(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.key_vault_uri(data, documentdb.DatabaseAccountKindMongoDB, documentdb.DefaultConsistencyLevelStrong),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_customerManagedKeyWithIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVaultKeyUriWithSystemAssignedIdentity(data, documentdb.DatabaseAccountKindMongoDB, documentdb.DefaultConsistencyLevelStrong),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultKeyUriWithSystemAssignedAndUserAssignedIdentity(data, documentdb.DatabaseAccountKindMongoDB, documentdb.DefaultConsistencyLevelStrong),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultKeyUriWithUserAssignedIdentity(data, documentdb.DatabaseAccountKindMongoDB, documentdb.DefaultConsistencyLevelStrong),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultKeyUriWithSystemAssignedIdentity(data, documentdb.DatabaseAccountKindMongoDB, documentdb.DefaultConsistencyLevelStrong),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_keyVaultUriUpdateConsistancy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.key_vault_uri(data, documentdb.DatabaseAccountKindMongoDB, documentdb.DefaultConsistencyLevelStrong),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.key_vault_uri(data, documentdb.DatabaseAccountKindMongoDB, documentdb.DefaultConsistencyLevelSession),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelSession, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_updateTagsWithUserAssignedDefaultIdentity(t *testing.T) {
	// Regression test case for issue #22466
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.updateTagWithUserAssignedDefaultIdentity(data, "Production"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateTagWithUserAssignedDefaultIdentity(data, "Sandbox"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_updateDefaultIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.defaultIdentity(data, documentdb.DatabaseAccountKindGlobalDocumentDB, "FirstPartyIdentity", documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateDefaultIdentity(data, documentdb.DatabaseAccountKindGlobalDocumentDB, `"SystemAssignedIdentity"`, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateDefaultIdentityUserAssigned(data, documentdb.DatabaseAccountKindGlobalDocumentDB, `join("=", ["UserAssignedIdentity", azurerm_user_assigned_identity.test.id])`, documentdb.DefaultConsistencyLevelEventual, "UserAssigned"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateDefaultIdentityUserAssigned(data, documentdb.DatabaseAccountKindGlobalDocumentDB, `join("=", ["UserAssignedIdentity", azurerm_user_assigned_identity.test.id])`, documentdb.DefaultConsistencyLevelEventual, "SystemAssigned, UserAssigned"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_userAssignedIdentityMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleUserAssignedIdentityBaseState(data),
		},
		data.ImportStep(),
		{
			Config: r.multipleUserAssignedIdentity(data, "azurerm_user_assigned_identity.test.id"),
		},
		data.ImportStep(),
		{
			Config: r.multipleUserAssignedIdentity(data, "azurerm_user_assigned_identity.test2.id"),
		},
		data.ImportStep(),
		{
			Config: r.multipleUserAssignedIdentity(data, "azurerm_user_assigned_identity.test.id"),
		},
		data.ImportStep(),
		{
			Config: r.multipleUserAssignedIdentityBaseState(data),
		},
		data.ImportStep(),
	})
}

//nolint:unparam
func testAccCosmosDBAccount_basicWith(t *testing.T, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, kind, consistency),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, consistency, 1),
			),
		},
		data.ImportStep(),
	})
}

func testAccCosmosDBAccount_basicDocumentDbWith(t *testing.T, consistency documentdb.DefaultConsistencyLevel) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, documentdb.DatabaseAccountKindGlobalDocumentDB, consistency),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, consistency, 1),
				checkAccCosmosDBAccount_sql(data),
			),
		},
		data.ImportStep(),
	})
}

func testAccCosmosDBAccount_basicMongoDBWith(t *testing.T, consistency documentdb.DefaultConsistencyLevel) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMongoDB(data, consistency),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, consistency, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "GlobalDocumentDB", documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		{
			Config:      r.requiresImport(data, documentdb.DefaultConsistencyLevelEventual),
			ExpectError: acceptance.RequiresImportError("azurerm_cosmosdb_account"),
		},
	})
}

func TestAccCosmosDBAccount_updateConsistency_global(t *testing.T) {
	testAccCosmosDBAccount_updateConsistency(t, documentdb.DatabaseAccountKindGlobalDocumentDB)
}

func TestAccCosmosDBAccount_updateConsistency_mongo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMongoDB(data, documentdb.DefaultConsistencyLevelStrong),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelStrong, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistencyMongoDB(data, documentdb.DefaultConsistencyLevelStrong, 8, 880),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelStrong, 1),
		},
		data.ImportStep(),
		{
			Config: r.basicMongoDB(data, documentdb.DefaultConsistencyLevelBoundedStaleness),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelBoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistencyMongoDB(data, documentdb.DefaultConsistencyLevelBoundedStaleness, 7, 770),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelBoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistencyMongoDB(data, documentdb.DefaultConsistencyLevelBoundedStaleness, 77, 700),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelBoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.basicMongoDB(data, documentdb.DefaultConsistencyLevelConsistentPrefix),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelConsistentPrefix, 1),
		},
		data.ImportStep(),
	})
}

func testAccCosmosDBAccount_updateConsistency(t *testing.T, kind documentdb.DatabaseAccountKind) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, kind, documentdb.DefaultConsistencyLevelStrong),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelStrong, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistency(data, kind, documentdb.DefaultConsistencyLevelStrong, 8, 880),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelStrong, 1),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, kind, documentdb.DefaultConsistencyLevelBoundedStaleness),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelBoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistency(data, kind, documentdb.DefaultConsistencyLevelBoundedStaleness, 7, 770),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelBoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistency(data, kind, documentdb.DefaultConsistencyLevelBoundedStaleness, 77, 700),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelBoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, kind, documentdb.DefaultConsistencyLevelConsistentPrefix),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelConsistentPrefix, 1),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_complete_mongo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeMongoDB(data, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 3),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_complete_global(t *testing.T) {
	testAccCosmosDBAccount_completeWith(t, documentdb.DatabaseAccountKindGlobalDocumentDB)
}

func TestAccCosmosDBAccount_complete_parse(t *testing.T) {
	testAccCosmosDBAccount_completeWith(t, documentdb.DatabaseAccountKindParse)
}

func testAccCosmosDBAccount_completeWith(t *testing.T, kind documentdb.DatabaseAccountKind) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, kind, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 3),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_complete_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeTags(data, documentdb.DatabaseAccountKindParse, documentdb.DefaultConsistencyLevelEventual),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_completeZoneRedundant_mongo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	// Limited regional availability
	data.Locations.Primary = "westeurope"
	data.Locations.Secondary = "northeurope"
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zoneRedundantMongoDB(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_completeZoneRedundant_global(t *testing.T) {
	testAccCosmosDBAccount_zoneRedundantWith(t, documentdb.DatabaseAccountKindGlobalDocumentDB)
}

func TestAccCosmosDBAccount_completeZoneRedundant_parse(t *testing.T) {
	testAccCosmosDBAccount_zoneRedundantWith(t, documentdb.DatabaseAccountKindParse)
}

func testAccCosmosDBAccount_zoneRedundantWith(t *testing.T, kind documentdb.DatabaseAccountKind) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	// Limited regional availability
	data.Locations.Primary = "westeurope"
	data.Locations.Secondary = "northeurope"
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zoneRedundant(data, kind),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_zoneRedundant_update_mongo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	// Limited regional availability
	data.Locations.Primary = "westeurope"
	data.Locations.Secondary = "northeurope"
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMongoDB(data, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.zoneRedundantMongoDBUpdate(data, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 2),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_update_mongo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMongoDB(data, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeMongoDB(data, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 3),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdatedMongoDB(data, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 3),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdatedMongoDB_RemoveDisableRateLimitingResponsesCapability(data, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 3),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithResourcesMongoDB(data, documentdb.DefaultConsistencyLevelEventual),
			Check:  acceptance.ComposeAggregateTestCheckFunc(
			// checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_update_global(t *testing.T) {
	testAccCosmosDBAccount_updateWith(t, documentdb.DatabaseAccountKindGlobalDocumentDB)
}

func TestAccCosmosDBAccount_update_parse(t *testing.T) {
	testAccCosmosDBAccount_updateWith(t, documentdb.DatabaseAccountKindParse)
}

func testAccCosmosDBAccount_updateWith(t *testing.T, kind documentdb.DatabaseAccountKind) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, kind, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, kind, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 3),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdated(data, kind, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 3),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdated_RemoveDisableRateLimitingResponsesCapabilities(data, kind, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 3),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithResources(data, kind, documentdb.DefaultConsistencyLevelEventual),
			Check:  acceptance.ComposeAggregateTestCheckFunc(
			// checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_capabilities_EnableAggregationPipeline(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, documentdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableAggregationPipeline"})
}

func TestAccCosmosDBAccount_capabilities_EnableCassandra(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, documentdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableCassandra"})
}

func TestAccCosmosDBAccount_capabilities_EnableGremlin(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, documentdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableGremlin"})
}

func TestAccCosmosDBAccount_capabilities_EnableTable(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, documentdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableTable"})
}

func TestAccCosmosDBAccount_capabilities_EnableServerless(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, documentdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableServerless"})
}

func TestAccCosmosDBAccount_capabilities_EnableMongo(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, documentdb.DatabaseAccountKindMongoDB, []string{"EnableMongo"})
}

func TestAccCosmosDBAccount_capabilities_MongoDBv34(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, documentdb.DatabaseAccountKindMongoDB, []string{"EnableMongo", "MongoDBv3.4"})
}

func TestAccCosmosDBAccount_capabilities_MongoDBv34_NoEnableMongo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.capabilities(data, documentdb.DatabaseAccountKindMongoDB, []string{"MongoDBv3.4"}),
			ExpectError: regexp.MustCompile("capability EnableMongo must be enabled if MongoDBv3.4 is also enabled"),
		},
	})
}

func TestAccCosmosDBAccount_capabilities_mongoEnableDocLevelTTL(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, documentdb.DatabaseAccountKindMongoDB, []string{"EnableMongo", "mongoEnableDocLevelTTL"})
}

func TestAccCosmosDBAccount_capabilities_DisableRateLimitingResponses(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, documentdb.DatabaseAccountKindMongoDB, []string{"EnableMongo", "DisableRateLimitingResponses"})
}

func TestAccCosmosDBAccount_capabilities_AllowSelfServeUpgradeToMongo36(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, documentdb.DatabaseAccountKindMongoDB, []string{"EnableMongo", "AllowSelfServeUpgradeToMongo36"})
}

func testAccCosmosDBAccount_capabilitiesWith(t *testing.T, kind documentdb.DatabaseAccountKind, capabilities []string) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.capabilities(data, kind, capabilities),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_capabilitiesAdd(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.capabilities(data, documentdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableCassandra"}),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.capabilities(data, documentdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableCassandra", "EnableAggregationPipeline"}),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_capabilitiesUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.capabilities(data, documentdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableCassandra"}),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.capabilities(data, documentdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableCassandra", "DisableRateLimitingResponses", "AllowSelfServeUpgradeToMongo36", "EnableAggregationPipeline", "MongoDBv3.4", "mongoEnableDocLevelTTL"}),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.capabilities(data, documentdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableCassandra", "AllowSelfServeUpgradeToMongo36", "EnableAggregationPipeline", "MongoDBv3.4", "mongoEnableDocLevelTTL"}),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_geoLocationsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "GlobalDocumentDB", documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.geoLocationUpdate(data, "GlobalDocumentDB", documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 2),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "GlobalDocumentDB", documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_freeTier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.freeTier(data, "GlobalDocumentDB", documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
				check.That(data.ResourceName).Key("enable_free_tier").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_analyticalStorage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.analyticalStorage(data, "MongoDB", documentdb.DefaultConsistencyLevelStrong, false),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelStrong, 1),
				check.That(data.ResourceName).Key("analytical_storage_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.analyticalStorage(data, "MongoDB", documentdb.DefaultConsistencyLevelStrong, true),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelStrong, 1),
				check.That(data.ResourceName).Key("analytical_storage_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_updateAnalyticalStorage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateAnalyticalStorage(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.AnalyticalStorageSchemaTypeWellDefined, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateAnalyticalStorage(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.AnalyticalStorageSchemaTypeFullFidelity, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_updateCapacity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateCapacity(data, documentdb.DatabaseAccountKindGlobalDocumentDB, -1, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateCapacity(data, documentdb.DatabaseAccountKindGlobalDocumentDB, 200, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_vNetFilters(t *testing.T) {
	if !features.FourPointOhBeta() {
		t.Skip("this test requires 4.0 mode")
	}
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vNetFilters(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("is_virtual_network_filter_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("virtual_network_rule.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

// todo remove for 4.0
func TestAccCosmosDBAccount_vNetFiltersThreePointOh(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skip("this test requires 3.0 mode")
	}
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vNetFiltersThreePointOh(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("is_virtual_network_filter_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("virtual_network_rule.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMongoDB(data, documentdb.DefaultConsistencyLevelSession),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.systemAssignedUserAssignedIdentity(data, documentdb.DefaultConsistencyLevelSession),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned, UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicMongoDB(data, documentdb.DefaultConsistencyLevelSession),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_backup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.type").HasValue("Periodic"),
				check.That(data.ResourceName).Key("backup.0.interval_in_minutes").HasValue("240"),
				check.That(data.ResourceName).Key("backup.0.retention_in_hours").HasValue("8"),
				check.That(data.ResourceName).Key("backup.0.storage_redundancy").HasValue("Geo"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithBackupPeriodic(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithBackupPeriodicUpdate(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.type").HasValue("Periodic"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_backupPeriodicToContinuous(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithBackupPeriodic(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithBackupContinuous(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_backupContinuous(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithBackupContinuous(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_backupPeriodicToContinuousUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithBackupPeriodic(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithBackupContinuousUpdate(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_networkBypass(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithNetworkBypass(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithoutNetworkBypass(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_mongoVersion32(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMongoDBVersion32(data, documentdb.DefaultConsistencyLevelSession),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_mongoVersion36(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMongoDBVersion36(data, documentdb.DefaultConsistencyLevelSession),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_mongoVersion40(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMongoDBVersion40(data, documentdb.DefaultConsistencyLevelSession),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_mongoVersion42(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMongoDBVersion42(data, documentdb.DefaultConsistencyLevelSession),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_mongoVersionUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMongoDBVersion32(data, documentdb.DefaultConsistencyLevelSession),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicMongoDBVersion36(data, documentdb.DefaultConsistencyLevelSession),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicMongoDBVersion40(data, documentdb.DefaultConsistencyLevelSession),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicMongoDBVersion36(data, documentdb.DefaultConsistencyLevelSession),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_localAuthenticationDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("local_authentication_disabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithLocalAuthenticationDisabled(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("local_authentication_disabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_defaultCreateMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.defaultCreateMode(data, documentdb.DatabaseAccountKindGlobalDocumentDB, documentdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_restoreCreateMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.restoreCreateMode(data, documentdb.DatabaseAccountKindMongoDB, documentdb.DefaultConsistencyLevelSession),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.DefaultConsistencyLevelSession, 1),
			),
		},
		data.ImportStep(),
	})
}

// todo remove for 4.0
func TestAccCosmosDBAccount_ipRangeFiltersThreePointOh(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skip("this test requires 3.0 mode")
	}
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipRangeFiltersThreePointOh(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ipRangeFiltersUpdatedThreePointOh(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.vNetFiltersThreePointOh(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_ipRangeFilters(t *testing.T) {
	if !features.FourPointOhBeta() {
		t.Skip("this test requires 4.0 mode")
	}
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipRangeFilters(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ipRangeFiltersUpdated(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.vNetFilters(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_withoutMaxAgeInSeconds(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withoutMaxAgeInSeconds(data, documentdb.DatabaseAccountKindParse, documentdb.DefaultConsistencyLevelEventual),
		},
		data.ImportStep(),
	})
}

func (t CosmosDBAccountResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DatabaseAccountID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cosmos.DatabaseClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Cosmos Database (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (CosmosDBAccountResource) basic(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
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
  kind                = "%s"

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(consistency))
}

func (CosmosDBAccountResource) basicMongoDB(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
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
  kind                = "MongoDB"

  capabilities {
    name = "EnableMongo"
  }

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(consistency))
}

func (r CosmosDBAccountResource) requiresImport(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_account" "import" {
  name                = azurerm_cosmosdb_account.test.name
  location            = azurerm_cosmosdb_account.test.location
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  offer_type          = azurerm_cosmosdb_account.test.offer_type

  consistency_policy {
    consistency_level = azurerm_cosmosdb_account.test.consistency_policy[0].consistency_level
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, r.basic(data, "GlobalDocumentDB", consistency))
}

func (CosmosDBAccountResource) consistency(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel, interval, staleness int) string {
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
  kind                = "%s"

  consistency_policy {
    consistency_level       = "%s"
    max_interval_in_seconds = %d
    max_staleness_prefix    = %d
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(consistency), interval, staleness)
}

func (CosmosDBAccountResource) consistencyMongoDB(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel, interval, staleness int) string {
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
  kind                = "MongoDB"

  capabilities {
    name = "EnableMongo"
  }

  consistency_policy {
    consistency_level       = "%s"
    max_interval_in_seconds = %d
    max_staleness_prefix    = %d
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(consistency), interval, staleness)
}

func (CosmosDBAccountResource) completePreReqs(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  dns_servers         = ["10.0.0.4", "10.0.0.5"]
}

resource "azurerm_subnet" "subnet1" {
  name                 = "acctest-SN1-%[1]d-1"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
  service_endpoints    = ["Microsoft.AzureCosmosDB"]
}

resource "azurerm_subnet" "subnet2" {
  name                 = "acctest-SN2-%[1]d-2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  service_endpoints    = ["Microsoft.AzureCosmosDB"]
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CosmosDBAccountResource) complete(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%[3]s"

  consistency_policy {
    consistency_level       = "%[4]s"
    max_interval_in_seconds = 300
    max_staleness_prefix    = 170000
  }

  is_virtual_network_filter_enabled = true

  virtual_network_rule {
    id = azurerm_subnet.subnet1.id
  }

  virtual_network_rule {
    id = azurerm_subnet.subnet2.id
  }

  enable_multiple_write_locations = true

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  geo_location {
    location          = "%[5]s"
    failover_priority = 1
  }

  geo_location {
    location          = "%[6]s"
    failover_priority = 2
  }

  cors_rule {
    allowed_origins    = ["http://www.example.com"]
    exposed_headers    = ["x-tempo-*"]
    allowed_headers    = ["x-tempo-*"]
    allowed_methods    = ["GET", "PUT"]
    max_age_in_seconds = 500
  }

  access_key_metadata_writes_enabled    = false
  network_acl_bypass_for_azure_services = true
}
`, r.completePreReqs(data), data.RandomInteger, string(kind), string(consistency), data.Locations.Secondary, data.Locations.Ternary)
}

func (r CosmosDBAccountResource) completeTags(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%[3]s"

  consistency_policy {
    consistency_level       = "%[4]s"
    max_interval_in_seconds = 300
    max_staleness_prefix    = 170000
  }

  is_virtual_network_filter_enabled = true

  virtual_network_rule {
    id = azurerm_subnet.subnet1.id
  }

  virtual_network_rule {
    id = azurerm_subnet.subnet2.id
  }

  enable_multiple_write_locations = true

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  geo_location {
    location          = "%[5]s"
    failover_priority = 1
  }

  geo_location {
    location          = "%[6]s"
    failover_priority = 2
  }

  cors_rule {
    allowed_origins    = ["http://www.example.com"]
    exposed_headers    = ["x-tempo-*"]
    allowed_headers    = ["x-tempo-*"]
    allowed_methods    = ["GET", "PUT"]
    max_age_in_seconds = 500
  }
  access_key_metadata_writes_enabled    = false
  network_acl_bypass_for_azure_services = true

  tags = {
    ENV = "Test"
  }
}
`, r.completePreReqs(data), data.RandomInteger, string(kind), string(consistency), data.Locations.Secondary, data.Locations.Ternary)
}

func (r CosmosDBAccountResource) completeMongoDB(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "MongoDB"

  capabilities {
    name = "EnableMongo"
  }

  consistency_policy {
    consistency_level       = "%[3]s"
    max_interval_in_seconds = 300
    max_staleness_prefix    = 170000
  }

  is_virtual_network_filter_enabled = true

  virtual_network_rule {
    id = azurerm_subnet.subnet1.id
  }

  virtual_network_rule {
    id = azurerm_subnet.subnet2.id
  }

  enable_multiple_write_locations = true

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  geo_location {
    location          = "%[4]s"
    failover_priority = 1
  }

  geo_location {
    location          = "%[5]s"
    failover_priority = 2
  }

  cors_rule {
    allowed_origins    = ["http://www.example.com"]
    exposed_headers    = ["x-tempo-*"]
    allowed_headers    = ["x-tempo-*"]
    allowed_methods    = ["GET", "PUT"]
    max_age_in_seconds = 500
  }

  access_key_metadata_writes_enabled    = false
  network_acl_bypass_for_azure_services = true
}
`, r.completePreReqs(data), data.RandomInteger, string(consistency), data.Locations.Secondary, data.Locations.Ternary)
}

func (CosmosDBAccountResource) zoneRedundant(data acceptance.TestData, kind documentdb.DatabaseAccountKind) string {
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
  kind                = "%s"

  enable_multiple_write_locations = true

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 300
    max_staleness_prefix    = 100000
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  geo_location {
    location          = "%s"
    failover_priority = 1
    zone_redundant    = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), data.Locations.Secondary)
}

func (CosmosDBAccountResource) zoneRedundantMongoDB(data acceptance.TestData) string {
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
  kind                = "MongoDB"

  capabilities {
    name = "EnableMongo"
  }

  enable_multiple_write_locations = true

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 300
    max_staleness_prefix    = 100000
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  geo_location {
    location          = "%s"
    failover_priority = 1
    zone_redundant    = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Secondary)
}

func (r CosmosDBAccountResource) completeUpdated(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%[3]s"

  capabilities {
    name = "DisableRateLimitingResponses"
  }

  capabilities {
    name = "AllowSelfServeUpgradeToMongo36"
  }

  capabilities {
    name = "EnableAggregationPipeline"
  }

  capabilities {
    name = "mongoEnableDocLevelTTL"
  }

  consistency_policy {
    consistency_level       = "%[4]s"
    max_interval_in_seconds = 360
    max_staleness_prefix    = 170000
  }

  is_virtual_network_filter_enabled = true

  virtual_network_rule {
    id = azurerm_subnet.subnet2.id
  }

  enable_multiple_write_locations = true

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  geo_location {
    location          = "%[5]s"
    failover_priority = 1
  }

  geo_location {
    location          = "%[6]s"
    failover_priority = 2
  }

  cors_rule {
    allowed_origins    = ["http://www.example.com", "http://www.test.com"]
    exposed_headers    = ["x-tempo-*", "x-method-*"]
    allowed_headers    = ["*"]
    allowed_methods    = ["GET"]
    max_age_in_seconds = 2147483647
  }

  access_key_metadata_writes_enabled = true
}
`, r.completePreReqs(data), data.RandomInteger, string(kind), string(consistency), data.Locations.Secondary, data.Locations.Ternary)
}

func (r CosmosDBAccountResource) completeUpdated_RemoveDisableRateLimitingResponsesCapabilities(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%[3]s"

  capabilities {
    name = "AllowSelfServeUpgradeToMongo36"
  }

  capabilities {
    name = "EnableAggregationPipeline"
  }

  capabilities {
    name = "mongoEnableDocLevelTTL"
  }

  consistency_policy {
    consistency_level       = "%[4]s"
    max_interval_in_seconds = 360
    max_staleness_prefix    = 170000
  }

  is_virtual_network_filter_enabled = true

  virtual_network_rule {
    id = azurerm_subnet.subnet2.id
  }

  enable_multiple_write_locations = true

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  geo_location {
    location          = "%[5]s"
    failover_priority = 1
  }

  geo_location {
    location          = "%[6]s"
    failover_priority = 2
  }

  cors_rule {
    allowed_origins    = ["http://www.example.com", "http://www.test.com"]
    exposed_headers    = ["x-tempo-*", "x-method-*"]
    allowed_headers    = ["*"]
    allowed_methods    = ["GET"]
    max_age_in_seconds = 2147483647
  }

  access_key_metadata_writes_enabled = true
}
`, r.completePreReqs(data), data.RandomInteger, string(kind), string(consistency), data.Locations.Secondary, data.Locations.Ternary)
}

func (r CosmosDBAccountResource) completeUpdatedMongoDB(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "MongoDB"

  capabilities {
    name = "EnableMongo"
  }

  capabilities {
    name = "DisableRateLimitingResponses"
  }

  capabilities {
    name = "AllowSelfServeUpgradeToMongo36"
  }

  capabilities {
    name = "EnableAggregationPipeline"
  }

  capabilities {
    name = "MongoDBv3.4"
  }

  capabilities {
    name = "mongoEnableDocLevelTTL"
  }

  consistency_policy {
    consistency_level       = "%[3]s"
    max_interval_in_seconds = 360
    max_staleness_prefix    = 170000
  }

  is_virtual_network_filter_enabled = true

  virtual_network_rule {
    id = azurerm_subnet.subnet2.id
  }

  enable_multiple_write_locations = true

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  geo_location {
    location          = "%[4]s"
    failover_priority = 1
  }

  geo_location {
    location          = "%[5]s"
    failover_priority = 2
  }

  cors_rule {
    allowed_origins    = ["http://www.example.com", "http://www.test.com"]
    exposed_headers    = ["x-tempo-*", "x-method-*"]
    allowed_headers    = ["*"]
    allowed_methods    = ["GET"]
    max_age_in_seconds = 2147483647
  }
  access_key_metadata_writes_enabled = true
}
`, r.completePreReqs(data), data.RandomInteger, string(consistency), data.Locations.Secondary, data.Locations.Ternary)
}

func (r CosmosDBAccountResource) completeUpdatedMongoDB_RemoveDisableRateLimitingResponsesCapability(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "MongoDB"

  capabilities {
    name = "EnableMongo"
  }

  capabilities {
    name = "AllowSelfServeUpgradeToMongo36"
  }

  capabilities {
    name = "EnableAggregationPipeline"
  }

  capabilities {
    name = "MongoDBv3.4"
  }

  capabilities {
    name = "mongoEnableDocLevelTTL"
  }

  consistency_policy {
    consistency_level       = "%[3]s"
    max_interval_in_seconds = 360
    max_staleness_prefix    = 170000
  }

  is_virtual_network_filter_enabled = true

  virtual_network_rule {
    id = azurerm_subnet.subnet2.id
  }

  enable_multiple_write_locations = true

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  geo_location {
    location          = "%[4]s"
    failover_priority = 1
  }

  geo_location {
    location          = "%[5]s"
    failover_priority = 2
  }

  cors_rule {
    allowed_origins    = ["http://www.example.com", "http://www.test.com"]
    exposed_headers    = ["x-tempo-*", "x-method-*"]
    allowed_headers    = ["*"]
    allowed_methods    = ["GET"]
    max_age_in_seconds = 2147483647
  }
  access_key_metadata_writes_enabled = true
}
`, r.completePreReqs(data), data.RandomInteger, string(consistency), data.Locations.Secondary, data.Locations.Ternary)
}

func (r CosmosDBAccountResource) basicWithResources(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%s"

  capabilities {
    name = "AllowSelfServeUpgradeToMongo36"
  }

  capabilities {
    name = "EnableAggregationPipeline"
  }

  capabilities {
    name = "mongoEnableDocLevelTTL"
  }

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, r.completePreReqs(data), data.RandomInteger, string(kind), string(consistency))
}

func (r CosmosDBAccountResource) basicWithResourcesMongoDB(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "MongoDB"

  capabilities {
    name = "EnableMongo"
  }

  capabilities {
    name = "AllowSelfServeUpgradeToMongo36"
  }

  capabilities {
    name = "EnableAggregationPipeline"
  }

  capabilities {
    name = "MongoDBv3.4"
  }

  capabilities {
    name = "mongoEnableDocLevelTTL"
  }

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, r.completePreReqs(data), data.RandomInteger, string(consistency))
}

func (CosmosDBAccountResource) capabilities(data acceptance.TestData, kind documentdb.DatabaseAccountKind, capabilities []string) string {
	capeTf := ""
	for _, c := range capabilities {
		capeTf += fmt.Sprintf("capabilities {name = \"%s\"}\n", c)
	}

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
  kind                = "%s"

  consistency_policy {
    consistency_level = "Strong"
  }

  %s

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), capeTf)
}

func (CosmosDBAccountResource) geoLocationUpdate(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
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
  kind                = "%s"

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  geo_location {
    location          = "%s"
    failover_priority = 1
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(consistency), data.Locations.Secondary)
}

func (CosmosDBAccountResource) zoneRedundantMongoDBUpdate(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
variable "geo_location" {
  type = list(object({
    location          = string
    failover_priority = string
    zone_redundant    = bool
  }))
  default = [
    {
      location          = "%s"
      failover_priority = 0
      zone_redundant    = false
    },
    {
      location          = "%s"
      failover_priority = 1
      zone_redundant    = true
    }
  ]
}

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
  kind                = "MongoDB"

  capabilities {
    name = "EnableMongo"
  }

  enable_multiple_write_locations = true
  enable_automatic_failover       = true

  consistency_policy {
    consistency_level = "%s"
  }

  dynamic "geo_location" {
    for_each = var.geo_location
    content {
      location          = geo_location.value.location
      failover_priority = geo_location.value.failover_priority
      zone_redundant    = geo_location.value.zone_redundant
    }
  }
}
`, data.Locations.Primary, data.Locations.Secondary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(consistency))
}

func (CosmosDBAccountResource) vNetFiltersPreReqs(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  dns_servers         = ["10.0.0.4", "10.0.0.5"]
}

resource "azurerm_subnet" "subnet1" {
  name                                           = "acctest-SN1-%[1]d-1"
  resource_group_name                            = azurerm_resource_group.test.name
  virtual_network_name                           = azurerm_virtual_network.test.name
  address_prefixes                               = ["10.0.1.0/24"]
  enforce_private_link_endpoint_network_policies = false
  enforce_private_link_service_network_policies  = false
}

resource "azurerm_subnet" "subnet2" {
  name                                           = "acctest-SN2-%[1]d-2"
  resource_group_name                            = azurerm_resource_group.test.name
  virtual_network_name                           = azurerm_virtual_network.test.name
  address_prefixes                               = ["10.0.2.0/24"]
  service_endpoints                              = ["Microsoft.AzureCosmosDB"]
  enforce_private_link_endpoint_network_policies = false
  enforce_private_link_service_network_policies  = false
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CosmosDBAccountResource) vNetFilters(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  enable_multiple_write_locations = false
  enable_automatic_failover       = false

  consistency_policy {
    consistency_level       = "Eventual"
    max_interval_in_seconds = 5
    max_staleness_prefix    = 100
  }

  is_virtual_network_filter_enabled = true
  ip_range_filter                   = []

  virtual_network_rule {
    id                                   = azurerm_subnet.subnet1.id
    ignore_missing_vnet_service_endpoint = true
  }

  virtual_network_rule {
    id                                   = azurerm_subnet.subnet2.id
    ignore_missing_vnet_service_endpoint = false
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, r.vNetFiltersPreReqs(data), data.RandomInteger)
}

func (r CosmosDBAccountResource) vNetFiltersThreePointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  enable_multiple_write_locations = false
  enable_automatic_failover       = false

  consistency_policy {
    consistency_level       = "Eventual"
    max_interval_in_seconds = 5
    max_staleness_prefix    = 100
  }

  is_virtual_network_filter_enabled = true
  ip_range_filter                   = ""

  virtual_network_rule {
    id                                   = azurerm_subnet.subnet1.id
    ignore_missing_vnet_service_endpoint = true
  }

  virtual_network_rule {
    id                                   = azurerm_subnet.subnet2.id
    ignore_missing_vnet_service_endpoint = false
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, r.vNetFiltersPreReqs(data), data.RandomInteger)
}

func (CosmosDBAccountResource) freeTier(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
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
  kind                = "%s"

  enable_free_tier = true

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(consistency))
}

func (CosmosDBAccountResource) analyticalStorage(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel, enableAnalyticalStorage bool) string {
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
  kind                = "%s"

  analytical_storage_enabled = %t

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), enableAnalyticalStorage, string(consistency))
}

func (CosmosDBAccountResource) mongoAnalyticalStorage(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
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
  kind                = "MongoDB"

  analytical_storage_enabled = true

  consistency_policy {
    consistency_level = "%s"
  }

  capabilities {
    name = "EnableMongo"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(consistency))
}

func checkAccCosmosDBAccount_basic(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel, locationCount int) acceptance.TestCheckFunc {
	return acceptance.ComposeTestCheckFunc(
		check.That(data.ResourceName).Key("name").Exists(),
		check.That(data.ResourceName).Key("resource_group_name").Exists(),
		check.That(data.ResourceName).Key("location").HasValue(azure.NormalizeLocation(data.Locations.Primary)),
		check.That(data.ResourceName).Key("tags.%").HasValue("0"),
		check.That(data.ResourceName).Key("offer_type").HasValue(string(documentdb.DatabaseAccountOfferTypeStandard)),
		check.That(data.ResourceName).Key("consistency_policy.0.consistency_level").HasValue(string(consistency)),
		check.That(data.ResourceName).Key("geo_location.#").HasValue(strconv.Itoa(locationCount)),
		check.That(data.ResourceName).Key("endpoint").Exists(),
		check.That(data.ResourceName).Key("read_endpoints.#").HasValue(strconv.Itoa(locationCount)),
		check.That(data.ResourceName).Key("primary_key").Exists(),
		check.That(data.ResourceName).Key("secondary_key").Exists(),
		check.That(data.ResourceName).Key("primary_readonly_key").Exists(),
		check.That(data.ResourceName).Key("secondary_readonly_key").Exists(),
	)
}

func checkAccCosmosDBAccount_sql(data acceptance.TestData) acceptance.TestCheckFunc {
	return acceptance.ComposeTestCheckFunc(
		check.That(data.ResourceName).Key("primary_sql_connection_string").Exists(),
		check.That(data.ResourceName).Key("secondary_sql_connection_string").Exists(),
		check.That(data.ResourceName).Key("primary_readonly_sql_connection_string").Exists(),
		check.That(data.ResourceName).Key("secondary_readonly_sql_connection_string").Exists(),
	)
}

func (CosmosDBAccountResource) network_access_enabled(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                          = "acctest-ca-%d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  offer_type                    = "Standard"
  kind                          = "%s"
  public_network_access_enabled = true

  capabilities {
    name = "EnableMongo"
  }

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(consistency))
}

func (CosmosDBAccountResource) key_vault_uri(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

provider "azuread" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

data "azuread_service_principal" "cosmosdb" {
  display_name = "Azure Cosmos DB"
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  purge_protection_enabled   = true
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "List",
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Update",
      "GetRotationPolicy",
    ]

    secret_permissions = [
      "Get",
      "Delete",
      "Set",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azuread_service_principal.cosmosdb.id

    key_permissions = [
      "List",
      "Create",
      "Delete",
      "Get",
      "Update",
      "UnwrapKey",
      "WrapKey",
      "GetRotationPolicy",
    ]

    secret_permissions = [
      "Get",
      "Delete",
      "Set",
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%s"
  key_vault_key_id    = azurerm_key_vault_key.test.versionless_id

  capabilities {
    name = "EnableMongo"
  }

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger, string(kind), string(consistency))
}

func (CosmosDBAccountResource) keyVaultKeyUriWithSystemAssignedIdentity(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

provider "azuread" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

data "azurerm_client_config" "current" {}

data "azuread_service_principal" "cosmosdb" {
  display_name = "Azure Cosmos DB"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctest-user-example"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  purge_protection_enabled   = true
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "List",
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Update",
      "GetRotationPolicy",
    ]

    secret_permissions = [
      "Get",
      "Delete",
      "Set",
    ]
  }

  access_policy {
    tenant_id = azurerm_user_assigned_identity.test.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    key_permissions = [
      "List",
      "Create",
      "Delete",
      "Get",
      "Update",
      "UnwrapKey",
      "WrapKey",
      "GetRotationPolicy",
    ]

    secret_permissions = [
      "Get",
      "Delete",
      "Set",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azuread_service_principal.cosmosdb.id

    key_permissions = [
      "List",
      "Create",
      "Delete",
      "Get",
      "Update",
      "UnwrapKey",
      "WrapKey",
      "GetRotationPolicy",
    ]

    secret_permissions = [
      "Get",
      "Delete",
      "Set",
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%s"
  key_vault_key_id    = azurerm_key_vault_key.test.versionless_id

  capabilities {
    name = "EnableMongo"
  }

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger, string(kind), string(consistency))
}

func (CosmosDBAccountResource) keyVaultKeyUriWithSystemAssignedAndUserAssignedIdentity(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

provider "azuread" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

data "azurerm_client_config" "current" {}

data "azuread_service_principal" "cosmosdb" {
  display_name = "Azure Cosmos DB"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctest-user-example"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  purge_protection_enabled   = true
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "List",
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Update",
      "GetRotationPolicy",
    ]

    secret_permissions = [
      "Get",
      "Delete",
      "Set",
    ]
  }

  access_policy {
    tenant_id = azurerm_user_assigned_identity.test.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    key_permissions = [
      "List",
      "Create",
      "Delete",
      "Get",
      "Update",
      "UnwrapKey",
      "WrapKey",
      "GetRotationPolicy",
    ]

    secret_permissions = [
      "Get",
      "Delete",
      "Set",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azuread_service_principal.cosmosdb.id

    key_permissions = [
      "List",
      "Create",
      "Delete",
      "Get",
      "Update",
      "UnwrapKey",
      "WrapKey",
      "GetRotationPolicy",
    ]

    secret_permissions = [
      "Get",
      "Delete",
      "Set",
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%s"
  key_vault_key_id    = azurerm_key_vault_key.test.versionless_id

  capabilities {
    name = "EnableMongo"
  }

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger, string(kind), string(consistency))
}

func (CosmosDBAccountResource) keyVaultKeyUriWithUserAssignedIdentity(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

provider "azuread" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

data "azurerm_client_config" "current" {}

data "azuread_service_principal" "cosmosdb" {
  display_name = "Azure Cosmos DB"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctest-user-example"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  purge_protection_enabled   = true
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "List",
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Update",
      "GetRotationPolicy",
    ]

    secret_permissions = [
      "Get",
      "Delete",
      "Set",
    ]
  }

  access_policy {
    tenant_id = azurerm_user_assigned_identity.test.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    key_permissions = [
      "List",
      "Create",
      "Delete",
      "Get",
      "Update",
      "UnwrapKey",
      "WrapKey",
      "GetRotationPolicy",
    ]

    secret_permissions = [
      "Get",
      "Delete",
      "Set",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azuread_service_principal.cosmosdb.id

    key_permissions = [
      "List",
      "Create",
      "Delete",
      "Get",
      "Update",
      "UnwrapKey",
      "WrapKey",
      "GetRotationPolicy",
    ]

    secret_permissions = [
      "Get",
      "Delete",
      "Set",
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%s"
  key_vault_key_id    = azurerm_key_vault_key.test.versionless_id

  capabilities {
    name = "EnableMongo"
  }

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger, string(kind), string(consistency))
}

func (CosmosDBAccountResource) systemAssignedUserAssignedIdentity(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctest-user-example"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "MongoDB"

  capabilities {
    name = "EnableMongo"
  }

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(consistency))
}

func (CosmosDBAccountResource) multipleUserAssignedIdentity(data acceptance.TestData, identityResource string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctest-user-example"
}

resource "azurerm_user_assigned_identity" "test2" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctest-user-example-two"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "MongoDB"

  default_identity_type = join("=", ["UserAssignedIdentity", %[4]s])

  capabilities {
    name = "EnableMongo"
  }

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 5
    max_staleness_prefix    = 100
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [%[4]s]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, identityResource)
}

func (CosmosDBAccountResource) multipleUserAssignedIdentityBaseState(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctest-user-example"
}

resource "azurerm_user_assigned_identity" "test2" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctest-user-example-two"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "MongoDB"

  capabilities {
    name = "EnableMongo"
  }

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 5
    max_staleness_prefix    = 100
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CosmosDBAccountResource) basicWithBackupPeriodic(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
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
  kind                = "%s"

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  backup {
    type                = "Periodic"
    interval_in_minutes = 120
    retention_in_hours  = 10
    storage_redundancy  = "Geo"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(consistency))
}

func (CosmosDBAccountResource) basicWithBackupPeriodicUpdate(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
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
  kind                = "%s"

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  backup {
    type                = "Periodic"
    interval_in_minutes = 60
    retention_in_hours  = 8
    storage_redundancy  = "Local"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(consistency))
}

func (CosmosDBAccountResource) basicWithBackupContinuous(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
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
  kind                = "%s"

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  backup {
    type = "Continuous"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(consistency))
}

func (CosmosDBAccountResource) basicWithBackupContinuousUpdate(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
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
  kind                = "%s"

  is_virtual_network_filter_enabled = true

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  backup {
    type = "Continuous"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(consistency))
}

func (CosmosDBAccountResource) basicWithNetworkBypassTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%[1]d"
  storage_account_id = azurerm_storage_account.test.id
}

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%[1]d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r CosmosDBAccountResource) basicWithNetworkBypass(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%s"

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  network_acl_bypass_for_azure_services = true
  network_acl_bypass_ids                = [azurerm_synapse_workspace.test.id]
}
`, r.basicWithNetworkBypassTemplate(data), data.RandomInteger, string(kind), string(consistency))
}

func (r CosmosDBAccountResource) basicWithoutNetworkBypass(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%s"

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

}
`, r.basicWithNetworkBypassTemplate(data), data.RandomInteger, string(kind), string(consistency))
}

func (CosmosDBAccountResource) basicMongoDBVersion32(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                 = "acctest-ca-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  offer_type           = "Standard"
  kind                 = "MongoDB"
  mongo_server_version = "3.2"

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  lifecycle {
    ignore_changes = [
      capabilities
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(consistency))
}

func (CosmosDBAccountResource) basicMongoDBVersion36(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                 = "acctest-ca-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  offer_type           = "Standard"
  kind                 = "MongoDB"
  mongo_server_version = "3.6"

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  lifecycle {
    ignore_changes = [
      capabilities
    ]
  }

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(consistency))
}

func (CosmosDBAccountResource) basicMongoDBVersion40(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                 = "acctest-ca-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  offer_type           = "Standard"
  kind                 = "MongoDB"
  mongo_server_version = "4.0"

  capabilities {
    name = "EnableMongo"
  }

  capabilities {
    name = "EnableMongo16MBDocumentSupport"
  }

  capabilities {
    name = "EnableMongoRetryableWrites"
  }

  capabilities {
    name = "EnableMongoRoleBasedAccessControl"
  }

  capabilities {
    name = "EnableUniqueCompoundNestedDocs"
  }

  capabilities {
    name = "EnableTtlOnCustomPath"
  }

  capabilities {
    name = "EnablePartialUniqueIndex"
  }

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  lifecycle {
    ignore_changes = [
      capabilities
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(consistency))
}

func (CosmosDBAccountResource) basicMongoDBVersion42(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                 = "acctest-ca-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  offer_type           = "Standard"
  kind                 = "MongoDB"
  mongo_server_version = "4.2"

  capabilities {
    name = "EnableMongo"
  }

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  lifecycle {
    ignore_changes = [
      capabilities
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(consistency))
}

func (CosmosDBAccountResource) basicWithLocalAuthenticationDisabled(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
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
  kind                = "%s"

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  local_authentication_disabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(consistency))
}

func (CosmosDBAccountResource) updateAnalyticalStorage(data acceptance.TestData, kind documentdb.DatabaseAccountKind, schemaType documentdb.AnalyticalStorageSchemaType, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                       = "acctest-ca-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  offer_type                 = "Standard"
  kind                       = "%s"
  analytical_storage_enabled = false

  analytical_storage {
    schema_type = "%s"
  }

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(schemaType), string(consistency))
}

func (CosmosDBAccountResource) updateCapacity(data acceptance.TestData, kind documentdb.DatabaseAccountKind, totalThroughputLimit int, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                       = "acctest-ca-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  offer_type                 = "Standard"
  kind                       = "%s"
  analytical_storage_enabled = false

  capacity {
    total_throughput_limit = %d
  }

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), totalThroughputLimit, string(consistency))
}

func (CosmosDBAccountResource) defaultIdentity(data acceptance.TestData, kind documentdb.DatabaseAccountKind, defaultIdentity string, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                  = "acctest-ca-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  offer_type            = "Standard"
  kind                  = "%s"
  default_identity_type = "%s"

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), defaultIdentity, string(consistency))
}

func (CosmosDBAccountResource) updateDefaultIdentity(data acceptance.TestData, kind documentdb.DatabaseAccountKind, defaultIdentity string, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                  = "acctest-ca-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  offer_type            = "Standard"
  kind                  = "%s"
  default_identity_type = %s

  identity {
    type = "SystemAssigned"
  }

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), defaultIdentity, string(consistency))
}

func (CosmosDBAccountResource) updateDefaultIdentityUserAssigned(data acceptance.TestData, kind documentdb.DatabaseAccountKind, defaultIdentity string, consistency documentdb.DefaultConsistencyLevel, identityType string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctest-user-example"
}

resource "azurerm_cosmosdb_account" "test" {
  name                  = "acctest-ca-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  offer_type            = "Standard"
  kind                  = "%s"
  default_identity_type = %s

  identity {
    type         = "%s"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), defaultIdentity, identityType, string(consistency))
}

func (CosmosDBAccountResource) defaultCreateMode(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
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
  kind                = "%s"
  create_mode         = "Default"

  consistency_policy {
    consistency_level = "%s"
  }

  backup {
    type = "Continuous"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(consistency))
}

func (CosmosDBAccountResource) restoreCreateMode(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test1" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "MongoDB"

  capabilities {
    name = "EnableMongo"
  }

  consistency_policy {
    consistency_level = "Eventual"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  backup {
    type = "Continuous"
  }
}

resource "azurerm_cosmosdb_mongo_database" "test" {
  name                = "acctest-mongodb-%d"
  resource_group_name = azurerm_cosmosdb_account.test1.resource_group_name
  account_name        = azurerm_cosmosdb_account.test1.name
}

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-mongodb-coll-%d"
  resource_group_name = azurerm_cosmosdb_mongo_database.test.resource_group_name
  account_name        = azurerm_cosmosdb_mongo_database.test.account_name
  database_name       = azurerm_cosmosdb_mongo_database.test.name

  index {
    keys   = ["_id"]
    unique = true
  }
}

data "azurerm_cosmosdb_restorable_database_accounts" "test" {
  name     = azurerm_cosmosdb_account.test1.name
  location = azurerm_resource_group.test.location
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%s"

  capabilities {
    name = "EnableMongo"
  }

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  backup {
    type = "Continuous"
  }

  create_mode = "Restore"

  restore {
    source_cosmosdb_account_id = data.azurerm_cosmosdb_restorable_database_accounts.test.accounts[0].id
    restore_timestamp_in_utc   = timeadd(timestamp(), "-1s")

    database {
      name             = azurerm_cosmosdb_mongo_database.test.name
      collection_names = [azurerm_cosmosdb_mongo_collection.test.name]
    }
  }

  // As "restore_timestamp_in_utc" is retrieved dynamically, so it would cause diff when tf plan. So we have to ignore it here.
  lifecycle {
    ignore_changes = [
      restore.0.restore_timestamp_in_utc
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, string(kind), string(consistency))
}

func (r CosmosDBAccountResource) ipRangeFilters(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  enable_multiple_write_locations = false
  enable_automatic_failover       = false

  consistency_policy {
    consistency_level       = "Eventual"
    max_interval_in_seconds = 5
    max_staleness_prefix    = 100
  }

  is_virtual_network_filter_enabled = true
  ip_range_filter                   = ["55.0.1.0/24"]

  virtual_network_rule {
    id                                   = azurerm_subnet.subnet1.id
    ignore_missing_vnet_service_endpoint = true
  }

  virtual_network_rule {
    id                                   = azurerm_subnet.subnet2.id
    ignore_missing_vnet_service_endpoint = false
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, r.vNetFiltersPreReqs(data), data.RandomInteger)
}

func (r CosmosDBAccountResource) ipRangeFiltersUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  enable_multiple_write_locations = false
  enable_automatic_failover       = false

  consistency_policy {
    consistency_level       = "Eventual"
    max_interval_in_seconds = 5
    max_staleness_prefix    = 100
  }

  is_virtual_network_filter_enabled = true
  ip_range_filter                   = ["55.0.1.0/24", "55.0.2.0/24"]

  virtual_network_rule {
    id                                   = azurerm_subnet.subnet1.id
    ignore_missing_vnet_service_endpoint = true
  }

  virtual_network_rule {
    id                                   = azurerm_subnet.subnet2.id
    ignore_missing_vnet_service_endpoint = false
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, r.vNetFiltersPreReqs(data), data.RandomInteger)
}

func (r CosmosDBAccountResource) ipRangeFiltersThreePointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  enable_multiple_write_locations = false
  enable_automatic_failover       = false

  consistency_policy {
    consistency_level       = "Eventual"
    max_interval_in_seconds = 5
    max_staleness_prefix    = 100
  }

  is_virtual_network_filter_enabled = true
  ip_range_filter                   = "55.0.1.0/24"

  virtual_network_rule {
    id                                   = azurerm_subnet.subnet1.id
    ignore_missing_vnet_service_endpoint = true
  }

  virtual_network_rule {
    id                                   = azurerm_subnet.subnet2.id
    ignore_missing_vnet_service_endpoint = false
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, r.vNetFiltersPreReqs(data), data.RandomInteger)
}

func (r CosmosDBAccountResource) ipRangeFiltersUpdatedThreePointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  enable_multiple_write_locations = false
  enable_automatic_failover       = false

  consistency_policy {
    consistency_level       = "Eventual"
    max_interval_in_seconds = 5
    max_staleness_prefix    = 100
  }

  is_virtual_network_filter_enabled = true
  ip_range_filter                   = "55.0.1.0/24,55.0.2.0/24"

  virtual_network_rule {
    id                                   = azurerm_subnet.subnet1.id
    ignore_missing_vnet_service_endpoint = true
  }

  virtual_network_rule {
    id                                   = azurerm_subnet.subnet2.id
    ignore_missing_vnet_service_endpoint = false
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, r.vNetFiltersPreReqs(data), data.RandomInteger)
}

func (CosmosDBAccountResource) updateTagWithUserAssignedDefaultIdentity(data acceptance.TestData, tag string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%[1]d"
  location = "westeurope"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctest-user-example"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 5
    max_staleness_prefix    = 100
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  default_identity_type = join("=", ["UserAssignedIdentity", azurerm_user_assigned_identity.test.id])

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    environment  = "%[2]s",
    created_date = "2023-07-18"
  }
}
`, data.RandomInteger, tag)
}

func (r CosmosDBAccountResource) withoutMaxAgeInSeconds(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
%[1]s
resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%[3]s"
  consistency_policy {
    consistency_level       = "%[4]s"
    max_interval_in_seconds = 300
    max_staleness_prefix    = 170000
  }
  is_virtual_network_filter_enabled = true
  virtual_network_rule {
    id = azurerm_subnet.subnet1.id
  }
  virtual_network_rule {
    id = azurerm_subnet.subnet2.id
  }
  enable_multiple_write_locations = true
  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
  geo_location {
    location          = "%[5]s"
    failover_priority = 1
  }
  geo_location {
    location          = "%[6]s"
    failover_priority = 2
  }
  cors_rule {
    allowed_origins = ["http://www.example.com"]
    exposed_headers = ["x-tempo-*"]
    allowed_headers = ["x-tempo-*"]
    allowed_methods = ["GET", "PUT"]
  }
  access_key_metadata_writes_enabled    = false
  network_acl_bypass_for_azure_services = true
}
`, r.completePreReqs(data), data.RandomInteger, string(kind), string(consistency), data.Locations.Secondary, data.Locations.Ternary)
}
