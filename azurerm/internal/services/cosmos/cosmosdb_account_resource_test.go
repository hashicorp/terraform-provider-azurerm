package cosmos_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-01-15/documentdb"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type CosmosDBAccountResource struct {
}

func TestAccCosmosDBAccount_basic_global_boundedStaleness(t *testing.T) {
	testAccCosmosDBAccount_basicWith(t, documentdb.GlobalDocumentDB, documentdb.BoundedStaleness)
}

func TestAccCosmosDBAccount_basic_global_consistentPrefix(t *testing.T) {
	testAccCosmosDBAccount_basicWith(t, documentdb.GlobalDocumentDB, documentdb.ConsistentPrefix)
}

func TestAccCosmosDBAccount_basic_global_eventual(t *testing.T) {
	testAccCosmosDBAccount_basicWith(t, documentdb.GlobalDocumentDB, documentdb.Eventual)
}

func TestAccCosmosDBAccount_basic_global_session(t *testing.T) {
	testAccCosmosDBAccount_basicWith(t, documentdb.GlobalDocumentDB, documentdb.Session)
}

func TestAccCosmosDBAccount_basic_global_strong(t *testing.T) {
	testAccCosmosDBAccount_basicWith(t, documentdb.GlobalDocumentDB, documentdb.Strong)
}

func TestAccCosmosDBAccount_basic_mongo_boundedStaleness(t *testing.T) {
	testAccCosmosDBAccount_basicMongoDBWith(t, documentdb.BoundedStaleness)
}

func TestAccCosmosDBAccount_basic_mongo_consistentPrefix(t *testing.T) {
	testAccCosmosDBAccount_basicMongoDBWith(t, documentdb.ConsistentPrefix)
}

func TestAccCosmosDBAccount_basic_mongo_eventual(t *testing.T) {
	testAccCosmosDBAccount_basicMongoDBWith(t, documentdb.Eventual)
}

func TestAccCosmosDBAccount_basic_mongo_session(t *testing.T) {
	testAccCosmosDBAccount_basicMongoDBWith(t, documentdb.Session)
}

func TestAccCosmosDBAccount_basic_mongo_strong(t *testing.T) {
	testAccCosmosDBAccount_basicMongoDBWith(t, documentdb.Strong)
}

func TestAccCosmosDBAccount_basic_parse_boundedStaleness(t *testing.T) {
	testAccCosmosDBAccount_basicWith(t, documentdb.Parse, documentdb.BoundedStaleness)
}

func TestAccCosmosDBAccount_basic_parse_consistentPrefix(t *testing.T) {
	testAccCosmosDBAccount_basicWith(t, documentdb.Parse, documentdb.ConsistentPrefix)
}

func TestAccCosmosDBAccount_basic_parse_eventual(t *testing.T) {
	testAccCosmosDBAccount_basicWith(t, documentdb.Parse, documentdb.Eventual)
}

func TestAccCosmosDBAccount_basic_parse_session(t *testing.T) {
	testAccCosmosDBAccount_basicWith(t, documentdb.Parse, documentdb.Session)
}

func TestAccCosmosDBAccount_basic_parse_strong(t *testing.T) {
	testAccCosmosDBAccount_basicWith(t, documentdb.Parse, documentdb.Strong)
}

func TestAccCosmosDBAccount_public_network_access_enabled(t *testing.T) {
	testAccCosmosDBAccount_public_network_access_enabled(t, documentdb.MongoDB, documentdb.Strong)
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
			Config: r.key_vault_uri(data, documentdb.MongoDB, documentdb.Strong),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Strong, 1),
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
			Config: r.key_vault_uri(data, documentdb.MongoDB, documentdb.Strong),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Strong, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.key_vault_uri(data, documentdb.MongoDB, documentdb.Session),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Session, 1),
			),
		},
		data.ImportStep(),
	})
}

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
			Config: r.basic(data, "GlobalDocumentDB", documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 1),
			),
		},
		{
			Config:      r.requiresImport(data, documentdb.Eventual),
			ExpectError: acceptance.RequiresImportError("azurerm_cosmosdb_account"),
		},
	})
}

func TestAccCosmosDBAccount_updateConsistency_global(t *testing.T) {
	testAccCosmosDBAccount_updateConsistency(t, documentdb.GlobalDocumentDB)
}

