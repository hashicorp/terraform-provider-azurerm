package managedredis_test

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
`, ManagedRedisResource{}.update(data))
}
