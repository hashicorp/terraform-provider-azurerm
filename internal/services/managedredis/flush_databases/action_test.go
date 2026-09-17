// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package flush_databases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

type ManagedRedisFlushDatabasesAction struct{}

func TestAccManagedRedisFlushDatabasesAction_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_databases_flush", "test")
	a := ManagedRedisFlushDatabasesAction{}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: a.basic(data),
				Check:  nil, // TODO
			},
		},
	})
}

func TestAccManagedRedisFlushDatabasesAction_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_databases_flush", "test")
	a := ManagedRedisFlushDatabasesAction{}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: a.complete(data),
				Check:  nil, // TODO
			},
		},
	})
}

func (r *ManagedRedisFlushDatabasesAction) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

action "azurerm_managed_redis_databases_flush" "test" {
  config {
    managed_redis_database_id = azurerm_managed_redis.test.default_database[0].id
  }
}
`, r.template(data))
}

func (r *ManagedRedisFlushDatabasesAction) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

action "azurerm_managed_redis_databases_flush" "test" {
  config {
    managed_redis_database_id = azurerm_managed_redis.test.default_database[0].id
    linked_database_ids       = [azurerm_managed_redis.test.default_database[0].id]
  }
}
`, r.template(data))
}

func (r *ManagedRedisFlushDatabasesAction) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "terraform_data" "trigger" {
  input = azurerm_managed_redis.test.id
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.azurerm_managed_redis_databases_flush.test]
    }
  }
}
`, testClusterHcl(data, "Balanced_B3"))
}

func testClusterHcl(data acceptance.TestData, skuName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-managedRedis-%[1]d"
  location = "%[2]s"
}

resource "azurerm_managed_redis" "test" {
  name                = "acctest-amr-%[1]d"
  resource_group_name = azurerm_resource_group.test.name

  location = "%[2]s"
  sku_name = "%[3]s"

  default_database {
    access_keys_authentication_enabled = false
    client_protocol                    = "Plaintext"
    clustering_policy                  = "EnterpriseCluster"
    eviction_policy                    = "NoEviction"
    geo_replication_group_name         = "acctest-amr-georep-%[1]d"

    module {
      name = "RediSearch"
      args = ""
    }

    module {
      name = "RedisJSON"
      args = ""
    }
  }

  high_availability_enabled = true

  tags = {
    ENV    = "Test",
    Method = "Update"
  }
}
`, data.RandomInteger, data.Locations.Primary, skuName)
}
