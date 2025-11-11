// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitive_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/accountconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CognitiveAccountConnectionTestResource struct{}

func TestAccCognitiveAccountConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("api_key", "metadata"),
	})
}

func TestAccCognitiveAccountConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

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

func TestAccCognitiveAccountConnection_apiKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.apiKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("api_key", "metadata"),
	})
}

func TestAccCognitiveAccountConnection_oauth2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.oauth2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("oauth2", "metadata"),
	})
}

func TestAccCognitiveAccountConnection_customKeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customKeys(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_keys", "metadata"),
	})
}

func TestAccCognitiveAccountConnection_AAD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.AAD(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("metadata"),
	})
}

func TestAccCognitiveAccountConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("api_key", "metadata"),
		{
			Config: r.apiKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("api_key", "metadata"),
	})
}

func TestAccCognitiveAccountConnection_updateStorageBlob(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.oauth2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("oauth2", "metadata"),
		{
			Config: r.oauth2Updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("oauth2", "metadata"),
		{
			Config: r.AAD(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("metadata"),
	})
}

func TestAccCognitiveAccountConnection_authTypeValidation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection", "test")
	r := CognitiveAccountConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.invalidAuthTypeApiKeyMismatch(data),
			ExpectError: regexp.MustCompile("when `auth_type` is `ApiKey`, `api_key` must be specified"),
		},
		{
			Config:      r.invalidAuthTypeAADWithOtherAuth(data),
			ExpectError: regexp.MustCompile("when `auth_type` is `AAD`, no other auth configuration blocks should be specified"),
		},
	})
}

func (r CognitiveAccountConnectionTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := accountconnectionresource.ParseConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cognitive.AccountConnectionResourceClient.AccountConnectionsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r CognitiveAccountConnectionTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
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

resource "azurerm_cognitive_account" "openai" {
  name                = "acctest-openai-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "OpenAI"
  sku_name            = "S0"

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
  name                  = "acctestcsc%[3]s"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r CognitiveAccountConnectionTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type            = "ApiKey"
  category             = "AzureOpenAI"
  target               = azurerm_cognitive_account.openai.endpoint
  api_key              = azurerm_cognitive_account.openai.primary_access_key

  metadata = {
    apiType    = "Azure"
    resourceId = azurerm_cognitive_account.openai.id
    location   = azurerm_cognitive_account.openai.location
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account_connection" "import" {
  name                 = azurerm_cognitive_account_connection.test.name
  cognitive_account_id = azurerm_cognitive_account_connection.test.cognitive_account_id
  auth_type            = azurerm_cognitive_account_connection.test.auth_type
  category             = azurerm_cognitive_account_connection.test.category
  target               = azurerm_cognitive_account_connection.test.target
  api_key              = azurerm_cognitive_account_connection.test.api_key
  metadata = {
    apiType    = azurerm_cognitive_account_connection.test.metadata.apiType
    resourceId = azurerm_cognitive_account_connection.test.metadata.resourceId
    location   = azurerm_cognitive_account_connection.test.metadata.location
  }
}
`, r.basic(data))
}

func (r CognitiveAccountConnectionTestResource) apiKey(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account" "openai2" {
  name                = "acctest-openai2-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "OpenAI"
  sku_name            = "S0"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type            = "ApiKey"
  category             = "AzureOpenAI"
  target               = azurerm_cognitive_account.openai2.endpoint
  api_key              = azurerm_cognitive_account.openai2.primary_access_key

  metadata = {
    apiType    = "Azure"
    resourceId = azurerm_cognitive_account.openai2.id
    location   = azurerm_cognitive_account.openai2.location
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionTestResource) oauth2(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-%d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type            = "OAuth2"
  category             = "AzureBlob"
  target               = azurerm_storage_account.test.primary_blob_endpoint
  metadata = {
    containerName = azurerm_storage_container.test.name
    accountName   = azurerm_storage_account.test.name
  }
  oauth2 {
    auth_url = "https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000/oauth2/v2.0/token"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionTestResource) oauth2Updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-%d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type            = "OAuth2"
  category             = "AzureBlob"
  target               = azurerm_storage_account.test.primary_blob_endpoint
  metadata = {
    containerName = azurerm_storage_container.test.name
    accountName   = azurerm_storage_account.test.name
  }
  oauth2 {
    auth_url        = "https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000/oauth2/v2.0/token"
    client_id       = "00000000-0000-0000-0000-000000000000"
    client_secret   = "placeHolderClientSecret"
    tenant_id       = "00000000-0000-0000-0000-000000000000"
    developer_token = "placeHolderDevToken"
    refresh_token   = "placeRefreshToken"
    username        = "placeHolderUsername"
    password        = "placeHolderPassword"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionTestResource) customKeys(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-%d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type            = "CustomKeys"
  category             = "CustomKeys"
  target               = azurerm_cognitive_account.openai.endpoint

  metadata = {
    apiType    = "Azure"
    resourceId = azurerm_cognitive_account.openai.id
    location   = azurerm_cognitive_account.openai.location
  }

  custom_keys = {
    primaryKey   = azurerm_cognitive_account.openai.primary_access_key
    secondaryKey = azurerm_cognitive_account.openai.secondary_access_key
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionTestResource) AAD(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-%d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type            = "AAD"
  category             = "AzureBlob"
  target               = azurerm_storage_account.test.primary_blob_endpoint

  metadata = {
    accountName   = azurerm_storage_account.test.name
    containerName = azurerm_storage_container.test.name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionTestResource) invalidAuthTypeApiKeyMismatch(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-invalid"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type            = "ApiKey"
  category             = "AzureOpenAI"
  target               = azurerm_cognitive_account.openai.endpoint

  metadata = {
    apiType    = "Azure"
    resourceId = azurerm_cognitive_account.openai.id
    location   = azurerm_cognitive_account.openai.location
  }

  # Missing api_key field - should cause validation error
}
`, r.template(data))
}

func (r CognitiveAccountConnectionTestResource) invalidAuthTypeAADWithOtherAuth(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cognitive_account_connection" "test" {
  name                 = "acctest-conn-invalid"
  cognitive_account_id = azurerm_cognitive_account.test.id
  auth_type            = "AAD"
  category             = "AzureBlob"
  target               = azurerm_storage_account.test.primary_blob_endpoint

  metadata = {
    accountName   = azurerm_storage_account.test.name
    containerName = azurerm_storage_container.test.name
  }

  # Should not specify other auth blocks when using AAD
  api_key = "should-not-be-here"
}
`, r.template(data))
}
