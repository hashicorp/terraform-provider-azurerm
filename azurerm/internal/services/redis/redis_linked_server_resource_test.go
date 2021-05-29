package redis_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/redis/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type RedisLinkedServerResource struct {
}

func TestAccRedisLinkedServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_linked_server", "test")
	r := RedisLinkedServerResource{}

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

func TestAccRedisLinkedServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_linked_server", "test")
	r := RedisLinkedServerResource{}
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

func (t RedisLinkedServerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.LinkedServerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Redis.LinkedServerClient.Get(ctx, id.ResourceGroup, id.RediName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Redis Linked Server (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.LinkedServerProperties != nil), nil
}

func (RedisLinkedServerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "pri" {
  name     = "acctestRG-redis-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "pri" {
  name                = "acctestRedispri%d"
  location            = azurerm_resource_group.pri.location
  resource_group_name = azurerm_resource_group.pri.name
  capacity            = 1
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false

  redis_configuration {
    maxmemory_reserved = 2
    maxmemory_delta    = 2
    maxmemory_policy   = "allkeys-lru"
  }
}

resource "azurerm_resource_group" "sec" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "sec" {
  name                = "acctestRedissec%d"
  location            = azurerm_resource_group.sec.location
  resource_group_name = azurerm_resource_group.sec.name
  capacity            = 1
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false

  redis_configuration {
    maxmemory_reserved = 2
    maxmemory_delta    = 2
    maxmemory_policy   = "allkeys-lru"
  }
}

resource "azurerm_redis_linked_server" "test" {
  target_redis_cache_name     = azurerm_redis_cache.pri.name
  resource_group_name         = azurerm_redis_cache.pri.resource_group_name
  linked_redis_cache_id       = azurerm_redis_cache.sec.id
  linked_redis_cache_location = azurerm_redis_cache.sec.location
  server_role                 = "Secondary"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger,
		data.RandomInteger, data.Locations.Secondary, data.RandomInteger)
}

func (r RedisLinkedServerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_redis_linked_server" "import" {
  target_redis_cache_name     = azurerm_redis_linked_server.test.target_redis_cache_name
  resource_group_name         = azurerm_redis_linked_server.test.resource_group_name
  linked_redis_cache_id       = azurerm_redis_linked_server.test.linked_redis_cache_id
  linked_redis_cache_location = azurerm_redis_linked_server.test.linked_redis_cache_location
  server_role                 = azurerm_redis_linked_server.test.server_role
}
`, r.basic(data))
}
