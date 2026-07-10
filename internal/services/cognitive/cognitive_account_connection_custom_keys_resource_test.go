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

type CognitiveAccountConnectionCustomKeysResource struct{}

func TestAccCognitiveAccountConnectionCustomKeys_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_custom_keys", "test")
	r := CognitiveAccountConnectionCustomKeysResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_keys"),
	})
}

func TestAccCognitiveAccountConnectionCustomKeys_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_custom_keys", "test")
	r := CognitiveAccountConnectionCustomKeysResource{}

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

func TestAccCognitiveAccountConnectionCustomKeys_withMetadata(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_custom_keys", "test")
	r := CognitiveAccountConnectionCustomKeysResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withMetadata(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_keys", "metadata"),
	})
}

func TestAccCognitiveAccountConnectionCustomKeys_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_custom_keys", "test")
	r := CognitiveAccountConnectionCustomKeysResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_keys", "metadata"),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_keys", "metadata"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_keys", "metadata"),
	})
}

func TestAccCognitiveAccountConnectionCustomKeys_remoteA2ACategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_custom_keys", "test")
	r := CognitiveAccountConnectionCustomKeysResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.remoteA2ACategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionCustomKeys_remoteToolCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_custom_keys", "test")
	r := CognitiveAccountConnectionCustomKeysResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.remoteToolCategory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccCognitiveAccountConnectionCustomKeys_importMismatchedAuthTypeFails(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_custom_keys", "test")
	r := CognitiveAccountConnectionCustomKeysResource{}

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

func (r CognitiveAccountConnectionCustomKeysResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r CognitiveAccountConnectionCustomKeysResource) template(data acceptance.TestData) string {
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

func (r CognitiveAccountConnectionCustomKeysResource) basic(data acceptance.TestData) string {
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

resource "azurerm_cognitive_account_connection_custom_keys" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "CustomKeys"
  target               = azurerm_cognitive_account.openai.endpoint

  custom_keys = {
    primaryKey   = azurerm_cognitive_account.openai.primary_access_key
    secondaryKey = azurerm_cognitive_account.openai.secondary_access_key
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionCustomKeysResource) withMetadata(data acceptance.TestData) string {
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

resource "azurerm_cognitive_account_connection_custom_keys" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "CustomKeys"
  target               = azurerm_cognitive_account.openai.endpoint

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_cognitive_account.openai.id
    Location   = azurerm_cognitive_account.openai.location
  }

  custom_keys = {
    primaryKey   = azurerm_cognitive_account.openai.primary_access_key
    secondaryKey = azurerm_cognitive_account.openai.secondary_access_key
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionCustomKeysResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account_connection_custom_keys" "import" {
  name                 = azurerm_cognitive_account_connection_custom_keys.test.name
  cognitive_account_id = azurerm_cognitive_account_connection_custom_keys.test.cognitive_account_id
  category             = azurerm_cognitive_account_connection_custom_keys.test.category
  target               = azurerm_cognitive_account_connection_custom_keys.test.target
  custom_keys          = azurerm_cognitive_account_connection_custom_keys.test.custom_keys
}
`, r.basic(data))
}

func (r CognitiveAccountConnectionCustomKeysResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account" "openai2" {
  name                = "acctest-cogacc-openai2-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "OpenAI"
  sku_name            = "S0"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account_connection_custom_keys" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "CustomKeys"
  target               = azurerm_cognitive_account.openai2.endpoint

  metadata = {
    ApiType    = "Azure"
    ResourceId = azurerm_cognitive_account.openai2.id
    Location   = azurerm_cognitive_account.openai2.location
  }

  custom_keys = {
    primaryKey   = azurerm_cognitive_account.openai2.primary_access_key
    secondaryKey = azurerm_cognitive_account.openai2.secondary_access_key
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionCustomKeysResource) remoteA2ACategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_custom_keys" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "RemoteA2A"
  target               = "https://a2a.example.com/"

  custom_keys = {
    apiKey = "test-api-key"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionCustomKeysResource) remoteToolCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_custom_keys" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "RemoteTool"
  target               = "https://tool.example.com/"

  custom_keys = {
    apiKey = "test-api-key"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionCustomKeysResource) importMismatchedAuthType(data acceptance.TestData) string {
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
