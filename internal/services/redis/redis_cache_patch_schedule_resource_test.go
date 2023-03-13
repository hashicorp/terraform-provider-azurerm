package redis_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2022-06-01/patchschedules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type RedisCachePatchScheduleResource struct{}

func TestAccRedisCachePatchSchedule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache_patch_schedule", "test")
	r := RedisCachePatchScheduleResource{}

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

func TestAccRedisCachePatchSchedule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache_patch_schedule", "test")
	r := RedisCachePatchScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.patchSchedule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccRedisCachePatchSchedule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache_patch_schedule", "test")
	r := RedisCachePatchScheduleResource{}

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

func (r RedisCachePatchScheduleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := patchschedules.ParseRediID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Redis.PatchSchedules.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r RedisCachePatchScheduleResource) basic(data acceptance.TestData) string {
	temp := r.template(data)
	return fmt.Sprintf(`
%s

resource azurerm_redis_cache_patch_schedule test {
  redis_cache_id = azurerm_redis_cache.test.id
  patch_schedule {
	day_of_week        = "Tuesday"
	start_hour_utc     = 8
	maintenance_window = "PT7H"
  }
}

`, temp)
}

func (r RedisCachePatchScheduleResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_redis_cache_patch_schedule" "import" {
  redis_cache_id = azurerm_redis_cache_patch_schedule.test.redis_cache_id
  patch_schedule {
    day_of_week        = "Tuesday"
    start_hour_utc     = 8
    maintenance_window = "PT7H"
  }
}
`, template)
}

func (r RedisCachePatchScheduleResource) patchSchedule(data acceptance.TestData) string {
	temp := r.template(data)
	return fmt.Sprintf(`%s

resource azurerm_redis_cache_patch_schedule test { 
  redis_cache_id = azurerm_redis_cache.test.id

  patch_schedule {
    day_of_week        = "Tuesday"
    start_hour_utc     =    20
  }

  patch_schedule {
    day_of_week        = "Wednesday"
    start_hour_utc     = 9
  }
}

`, temp)
}

func (RedisCachePatchScheduleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 1
  family              = "C"
  sku_name            = "Basic"
  enable_non_ssl_port = false
  minimum_tls_version = "1.2"

  redis_configuration {
  }

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
