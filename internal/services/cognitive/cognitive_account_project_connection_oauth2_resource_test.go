// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/projectconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CognitiveAccountProjectConnectionOAuth2Resource struct{}

func TestAccCognitiveAccountProjectConnectionOAuth2_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_oauth2", "test")
	r := CognitiveAccountProjectConnectionOAuth2Resource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("oauth2", "metadata"),
	})
}

func TestAccCognitiveAccountProjectConnectionOAuth2_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_oauth2", "test")
	r := CognitiveAccountProjectConnectionOAuth2Resource{}

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

func TestAccCognitiveAccountProjectConnectionOAuth2_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_oauth2", "test")
	r := CognitiveAccountProjectConnectionOAuth2Resource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("oauth2", "metadata"),
		{
			Config: r.anotherContainer(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("oauth2", "metadata"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("oauth2", "metadata"),
	})
}

func (r CognitiveAccountProjectConnectionOAuth2Resource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := projectconnectionresource.ParseProjectConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cognitive.ProjectConnectionResourceClient.ProjectConnectionsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r CognitiveAccountProjectConnectionOAuth2Resource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-apc-%[1]d"
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

resource "azurerm_cognitive_account_project" "test" {
  name                 = "acctest-%[1]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  location             = azurerm_resource_group.test.location
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
  name                  = "acctestsc%[3]s"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r CognitiveAccountProjectConnectionOAuth2Resource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_project_connection_oauth2" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_project_id = azurerm_cognitive_account_project.test.id
  category             = "AzureBlob"
  target               = azurerm_storage_account.test.primary_blob_endpoint

  metadata = {
    containerName = azurerm_storage_container.test.name
    accountName   = azurerm_storage_account.test.name
  }

  oauth2 {
    authentication_url = "https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000/oauth2/v2.0/token"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountProjectConnectionOAuth2Resource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account_project_connection_oauth2" "import" {
  name                 = azurerm_cognitive_account_project_connection_oauth2.test.name
  cognitive_project_id = azurerm_cognitive_account_project_connection_oauth2.test.cognitive_project_id
  category             = azurerm_cognitive_account_project_connection_oauth2.test.category
  target               = azurerm_cognitive_account_project_connection_oauth2.test.target

  metadata = {
    containerName = azurerm_cognitive_account_project_connection_oauth2.test.metadata.containerName
    accountName   = azurerm_cognitive_account_project_connection_oauth2.test.metadata.accountName
  }

  oauth2 {
    authentication_url = "https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000/oauth2/v2.0/token"
  }
}
`, r.basic(data))
}

func (r CognitiveAccountProjectConnectionOAuth2Resource) anotherContainer(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_storage_container" "test2" {
  name                  = "acctestsc2%[2]s"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_cognitive_account_project_connection_oauth2" "test" {
  name                 = "acctest-conn-%[3]d"
  cognitive_project_id = azurerm_cognitive_account_project.test.id
  category             = "AzureBlob"
  target               = azurerm_storage_account.test.primary_blob_endpoint

  metadata = {
    containerName = azurerm_storage_container.test2.name
    accountName   = azurerm_storage_account.test.name
  }

  oauth2 {
    authentication_url = "https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000/oauth2/v2.0/token"
  }
}
`, r.template(data), data.RandomString, data.RandomInteger)
}

func (r CognitiveAccountProjectConnectionOAuth2Resource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_project_connection_oauth2" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_project_id = azurerm_cognitive_account_project.test.id
  category             = "AzureBlob"
  target               = azurerm_storage_account.test.primary_blob_endpoint

  metadata = {
    containerName = azurerm_storage_container.test.name
    accountName   = azurerm_storage_account.test.name
  }

  oauth2 {
    authentication_url = "https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000/oauth2/v2.0/token"
    client_id          = "00000000-0000-0000-0000-000000000000"
    client_secret      = "placeHolderClientSecret"
    tenant_id          = "00000000-0000-0000-0000-000000000000"
    developer_token    = "placeHolderDevToken"
    refresh_token      = "placeRefreshToken"
    username           = "placeHolderUsername"
    password           = "placeHolderPassword"
  }
}
`, r.template(data), data.RandomInteger)
}
