// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/accountconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CognitiveAccountConnectionAccountManagedIdentityResource struct{}

func TestAccCognitiveAccountConnectionAccountManagedIdentity_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_account_managed_identity", "test")
	r := CognitiveAccountConnectionAccountManagedIdentityResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("metadata"),
	})
}

func TestAccCognitiveAccountConnectionAccountManagedIdentity_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_account_managed_identity", "test")
	r := CognitiveAccountConnectionAccountManagedIdentityResource{}

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

func TestAccCognitiveAccountConnectionAccountManagedIdentity_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_account_managed_identity", "test")
	r := CognitiveAccountConnectionAccountManagedIdentityResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("metadata"),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("metadata"),
	})
}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := accountconnectionresource.ParseConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cognitive.AccountConnectionResourceClient.AccountConnectionsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-ami-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cognitive_account" "test" {
  name                       = "acctest-aiservices-%[1]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "acctestaiservices-%[1]d"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, "Australia East")
}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                      = "acctkv%[3]d"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  tenant_id                 = data.azurerm_client_config.current.tenant_id
  enable_rbac_authorization = true
  sku_name                  = "standard"
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_key_vault.test.id
  role_definition_name = "Key Vault Secrets Officer"
  principal_id         = azurerm_cognitive_account.test.identity[0].principal_id
}

resource "azurerm_cognitive_account_connection_account_managed_identity" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "AzureKeyVault"
  target               = azurerm_key_vault.test.id

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_key_vault.test.id
    location   = azurerm_key_vault.test.location
  }

  depends_on = [azurerm_role_assignment.test]
}
`, r.template(data), data.RandomInteger, data.RandomInteger%1000000)
}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account_connection_account_managed_identity" "import" {
  name                 = azurerm_cognitive_account_connection_account_managed_identity.test.name
  cognitive_account_id = azurerm_cognitive_account_connection_account_managed_identity.test.cognitive_account_id
  category             = azurerm_cognitive_account_connection_account_managed_identity.test.category
  target               = azurerm_cognitive_account_connection_account_managed_identity.test.target

  metadata = {
    ApiType    = azurerm_cognitive_account_connection_account_managed_identity.test.metadata.ApiType
    ResourceId = azurerm_cognitive_account_connection_account_managed_identity.test.metadata.ResourceId
    location   = azurerm_cognitive_account_connection_account_managed_identity.test.metadata.location
  }
}
`, r.basic(data))
}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test2" {
  name                      = "acctkw%[3]d"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  tenant_id                 = data.azurerm_client_config.current.tenant_id
  enable_rbac_authorization = true
  sku_name                  = "standard"
}

resource "azurerm_role_assignment" "test2" {
  scope                = azurerm_key_vault.test2.id
  role_definition_name = "Key Vault Secrets Officer"
  principal_id         = azurerm_cognitive_account.test.identity[0].principal_id
}

resource "azurerm_cognitive_account_connection_account_managed_identity" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "AzureKeyVault"
  target               = azurerm_key_vault.test2.id

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_key_vault.test2.id
    location   = azurerm_key_vault.test2.location
  }

  depends_on = [azurerm_role_assignment.test2]
}
`, r.template(data), data.RandomInteger, data.RandomInteger%1000000)
}
