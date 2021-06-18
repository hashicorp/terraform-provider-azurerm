package springcloud_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SpringCloudAppCosmosDbAssociationResource struct {
}

func TestAccSpringCloudAppCosmosDbAssociation_cassandra_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_cosmosdb_association", "test")
	r := SpringCloudAppCosmosDbAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cassandra_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cosmosdb_access_key"),
	})
}

func TestAccSpringCloudAppCosmosDbAssociation_cassandra_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_cosmosdb_association", "test")
	r := SpringCloudAppCosmosDbAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cassandra_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.cassandra_requiresImport),
	})
}

func TestAccSpringCloudAppCosmosDbAssociation_cassandra_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_cosmosdb_association", "test")
	r := SpringCloudAppCosmosDbAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cassandra_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cosmosdb_access_key"),
		{
			Config: r.cassandra_update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cosmosdb_access_key"),
	})
}

func TestAccSpringCloudAppCosmosDbAssociation_gremlin_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_cosmosdb_association", "test")
	r := SpringCloudAppCosmosDbAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gremlin_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cosmosdb_access_key"),
	})
}

func TestAccSpringCloudAppCosmosDbAssociation_gremlin_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_cosmosdb_association", "test")
	r := SpringCloudAppCosmosDbAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gremlin_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.gremlin_requiresImport),
	})
}

func TestAccSpringCloudAppCosmosDbAssociation_gremlin_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_cosmosdb_association", "test")
	r := SpringCloudAppCosmosDbAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gremlin_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cosmosdb_access_key"),
		{
			Config: r.gremlin_update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cosmosdb_access_key"),
	})
}

func TestAccSpringCloudAppCosmosDbAssociation_mongo_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_cosmosdb_association", "test")
	r := SpringCloudAppCosmosDbAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.mongo_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cosmosdb_access_key"),
	})
}

func TestAccSpringCloudAppCosmosDbAssociation_mongo_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_cosmosdb_association", "test")
	r := SpringCloudAppCosmosDbAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.mongo_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.mongo_requiresImport),
	})
}

func TestAccSpringCloudAppCosmosDbAssociation_mongo_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_cosmosdb_association", "test")
	r := SpringCloudAppCosmosDbAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.mongo_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cosmosdb_access_key"),
		{
			Config: r.mongo_update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cosmosdb_access_key"),
	})
}

func TestAccSpringCloudAppCosmosDbAssociation_sql_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_cosmosdb_association", "test")
	r := SpringCloudAppCosmosDbAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sql_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cosmosdb_access_key"),
	})
}

func TestAccSpringCloudAppCosmosDbAssociation_sql_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_cosmosdb_association", "test")
	r := SpringCloudAppCosmosDbAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sql_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.sql_requiresImport),
	})
}

func TestAccSpringCloudAppCosmosDbAssociation_sql_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_cosmosdb_association", "test")
	r := SpringCloudAppCosmosDbAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sql_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cosmosdb_access_key"),
		{
			Config: r.sql_update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cosmosdb_access_key"),
	})
}

func TestAccSpringCloudAppCosmosDbAssociation_table_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_cosmosdb_association", "test")
	r := SpringCloudAppCosmosDbAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.table_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cosmosdb_access_key"),
	})
}

func TestAccSpringCloudAppCosmosDbAssociation_table_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_cosmosdb_association", "test")
	r := SpringCloudAppCosmosDbAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.table_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.table_requiresImport),
	})
}

func TestAccSpringCloudAppCosmosDbAssociation_table_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_cosmosdb_association", "test")
	r := SpringCloudAppCosmosDbAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.table_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cosmosdb_access_key"),
		{
			Config: r.table_update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cosmosdb_access_key"),
	})
}

func (t SpringCloudAppCosmosDbAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SpringCloudAppAssociationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppPlatform.BindingsClient.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.BindingName)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (r SpringCloudAppCosmosDbAssociationResource) cassandra_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app_cosmosdb_association" "test" {
  name                             = "acctestscac-%d"
  spring_cloud_app_id              = azurerm_spring_cloud_app.test.id
  cosmosdb_account_id              = azurerm_cosmosdb_account.test.id
  api_type                         = "cassandra"
  cosmosdb_cassandra_keyspace_name = azurerm_cosmosdb_cassandra_keyspace.test.name
  cosmosdb_access_key              = azurerm_cosmosdb_account.test.primary_key
}
`, r.cassandra_template(data), data.RandomInteger)
}

func (r SpringCloudAppCosmosDbAssociationResource) cassandra_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app_cosmosdb_association" "import" {
  name                             = azurerm_spring_cloud_app_cosmosdb_association.test.name
  spring_cloud_app_id              = azurerm_spring_cloud_app_cosmosdb_association.test.spring_cloud_app_id
  cosmosdb_account_id              = azurerm_spring_cloud_app_cosmosdb_association.test.cosmosdb_account_id
  api_type                         = azurerm_spring_cloud_app_cosmosdb_association.test.api_type
  cosmosdb_cassandra_keyspace_name = azurerm_spring_cloud_app_cosmosdb_association.test.cosmosdb_cassandra_keyspace_name
  cosmosdb_access_key              = azurerm_spring_cloud_app_cosmosdb_association.test.cosmosdb_access_key
}
`, r.cassandra_basic(data))
}

func (r SpringCloudAppCosmosDbAssociationResource) cassandra_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_cassandra_keyspace" "update" {
  name                = "acctest-ck1-%d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_spring_cloud_app_cosmosdb_association" "test" {
  name                             = "acctestscac-%d"
  spring_cloud_app_id              = azurerm_spring_cloud_app.test.id
  cosmosdb_account_id              = azurerm_cosmosdb_account.test.id
  api_type                         = "cassandra"
  cosmosdb_cassandra_keyspace_name = azurerm_cosmosdb_cassandra_keyspace.update.name
  cosmosdb_access_key              = azurerm_cosmosdb_account.test.primary_key
}
`, r.cassandra_template(data), data.RandomInteger, data.RandomInteger)
}

func (r SpringCloudAppCosmosDbAssociationResource) gremlin_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app_cosmosdb_association" "test" {
  name                           = "acctestscac-%d"
  spring_cloud_app_id            = azurerm_spring_cloud_app.test.id
  cosmosdb_account_id            = azurerm_cosmosdb_account.test.id
  api_type                       = "gremlin"
  cosmosdb_gremlin_database_name = azurerm_cosmosdb_gremlin_database.test.name
  cosmosdb_gremlin_graph_name    = azurerm_cosmosdb_gremlin_graph.test.name
  cosmosdb_access_key            = azurerm_cosmosdb_account.test.primary_key
}
`, r.gremlin_template(data), data.RandomInteger)
}

func (r SpringCloudAppCosmosDbAssociationResource) gremlin_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app_cosmosdb_association" "import" {
  name                           = azurerm_spring_cloud_app_cosmosdb_association.test.name
  spring_cloud_app_id            = azurerm_spring_cloud_app_cosmosdb_association.test.spring_cloud_app_id
  cosmosdb_account_id            = azurerm_spring_cloud_app_cosmosdb_association.test.cosmosdb_account_id
  api_type                       = azurerm_spring_cloud_app_cosmosdb_association.test.api_type
  cosmosdb_gremlin_database_name = azurerm_spring_cloud_app_cosmosdb_association.test.cosmosdb_gremlin_database_name
  cosmosdb_gremlin_graph_name    = azurerm_spring_cloud_app_cosmosdb_association.test.cosmosdb_gremlin_graph_name
  cosmosdb_access_key            = azurerm_spring_cloud_app_cosmosdb_association.test.cosmosdb_access_key
}
`, r.gremlin_basic(data))
}

func (r SpringCloudAppCosmosDbAssociationResource) gremlin_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_gremlin_graph" "update" {
  name                = "acctest-CGG1-%d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_gremlin_database.test.name
  partition_key_path  = "/test"
  throughput          = 400

  index_policy {
    automatic      = true
    indexing_mode  = "Consistent"
    included_paths = ["/*"]
    excluded_paths = ["/\"_etag\"/?"]
  }
}

resource "azurerm_spring_cloud_app_cosmosdb_association" "test" {
  name                           = "acctestscac-%d"
  spring_cloud_app_id            = azurerm_spring_cloud_app.test.id
  cosmosdb_account_id            = azurerm_cosmosdb_account.test.id
  api_type                       = "gremlin"
  cosmosdb_gremlin_database_name = azurerm_cosmosdb_gremlin_database.test.name
  cosmosdb_gremlin_graph_name    = azurerm_cosmosdb_gremlin_graph.update.name
  cosmosdb_access_key            = azurerm_cosmosdb_account.test.primary_key
}
`, r.gremlin_template(data), data.RandomInteger, data.RandomInteger)
}

