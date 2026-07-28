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

type CognitiveAccountConnectionApiKeyResource struct{}

func TestAccCognitiveAccountConnectionApiKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("api_key"),
	})
}

func TestAccCognitiveAccountConnectionApiKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

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

func TestAccCognitiveAccountConnectionApiKey_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("api_key"),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("api_key", "metadata"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("api_key"),
	})
}

func TestAccCognitiveAccountConnectionApiKey_updateOptionalTarget(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.openAICategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("api_key"),
		{
			Config: r.openAIWithTarget(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("api_key"),
		{
			Config: r.openAICategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("api_key"),
	})
}

func TestAccCognitiveAccountConnectionApiKey_aiServicesCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.aiServicesCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionApiKey_azureOpenAICategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureOpenAICategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionApiKey_appInsightsCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appInsightsCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionApiKey_cognitiveSearchCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cognitiveSearchCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionApiKey_apiKeyCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.apiKeyCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionApiKey_groundingWithCustomSearchCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.groundingWithCustomSearchCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionApiKey_openAICategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.openAICategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionApiKey_serpCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serpCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionApiKey_serverlessCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serverlessCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionApiKey_apiManagementCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.apiManagementCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionApiKey_appConfigCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appConfigCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionApiKey_bingLLMSearchCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.bingLLMSearchCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionApiKey_cognitiveServiceCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cognitiveServiceCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionApiKey_groundingWithBingSearchCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.groundingWithBingSearchCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionApiKey_modelGatewayCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.modelGatewayCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionApiKey_pineconeCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.pineconeCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionApiKey_importMismatchedAuthTypeFails(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_api_key", "test")
	r := CognitiveAccountConnectionApiKeyResource{}

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

func (r CognitiveAccountConnectionApiKeyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r CognitiveAccountConnectionApiKeyResource) template(data acceptance.TestData) string {
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

func (r CognitiveAccountConnectionApiKeyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_api_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "ApiKey"
  target               = "https://api.example.com/"
  api_key              = "test-api-key"
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionApiKeyResource) azureOpenAICategory(data acceptance.TestData) string {
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

resource "azurerm_cognitive_account_connection_api_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "AzureOpenAI"
  target               = azurerm_cognitive_account.openai.endpoint
  api_key              = azurerm_cognitive_account.openai.primary_access_key

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_cognitive_account.openai.id
    Location   = azurerm_cognitive_account.openai.location
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionApiKeyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account_connection_api_key" "import" {
  name                 = azurerm_cognitive_account_connection_api_key.test.name
  cognitive_account_id = azurerm_cognitive_account_connection_api_key.test.cognitive_account_id
  category             = azurerm_cognitive_account_connection_api_key.test.category
  target               = azurerm_cognitive_account_connection_api_key.test.target
  api_key              = azurerm_cognitive_account_connection_api_key.test.api_key
}
`, r.basic(data))
}

func (r CognitiveAccountConnectionApiKeyResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_api_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "ApiKey"
  target               = "https://api2.example.com/"
  api_key              = "test-api-key-2"

  metadata = {
    Type = "updated"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionApiKeyResource) aiServicesCategory(data acceptance.TestData) string {
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

resource "azurerm_cognitive_account_connection_api_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "AIServices"
  target               = azurerm_cognitive_account.aiservices.endpoint
  api_key              = azurerm_cognitive_account.aiservices.primary_access_key

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_cognitive_account.aiservices.id
    Location   = azurerm_cognitive_account.aiservices.location
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionApiKeyResource) appInsightsCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_application_insights" "test" {
  name                = "acctestappi-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_cognitive_account_connection_api_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "AppInsights"
  target               = azurerm_application_insights.test.id
  api_key              = azurerm_application_insights.test.connection_string

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_application_insights.test.id
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionApiKeyResource) cognitiveSearchCategory(data acceptance.TestData) string {
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

resource "azurerm_cognitive_account_connection_api_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "CognitiveSearch"
  target               = "https://${azurerm_search_service.test.name}.search.windows.net/"
  api_key              = azurerm_search_service.test.primary_key

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_search_service.test.id
    Location   = azurerm_search_service.test.location
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionApiKeyResource) apiKeyCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_api_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "ApiKey"
  target               = "https://api.example.com/"
  api_key              = "test-api-key"

  metadata = {}
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionApiKeyResource) groundingWithCustomSearchCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_api_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "GroundingWithCustomSearch"
  target               = "https://api.bing.microsoft.com/"
  api_key              = "test-api-key"

  metadata = {
    ApiType    = "Azure"
    ResourceId = "/subscriptions/%[3]s/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.Bing/accounts/acctestbing%[2]d"
    Type       = "bing_custom_search"
  }
}
`, r.template(data), data.RandomInteger, data.Subscriptions.Primary)
}

func (r CognitiveAccountConnectionApiKeyResource) openAICategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_api_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "OpenAI"
  api_key              = "test-api-key"
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionApiKeyResource) serpCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_api_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "Serp"
  api_key              = "test-api-key"
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionApiKeyResource) openAIWithTarget(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_api_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "OpenAI"
  target               = "https://api.openai.com/"
  api_key              = "test-api-key-2"
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionApiKeyResource) serverlessCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_api_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "Serverless"
  target               = "https://acctest-%[2]d.example.models.ai.azure.com/"
  api_key              = "test-api-key"

  metadata = {}
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionApiKeyResource) apiManagementCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_api_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "ApiManagement"
  target               = "https://api-management.example.com/"
  api_key              = "test-api-key"
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionApiKeyResource) appConfigCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_api_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "AppConfig"
  target               = "https://app-config.example.com/"
  api_key              = "test-api-key"
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionApiKeyResource) bingLLMSearchCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_api_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "BingLLMSearch"
  target               = "https://api.bing.microsoft.com/"
  api_key              = "test-api-key"

  metadata = {
    Location = "%[3]s"
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r CognitiveAccountConnectionApiKeyResource) cognitiveServiceCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_api_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "CognitiveService"
  target               = "https://cognitive-service.example.com/"
  api_key              = "test-api-key"

  metadata = {
    Kind = "AIServices"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionApiKeyResource) groundingWithBingSearchCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_api_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "GroundingWithBingSearch"
  target               = "https://api.bing.microsoft.com/"
  api_key              = "test-api-key"
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionApiKeyResource) modelGatewayCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_api_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "ModelGateway"
  target               = "https://gateway.example.com/"
  api_key              = "test-api-key"
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionApiKeyResource) pineconeCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_api_key" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "Pinecone"
  api_key              = "test-api-key"
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionApiKeyResource) importMismatchedAuthType(data acceptance.TestData) string {
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
