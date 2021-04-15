package cosmos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-01-15/documentdb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type CosmosDBAccountMongoResource struct {
}

func TestAccCosmosDBAccountMongo_basic_boundedStaleness(t *testing.T) {
	testAccCosmosDBAccountMongo_basicWithConsistency(t, documentdb.BoundedStaleness)
}

func TestAccCosmosDBAccountMongo_basic_consistentPrefix(t *testing.T) {
	testAccCosmosDBAccountMongo_basicWithConsistency(t, documentdb.ConsistentPrefix)
}

func TestAccCosmosDBAccountMongo_basic_eventual(t *testing.T) {
	testAccCosmosDBAccountMongo_basicWithConsistency(t, documentdb.Eventual)
}

func TestAccCosmosDBAccountMongo_basic_session(t *testing.T) {
	testAccCosmosDBAccountMongo_basicWithConsistency(t, documentdb.Session)
}

func TestAccCosmosDBAccountMongo_basic_strong(t *testing.T) {
	testAccCosmosDBAccountMongo_basicWithConsistency(t, documentdb.Strong)
}

func TestAccCosmosDBAccountMongo_public_network_access_enabled(t *testing.T) {
	testAccCosmosDBAccountMongo_public_network_access_enabled(t, documentdb.Strong)
}

func testAccCosmosDBAccountMongo_public_network_access_enabled(t *testing.T, consistency documentdb.DefaultConsistencyLevel) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountMongoResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.network_access_enabled(data, consistency),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, consistency, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccountMongo_keyVaultUri(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountMongoResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.key_vault_uri(data, documentdb.Strong),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Strong, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccountMongo_keyVaultUriUpdateConsistancy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountMongoResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.key_vault_uri(data, documentdb.Strong),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Strong, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.key_vault_uri(data, documentdb.Session),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Session, 1),
			),
		},
		data.ImportStep(),
	})
}

func testAccCosmosDBAccountMongo_basicWithConsistency(t *testing.T, consistency documentdb.DefaultConsistencyLevel) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountMongoResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, consistency),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, consistency, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccountMongo_updateConsistency(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountMongoResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, documentdb.Strong),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.Strong, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistency(data, documentdb.Strong, 8, 880),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.Strong, 1),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, documentdb.BoundedStaleness),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.BoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistency(data, documentdb.BoundedStaleness, 7, 770),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.BoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.consistency(data, documentdb.BoundedStaleness, 77, 700),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.BoundedStaleness, 1),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, documentdb.ConsistentPrefix),
			Check:  checkAccCosmosDBAccount_basic(data, documentdb.ConsistentPrefix, 1),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccountMongo_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountMongoResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data, documentdb.Eventual),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 3),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccountMongo_zoneRedundant(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountMongoResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.zoneRedundant(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccCosmosDBAccountMongo_zoneRedundant_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountMongoResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, documentdb.Eventual),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.zoneRedundantUpdate(data, documentdb.Eventual),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 2),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccountMongo_zoneRedundant_update(t *testing.T) {
	testAccCosmosDBAccountMongo_zoneRedundant_update(t)
}

func TestAccCosmosDBAccountMongo_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountMongoResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, documentdb.Eventual),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 1),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, documentdb.Eventual),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 3),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdated(data, documentdb.Eventual),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 3),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithResources(data, documentdb.Eventual),
			Check:  resource.ComposeAggregateTestCheckFunc(
			// checkAccCosmosDBAccount_basic(data, documentdb.Eventual, 1),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDBAccountMongo_capabilities_EnableMongo(t *testing.T) {
	testAccCosmosDBAccountMongo_capabilitiesWith(t, []string{"EnableMongo"})
}

func TestAccCosmosDBAccountMongo_capabilities_mongoEnableDocLevelTTL(t *testing.T) {
	testAccCosmosDBAccountMongo_capabilitiesWith(t, []string{"EnableMongo", "mongoEnableDocLevelTTL"})
}

func TestAccCosmosDBAccountMongo_capabilities_DisableRateLimitingResponses(t *testing.T) {
	testAccCosmosDBAccountMongo_capabilitiesWith(t, []string{"EnableMongo", "DisableRateLimitingResponses"})
}

func TestAccCosmosDBAccountMongo_capabilities_AllowSelfServeUpgradeToMongo36(t *testing.T) {
	testAccCosmosDBAccountMongo_capabilitiesWith(t, []string{"EnableMongo", "AllowSelfServeUpgradeToMongo36"})
}

func testAccCosmosDBAccountMongo_capabilitiesWith(t *testing.T, capabilities []string) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountMongoResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.capabilities(data, capabilities),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.Strong, 1),
			),
		},
		data.ImportStep(),
	})
}

func (t CosmosDBAccountMongoResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
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

func (CosmosDBAccountMongoResource) basic(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
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

func (CosmosDBAccountMongoResource) consistency(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel, interval, staleness int) string {
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

  consistency_policy {
    consistency_level       = "%s"
    max_interval_in_seconds = %d
    max_staleness_prefix    = %d
  }

  capabilities {
    name = "EnableMongo"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(consistency), interval, staleness)
}

func (CosmosDBAccountMongoResource) completePreReqs(data acceptance.TestData) string {
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

func (r CosmosDBAccountMongoResource) complete(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "MongoDB"

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

  capabilities {
    name = "EnableMongo"
  }

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
}
`, r.completePreReqs(data), data.RandomInteger, string(consistency), data.Locations.Secondary, data.Locations.Ternary)
}

func (CosmosDBAccountMongoResource) zoneRedundant(data acceptance.TestData) string {
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

  enable_multiple_write_locations = true

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 300
    max_staleness_prefix    = 100000
  }

  capabilities {
    name = "EnableMongo"
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

func (r CosmosDBAccountMongoResource) completeUpdated(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "MongoDB"

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
}
`, r.completePreReqs(data), data.RandomInteger, string(consistency), data.Locations.Secondary, data.Locations.Ternary)
}

func (r CosmosDBAccountMongoResource) basicWithResources(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "MongoDB"

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
`, r.completePreReqs(data), data.RandomInteger, string(consistency))
}

func (CosmosDBAccountMongoResource) capabilities(data acceptance.TestData, capabilities []string) string {
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
  kind                = "MongoDB"

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

func (CosmosDBAccountMongoResource) mongo36WithCapabilities(data acceptance.TestData, capabilities []string) string {
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
  name                 = "acctest-ca-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  offer_type           = "Standard"
  kind                 = "MongoDB"
  mongo_server_version = "3.6"

  consistency_policy {
    consistency_level = "Strong"
  }

  capabilities {
    name = "EnableMongo"
  }

  %s

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, capeTf)
}

func (CosmosDBAccountMongoResource) geoLocationUpdate(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
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

  capabilities {
    name = "EnableMongo"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(consistency), data.Locations.Secondary)
}

func (CosmosDBAccountMongoResource) zoneRedundantUpdate(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
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

  enable_multiple_write_locations = true
  enable_automatic_failover       = true

  consistency_policy {
    consistency_level = "%s"
  }

  capabilities {
    name = "EnableMongo"
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

func (CosmosDBAccountMongoResource) network_access_enabled(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
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
  kind                          = "MongoDB"
  public_network_access_enabled = true

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

func (CosmosDBAccountMongoResource) key_vault_uri(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
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
  kind                = "MongoDB"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger, string(consistency))
}
