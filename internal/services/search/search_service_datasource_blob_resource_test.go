// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package search_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/datasources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SearchServiceDatasourceBlobResource struct{}

func TestAccSearchServiceDatasourceBlob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_search_service_datasource_blob", "test")
	r := SearchServiceDatasourceBlobResource{}

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

func TestAccSearchServiceDatasourceBlob_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_search_service_datasource_blob", "test")
	r := SearchServiceDatasourceBlobResource{}

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

func TestAccSearchServiceDatasourceBlob_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_search_service_datasource_blob", "test")
	r := SearchServiceDatasourceBlobResource{}

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

func TestAccSearchServiceDatasourceBlob_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_search_service_datasource_blob", "test")
	r := SearchServiceDatasourceBlobResource{}

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
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("connection_string"),
	})
}

func TestAccSearchServiceDatasourceBlob_encryptionKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_search_service_datasource_blob", "test")
	r := SearchServiceDatasourceBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withEncryptionKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("connection_string"),
	})
}

func TestAccSearchServiceDatasourceBlob_encryptionKeyWithAppCredentials(t *testing.T) {
	appClientId := os.Getenv("ARM_TEST_SEARCH_APP_CLIENT_ID")
	appClientSecret := os.Getenv("ARM_TEST_SEARCH_APP_CLIENT_SECRET")
	if appClientId == "" || appClientSecret == "" {
		t.Skip("Skipping: ARM_TEST_SEARCH_APP_CLIENT_ID and ARM_TEST_SEARCH_APP_CLIENT_SECRET must be set")
	}

	data := acceptance.BuildTestData(t, "azurerm_search_service_datasource_blob", "test")
	r := SearchServiceDatasourceBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withEncryptionKeyAndAppCredentials(data, appClientId, appClientSecret),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("connection_string", "encryption_key.0.application_secret"),
	})
}

func (r SearchServiceDatasourceBlobResource) Exists(ctx context.Context, c *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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
		return nil, fmt.Errorf("%s could not be retrieved: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (SearchServiceDatasourceBlobResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-search-%[1]d"
  location = "%[2]s"
}

resource "azurerm_search_service" "test" {
  name                = "acctestsearch%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "standard"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name               = "acctestsc%[1]d"
  storage_account_id = azurerm_storage_account.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r SearchServiceDatasourceBlobResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_search_service_datasource_blob" "test" {
  name              = "acctestds%d"
  search_service_id = azurerm_search_service.test.id
  container_name    = azurerm_storage_container.test.name
  connection_string = azurerm_storage_account.test.primary_connection_string
}
`, template, data.RandomInteger)
}

func (r SearchServiceDatasourceBlobResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_search_service_datasource_blob" "import" {
  name              = "acctestds%d"
  search_service_id = azurerm_search_service.test.id
  container_name    = azurerm_storage_container.test.name
  connection_string = azurerm_storage_account.test.primary_connection_string
}
`, template, data.RandomInteger)
}

func (r SearchServiceDatasourceBlobResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_search_service_datasource_blob" "test" {
  name                     = "acctestds%d"
  search_service_id        = azurerm_search_service.test.id
  container_name           = azurerm_storage_container.test.name
  container_query          = "/folder"
  connection_string        = azurerm_storage_account.test.primary_connection_string
  description              = "test description"
  soft_delete_column_name  = "IsDeleted"
  soft_delete_marker_value = "true"
}
`, template, data.RandomInteger)
}

func (r SearchServiceDatasourceBlobResource) encryptionKeyTemplate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv%s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  purge_protection_enabled   = true
  enable_rbac_authorization  = true
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
  name         = "acctestkvkey%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = [azurerm_role_assignment.current_user_keyvault]
}
`, template, data.RandomString, data.RandomString)
}

func (r SearchServiceDatasourceBlobResource) withEncryptionKey(data acceptance.TestData) string {
	template := r.encryptionKeyTemplate(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_search_service_datasource_blob" "test" {
  name              = "acctestds%d"
  search_service_id = azurerm_search_service.test.id
  container_name    = azurerm_storage_container.test.name
  connection_string = azurerm_storage_account.test.primary_connection_string

  encryption_key {
    key_name      = azurerm_key_vault_key.test.name
    key_version   = azurerm_key_vault_key.test.version
    key_vault_uri = azurerm_key_vault.test.vault_uri
  }

  depends_on = [azurerm_role_assignment.search_keyvault]
}
`, template, data.RandomInteger)
}

func (r SearchServiceDatasourceBlobResource) withEncryptionKeyAndAppCredentials(data acceptance.TestData, appClientId, appClientSecret string) string {
	template := r.encryptionKeyTemplate(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_role_assignment" "app_keyvault" {
  scope                = azurerm_key_vault.test.id
  role_definition_name = "Key Vault Crypto Service Encryption User"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_search_service_datasource_blob" "test" {
  name              = "acctestds%d"
  search_service_id = azurerm_search_service.test.id
  container_name    = azurerm_storage_container.test.name
  connection_string = azurerm_storage_account.test.primary_connection_string

  encryption_key {
    key_name           = azurerm_key_vault_key.test.name
    key_version        = azurerm_key_vault_key.test.version
    key_vault_uri      = azurerm_key_vault.test.vault_uri
    application_id     = %q
    application_secret = %q
  }

  depends_on = [azurerm_role_assignment.search_keyvault, azurerm_role_assignment.app_keyvault]
}
`, template, data.RandomInteger, appClientId, appClientSecret)
}
