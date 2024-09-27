// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package redis_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-03-01/redis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type RedisCacheAccessPolicyResource struct{}

func TestAccRedisCacheAccessPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache_access_policy", "test")
	r := RedisCacheAccessPolicyResource{}

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

func TestAccRedisCacheAccessPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache_access_policy", "test")
	r := RedisCacheAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedisCacheAccessPolicy_multi(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache_access_policy", "test")
	r := RedisCacheAccessPolicyResource{}
	accessPolicyTwo := "azurerm_redis_cache_access_policy.test2"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multi(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(accessPolicyTwo).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		data.ImportStep(accessPolicyTwo),
	})
}

func TestAccRedisCacheAccessPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache_access_policy", "test")
	r := RedisCacheAccessPolicyResource{}
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

func (t RedisCacheAccessPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := redis.ParseAccessPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Redis.Redis.AccessPolicyGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (RedisCacheAccessPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "test" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "test" {
  name                 = "acctestRedis-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  capacity             = 1
  family               = "C"
  sku_name             = "Basic"
  non_ssl_port_enabled = true
  minimum_tls_version  = "1.2"

  redis_configuration {
  }
}

resource "azurerm_redis_cache_access_policy" "test" {
  name           = "acctestRedisAccessPolicytest"
  redis_cache_id = azurerm_redis_cache.test.id
  permissions    = "+@read +@connection +cluster|info allkeys"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RedisCacheAccessPolicyResource) basicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "test" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "test" {
  name                 = "acctestRedis-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  capacity             = 1
  family               = "C"
  sku_name             = "Basic"
  non_ssl_port_enabled = true
  minimum_tls_version  = "1.2"

  redis_configuration {
  }
}

resource "azurerm_redis_cache_access_policy" "test" {
  name           = "acctestRedisAccessPolicytest"
  redis_cache_id = azurerm_redis_cache.test.id
  permissions    = "+@read +@write +@connection +cluster|info allkeys"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r RedisCacheAccessPolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_redis_cache_access_policy" "import" {
  name           = azurerm_redis_cache_access_policy.test.name
  redis_cache_id = azurerm_redis_cache_access_policy.test.redis_cache_id
  permissions    = azurerm_redis_cache_access_policy.test.permissions
}
`, r.basic(data))
}

func (r RedisCacheAccessPolicyResource) multi(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_redis_cache_access_policy" "test2" {
  name           = "acctestRedisAccessPolicytest2"
  redis_cache_id = azurerm_redis_cache.test.id
  permissions    = "+@read +@connection +cluster|info allkeys"
}
`, r.basic(data))
}
