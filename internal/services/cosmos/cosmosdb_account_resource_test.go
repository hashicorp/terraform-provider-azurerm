// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CosmosDBAccountResource struct{}

func TestAccCosmosDBAccount_basic_global_boundedStaleness(t *testing.T) {
	testAccCosmosDBAccount_basicDocumentDbWith(t, cosmosdb.DefaultConsistencyLevelBoundedStaleness)
}

func TestAccCosmosDBAccount_basic_global_consistentPrefix(t *testing.T) {
	testAccCosmosDBAccount_basicDocumentDbWith(t, cosmosdb.DefaultConsistencyLevelConsistentPrefix)
}

func TestAccCosmosDBAccount_basic_global_eventual(t *testing.T) {
	testAccCosmosDBAccount_basicDocumentDbWith(t, cosmosdb.DefaultConsistencyLevelEventual)
}

func TestAccCosmosDBAccount_basic_global_session(t *testing.T) {
	testAccCosmosDBAccount_basicDocumentDbWith(t, cosmosdb.DefaultConsistencyLevelSession)
}

func TestAccCosmosDBAccount_basic_global_strong(t *testing.T) {
	testAccCosmosDBAccount_basicDocumentDbWith(t, cosmosdb.DefaultConsistencyLevelStrong)
}

func TestAccCosmosDBAccount_basic_mongo_boundedStaleness(t *testing.T) {
	testAccCosmosDBAccount_basicMongoDBWith(t, cosmosdb.DefaultConsistencyLevelBoundedStaleness)
}

func TestAccCosmosDBAccount_basic_mongo_consistentPrefix(t *testing.T) {
	testAccCosmosDBAccount_basicMongoDBWith(t, cosmosdb.DefaultConsistencyLevelConsistentPrefix)
}

func TestAccCosmosDBAccount_basic_mongo_eventual(t *testing.T) {
	testAccCosmosDBAccount_basicMongoDBWith(t, cosmosdb.DefaultConsistencyLevelEventual)
}

func TestAccCosmosDBAccount_basic_mongo_session(t *testing.T) {
	testAccCosmosDBAccount_basicMongoDBWith(t, cosmosdb.DefaultConsistencyLevelSession)
}

func TestAccCosmosDBAccount_basic_mongo_strong(t *testing.T) {
	testAccCosmosDBAccount_basicMongoDBWith(t, cosmosdb.DefaultConsistencyLevelStrong)
}

func TestAccCosmosDBAccount_basic_mongo_strong_without_capability(t *testing.T) {
	testAccCosmosDBAccount_basicMongoDBWith(t, cosmosdb.DefaultConsistencyLevelStrong)
}

func TestAccCosmosDBAccount_basic_parse_boundedStaleness(t *testing.T) {
	testAccCosmosDBAccount_basicWith(t, cosmosdb.DatabaseAccountKindParse, cosmosdb.DefaultConsistencyLevelBoundedStaleness)
}

func TestAccCosmosDBAccount_basic_parse_consistentPrefix(t *testing.T) {
	testAccCosmosDBAccount_basicWith(t, cosmosdb.DatabaseAccountKindParse, cosmosdb.DefaultConsistencyLevelConsistentPrefix)
}

func TestAccCosmosDBAccount_basic_parse_eventual(t *testing.T) {
	testAccCosmosDBAccount_basicWith(t, cosmosdb.DatabaseAccountKindParse, cosmosdb.DefaultConsistencyLevelEventual)
}

func TestAccCosmosDBAccount_basic_parse_session(t *testing.T) {
	testAccCosmosDBAccount_basicWith(t, cosmosdb.DatabaseAccountKindParse, cosmosdb.DefaultConsistencyLevelSession)
}

func TestAccCosmosDBAccount_basic_parse_strong(t *testing.T) {
	testAccCosmosDBAccount_basicWith(t, cosmosdb.DatabaseAccountKindParse, cosmosdb.DefaultConsistencyLevelStrong)
}

func TestAccCosmosDBAccount_public_network_access_enabled(t *testing.T) {
	testAccCosmosDBAccount_public_network_access_enabled(t, cosmosdb.DatabaseAccountKindMongoDB, cosmosdb.DefaultConsistencyLevelStrong)
}

