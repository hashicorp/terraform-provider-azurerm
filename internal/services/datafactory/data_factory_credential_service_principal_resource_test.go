// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/credentials"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CredentialServicePrincipalResource struct{}

func TestAccDataFactoryCredentialServicePrincipal_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_credential_service_principal", "test")
	r := CredentialServicePrincipalResource{}

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

func TestAccDataFactoryCredentialServicePrincipal_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_credential_service_principal", "test")
	r := CredentialServicePrincipalResource{}
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

func TestAccDataFactoryCredentialServicePrincipal_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_credential_service_principal", "test")
	r := CredentialServicePrincipalResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryCredentialServicePrincipal_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_credential_service_principal", "test")
	r := CredentialServicePrincipalResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
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

func (r CredentialServicePrincipalResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := credentials.ParseCredentialID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataFactory.Credentials.CredentialOperationsGet(ctx, *id, credentials.DefaultCredentialOperationsGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r CredentialServicePrincipalResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%[2]d"
  location = "%[1]s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_key_vault" "test" {
  name                       = "kv%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Get",
    ]

    secret_permissions = [
      "Set",
      "Get",
      "Delete",
      "Purge",
      "Recover"
    ]
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "acctestkvsecret"
  value        = "fakedsecret"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_data_factory_linked_service_key_vault" "test" {
  name            = "acctestlskv%[2]d"
  data_factory_id = azurerm_data_factory.test.id
  key_vault_id    = azurerm_key_vault.test.id
}

resource "azurerm_data_factory_linked_service_key_vault" "test2" {
  name            = "acctestlskv2%[2]d"
  data_factory_id = azurerm_data_factory.test.id
  key_vault_id    = azurerm_key_vault.test.id
}

resource "azurerm_key_vault_secret" "test2" {
  name         = "anothersecret"
  value        = "fakedsecret"
  key_vault_id = azurerm_key_vault.test.id
}
`, data.Locations.Primary, data.RandomInteger)
}

func (r CredentialServicePrincipalResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_credential_service_principal" "test" {
  name                 = "credential%d"
  data_factory_id      = azurerm_data_factory.test.id
  tenant_id            = data.azurerm_client_config.current.tenant_id
  service_principal_id = data.azurerm_client_config.current.object_id
  service_principal_key {
    linked_service_name = azurerm_data_factory_linked_service_key_vault.test.name
    secret_name         = azurerm_key_vault_secret.test.name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CredentialServicePrincipalResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_credential_service_principal" "import" {
  name                 = azurerm_data_factory_credential_service_principal.test.name
  data_factory_id      = azurerm_data_factory_credential_service_principal.test.data_factory_id
  tenant_id            = data.azurerm_client_config.current.tenant_id
  service_principal_id = data.azurerm_client_config.current.object_id
  service_principal_key {
    linked_service_name = azurerm_data_factory_linked_service_key_vault.test.name
    secret_name         = azurerm_key_vault_secret.test.name
  }
}
`, config)
}

func (r CredentialServicePrincipalResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_factory_credential_service_principal" "test" {
  name                 = "credential%[2]d"
  description          = "UPDATED DESCRIPTION"
  data_factory_id      = azurerm_data_factory.test.id
  tenant_id            = "00000000-0000-0000-0000-000000000000"
  service_principal_id = "00000000-0000-0000-0000-000000000000"
  service_principal_key {
    linked_service_name = azurerm_data_factory_linked_service_key_vault.test2.name
    secret_name         = azurerm_key_vault_secret.test2.name
    secret_version      = azurerm_key_vault_secret.test2.version
  }
  annotations = ["1", "2"]
}
`, r.template(data), data.RandomInteger)
}
