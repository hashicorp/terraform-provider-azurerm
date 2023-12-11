// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/cache"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApimanagementRedisCacheResource struct{}

func TestAccApiManagementRedisCache_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_redis_cache", "test")
	r := ApimanagementRedisCacheResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("connection_string"),
	})
}

func TestAccApiManagementRedisCache_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_redis_cache", "test")
	r := ApimanagementRedisCacheResource{}
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

func TestAccApiManagementRedisCache_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_redis_cache", "test")
	r := ApimanagementRedisCacheResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("connection_string"),
	})
}

func TestAccApiManagementRedisCache_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_redis_cache", "test")
	r := ApimanagementRedisCacheResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("connection_string"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("connection_string"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("connection_string"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("connection_string"),
	})
}

func (r ApimanagementRedisCacheResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := cache.ParseCacheID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.ApiManagement.CacheClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %q %+v", id, err)
	}
	return pointer.To(true), nil
}

func (r ApimanagementRedisCacheResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-apim-%[1]d"
  location = "%[2]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%[1]d"
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

resource "azurerm_redis_cache" "test2" {
  name                = "acctestRedis2-%[1]d"
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
`, data.RandomInteger, data.Locations.Primary)
}

func (r ApimanagementRedisCacheResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_redis_cache" "test" {
  name              = "acctest-Redis-Cache-%d"
  api_management_id = azurerm_api_management.test.id
  connection_string = azurerm_redis_cache.test.primary_connection_string
}
`, r.template(data), data.RandomInteger)
}

func (r ApimanagementRedisCacheResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_redis_cache" "import" {
  name              = azurerm_api_management_redis_cache.test.name
  api_management_id = azurerm_api_management.test.id
  connection_string = azurerm_redis_cache.test.primary_connection_string
}
`, r.basic(data))
}

func (r ApimanagementRedisCacheResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_redis_cache" "test" {
  name              = "acctest-Redis-Cache-%d"
  api_management_id = azurerm_api_management.test.id
  connection_string = azurerm_redis_cache.test.primary_connection_string
  description       = "Redis cache instances"
  redis_cache_id    = azurerm_redis_cache.test.id
  cache_location    = "%s"
}
`, template, data.RandomInteger, data.Locations.Secondary)
}

func (r ApimanagementRedisCacheResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_redis_cache" "test" {
  name              = "acctest-Redis-Cache-%d"
  api_management_id = azurerm_api_management.test.id
  connection_string = azurerm_redis_cache.test2.primary_connection_string
  description       = "Redis cache Update"
  redis_cache_id    = azurerm_redis_cache.test2.id
  cache_location    = "%s"
}
`, template, data.RandomInteger, data.Locations.Ternary)
}