func testAccCosmosDBAccount_public_network_access_enabled(t *testing.T, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) {
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
			Config: r.key_vault_uri(data, cosmosdb.DatabaseAccountKindMongoDB, cosmosdb.DefaultConsistencyLevelStrong),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_ManagedHSMUri(t *testing.T) {
	if os.Getenv("ARM_TEST_HSM_KEY") == "" {
		t.Skip("Skipping as ARM_TEST_HSM_KEY is not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedHSMKey(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
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
			Config: r.keyVaultKeyUriWithSystemAssignedIdentity(data, cosmosdb.DatabaseAccountKindMongoDB, cosmosdb.DefaultConsistencyLevelStrong),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultKeyUriWithSystemAssignedAndUserAssignedIdentity(data, cosmosdb.DatabaseAccountKindMongoDB, cosmosdb.DefaultConsistencyLevelStrong),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultKeyUriWithUserAssignedIdentity(data, cosmosdb.DatabaseAccountKindMongoDB, cosmosdb.DefaultConsistencyLevelStrong),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultKeyUriWithSystemAssignedIdentity(data, cosmosdb.DatabaseAccountKindMongoDB, cosmosdb.DefaultConsistencyLevelStrong),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_updateMongoDBVersionCapabilities(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMongoDB(data, cosmosdb.DefaultConsistencyLevelStrong),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateMongoDBVersionCapabilities(data, cosmosdb.DefaultConsistencyLevelStrong),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
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
			Config: r.key_vault_uri(data, cosmosdb.DatabaseAccountKindMongoDB, cosmosdb.DefaultConsistencyLevelStrong),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.key_vault_uri(data, cosmosdb.DatabaseAccountKindMongoDB, cosmosdb.DefaultConsistencyLevelSession),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelSession, 1),
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

func TestAccCosmosDBAccount_minimalTlsVersion(t *testing.T) {
	if features.FivePointOh() {
		t.Skipf("There is no more available values for `minimal_tls_version` to test.")
	}
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMinimalTlsVersion(data, cosmosdb.MinimalTlsVersionTls),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("minimal_tls_version").HasValue("Tls"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicMinimalTlsVersion(data, cosmosdb.MinimalTlsVersionTlsOneOne),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("minimal_tls_version").HasValue("Tls11"),
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
			Config: r.basic(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.defaultIdentity(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, "FirstPartyIdentity", cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateDefaultIdentity(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, `"SystemAssignedIdentity"`, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateDefaultIdentityUserAssigned(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, `join("=", ["UserAssignedIdentity", azurerm_user_assigned_identity.test.id])`, cosmosdb.DefaultConsistencyLevelEventual, "UserAssigned"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateDefaultIdentityUserAssigned(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, `join("=", ["UserAssignedIdentity", azurerm_user_assigned_identity.test.id])`, cosmosdb.DefaultConsistencyLevelEventual, "SystemAssigned, UserAssigned"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
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
func testAccCosmosDBAccount_basicWith(t *testing.T, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) {
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

func testAccCosmosDBAccount_basicDocumentDbWith(t *testing.T, consistency cosmosdb.DefaultConsistencyLevel) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, consistency),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, consistency, 1),
				checkAccCosmosDBAccount_sql(data),
			),
		},
		data.ImportStep(),
	})
}

func testAccCosmosDBAccount_basicMongoDBWith(t *testing.T, consistency cosmosdb.DefaultConsistencyLevel) {
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
			Config: r.basic(data, "GlobalDocumentDB", cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		{
			Config:      r.requiresImport(data, cosmosdb.DefaultConsistencyLevelEventual),
			ExpectError: acceptance.RequiresImportError("azurerm_cosmosdb_account"),
		},
	})
}

func TestAccCosmosDBAccount_updateConsistency_global(t *testing.T) {
	testAccCosmosDBAccount_updateConsistency(t, cosmosdb.DatabaseAccountKindGlobalDocumentDB)
}

func TestAccCosmosDBAccount_updateConsistency_mongo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMongoDB(data, cosmosdb.DefaultConsistencyLevelStrong),
			Check:  checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistencyMongoDB(data, cosmosdb.DefaultConsistencyLevelStrong, 8, 880),
			Check:  checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
		},
		data.ImportStep(),
		{
			Config: r.basicMongoDB(data, cosmosdb.DefaultConsistencyLevelBoundedStaleness),
			Check:  checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelBoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistencyMongoDB(data, cosmosdb.DefaultConsistencyLevelBoundedStaleness, 7, 770),
			Check:  checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelBoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistencyMongoDB(data, cosmosdb.DefaultConsistencyLevelBoundedStaleness, 77, 700),
			Check:  checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelBoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.basicMongoDB(data, cosmosdb.DefaultConsistencyLevelConsistentPrefix),
			Check:  checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelConsistentPrefix, 1),
		},
		data.ImportStep(),
	})
}