func (r SpringCloudAppCosmosDbAssociationResource) mongo_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app_cosmosdb_association" "test" {
  name                         = "acctestscac-%d"
  spring_cloud_app_id          = azurerm_spring_cloud_app.test.id
  cosmosdb_account_id          = azurerm_cosmosdb_account.test.id
  api_type                     = "mongo"
  cosmosdb_mongo_database_name = azurerm_cosmosdb_mongo_database.test.name
  cosmosdb_access_key          = azurerm_cosmosdb_account.test.primary_key
}
`, r.mongo_template(data), data.RandomInteger)
}

func (r SpringCloudAppCosmosDbAssociationResource) mongo_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app_cosmosdb_association" "import" {
  name                         = azurerm_spring_cloud_app_cosmosdb_association.test.name
  spring_cloud_app_id          = azurerm_spring_cloud_app_cosmosdb_association.test.spring_cloud_app_id
  cosmosdb_account_id          = azurerm_spring_cloud_app_cosmosdb_association.test.cosmosdb_account_id
  api_type                     = azurerm_spring_cloud_app_cosmosdb_association.test.api_type
  cosmosdb_mongo_database_name = azurerm_spring_cloud_app_cosmosdb_association.test.cosmosdb_mongo_database_name
  cosmosdb_access_key          = azurerm_spring_cloud_app_cosmosdb_association.test.cosmosdb_access_key
}
`, r.mongo_basic(data))
}

func (r SpringCloudAppCosmosDbAssociationResource) mongo_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_mongo_database" "update" {
  name                = "acctest-mongo1-%d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_spring_cloud_app_cosmosdb_association" "test" {
  name                         = "acctestscac-%d"
  spring_cloud_app_id          = azurerm_spring_cloud_app.test.id
  cosmosdb_account_id          = azurerm_cosmosdb_account.test.id
  api_type                     = "mongo"
  cosmosdb_mongo_database_name = azurerm_cosmosdb_mongo_database.update.name
  cosmosdb_access_key          = azurerm_cosmosdb_account.test.primary_key
}
`, r.mongo_template(data), data.RandomInteger, data.RandomInteger)
}

func (r SpringCloudAppCosmosDbAssociationResource) sql_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app_cosmosdb_association" "test" {
  name                       = "acctestscac-%d"
  spring_cloud_app_id        = azurerm_spring_cloud_app.test.id
  cosmosdb_account_id        = azurerm_cosmosdb_account.test.id
  api_type                   = "sql"
  cosmosdb_sql_database_name = azurerm_cosmosdb_sql_database.test.name
  cosmosdb_access_key        = azurerm_cosmosdb_account.test.primary_key
}
`, r.sql_template(data), data.RandomInteger)
}

func (r SpringCloudAppCosmosDbAssociationResource) sql_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app_cosmosdb_association" "import" {
  name                       = azurerm_spring_cloud_app_cosmosdb_association.test.name
  spring_cloud_app_id        = azurerm_spring_cloud_app_cosmosdb_association.test.spring_cloud_app_id
  cosmosdb_account_id        = azurerm_spring_cloud_app_cosmosdb_association.test.cosmosdb_account_id
  api_type                   = azurerm_spring_cloud_app_cosmosdb_association.test.api_type
  cosmosdb_sql_database_name = azurerm_spring_cloud_app_cosmosdb_association.test.cosmosdb_sql_database_name
  cosmosdb_access_key        = azurerm_spring_cloud_app_cosmosdb_association.test.cosmosdb_access_key
}
`, r.sql_basic(data))
}

func (r SpringCloudAppCosmosDbAssociationResource) sql_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_sql_database" "update" {
  name                = "acctest-sql1-%d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_spring_cloud_app_cosmosdb_association" "test" {
  name                       = "acctestscac-%d"
  spring_cloud_app_id        = azurerm_spring_cloud_app.test.id
  cosmosdb_account_id        = azurerm_cosmosdb_account.test.id
  api_type                   = "sql"
  cosmosdb_sql_database_name = azurerm_cosmosdb_sql_database.update.name
  cosmosdb_access_key        = azurerm_cosmosdb_account.test.primary_key
}
`, r.sql_template(data), data.RandomInteger, data.RandomInteger)
}

