// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/datasources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SearchServiceDatasourceDataLakeStorageGen2Resource struct{}

func TestAccSearchServiceDatasourceDataLakeStorageGen2_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_search_service_datasource_data_lake_storage_gen2", "test")
	r := SearchServiceDatasourceDataLakeStorageGen2Resource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("connection_string_wo_version"),
	})
}

func TestAccSearchServiceDatasourceDataLakeStorageGen2_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_search_service_datasource_data_lake_storage_gen2", "test")
	r := SearchServiceDatasourceDataLakeStorageGen2Resource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSearchServiceDatasourceDataLakeStorageGen2_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_search_service_datasource_data_lake_storage_gen2", "test")
	r := SearchServiceDatasourceDataLakeStorageGen2Resource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("connection_string_wo_version"),
	})
}

func TestAccSearchServiceDatasourceDataLakeStorageGen2_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_search_service_datasource_data_lake_storage_gen2", "test")
	r := SearchServiceDatasourceDataLakeStorageGen2Resource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("connection_string_wo_version"),
		{
			Config: r.updated(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("connection_string_wo_version"),
		{
			Config: r.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("connection_string_wo_version"),
	})
}

func (r SearchServiceDatasourceDataLakeStorageGen2Resource) Exists(ctx context.Context, c *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := datasources.ParseDatasourceID(state.ID)
	if err != nil {
		return nil, err
	}

	client := c.Search.SearchDataPlaneClient.DataSources.Clone(id.BaseURI)
	resp, err := client.Get(ctx, *id, datasources.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if resp.Model == nil {
		return nil, fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if resp.Model.Type != datasources.SearchIndexerDataSourceTypeAdlsgenTwo {
		return nil, fmt.Errorf("%s has type `%s`, expected `adlsgen2`", id, resp.Model.Type)
	}
	return pointer.To(true), nil
}

func (r SearchServiceDatasourceDataLakeStorageGen2Resource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-search-%[1]d"
  location = "%[2]s"
}

resource "azurerm_search_service" "test" {
  name                        = "acctestsearch%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku                         = "standard"
  authentication_failure_mode = "http403"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_role_assignment" "current_user_search" {
  scope                = azurerm_search_service.test.id
  role_definition_name = "Search Service Contributor"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  is_hns_enabled           = true
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctestfs%[1]d"
  storage_account_id = azurerm_storage_account.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r SearchServiceDatasourceDataLakeStorageGen2Resource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_search_service_datasource_data_lake_storage_gen2" "test" {
  name                         = "acctestds%d"
  search_service_endpoint      = azurerm_search_service.test.endpoint
  filesystem_name              = azurerm_storage_data_lake_gen2_filesystem.test.name
  connection_string_wo         = azurerm_storage_account.test.primary_connection_string
  connection_string_wo_version = 1

  depends_on = [azurerm_role_assignment.current_user_search]
}
`, r.template(data), data.RandomInteger)
}

func (r SearchServiceDatasourceDataLakeStorageGen2Resource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_search_service_datasource_data_lake_storage_gen2" "import" {
  name                         = azurerm_search_service_datasource_data_lake_storage_gen2.test.name
  search_service_endpoint      = azurerm_search_service.test.endpoint
  filesystem_name              = azurerm_storage_data_lake_gen2_filesystem.test.name
  connection_string_wo         = azurerm_storage_account.test.primary_connection_string
  connection_string_wo_version = 1
}
`, r.basic(data))
}

func (r SearchServiceDatasourceDataLakeStorageGen2Resource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_data_lake_gen2_filesystem" "updated" {
  name               = "acctestfs%[2]d-updated"
  storage_account_id = azurerm_storage_account.test.id
}

resource "azurerm_search_service_datasource_data_lake_storage_gen2" "test" {
  name                         = "acctestds%[2]d"
  search_service_endpoint      = azurerm_search_service.test.endpoint
  filesystem_name              = azurerm_storage_data_lake_gen2_filesystem.updated.name
  container_query              = "documents"
  connection_string_wo         = azurerm_storage_account.test.secondary_connection_string
  connection_string_wo_version = 2
  description                  = "updated description"
  soft_delete_column_name      = "IsDeleted"
  soft_delete_marker_value     = "true"

  depends_on = [azurerm_role_assignment.current_user_search]
}
`, r.template(data), data.RandomInteger)
}

func (r SearchServiceDatasourceDataLakeStorageGen2Resource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv%[2]s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  purge_protection_enabled   = true
  rbac_authorization_enabled = true
}

resource "azurerm_role_assignment" "current_user_keyvault" {
  scope                = azurerm_key_vault.test.id
  role_definition_name = "Key Vault Administrator"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_role_assignment" "search_keyvault" {
  scope                = azurerm_key_vault.test.id
  role_definition_name = "Key Vault Crypto Service Encryption User"
  principal_id         = azurerm_search_service.test.identity[0].principal_id
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctestkvkey%[2]s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [azurerm_role_assignment.current_user_keyvault]
}

resource "azurerm_search_service_datasource_data_lake_storage_gen2" "test" {
  name                         = "acctestds%[3]d"
  search_service_endpoint      = azurerm_search_service.test.endpoint
  filesystem_name              = azurerm_storage_data_lake_gen2_filesystem.test.name
  container_query              = "documents"
  connection_string_wo         = azurerm_storage_account.test.primary_connection_string
  connection_string_wo_version = 1
  description                  = "complete description"
  soft_delete_column_name      = "IsDeleted"
  soft_delete_marker_value     = "true"

  encryption_key {
    key_vault_key_id = azurerm_key_vault_key.test.id
  }

  depends_on = [azurerm_role_assignment.current_user_search, azurerm_role_assignment.search_keyvault]
}
`, r.template(data), data.RandomString, data.RandomInteger)
}