func TestAccCosmosDBAccount_updateConsistency_mongo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMongoDB(data, documentdb.Strong),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.Strong, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistencyMongoDB(data, documentdb.Strong, 8, 880),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.Strong, 1),
		},
		data.ImportStep(),
		{
			Config: r.basicMongoDB(data, documentdb.BoundedStaleness),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.BoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistencyMongoDB(data, documentdb.BoundedStaleness, 7, 770),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.BoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistencyMongoDB(data, documentdb.BoundedStaleness, 77, 700),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.BoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.basicMongoDB(data, documentdb.ConsistentPrefix),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.ConsistentPrefix, 1),
		},
		data.ImportStep(),
	})
}

func testAccCosmosDBAccount_updateConsistency(t *testing.T, kind documentdb.DatabaseAccountKind) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, kind, documentdb.Strong),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.Strong, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistency(data, kind, documentdb.Strong, 8, 880),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.Strong, 1),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, kind, documentdb.BoundedStaleness),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.BoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistency(data, kind, documentdb.BoundedStaleness, 7, 770),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.BoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistency(data, kind, documentdb.BoundedStaleness, 77, 700),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.BoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, kind, documentdb.ConsistentPrefix),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.ConsistentPrefix, 1),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_complete_mongo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeMongoDB(data, documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 3),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_complete_global(t *testing.T) {
	testAccCosmosDBAccount_completeWith(t, documentdb.GlobalDocumentDB)
}

func TestAccCosmosDBAccount_complete_parse(t *testing.T) {
	testAccCosmosDBAccount_completeWith(t, documentdb.Parse)
}

func testAccCosmosDBAccount_completeWith(t *testing.T, kind documentdb.DatabaseAccountKind) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, kind, documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 3),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_completeZoneRedundant_mongo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
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
	testAccCosmosDBAccount_zoneRedundantWith(t, documentdb.GlobalDocumentDB)
}

func TestAccCosmosDBAccount_completeZoneRedundant_parse(t *testing.T) {
	testAccCosmosDBAccount_zoneRedundantWith(t, documentdb.Parse)
}

func testAccCosmosDBAccount_zoneRedundantWith(t *testing.T, kind documentdb.DatabaseAccountKind) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
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
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMongoDB(data, documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.zoneRedundantMongoDBUpdate(data, documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 2),
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
			Config: r.basicMongoDB(data, documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeMongoDB(data, documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 3),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdatedMongoDB(data, documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 3),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithResourcesMongoDB(data, documentdb.Eventual),
			Check:  acceptance.ComposeAggregateTestCheckFunc(
			// checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_update_global(t *testing.T) {
	testAccCosmosDBAccount_updateWith(t, documentdb.GlobalDocumentDB)
}

func TestAccCosmosDBAccount_update_parse(t *testing.T) {
	testAccCosmosDBAccount_updateWith(t, documentdb.Parse)
}

func testAccCosmosDBAccount_updateWith(t *testing.T, kind documentdb.DatabaseAccountKind) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, kind, documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, kind, documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 3),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdated(data, kind, documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 3),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithResources(data, kind, documentdb.Eventual),
			Check:  acceptance.ComposeAggregateTestCheckFunc(
			// checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccount_capabilities_EnableAggregationPipeline(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, documentdb.GlobalDocumentDB, []string{"EnableAggregationPipeline"})
}

func TestAccCosmosDBAccount_capabilities_EnableCassandra(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, documentdb.GlobalDocumentDB, []string{"EnableCassandra"})
}

func TestAccCosmosDBAccount_capabilities_EnableGremlin(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, documentdb.GlobalDocumentDB, []string{"EnableGremlin"})
}

func TestAccCosmosDBAccount_capabilities_EnableTable(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, documentdb.GlobalDocumentDB, []string{"EnableTable"})
}

func TestAccCosmosDBAccount_capabilities_EnableServerless(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, documentdb.GlobalDocumentDB, []string{"EnableServerless"})
}

func TestAccCosmosDBAccount_capabilities_EnableMongo(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, documentdb.MongoDB, []string{"EnableMongo"})
}

func TestAccCosmosDBAccount_capabilities_MongoDBv34(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, documentdb.MongoDB, []string{"EnableMongo", "MongoDBv3.4"})
}

func TestAccCosmosDBAccount_capabilities_mongoEnableDocLevelTTL(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, documentdb.MongoDB, []string{"EnableMongo", "mongoEnableDocLevelTTL"})
}

