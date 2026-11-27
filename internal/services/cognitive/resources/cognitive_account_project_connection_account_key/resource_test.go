// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_project_connection_account_key_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/projectconnectionresource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CognitiveAccountProjectConnectionAccountKeyResource struct{}

func TestAccCognitiveAccountProjectConnectionAccountKey_regressionTest(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_account_key", "test")
	r := CognitiveAccountProjectConnectionAccountKeyResource{}
	data.ResourceRegressionTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
		},
	}, "")
}

func TestAccCognitiveAccountProjectConnectionAccountKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_account_key", "test")
	r := CognitiveAccountProjectConnectionAccountKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("account_key"),
	})
}

func TestAccCognitiveAccountProjectConnectionAccountKey_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_account_key", "test")
	r := CognitiveAccountProjectConnectionAccountKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("account_key"),
	})
}

func TestAccCognitiveAccountProjectConnectionAccountKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_account_key", "test")
	r := CognitiveAccountProjectConnectionAccountKeyResource{}

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

func TestAccCognitiveAccountProjectConnectionAccountKey_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_account_key", "test")
	r := CognitiveAccountProjectConnectionAccountKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("account_key"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("account_key"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("account_key"),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("account_key"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("account_key"),
	})
}

func TestAccCognitiveAccountProjectConnectionAccountKey_storageAccountCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_account_key", "test")
	r := CognitiveAccountProjectConnectionAccountKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("account_key"),
	})
}

func TestAccCognitiveAccountProjectConnectionAccountKey_metadataRequiresResourceId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_account_key", "test")
	r := CognitiveAccountProjectConnectionAccountKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.metadataWithoutResourceId(data),
			ExpectError: regexp.MustCompile("`metadata` must include a non-empty `ResourceId` value"),
		},
	})
}

func TestAccCognitiveAccountProjectConnectionAccountKey_importMismatchedAuthTypeFails(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_account_key", "test")
	r := CognitiveAccountProjectConnectionAccountKeyResource{}

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
				rs, ok := s.RootModule().Resources["azurerm_cognitive_account_project_connection_entra_id.wrong"]
				if !ok {
					return "", fmt.Errorf("resource `%s` not found in state", "azurerm_cognitive_account_project_connection_entra_id.wrong")
				}
				return rs.Primary.ID, nil
			},
			ExpectError: regexp.MustCompile("cannot be managed by"),
		},
	})
}

func (r CognitiveAccountProjectConnectionAccountKeyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r CognitiveAccountProjectConnectionAccountKeyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-pc-%[1]d"
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

resource "azurerm_cognitive_account_project" "test" {
  name                 = "acctest-proj-%[1]d"
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

resource "azurerm_storage_account" "test2" {
  name                     = "acctestsb%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r CognitiveAccountProjectConnectionAccountKeyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_project_connection_account_key" "test" {
  name                         = "acctest-conn-%[2]d"
  cognitive_account_project_id = azurerm_cognitive_account_project.test.id
  category                     = "AzureStorageAccount"
  target                       = azurerm_storage_account.test.primary_blob_endpoint
  account_key                  = azurerm_storage_account.test.primary_access_key

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_storage_account.test.id
    Location   = azurerm_storage_account.test.location
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountProjectConnectionAccountKeyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account_project_connection_account_key" "import" {
  name                         = azurerm_cognitive_account_project_connection_account_key.test.name
  cognitive_account_project_id = azurerm_cognitive_account_project_connection_account_key.test.cognitive_account_project_id
  category                     = azurerm_cognitive_account_project_connection_account_key.test.category
  target                       = azurerm_cognitive_account_project_connection_account_key.test.target
  account_key                  = azurerm_cognitive_account_project_connection_account_key.test.account_key

  metadata = {
    ApiType    = azurerm_cognitive_account_project_connection_account_key.test.metadata.ApiType
    ResourceId = azurerm_cognitive_account_project_connection_account_key.test.metadata.ResourceId
    Location   = azurerm_cognitive_account_project_connection_account_key.test.metadata.Location
  }
}
`, r.basic(data))
}

func (r CognitiveAccountProjectConnectionAccountKeyResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_project_connection_account_key" "test" {
  name                         = "acctest-conn-%[2]d"
  cognitive_account_project_id = azurerm_cognitive_account_project.test.id
  category                     = "AzureStorageAccount"
  target                       = azurerm_storage_account.test.primary_blob_endpoint
  account_key                  = azurerm_storage_account.test.primary_access_key

  metadata = {
    ApiType     = "Azure"
    ResourceId  = azurerm_storage_account.test.id
    Location    = azurerm_storage_account.test.location
    CustomLabel = "acceptance-test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountProjectConnectionAccountKeyResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_project_connection_account_key" "test" {
  name                         = "acctest-conn-%[2]d"
  cognitive_account_project_id = azurerm_cognitive_account_project.test.id
  category                     = "AzureStorageAccount"
  target                       = azurerm_storage_account.test2.primary_blob_endpoint
  account_key                  = azurerm_storage_account.test2.primary_access_key

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_storage_account.test2.id
    Location   = azurerm_storage_account.test2.location
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountProjectConnectionAccountKeyResource) metadataWithoutResourceId(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_project_connection_account_key" "test" {
  name                         = "acctest-conn-%[2]d"
  cognitive_account_project_id = azurerm_cognitive_account_project.test.id
  category                     = "AzureStorageAccount"
  target                       = "https://example.com"
  account_key                  = "example"

  metadata = {
    ApiType = "Azure"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountProjectConnectionAccountKeyResource) importMismatchedAuthType(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account_project_connection_entra_id" "wrong" {
  name                         = "acctest-wrong-%[2]d"
  cognitive_account_project_id = azurerm_cognitive_account_project.test.id
  category                     = "Databricks"
  target                       = "https://wrong-%[2]d.databricks.net/"
}
`, r.basic(data), data.RandomInteger)
}
