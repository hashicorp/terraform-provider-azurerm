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

type CognitiveAccountConnectionAccountKeyResource struct{}

func TestAccCognitiveAccountConnectionAccountKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_account_key", "test")
	r := CognitiveAccountConnectionAccountKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("account_key", "metadata"),
	})
}

func TestAccCognitiveAccountConnectionAccountKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_account_key", "test")
	r := CognitiveAccountConnectionAccountKeyResource{}

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

func TestAccCognitiveAccountConnectionAccountKey_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_account_key", "test")
	r := CognitiveAccountConnectionAccountKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("account_key", "metadata"),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("account_key", "metadata"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("account_key", "metadata"),
	})
}

func TestAccCognitiveAccountConnectionAccountKey_importMismatchedAuthTypeFails(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_account_key", "test")
	r := CognitiveAccountConnectionAccountKeyResource{}

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

func (r CognitiveAccountConnectionAccountKeyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r CognitiveAccountConnectionAccountKeyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-ak-%[1]d"
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

func (r CognitiveAccountConnectionAccountKeyResource) basic(data acceptance.TestData) string {
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

resource "azurerm_cognitive_account_connection_account_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "AzureStorageAccount"
  target               = azurerm_storage_account.test.primary_blob_endpoint
  account_key          = azurerm_storage_account.test.primary_access_key

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_storage_account.test.id
    Location   = azurerm_storage_account.test.location
  }
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r CognitiveAccountConnectionAccountKeyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account_connection_account_key" "import" {
  name                 = azurerm_cognitive_account_connection_account_key.test.name
  cognitive_account_id = azurerm_cognitive_account_connection_account_key.test.cognitive_account_id
  category             = azurerm_cognitive_account_connection_account_key.test.category
  target               = azurerm_cognitive_account_connection_account_key.test.target
  account_key          = azurerm_cognitive_account_connection_account_key.test.account_key

  metadata = {
    ApiType    = azurerm_cognitive_account_connection_account_key.test.metadata.ApiType
    ResourceId = azurerm_cognitive_account_connection_account_key.test.metadata.ResourceId
    Location   = azurerm_cognitive_account_connection_account_key.test.metadata.Location
  }
}
`, r.basic(data))
}

func (r CognitiveAccountConnectionAccountKeyResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_storage_account" "test2" {
  name                     = "acctestb%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_cognitive_account_connection_account_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "AzureStorageAccount"
  target               = azurerm_storage_account.test2.primary_blob_endpoint
  account_key          = azurerm_storage_account.test2.primary_access_key

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_storage_account.test2.id
    Location   = azurerm_storage_account.test2.location
  }
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r CognitiveAccountConnectionAccountKeyResource) importMismatchedAuthType(data acceptance.TestData) string {
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