func testAccCosmosDBAccount_updateConsistency(t *testing.T, kind cosmosdb.DatabaseAccountKind) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, kind, cosmosdb.DefaultConsistencyLevelStrong),
			Check:  checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistency(data, kind, false, cosmosdb.DefaultConsistencyLevelStrong, 8, 880),
			Check:  checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, kind, cosmosdb.DefaultConsistencyLevelBoundedStaleness),
			Check:  checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelBoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistency(data, kind, true, cosmosdb.DefaultConsistencyLevelBoundedStaleness, 7, 770),
			Check:  checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelBoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistency(data, kind, false, cosmosdb.DefaultConsistencyLevelBoundedStaleness, 77, 700),
			Check:  checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelBoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, kind, cosmosdb.DefaultConsistencyLevelConsistentPrefix),
			Check:  checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelConsistentPrefix, 1),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_complete_mongo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeMongoDB(data, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 3),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_complete_global(t *testing.T) {
	testAccCosmosDBAccount_completeWith(t, cosmosdb.DatabaseAccountKindGlobalDocumentDB)
}

func TestAccCosmosDBAccount_complete_parse(t *testing.T) {
	testAccCosmosDBAccount_completeWith(t, cosmosdb.DatabaseAccountKindParse)
}

func testAccCosmosDBAccount_completeWith(t *testing.T, kind cosmosdb.DatabaseAccountKind) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, kind, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 3),
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
			Config: r.completeTags(data, cosmosdb.DatabaseAccountKindParse, cosmosdb.DefaultConsistencyLevelEventual),
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
	testAccCosmosDBAccount_zoneRedundantWith(t, cosmosdb.DatabaseAccountKindGlobalDocumentDB)
}

func TestAccCosmosDBAccount_completeZoneRedundant_parse(t *testing.T) {
	testAccCosmosDBAccount_zoneRedundantWith(t, cosmosdb.DatabaseAccountKindParse)
}

func testAccCosmosDBAccount_zoneRedundantWith(t *testing.T, kind cosmosdb.DatabaseAccountKind) {
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
			Config: r.basicMongoDB(data, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.zoneRedundantMongoDBUpdate(data, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 2),
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
			Config: r.basicMongoDB(data, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeMongoDB(data, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 3),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdatedMongoDB(data, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 3),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdatedMongoDB_RemoveDisableRateLimitingResponsesCapability(data, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 3),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithResourcesMongoDB(data, cosmosdb.DefaultConsistencyLevelEventual),
			Check:  acceptance.ComposeAggregateTestCheckFunc(
			// checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_update_global(t *testing.T) {
	testAccCosmosDBAccount_updateWith(t, cosmosdb.DatabaseAccountKindGlobalDocumentDB)
}

func TestAccCosmosDBAccount_update_parse(t *testing.T) {
	testAccCosmosDBAccount_updateWith(t, cosmosdb.DatabaseAccountKindParse)
}

func testAccCosmosDBAccount_updateWith(t *testing.T, kind cosmosdb.DatabaseAccountKind) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, kind, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, kind, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 3),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdated(data, kind, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 3),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdated_RemoveDisableRateLimitingResponsesCapabilities(data, kind, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 3),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithResources(data, kind, cosmosdb.DefaultConsistencyLevelEventual),
			Check:  acceptance.ComposeAggregateTestCheckFunc(
			// checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_capabilities_EnableAggregationPipeline(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, cosmosdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableAggregationPipeline"})
}

func TestAccCosmosDBAccount_capabilities_EnableCassandra(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, cosmosdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableCassandra"})
}

func TestAccCosmosDBAccount_capabilities_EnableGremlin(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, cosmosdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableGremlin"})
}

func TestAccCosmosDBAccount_capabilities_EnableTable(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, cosmosdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableTable"})
}

func TestAccCosmosDBAccount_capabilities_EnableServerless(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, cosmosdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableServerless"})
}

func TestAccCosmosDBAccount_capabilities_EnableNoSQLVectorSearch(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, cosmosdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableNoSQLVectorSearch"})
}

func TestAccCosmosDBAccount_capabilities_EnableNoSQLFullTextSearch(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, cosmosdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableNoSQLFullTextSearch"})
}

func TestAccCosmosDBAccount_capabilities_EnableMongo(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, cosmosdb.DatabaseAccountKindMongoDB, []string{"EnableMongo"})
}

func TestAccCosmosDBAccount_capabilities_MongoDBv34(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, cosmosdb.DatabaseAccountKindMongoDB, []string{"EnableMongo", "MongoDBv3.4"})
}

func TestAccCosmosDBAccount_capabilities_MongoDBv34_NoEnableMongo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.capabilities(data, cosmosdb.DatabaseAccountKindMongoDB, []string{"MongoDBv3.4"}),
			ExpectError: regexp.MustCompile("capability EnableMongo must be enabled if MongoDBv3.4 is also enabled"),
		},
	})
}

