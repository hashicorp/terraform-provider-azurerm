package redisenterprise_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redisenterprise/sdk/2022-01-01/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type RedisenterpriseGeoDatabaseResource struct{}

func TestRedisEnterpriseGeoDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_enterprise_geo_database", "test")
	r := RedisenterpriseGeoDatabaseResource{}
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

func TestRedisEnterpriseGeoDatabase_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_enterprise_geo_database", "test")
	r := RedisenterpriseGeoDatabaseResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestRedisEnterpriseGeoDatabase_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_enterprise_geo_database", "test")
	r := RedisenterpriseGeoDatabaseResource{}
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

func TestRedisEnterpriseGeoDatabase_unlink(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_enterprise_geo_database", "test")
	r := RedisenterpriseGeoDatabaseResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
		{
			Config: r.unlinkDatabase(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}
func (r RedisenterpriseGeoDatabaseResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := databases.ParseDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.RedisEnterprise.GeoDatabaseClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (r RedisenterpriseGeoDatabaseResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_redis_enterprise_geo_database" "test" {
  cluster_id          = azurerm_redis_enterprise_cluster.test.id

  client_protocol            = "Encrypted"
  clustering_policy          = "EnterpriseCluster"
  eviction_policy            = "NoEviction"

  linked_database_id = [
    "${azurerm_redis_enterprise_cluster.test.id}/databases/default",
    "${azurerm_redis_enterprise_cluster.test1.id}/databases/default",
    "${azurerm_redis_enterprise_cluster.test2.id}/databases/default"
  ]

  linked_database_group_nickname = "tftestGeoGroup"
}
`, template)
}

func (r RedisenterpriseGeoDatabaseResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_redis_enterprise_geo_database" "import" {
  cluster_id = azurerm_redis_enterprise_geo_database.test.cluster_id

  client_protocol   = azurerm_redis_enterprise_geo_database.test.client_protocol
  clustering_policy = azurerm_redis_enterprise_geo_database.test.clustering_policy
  eviction_policy   = azurerm_redis_enterprise_geo_database.test.eviction_policy

  linked_database_id = [
    "${azurerm_redis_enterprise_cluster.test.id}/databases/default",
    "${azurerm_redis_enterprise_cluster.test1.id}/databases/default",
    "${azurerm_redis_enterprise_cluster.test2.id}/databases/default"
  ]
}
`, config)
}

func (r RedisenterpriseGeoDatabaseResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_redis_enterprise_geo_database" "test" {
  cluster_id = azurerm_redis_enterprise_cluster.test.id

  client_protocol            = "Encrypted"
  clustering_policy          = "EnterpriseCluster"
  eviction_policy            = "NoEviction"
  redi_search_module_enabled = true
  redi_search_module_args    = ""
  port                       = 10000

  linked_database_id = [
    "${azurerm_redis_enterprise_cluster.test.id}/databases/default",
    "${azurerm_redis_enterprise_cluster.test1.id}/databases/default",
    "${azurerm_redis_enterprise_cluster.test2.id}/databases/default",
  ]

  linked_database_group_nickname = "tftestGeoGroup"
}
`, template)
}

func (r RedisenterpriseGeoDatabaseResource) unlinkDatabase(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_redis_enterprise_geo_database" "test" {
  cluster_id = azurerm_redis_enterprise_cluster.test.id

  client_protocol   = "Encrypted"
  clustering_policy = "EnterpriseCluster"
  eviction_policy   = "NoEviction"

  linked_database_id = [
    "${azurerm_redis_enterprise_cluster.test.id}/databases/default",
    "${azurerm_redis_enterprise_cluster.test1.id}/databases/default",
  ]

  force_unlink_database_id = [
    "${azurerm_redis_enterprise_cluster.test1.id}/databases/default"
  ]
  linked_database_group_nickname = "tftestGeoGroup"
}
`, template)
}

func (r RedisenterpriseGeoDatabaseResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-redisEnterprise-%d"
  location = "%s"
}

resource "azurerm_redis_enterprise_cluster" "test" {
  name                = "acctest-geo-1"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Enterprise_E20-2"
}

resource "azurerm_redis_enterprise_cluster" "test1" {
  name                = "acctest-geo-2"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Enterprise_E20-2"
}

resource "azurerm_redis_enterprise_cluster" "test2" {
  name                = "acctest-geo-3"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Enterprise_E20-2"
}


`, data.RandomInteger, data.Locations.Primary)
}
