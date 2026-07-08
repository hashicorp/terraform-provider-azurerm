// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/accountconnectionresource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("metadata"),
	})
}

func TestAccCognitiveAccountConnectionAccountManagedIdentity_aiServicesCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_account_managed_identity", "test")
	r := CognitiveAccountConnectionAccountManagedIdentityResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.aiServicesCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionAccountManagedIdentity_azureOpenAICategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_account_managed_identity", "test")
	r := CognitiveAccountConnectionAccountManagedIdentityResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureOpenAICategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionAccountManagedIdentity_azureKeyVaultCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_account_managed_identity", "test")
	r := CognitiveAccountConnectionAccountManagedIdentityResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureKeyVaultCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionAccountManagedIdentity_storageAccountCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_account_managed_identity", "test")
	r := CognitiveAccountConnectionAccountManagedIdentityResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccountCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionAccountManagedIdentity_importMismatchedAuthTypeFails(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_account_managed_identity", "test")
	r := CognitiveAccountConnectionAccountManagedIdentityResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.importMismatchedAuthType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			ResourceName: data.ResourceName,
			ImportState:  true,
			ImportStateIdFunc: func(s *terraform.State) (string, error) {
				rs, ok := s.RootModule().Resources["azurerm_cognitive_account_connection_entra_id.wrong"]
				if !ok {
					return "", fmt.Errorf("resource `%s` not found in state", "azurerm_cognitive_account_connection_entra_id.wrong")
				}
				return rs.Primary.ID, nil
			},
			ExpectError: regexp.MustCompile("cannot be managed by"),
		},
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
  name                       = "acctest-cogacc-%[1]d"
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
`, data.RandomInteger, data.Locations.Primary)
}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account" "aiservices" {
  name                       = "acctest-cogacc-2-%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "acctestaisvc2-%[2]d"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account_connection_account_managed_identity" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "AIServices"
  target               = azurerm_cognitive_account.aiservices.endpoint

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_cognitive_account.aiservices.id
    Location   = azurerm_cognitive_account.aiservices.location
  }
}
`, r.template(data), data.RandomInteger)
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
    Location   = azurerm_cognitive_account_connection_account_managed_identity.test.metadata.Location
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

resource "azurerm_cognitive_account" "aiservices2" {
  name                       = "acctest-cogacc-ai2-%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "acctestaisvc3-%[2]d"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account_connection_account_managed_identity" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "AIServices"
  target               = azurerm_cognitive_account.aiservices2.endpoint

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_cognitive_account.aiservices2.id
    Location   = azurerm_cognitive_account.aiservices2.location
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) aiServicesCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account" "aiservices" {
  name                       = "acctest-cogacc-ai-%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "acctestaisvc2-%[2]d"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account_connection_account_managed_identity" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "AIServices"
  target               = azurerm_cognitive_account.aiservices.endpoint

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_cognitive_account.aiservices.id
    Location   = azurerm_cognitive_account.aiservices.location
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) azureOpenAICategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account" "openai" {
  name                = "acctest-cogacc-openai-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "OpenAI"
  sku_name            = "S0"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account_connection_account_managed_identity" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "AzureOpenAI"
  target               = azurerm_cognitive_account.openai.endpoint

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_cognitive_account.openai.id
    Location   = azurerm_cognitive_account.openai.location
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) azureKeyVaultCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                      = "acctkv%[3]s"
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
  target               = azurerm_key_vault.test.vault_uri

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_key_vault.test.id
    Location   = azurerm_key_vault.test.location
  }

  depends_on = [azurerm_role_assignment.test]
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) storageAccountCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctesta%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_cognitive_account_connection_account_managed_identity" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "AzureStorageAccount"
  target               = azurerm_storage_account.test.primary_blob_endpoint

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_storage_account.test.id
    Location   = azurerm_storage_account.test.location
  }
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) importMismatchedAuthType(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account_connection_entra_id" "wrong" {
  name                 = "acctest-wrong-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "Databricks"
  target               = "https://wrong-%[2]d.databricks.net/"
}
`, r.basic(data), data.RandomInteger)
}
