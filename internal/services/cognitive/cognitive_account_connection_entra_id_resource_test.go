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

type CognitiveAccountConnectionEntraIdResource struct{}

func TestAccCognitiveAccountConnectionEntraID_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_entra_id", "test")
	r := CognitiveAccountConnectionEntraIdResource{}

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

func TestAccCognitiveAccountConnectionEntraID_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_entra_id", "test")
	r := CognitiveAccountConnectionEntraIdResource{}

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

func TestAccCognitiveAccountConnectionEntraID_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_entra_id", "test")
	r := CognitiveAccountConnectionEntraIdResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
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
		data.ImportStep(),
	})
}

func TestAccCognitiveAccountConnectionEntraID_aiServicesCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_entra_id", "test")
	r := CognitiveAccountConnectionEntraIdResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.aiServicesCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionEntraID_azureOpenAICategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_entra_id", "test")
	r := CognitiveAccountConnectionEntraIdResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureOpenAICategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionEntraID_storageAccountCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_entra_id", "test")
	r := CognitiveAccountConnectionEntraIdResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccountCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionEntraID_cognitiveSearchCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_entra_id", "test")
	r := CognitiveAccountConnectionEntraIdResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cognitiveSearchCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionEntraID_cosmosDbCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_entra_id", "test")
	r := CognitiveAccountConnectionEntraIdResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cosmosDbCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionEntraID_databricksCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_entra_id", "test")
	r := CognitiveAccountConnectionEntraIdResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.databricksCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionEntraID_apiManagementCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_entra_id", "test")
	r := CognitiveAccountConnectionEntraIdResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.apiManagementCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionEntraID_appConfigCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_entra_id", "test")
	r := CognitiveAccountConnectionEntraIdResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appConfigCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionEntraID_cognitiveServiceCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_entra_id", "test")
	r := CognitiveAccountConnectionEntraIdResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cognitiveServiceCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionEntraID_managedOnlineEndpointCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_entra_id", "test")
	r := CognitiveAccountConnectionEntraIdResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedOnlineEndpointCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionEntraID_microsoftFabricCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_entra_id", "test")
	r := CognitiveAccountConnectionEntraIdResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.microsoftFabricCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionEntraID_sharepointCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_entra_id", "test")
	r := CognitiveAccountConnectionEntraIdResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sharepointCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionEntraID_importMismatchedAuthTypeFails(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_entra_id", "test")
	r := CognitiveAccountConnectionEntraIdResource{}

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
				rs, ok := s.RootModule().Resources["azurerm_cognitive_account_connection_api_key.wrong"]
				if !ok {
					return "", fmt.Errorf("resource `%s` not found in state", "azurerm_cognitive_account_connection_api_key.wrong")
				}
				return rs.Primary.ID, nil
			},
			ExpectError: regexp.MustCompile("cannot be managed by"),
		},
	})
}

func (r CognitiveAccountConnectionEntraIdResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r CognitiveAccountConnectionEntraIdResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-ac-%[1]d"
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

func (r CognitiveAccountConnectionEntraIdResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_entra_id" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "Databricks"
  target               = "https://workspace-%[2]d.databricks.net/"
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionEntraIdResource) azureOpenAICategory(data acceptance.TestData) string {
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

resource "azurerm_cognitive_account_connection_entra_id" "test" {
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

func (r CognitiveAccountConnectionEntraIdResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account_connection_entra_id" "import" {
  name                 = azurerm_cognitive_account_connection_entra_id.test.name
  cognitive_account_id = azurerm_cognitive_account_connection_entra_id.test.cognitive_account_id
  category             = azurerm_cognitive_account_connection_entra_id.test.category
  target               = azurerm_cognitive_account_connection_entra_id.test.target
}
`, r.basic(data))
}

func (r CognitiveAccountConnectionEntraIdResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_entra_id" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "Databricks"
  target               = "https://workspace2-%[2]d.databricks.net/"

  metadata = {
    "ApiType" : "Azure",
    "azure_databricks_connection_type" : "job",
    "job_id" : "123"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionEntraIdResource) aiServicesCategory(data acceptance.TestData) string {
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

resource "azurerm_cognitive_account_connection_entra_id" "test" {
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

func (r CognitiveAccountConnectionEntraIdResource) storageAccountCategory(data acceptance.TestData) string {
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

resource "azurerm_cognitive_account_connection_entra_id" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "AzureStorageAccount"
  target               = azurerm_storage_account.test.primary_blob_endpoint

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_storage_account.test.id
  }
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r CognitiveAccountConnectionEntraIdResource) cognitiveSearchCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_search_service" "test" {
  name                = "acctestsearchsvc%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "standard"
}

resource "azurerm_cognitive_account_connection_entra_id" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "CognitiveSearch"
  target               = "https://${azurerm_search_service.test.name}.search.windows.net/"

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_search_service.test.id
    Type       = "azure_ai_search"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionEntraIdResource) cosmosDbCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                         = "acctestcdb%[3]s"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  offer_type                   = "Standard"
  kind                         = "GlobalDocumentDB"
  local_authentication_enabled = false

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 10
    max_staleness_prefix    = 200
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

resource "azurerm_cognitive_account_connection_entra_id" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "CosmosDb"
  target               = azurerm_cosmosdb_account.test.endpoint

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_cosmosdb_account.test.id
    Location   = azurerm_cosmosdb_account.test.location
  }
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r CognitiveAccountConnectionEntraIdResource) databricksCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_entra_id" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "Databricks"
  target               = "https://workspace%[2]d.databricks.net/"
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionEntraIdResource) apiManagementCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_entra_id" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "ApiManagement"
  target               = "https://api-management.example.com/"
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionEntraIdResource) appConfigCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_entra_id" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "AppConfig"
  target               = "https://app-config.example.com/"
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionEntraIdResource) cognitiveServiceCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_entra_id" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "CognitiveService"
  target               = "https://cognitive-service.example.com/"

  metadata = {
    Kind = "AIServices"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionEntraIdResource) managedOnlineEndpointCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_entra_id" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "ManagedOnlineEndpoint"
  target               = "https://endpoint.example.com/"
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionEntraIdResource) microsoftFabricCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_entra_id" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "MicrosoftFabric"
  target               = "https://fabric.microsoft.com/"
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionEntraIdResource) sharepointCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_entra_id" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "Sharepoint"
  target               = "https://contoso.sharepoint.com/"
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionEntraIdResource) importMismatchedAuthType(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account_connection_api_key" "wrong" {
  name                 = "acctest-wrong-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "ApiKey"
  target               = "https://api.example.com/"
  api_key              = "test-api-key"
}
`, r.basic(data), data.RandomInteger)
}