func TestAccCosmosDBAccount_capabilities_DisableRateLimitingResponses(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, documentdb.MongoDB, []string{"EnableMongo", "DisableRateLimitingResponses"})
}

func TestAccCosmosDBAccount_capabilities_AllowSelfServeUpgradeToMongo36(t *testing.T) {
	testAccCosmosDBAccount_capabilitiesWith(t, documentdb.MongoDB, []string{"EnableMongo", "AllowSelfServeUpgradeToMongo36"})
}

func testAccCosmosDBAccount_capabilitiesWith(t *testing.T, kind documentdb.DatabaseAccountKind, capabilities []string) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.capabilities(data, kind, capabilities),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Strong, 1),
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
			Config: r.capabilities(data, documentdb.GlobalDocumentDB, []string{"EnableCassandra"}),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Strong, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.capabilities(data, documentdb.GlobalDocumentDB, []string{"EnableCassandra", "EnableAggregationPipeline"}),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Strong, 1),
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
			Config: r.capabilities(data, documentdb.GlobalDocumentDB, []string{"EnableCassandra"}),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Strong, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.capabilities(data, documentdb.GlobalDocumentDB, []string{"EnableTable", "EnableAggregationPipeline"}),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Strong, 1),
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
			Config: r.basic(data, "GlobalDocumentDB", documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.geoLocationUpdate(data, "GlobalDocumentDB", documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 2),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "GlobalDocumentDB", documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 1),
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
			Config: r.freeTier(data, "GlobalDocumentDB", documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 1),
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
			Config: r.analyticalStorage(data, "GlobalDocumentDB", documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 1),
				check.That(data.ResourceName).Key("analytical_storage_enabled").HasValue("true"),
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
			Config: r.basicMongoDB(data, documentdb.Session),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.systemAssignedIdentity(data, documentdb.Session),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicMongoDB(data, documentdb.Session),
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
			Config: r.basic(data, documentdb.GlobalDocumentDB, documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.type").HasValue("Periodic"),
				check.That(data.ResourceName).Key("backup.0.interval_in_minutes").HasValue("240"),
				check.That(data.ResourceName).Key("backup.0.retention_in_hours").HasValue("8"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithBackupPeriodic(data, documentdb.GlobalDocumentDB, documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithBackupPeriodicUpdate(data, documentdb.GlobalDocumentDB, documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.type").HasValue("Periodic"),
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
			Config: r.basicWithBackupContinuous(data, documentdb.GlobalDocumentDB, documentdb.Eventual),
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
			Config: r.basic(data, documentdb.GlobalDocumentDB, documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithNetworkBypass(data, documentdb.GlobalDocumentDB, documentdb.Eventual),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithoutNetworkBypass(data, documentdb.GlobalDocumentDB, documentdb.Eventual),
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
			Config: r.basicMongoDBVersion32(data, documentdb.Session),
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
			Config: r.basicMongoDBVersion40(data, documentdb.Session),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
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
  address_prefix       = "10.0.1.0/24"
  service_endpoints    = ["Microsoft.AzureCosmosDB"]
}

resource "azurerm_subnet" "subnet2" {
  name                 = "acctest-SN2-%[1]d-2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
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
    max_age_in_seconds = "500"
  }

  access_key_metadata_writes_enabled    = false
  network_acl_bypass_for_azure_services = true
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
    max_age_in_seconds = "500"
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
    max_age_in_seconds = "2000000000"
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
    max_age_in_seconds = "2000000000"
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

func (CosmosDBAccountResource) analyticalStorage(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
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

  analytical_storage_enabled = true

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
		check.That(data.ResourceName).Key("offer_type").HasValue(string(documentdb.Standard)),
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
      purge_soft_delete_on_destroy = false
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
  soft_delete_enabled        = true
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "list",
      "create",
      "delete",
      "get",
      "purge",
      "update",
    ]

    secret_permissions = [
      "get",
      "delete",
      "set",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azuread_service_principal.cosmosdb.id

    key_permissions = [
      "list",
      "create",
      "delete",
      "get",
      "update",
      "unwrapKey",
      "wrapKey",
    ]

    secret_permissions = [
      "get",
      "delete",
      "set",
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

func (CosmosDBAccountResource) systemAssignedIdentity(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
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

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(consistency))
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