func TestAccCosmosDBAccount_capabilities_mongoEnableDocLevelTTL(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, cosmosdb.DatabaseAccountKindMongoDB, []string{"EnableMongo", "mongoEnableDocLevelTTL"})
}

func TestAccCosmosDBAccount_capabilities_DisableRateLimitingResponses(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, cosmosdb.DatabaseAccountKindMongoDB, []string{"EnableMongo", "DisableRateLimitingResponses"})
}

func TestAccCosmosDBAccount_capabilities_AllowSelfServeUpgradeToMongo36(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, cosmosdb.DatabaseAccountKindMongoDB, []string{"EnableMongo", "AllowSelfServeUpgradeToMongo36"})
}

func testAccCosmosDBAccount_capabilitiesWith(t *testing.T, kind cosmosdb.DatabaseAccountKind, capabilities []string) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.capabilities(data, kind, capabilities),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
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
			Config: r.capabilities(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableCassandra"}),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.capabilities(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableCassandra", "EnableAggregationPipeline", "DeleteAllItemsByPartitionKey"}),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
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
			Config: r.capabilities(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableCassandra"}),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.capabilities(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableCassandra", "DisableRateLimitingResponses", "AllowSelfServeUpgradeToMongo36", "EnableAggregationPipeline", "MongoDBv3.4", "mongoEnableDocLevelTTL"}),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.capabilities(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableCassandra", "AllowSelfServeUpgradeToMongo36", "EnableAggregationPipeline", "MongoDBv3.4", "mongoEnableDocLevelTTL"}),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
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
			Config: r.basic(data, "GlobalDocumentDB", cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.geoLocationUpdate(data, "GlobalDocumentDB", cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 2),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "GlobalDocumentDB", cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
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
			Config: r.freeTier(data, "GlobalDocumentDB", cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
				check.That(data.ResourceName).Key("free_tier_enabled").HasValue("true"),
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
			Config: r.analyticalStorage(data, "MongoDB", cosmosdb.DefaultConsistencyLevelStrong, false),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
				check.That(data.ResourceName).Key("analytical_storage_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.analyticalStorage(data, "MongoDB", cosmosdb.DefaultConsistencyLevelStrong, true),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
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
			Config: r.basic(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateAnalyticalStorage(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.AnalyticalStorageSchemaTypeWellDefined, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateAnalyticalStorage(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.AnalyticalStorageSchemaTypeFullFidelity, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
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
			Config: r.basic(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateCapacity(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, -1, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateCapacity(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, 200, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_vNetFilters(t *testing.T) {
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

func TestAccCosmosDBAccount_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMongoDB(data, cosmosdb.DefaultConsistencyLevelSession),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.systemAssignedUserAssignedIdentity(data, cosmosdb.DefaultConsistencyLevelSession),
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
			Config: r.basicMongoDB(data, cosmosdb.DefaultConsistencyLevelSession),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_storageRedundancyUndefined(t *testing.T) {
	// Regression test for MSFT IcM where the SDK is supplied a 'nil' pointer for the
	// 'storage_redundancy' field, the new transport layer would send an 'empty' string
	// instead of omitting the field from the PUT call which would result in the API
	// returning an error...
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageRedundancyUndefined(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.type").HasValue("Periodic"),
				check.That(data.ResourceName).Key("backup.0.interval_in_minutes").HasValue("120"),
				check.That(data.ResourceName).Key("backup.0.retention_in_hours").HasValue("10"),
				check.That(data.ResourceName).Key("backup.0.storage_redundancy").HasValue("Geo"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_backupOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
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
			Config: r.basicWithBackupPeriodic(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.type").HasValue("Periodic"),
				check.That(data.ResourceName).Key("backup.0.interval_in_minutes").HasValue("120"),
				check.That(data.ResourceName).Key("backup.0.retention_in_hours").HasValue("10"),
				check.That(data.ResourceName).Key("backup.0.storage_redundancy").HasValue("Geo"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithBackupPeriodicUpdate(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.type").HasValue("Periodic"),
				check.That(data.ResourceName).Key("backup.0.interval_in_minutes").HasValue("60"),
				check.That(data.ResourceName).Key("backup.0.retention_in_hours").HasValue("8"),
				check.That(data.ResourceName).Key("backup.0.storage_redundancy").HasValue("Local"),
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
			Config: r.basicWithBackupPeriodic(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithBackupContinuous(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual, cosmosdb.ContinuousTierContinuousSevenDays),
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
			Config: r.basicWithBackupContinuous(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual, cosmosdb.ContinuousTierContinuousSevenDays),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithBackupContinuous(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual, cosmosdb.ContinuousTierContinuousThreeZeroDays),
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
			Config: r.basicWithBackupPeriodic(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithBackupContinuousUpdate(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
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
			Config: r.basic(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithNetworkBypass(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithoutNetworkBypass(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
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
			Config: r.basicMongoDBVersion32(data, cosmosdb.DefaultConsistencyLevelSession),
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
			Config: r.basicMongoDBVersion36(data, cosmosdb.DefaultConsistencyLevelSession),
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
			Config: r.basicMongoDBVersion40(data, cosmosdb.DefaultConsistencyLevelSession),
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
			Config: r.basicMongoDBVersion(data, cosmosdb.DefaultConsistencyLevelSession, "4.2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_mongoVersion50(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMongoDBVersion(data, cosmosdb.DefaultConsistencyLevelStrong, "5.0"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_mongoVersion60(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMongoDBVersion(data, cosmosdb.DefaultConsistencyLevelSession, "6.0"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_mongoVersion70(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMongoDBVersion(data, cosmosdb.DefaultConsistencyLevelSession, "7.0"),
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
			Config: r.basicMongoDBVersion32(data, cosmosdb.DefaultConsistencyLevelSession),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicMongoDBVersion36(data, cosmosdb.DefaultConsistencyLevelSession),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicMongoDBVersion40(data, cosmosdb.DefaultConsistencyLevelSession),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicMongoDBVersion36(data, cosmosdb.DefaultConsistencyLevelSession),
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
			Config: r.basic(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("local_authentication_disabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithLocalAuthenticationDisabled(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("local_authentication_disabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_updateBurstCapacity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("burst_capacity_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithBurstCapacityEnabled(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("burst_capacity_enabled").HasValue("true"),
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
			Config: r.defaultCreateMode(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelEventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelEventual, 1),
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
			Config: r.restoreCreateMode(data, cosmosdb.DatabaseAccountKindMongoDB, cosmosdb.DefaultConsistencyLevelSession),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelSession, 1),
				check.That(data.ResourceName).Key("minimal_tls_version").HasValue("Tls12"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_tablesToRestore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tablesToRestore(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelStrong),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_gremlinDatabasesToRestore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gremlinDatabasesToRestore(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelStrong),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, cosmosdb.DefaultConsistencyLevelStrong, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_ipRangeFilters(t *testing.T) {
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
			Config: r.withoutMaxAgeInSeconds(data, cosmosdb.DatabaseAccountKindParse, cosmosdb.DefaultConsistencyLevelEventual),
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

	return pointer.To(resp.ID != nil), nil
}

func (CosmosDBAccountResource) basic(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (CosmosDBAccountResource) basicMinimalTlsVersion(data acceptance.TestData, tls cosmosdb.MinimalTlsVersion) string {
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
  minimal_tls_version = "%s"

  consistency_policy {
    consistency_level = "Eventual"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(tls))
}

func (CosmosDBAccountResource) basicMongoDB(data acceptance.TestData, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (r CosmosDBAccountResource) requiresImport(data acceptance.TestData, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (CosmosDBAccountResource) consistency(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, partitionMergeEnabled bool, consistency cosmosdb.DefaultConsistencyLevel, interval, staleness int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                    = "acctest-ca-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  offer_type              = "Standard"
  kind                    = "%s"
  partition_merge_enabled = %t

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), partitionMergeEnabled, string(consistency), interval, staleness)
}

func (CosmosDBAccountResource) consistencyMongoDB(data acceptance.TestData, consistency cosmosdb.DefaultConsistencyLevel, interval, staleness int) string {
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

func (r CosmosDBAccountResource) complete(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

  multiple_write_locations_enabled = true

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

func (r CosmosDBAccountResource) completeTags(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

  multiple_write_locations_enabled = true

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

func (r CosmosDBAccountResource) completeMongoDB(data acceptance.TestData, consistency cosmosdb.DefaultConsistencyLevel) string {
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

  multiple_write_locations_enabled = true

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

func (CosmosDBAccountResource) zoneRedundant(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind) string {
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

  multiple_write_locations_enabled = true

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

  multiple_write_locations_enabled = true

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

func (r CosmosDBAccountResource) completeUpdated(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

  multiple_write_locations_enabled = true

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

func (r CosmosDBAccountResource) completeUpdated_RemoveDisableRateLimitingResponsesCapabilities(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

  multiple_write_locations_enabled = true

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

func (r CosmosDBAccountResource) completeUpdatedMongoDB(data acceptance.TestData, consistency cosmosdb.DefaultConsistencyLevel) string {
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

  multiple_write_locations_enabled = true

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

func (r CosmosDBAccountResource) completeUpdatedMongoDB_RemoveDisableRateLimitingResponsesCapability(data acceptance.TestData, consistency cosmosdb.DefaultConsistencyLevel) string {
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

  multiple_write_locations_enabled = true

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

func (r CosmosDBAccountResource) basicWithResources(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (r CosmosDBAccountResource) basicWithResourcesMongoDB(data acceptance.TestData, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (CosmosDBAccountResource) capabilities(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, capabilities []string) string {
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

func (CosmosDBAccountResource) geoLocationUpdate(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (CosmosDBAccountResource) zoneRedundantMongoDBUpdate(data acceptance.TestData, consistency cosmosdb.DefaultConsistencyLevel) string {
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

  multiple_write_locations_enabled = true
  automatic_failover_enabled       = true

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
  name                                          = "acctest-SN1-%[1]d-1"
  resource_group_name                           = azurerm_resource_group.test.name
  virtual_network_name                          = azurerm_virtual_network.test.name
  address_prefixes                              = ["10.0.1.0/24"]
  private_endpoint_network_policies             = "Disabled"
  private_link_service_network_policies_enabled = false
}

resource "azurerm_subnet" "subnet2" {
  name                                          = "acctest-SN2-%[1]d-2"
  resource_group_name                           = azurerm_resource_group.test.name
  virtual_network_name                          = azurerm_virtual_network.test.name
  address_prefixes                              = ["10.0.2.0/24"]
  service_endpoints                             = ["Microsoft.AzureCosmosDB"]
  private_endpoint_network_policies             = "Disabled"
  private_link_service_network_policies_enabled = false
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

  multiple_write_locations_enabled = false
  automatic_failover_enabled       = false

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

func (CosmosDBAccountResource) freeTier(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

  free_tier_enabled = true

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

func (CosmosDBAccountResource) analyticalStorage(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel, enableAnalyticalStorage bool) string {
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

func (CosmosDBAccountResource) mongoAnalyticalStorage(data acceptance.TestData, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func checkAccCosmosDBAccount_basic(data acceptance.TestData, consistency cosmosdb.DefaultConsistencyLevel, locationCount int) acceptance.TestCheckFunc {
	return acceptance.ComposeTestCheckFunc(
		check.That(data.ResourceName).Key("name").Exists(),
		check.That(data.ResourceName).Key("resource_group_name").Exists(),
		check.That(data.ResourceName).Key("location").HasValue(azure.NormalizeLocation(data.Locations.Primary)),
		check.That(data.ResourceName).Key("tags.%").HasValue("0"),
		check.That(data.ResourceName).Key("offer_type").HasValue(string(cosmosdb.DatabaseAccountOfferTypeStandard)),
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

func (CosmosDBAccountResource) network_access_enabled(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (CosmosDBAccountResource) key_vault_uri(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (CosmosDBAccountResource) keyVaultKeyUriWithSystemAssignedIdentity(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (CosmosDBAccountResource) keyVaultKeyUriWithSystemAssignedAndUserAssignedIdentity(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (CosmosDBAccountResource) keyVaultKeyUriWithUserAssignedIdentity(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (CosmosDBAccountResource) managedHSMKey(data acceptance.TestData) string {
	// Purge Protection must be enabled to configure Managed HSM Key: https://learn.microsoft.com/en-us/azure/cosmos-db/how-to-setup-customer-managed-keys-mhsm#configure-your-azure-managed-hsm-key-vault
	// hsmTemplate := customermanagedkeys.ManagedHSMKeyTempalte(data.RandomInteger, data.RandomString, enablePurgeProtection, []string{"data.azuread_service_principal.cosmosdb.id"})
	raName1, _ := uuid.GenerateUUID()
	raName2, _ := uuid.GenerateUUID()
	raName3, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctest-user-example"
}

data "azuread_service_principal" "cosmosdb" {
  display_name = "Azure Cosmos DB"
}

resource "azurerm_key_vault" "test" {
  name                       = "acc%[1]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id
    certificate_permissions = [
      "Create",
      "Delete",
      "DeleteIssuers",
      "Get",
      "Purge",
      "Update"
    ]
  }
  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_certificate" "cert" {
  count        = 3
  name         = "acchsmcert${count.index}"
  key_vault_id = azurerm_key_vault.test.id
  certificate_policy {
    issuer_parameters {
      name = "Self"
    }
    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }
    lifetime_action {
      action {
        action_type = "AutoRenew"
      }
      trigger {
        days_before_expiry = 30
      }
    }
    secret_properties {
      content_type = "application/x-pkcs12"
    }
    x509_certificate_properties {
      extended_key_usage = []
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]
      subject            = "CN=hello-world"
      validity_in_months = 12
    }
  }
}

resource "azurerm_key_vault_managed_hardware_security_module" "test" {
  name                       = "kvHsm%[1]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  sku_name                   = "Standard_B1"
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  admin_object_ids           = [data.azurerm_client_config.current.object_id]
  purge_protection_enabled   = true
  soft_delete_retention_days = 7

  security_domain_key_vault_certificate_ids = [for cert in azurerm_key_vault_certificate.cert : cert.id]
  security_domain_quorum                    = 3
}

data "azurerm_key_vault_managed_hardware_security_module_role_definition" "crypto-officer" {
  name           = "515eb02d-2335-4d2d-92f2-b1cbdf9c3778"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
}

data "azurerm_key_vault_managed_hardware_security_module_role_definition" "crypto-user" {
  name           = "21dbd100-6940-42c2-9190-5d6cb909625b"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
}

data "azurerm_key_vault_managed_hardware_security_module_role_definition" "encrypt-user" {
  name           = "33413926-3206-4cdd-b39a-83574fe37a17"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "client1" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "%[3]s"
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.crypto-officer.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "client2" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "%[4]s"
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.crypto-user.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "racosmos" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "%[5]s"
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.encrypt-user.resource_manager_id
  principal_id       = data.azuread_service_principal.cosmosdb.object_id

  depends_on = [azurerm_key_vault_managed_hardware_security_module_key.test]
}

resource "azurerm_key_vault_managed_hardware_security_module_key" "test" {
  name           = "acctestHSMK-%[1]d"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  key_type       = "RSA-HSM"
  key_size       = 2048
  key_opts       = ["unwrapKey", "wrapKey"]

  depends_on = [
    azurerm_key_vault_managed_hardware_security_module_role_assignment.client1,
    azurerm_key_vault_managed_hardware_security_module_role_assignment.client2
  ]
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acc-ca-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "MongoDB"
  managed_hsm_key_id  = azurerm_key_vault_managed_hardware_security_module_key.test.id

  capabilities {
    name = "EnableMongo"
  }

  consistency_policy {
    consistency_level = "Strong"
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

  // depends_on = [azurerm_key_vault_managed_hardware_security_module_role_assignment.racosmos]
}
`, data.RandomInteger, data.Locations.Primary, raName1, raName2, raName3)
}

func (CosmosDBAccountResource) systemAssignedUserAssignedIdentity(data acceptance.TestData, consistency cosmosdb.DefaultConsistencyLevel) string {
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

  depends_on = [azurerm_user_assigned_identity.test, azurerm_user_assigned_identity.test2]
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

  depends_on = [azurerm_user_assigned_identity.test, azurerm_user_assigned_identity.test2]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CosmosDBAccountResource) basicWithBackupPeriodic(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (CosmosDBAccountResource) storageRedundancyUndefined(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(consistency))
}

func (CosmosDBAccountResource) basicWithBackupPeriodicUpdate(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (CosmosDBAccountResource) basicWithBackupContinuous(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel, tier cosmosdb.ContinuousTier) string {
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
    tier = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(consistency), string(tier))
}

func (CosmosDBAccountResource) basicWithBackupContinuousUpdate(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (r CosmosDBAccountResource) basicWithNetworkBypass(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (r CosmosDBAccountResource) basicWithoutNetworkBypass(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (CosmosDBAccountResource) basicMongoDBVersion32(data acceptance.TestData, consistency cosmosdb.DefaultConsistencyLevel) string {
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
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(consistency))
}

func (CosmosDBAccountResource) updateMongoDBVersionCapabilities(data acceptance.TestData, consistency cosmosdb.DefaultConsistencyLevel) string {
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

  capabilities {
    name = "EnableMongo16MBDocumentSupport"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(consistency))
}

func (CosmosDBAccountResource) basicMongoDBVersion36(data acceptance.TestData, consistency cosmosdb.DefaultConsistencyLevel) string {
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
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(consistency))
}

func (CosmosDBAccountResource) basicMongoDBVersion40(data acceptance.TestData, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (CosmosDBAccountResource) basicMongoDBVersion(data acceptance.TestData, consistency cosmosdb.DefaultConsistencyLevel, version string) string {
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
  mongo_server_version = "%s"

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version, string(consistency))
}

func (CosmosDBAccountResource) basicWithLocalAuthenticationDisabled(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (CosmosDBAccountResource) basicWithBurstCapacityEnabled(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

  burst_capacity_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(consistency))
}

func (CosmosDBAccountResource) updateAnalyticalStorage(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, schemaType cosmosdb.AnalyticalStorageSchemaType, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (CosmosDBAccountResource) updateCapacity(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, totalThroughputLimit int, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (CosmosDBAccountResource) defaultIdentity(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, defaultIdentity string, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (CosmosDBAccountResource) updateDefaultIdentity(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, defaultIdentity string, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (CosmosDBAccountResource) updateDefaultIdentityUserAssigned(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, defaultIdentity string, consistency cosmosdb.DefaultConsistencyLevel, identityType string) string {
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

func (CosmosDBAccountResource) defaultCreateMode(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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

func (CosmosDBAccountResource) restoreCreateMode(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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
  minimal_tls_version = "Tls12"

  capabilities {
    name = "EnableMongo"
  }

  consistency_policy {
    consistency_level = "Session"
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

  // indices can cause test to be inconsistent
  // I believe there is a bug within the azurerm_cosmosdb_mongo_collection that causes inconsistent results on read
  lifecycle {
    ignore_changes = [
      index
    ]
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

func (CosmosDBAccountResource) tablesToRestore(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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
  kind                = "GlobalDocumentDB"

  capabilities {
    name = "EnableTable"
  }

  consistency_policy {
    consistency_level = "Strong"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  backup {
    type = "Continuous"
  }
}

resource "azurerm_cosmosdb_table" "test" {
  name                = "acctest-sqltable-%d"
  resource_group_name = azurerm_cosmosdb_account.test1.resource_group_name
  account_name        = azurerm_cosmosdb_account.test1.name
}

resource "azurerm_cosmosdb_table" "test2" {
  name                = "acctest-sqltable2-%d"
  resource_group_name = azurerm_cosmosdb_account.test1.resource_group_name
  account_name        = azurerm_cosmosdb_account.test1.name

  depends_on = [azurerm_cosmosdb_table.test]
}

data "azurerm_cosmosdb_restorable_database_accounts" "test" {
  name     = azurerm_cosmosdb_account.test1.name
  location = azurerm_resource_group.test.location

  depends_on = [azurerm_cosmosdb_table.test2]
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%s"

  capabilities {
    name = "EnableTable"
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
    tables_to_restore          = [azurerm_cosmosdb_table.test.name, azurerm_cosmosdb_table.test2.name]
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

func (CosmosDBAccountResource) gremlinDatabasesToRestore(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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
  kind                = "GlobalDocumentDB"

  capabilities {
    name = "EnableGremlin"
  }

  consistency_policy {
    consistency_level = "Strong"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  backup {
    type = "Continuous"
  }
}

resource "azurerm_cosmosdb_gremlin_database" "test" {
  name                = "acctest-gremlindb-%d"
  resource_group_name = azurerm_cosmosdb_account.test1.resource_group_name
  account_name        = azurerm_cosmosdb_account.test1.name
}

resource "azurerm_cosmosdb_gremlin_graph" "test" {
  name                = "acctest-CGRPC-%d"
  resource_group_name = azurerm_cosmosdb_account.test1.resource_group_name
  account_name        = azurerm_cosmosdb_account.test1.name
  database_name       = azurerm_cosmosdb_gremlin_database.test.name
  partition_key_path  = "/test"
  throughput          = 400
}

resource "azurerm_cosmosdb_gremlin_graph" "test2" {
  name                = "acctest-CGRPC2-%d"
  resource_group_name = azurerm_cosmosdb_account.test1.resource_group_name
  account_name        = azurerm_cosmosdb_account.test1.name
  database_name       = azurerm_cosmosdb_gremlin_database.test.name
  partition_key_path  = "/test2"
  throughput          = 500

  depends_on = [azurerm_cosmosdb_gremlin_graph.test]
}

data "azurerm_cosmosdb_restorable_database_accounts" "test" {
  name     = azurerm_cosmosdb_account.test1.name
  location = azurerm_resource_group.test.location

  depends_on = [azurerm_cosmosdb_gremlin_graph.test2]
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%s"

  capabilities {
    name = "EnableGremlin"
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

    gremlin_database {
      name        = azurerm_cosmosdb_gremlin_database.test.name
      graph_names = [azurerm_cosmosdb_gremlin_graph.test.name, azurerm_cosmosdb_gremlin_graph.test2.name]
    }
  }

  // As "restore_timestamp_in_utc" is retrieved dynamically, so it would cause diff when tf plan. So we have to ignore it here.
  lifecycle {
    ignore_changes = [
      restore.0.restore_timestamp_in_utc
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, string(kind), string(consistency))
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

  multiple_write_locations_enabled = false
  automatic_failover_enabled       = false

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

  multiple_write_locations_enabled = false
  automatic_failover_enabled       = false

  consistency_policy {
    consistency_level       = "Eventual"
    max_interval_in_seconds = 5
    max_staleness_prefix    = 100
  }

  is_virtual_network_filter_enabled = true
  ip_range_filter                   = ["55.0.1.0/24", "55.0.2.0/24", "0.0.0.0"]

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

func (r CosmosDBAccountResource) withoutMaxAgeInSeconds(data acceptance.TestData, kind cosmosdb.DatabaseAccountKind, consistency cosmosdb.DefaultConsistencyLevel) string {
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
  multiple_write_locations_enabled = true
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