func (r SpringCloudAppCosmosDbAssociationResource) table_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app_cosmosdb_association" "test" {
  name                = "acctestscac-%d"
  spring_cloud_app_id = azurerm_spring_cloud_app.test.id
  cosmosdb_account_id = azurerm_cosmosdb_account.test.id
  api_type            = "table"
  cosmosdb_access_key = azurerm_cosmosdb_account.test.primary_key
}
`, r.table_template(data), data.RandomInteger)
}

func (r SpringCloudAppCosmosDbAssociationResource) table_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app_cosmosdb_association" "import" {
  name                = azurerm_spring_cloud_app_cosmosdb_association.test.name
  spring_cloud_app_id = azurerm_spring_cloud_app_cosmosdb_association.test.spring_cloud_app_id
  cosmosdb_account_id = azurerm_spring_cloud_app_cosmosdb_association.test.cosmosdb_account_id
  api_type            = azurerm_spring_cloud_app_cosmosdb_association.test.api_type
  cosmosdb_access_key = azurerm_spring_cloud_app_cosmosdb_association.test.cosmosdb_access_key
}
`, r.table_basic(data))
}

func (r SpringCloudAppCosmosDbAssociationResource) table_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app_cosmosdb_association" "test" {
  name                = "acctestscac-%d"
  spring_cloud_app_id = azurerm_spring_cloud_app.test.id
  cosmosdb_account_id = azurerm_cosmosdb_account.test.id
  api_type            = "table"
  cosmosdb_access_key = azurerm_cosmosdb_account.test.secondary_key
}
`, r.table_template(data), data.RandomInteger)
}

func (r SpringCloudAppCosmosDbAssociationResource) cassandra_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%d"
  location = "%s"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_spring_cloud_app" "test" {
  name                = "acctest-sca-%d"
  resource_group_name = azurerm_spring_cloud_service.test.resource_group_name
  service_name        = azurerm_spring_cloud_service.test.name
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

  capabilities {
    name = "EnableCassandra"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_cassandra_keyspace" "test" {
  name                = "acctest-ck-%d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r SpringCloudAppCosmosDbAssociationResource) gremlin_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%d"
  location = "%s"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_spring_cloud_app" "test" {
  name                = "acctest-sca-%d"
  resource_group_name = azurerm_spring_cloud_service.test.resource_group_name
  service_name        = azurerm_spring_cloud_service.test.name
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

  capabilities {
    name = "EnableGremlin"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_gremlin_database" "test" {
  name                = "acctest-CGD-%d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_cosmosdb_gremlin_graph" "test" {
  name                = "acctest-CGG-%d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_gremlin_database.test.name
  partition_key_path  = "/test"
  throughput          = 400

  index_policy {
    automatic      = true
    indexing_mode  = "Consistent"
    included_paths = ["/*"]
    excluded_paths = ["/\"_etag\"/?"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r SpringCloudAppCosmosDbAssociationResource) mongo_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%d"
  location = "%s"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_spring_cloud_app" "test" {
  name                = "acctest-sca-%d"
  resource_group_name = azurerm_spring_cloud_service.test.resource_group_name
  service_name        = azurerm_spring_cloud_service.test.name
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

  capabilities {
    name = "EnableMongo"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_mongo_database" "test" {
  name                = "acctest-mongo-%d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r SpringCloudAppCosmosDbAssociationResource) sql_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%d"
  location = "%s"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_spring_cloud_app" "test" {
  name                = "acctest-sca-%d"
  resource_group_name = azurerm_spring_cloud_service.test.resource_group_name
  service_name        = azurerm_spring_cloud_service.test.name
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

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "acctest-sql-%d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r SpringCloudAppCosmosDbAssociationResource) table_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%d"
  location = "%s"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_spring_cloud_app" "test" {
  name                = "acctest-sca-%d"
  resource_group_name = azurerm_spring_cloud_service.test.resource_group_name
  service_name        = azurerm_spring_cloud_service.test.name
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

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
