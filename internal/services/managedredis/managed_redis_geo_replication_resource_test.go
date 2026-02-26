// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managedredis_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const defaultDatabaseName = "default"

type ManagedRedisGeoReplicationResource struct{}

func TestAccManagedRedisGeoReplication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_geo_replication", "test")
	r := ManagedRedisGeoReplicationResource{}
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

func TestAccManagedRedisGeoReplication_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_geo_replication", "test")
	r := ManagedRedisGeoReplicationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.threeClustersOneAndTwoLinked(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("linked_managed_redis_ids.0").MatchesRegex(regexp.MustCompile(`redisEnterprise/acctest-amr2`)),
			),
		},
		data.ImportStep(),
		{
			Config: r.threeClustersOneAndThreeLinked(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("linked_managed_redis_ids.0").MatchesRegex(regexp.MustCompile(`redisEnterprise/acctest-amr3`)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedRedisGeoReplication_fourClustersLinkUnlink(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_geo_replication", "test")
	r := ManagedRedisGeoReplicationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		// refresh step
		{
			Config: r.fourClustersAllLinked(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("linked_managed_redis_ids.#").HasValue("3"),
			),
		},
		data.ImportStep(),
		// refresh step
		{
			Config: r.fourClustersOneAndTwoLinked(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("linked_managed_redis_ids.0").MatchesRegex(regexp.MustCompile(`redisEnterprise/acctest-amr2`)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedRedisGeoReplication_linkedManagedRedisIdsCannotContainSelf(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_geo_replication", "test")
	r := ManagedRedisGeoReplicationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.linkedManagedRedisIdsContainSelf(data),
			ExpectError: regexp.MustCompile("linked_managed_redis_ids cannot contain the same value as managed_redis_id"),
		},
	})
}

func (r ManagedRedisGeoReplicationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	clusterId, err := redisenterprise.ParseRedisEnterpriseID(state.ID)
	if err != nil {
		return nil, err
	}

	dbId := databases.NewDatabaseID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, defaultDatabaseName)

	resp, err := client.ManagedRedis.DatabaseClient.Get(ctx, dbId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", dbId, err)
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if geoProps := props.GeoReplication; geoProps != nil {
				return pointer.To(len(pointer.From(geoProps.LinkedDatabases)) > 0), nil
			}
		}
	}

	return pointer.To(false), nil
}

func (r ManagedRedisGeoReplicationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-managedRedis-%[1]d"
  location = "%[2]s"
}

resource "azurerm_managed_redis" "amr1" {
  name                = "acctest-amr1-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[2]s"
  sku_name            = "Balanced_B3"

  default_database {
    geo_replication_group_name = "acctest-georep-%[1]d"
  }
}

resource "azurerm_managed_redis" "amr2" {
  name                = "acctest-amr2-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  sku_name            = "Balanced_B3"

  default_database {
    geo_replication_group_name = "acctest-georep-%[1]d"
  }
}

resource "azurerm_managed_redis" "amr3" {
  name                = "acctest-amr3-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[4]s"
  sku_name            = "Balanced_B3"

  default_database {
    geo_replication_group_name = "acctest-georep-%[1]d"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary, data.Locations.Ternary)
}

func (r ManagedRedisGeoReplicationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_managed_redis_geo_replication" "test" {
  managed_redis_id = azurerm_managed_redis.amr1.id
  linked_managed_redis_ids = [
    azurerm_managed_redis.amr2.id,
    azurerm_managed_redis.amr3.id,
  ]
}
`, r.template(data))
}

func (r ManagedRedisGeoReplicationResource) threeClustersOneAndTwoLinked(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_managed_redis_geo_replication" "test" {
  managed_redis_id = azurerm_managed_redis.amr1.id
  linked_managed_redis_ids = [
    azurerm_managed_redis.amr2.id,
  ]
}
`, r.template(data))
}

func (r ManagedRedisGeoReplicationResource) threeClustersOneAndThreeLinked(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_managed_redis_geo_replication" "test" {
  managed_redis_id = azurerm_managed_redis.amr1.id
  linked_managed_redis_ids = [
    azurerm_managed_redis.amr3.id,
  ]
}
`, r.template(data))
}

func (r ManagedRedisGeoReplicationResource) fourClustersAllLinked(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_managed_redis" "amr4" {
  name                = "acctest-amr4-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "westeurope"
  sku_name            = "Balanced_B3"

  default_database {
    geo_replication_group_name = "acctest-georep-%[2]d"
  }
}

resource "azurerm_managed_redis_geo_replication" "test" {
  managed_redis_id = azurerm_managed_redis.amr1.id
  linked_managed_redis_ids = [
    azurerm_managed_redis.amr2.id,
    azurerm_managed_redis.amr3.id,
    azurerm_managed_redis.amr4.id,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r ManagedRedisGeoReplicationResource) fourClustersOneAndTwoLinked(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_managed_redis" "amr4" {
  name                = "acctest-amr4-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "westeurope"
  sku_name            = "Balanced_B3"

  default_database {
    geo_replication_group_name = "acctest-georep-%[2]d"
  }
}

resource "azurerm_managed_redis_geo_replication" "test" {
  managed_redis_id = azurerm_managed_redis.amr1.id
  linked_managed_redis_ids = [
    azurerm_managed_redis.amr2.id,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r ManagedRedisGeoReplicationResource) linkedManagedRedisIdsContainSelf(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-managedRedis-%[1]d"
  location = "%[2]s"
}

resource "azurerm_managed_redis" "amr1" {
  name                = "acctest-amr1-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[2]s"
  sku_name            = "Balanced_B3"

  default_database {
    geo_replication_group_name = "acctest-georep-%[1]d"
  }
}

resource "azurerm_managed_redis" "amr2" {
  name                = "acctest-amr2-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  sku_name            = "Balanced_B3"

  default_database {
    geo_replication_group_name = "acctest-georep-%[1]d"
  }
}

resource "azurerm_managed_redis_geo_replication" "test" {
  managed_redis_id = azurerm_managed_redis.amr1.id
  linked_managed_redis_ids = [
    azurerm_managed_redis.amr1.id,
    azurerm_managed_redis.amr2.id,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}
