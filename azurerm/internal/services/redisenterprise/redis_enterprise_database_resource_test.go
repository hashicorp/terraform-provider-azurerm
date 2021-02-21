package redisenterprise_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/redisenterprise/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type RedisenterpriseDatabaseResource struct{}

func TestRedisEnterpriseDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_enterprise_database", "test")
	r := RedisenterpriseDatabaseResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestRedisEnterpriseDatabase_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_enterprise_database", "test")
	r := RedisenterpriseDatabaseResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestRedisEnterpriseDatabase_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_enterprise_database", "test")
	r := RedisenterpriseDatabaseResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestRedisEnterpriseDatabase_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_enterprise_database", "test")
	r := RedisenterpriseDatabaseResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestRedisEnterpriseDatabase_updateModules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_enterprise_database", "test")
	r := RedisenterpriseDatabaseResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateModules(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestRedisEnterpriseDatabase_updatePersistence(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_enterprise_database", "test")
	r := RedisenterpriseDatabaseResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updatePersistence(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r RedisenterpriseDatabaseResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.RedisEnterpriseDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.RedisEnterprise.DatabaseClient.Get(ctx, id.ResourceGroup, id.ClusterName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Redis Entrprise Database %q (Resource Group %q / clusterName %q): %+v", id.Name, id.ResourceGroup, id.ClusterName, err)
	}
	return utils.Bool(true), nil
}

func (r RedisenterpriseDatabaseResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-redisenterprise-%d"
  location = "%s"
}

resource "azurerm_redis_enterprise_cluster" "test" {
  name                = "acctest-rec-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku_name = "EnterpriseFlash_F300-3"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r RedisenterpriseDatabaseResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_redis_enterprise_database" "test" {
  name = "acctest-rd-%d"
  resource_group_name = azurerm_resource_group.test.name
  cluster_name = azurerm_redis_enterprise_cluster.test.name
}
`, template, data.RandomInteger)
}

func (r RedisenterpriseDatabaseResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_redis_enterprise_database" "import" {
  name = azurerm_redis_enterprise_database.test.name
  resource_group_name = azurerm_redis_enterprise_database.test.resource_group_name
  cluster_name = azurerm_redis_enterprise_database.test.cluster_name
}
`, config)
}

func (r RedisenterpriseDatabaseResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_redis_enterprise_database" "test" {
  name = "acctest-rd-%d"
  resource_group_name = azurerm_resource_group.test.name
  cluster_name = azurerm_redis_enterprise_cluster.test.name
  client_protocol = "Encrypted"
  clustering_policy = "EnterpriseCluster"
  eviction_policy = "AllKeysLRU"
  module {
    name = "RedisBloom"
    args = "ERROR_RATE 0.00 INITIAL_SIZE 400"
  }

  module {
    name = "RedisTimeSeries"
    args = "RETENTION_POLICY 20"
  }

  module {
    name = "RediSearch"
    args = ""
  }

  persistence {
    aof_enabled = true
    aof_frequency = "1s"
    rdb_enabled = false
    rdb_frequency = ""
  }
  port = 10000
}
`, template, data.RandomInteger)
}

func (r RedisenterpriseDatabaseResource) updateModules(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_redis_enterprise_database" "test" {
  name = "acctest-rd-%d"
  resource_group_name = azurerm_resource_group.test.name
  cluster_name = azurerm_redis_enterprise_cluster.test.name
  client_protocol = "Encrypted"
  clustering_policy = "EnterpriseCluster"
  eviction_policy = "AllKeysLRU"
  module {
    name = "RedisBloom"
    args = "ERROR_RATE 0.00 INITIAL_SIZE 400"
  }

  module {
    name = "RedisTimeSeries"
    args = "RETENTION_POLICY 20"
  }

  module {
    name = "RediSearch"
    args = ""
  }

  persistence {
    aof_enabled = true
    aof_frequency = "1s"
    rdb_enabled = false
    rdb_frequency = ""
  }
  port = 10000
}
`, template, data.RandomInteger)
}

func (r RedisenterpriseDatabaseResource) updatePersistence(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_redis_enterprise_database" "test" {
  name = "acctest-rd-%d"
  resource_group_name = azurerm_resource_group.test.name
  cluster_name = azurerm_redis_enterprise_cluster.test.name
  client_protocol = "Encrypted"
  clustering_policy = "EnterpriseCluster"
  eviction_policy = "AllKeysLRU"
  module {
    name = "RedisBloom"
    args = "ERROR_RATE 0.00 INITIAL_SIZE 400"
  }

  module {
    name = "RedisTimeSeries"
    args = "RETENTION_POLICY 20"
  }

  module {
    name = "RediSearch"
    args = ""
  }

  persistence {
    aof_enabled = true
    aof_frequency = "1s"
    rdb_enabled = false
    rdb_frequency = ""
  }
  port = 10000
}
`, template, data.RandomInteger)
}
