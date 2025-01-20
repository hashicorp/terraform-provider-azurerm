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

type RedisCacheAccessPolicyAssignmentResource struct{}

func TestAccRedisCacheAccessPolicyAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache_access_policy_assignment", "test")
	r := RedisCacheAccessPolicyAssignmentResource{}

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

func TestAccRedisCacheAccessPolicyAssignment_multi(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache_access_policy_assignment", "test")
	r := RedisCacheAccessPolicyAssignmentResource{}
	accessPolicyAssignmentTwo := "azurerm_redis_cache_access_policy_assignment.test2"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multi(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(accessPolicyAssignmentTwo).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		data.ImportStep(accessPolicyAssignmentTwo),
	})
}

func TestAccRedisCacheAccessPolicyAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache_access_policy_assignment", "test")
	r := RedisCacheAccessPolicyAssignmentResource{}
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

func (t RedisCacheAccessPolicyAssignmentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := redis.ParseAccessPolicyAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Redis.Redis.AccessPolicyAssignmentGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (RedisCacheAccessPolicyAssignmentResource) basic(data acceptance.TestData) string {
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

resource "azurerm_redis_cache_access_policy_assignment" "test" {
  name               = "acctestRedisAccessPolicyAssignmentTest"
  redis_cache_id     = azurerm_redis_cache.test.id
  access_policy_name = "Data Contributor"
  object_id          = data.azurerm_client_config.test.object_id
  object_id_alias    = "ServicePrincipal"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r RedisCacheAccessPolicyAssignmentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_redis_cache_access_policy_assignment" "import" {
  name               = azurerm_redis_cache_access_policy_assignment.test.name
  redis_cache_id     = azurerm_redis_cache_access_policy_assignment.test.redis_cache_id
  access_policy_name = azurerm_redis_cache_access_policy_assignment.test.access_policy_name
  object_id          = azurerm_redis_cache_access_policy_assignment.test.object_id
  object_id_alias    = azurerm_redis_cache_access_policy_assignment.test.object_id_alias
}
`, r.basic(data))
}

func (r RedisCacheAccessPolicyAssignmentResource) multi(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azuread" {}

resource "azurerm_redis_cache_access_policy" "test2" {
  name           = "acctestRedisAccessPolicytest2"
  redis_cache_id = azurerm_redis_cache.test.id
  permissions    = "+@read +@connection +cluster|info allkeys"
}

resource "azuread_group" "test2" {
  display_name     = "acctestredis"
  security_enabled = true
}

resource "azurerm_redis_cache_access_policy_assignment" "test2" {
  name               = "acctestRedisAccessPolicyAssignmentTest2"
  redis_cache_id     = azurerm_redis_cache.test.id
  access_policy_name = "Data Contributor"
  object_id          = azuread_group.test2.id
  object_id_alias    = "Group"
}
`, r.basic(data))
}
