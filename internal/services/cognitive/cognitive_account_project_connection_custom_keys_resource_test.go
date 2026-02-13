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

type CognitiveAccountProjectConnectionCustomKeysResource struct{}

func TestAccCognitiveAccountProjectConnectionCustomKeys_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_custom_keys", "test")
	r := CognitiveAccountProjectConnectionCustomKeysResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_keys", "metadata"),
	})
}

func TestAccCognitiveAccountProjectConnectionCustomKeys_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_custom_keys", "test")
	r := CognitiveAccountProjectConnectionCustomKeysResource{}

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

func TestAccCognitiveAccountProjectConnectionCustomKeys_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_custom_keys", "test")
	r := CognitiveAccountProjectConnectionCustomKeysResource{}

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
	})
}

func (r CognitiveAccountProjectConnectionCustomKeysResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r CognitiveAccountProjectConnectionCustomKeysResource) template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary)
}

func (r CognitiveAccountProjectConnectionCustomKeysResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_project_connection_custom_keys" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_project_id = azurerm_cognitive_account_project.test.id
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

func (r CognitiveAccountProjectConnectionCustomKeysResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account_project_connection_custom_keys" "import" {
  name                 = azurerm_cognitive_account_project_connection_custom_keys.test.name
  cognitive_project_id = azurerm_cognitive_account_project_connection_custom_keys.test.cognitive_project_id
  category             = azurerm_cognitive_account_project_connection_custom_keys.test.category
  target               = azurerm_cognitive_account_project_connection_custom_keys.test.target

  metadata = {
    apiType    = azurerm_cognitive_account_project_connection_custom_keys.test.metadata.apiType
    resourceId = azurerm_cognitive_account_project_connection_custom_keys.test.metadata.resourceId
    location   = azurerm_cognitive_account_project_connection_custom_keys.test.metadata.location
  }

  custom_keys = {
    primaryKey   = azurerm_cognitive_account_project_connection_custom_keys.test.custom_keys.primaryKey
    secondaryKey = azurerm_cognitive_account_project_connection_custom_keys.test.custom_keys.secondaryKey
  }
}
`, r.basic(data))
}

func (r CognitiveAccountProjectConnectionCustomKeysResource) updated(data acceptance.TestData) string {
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

resource "azurerm_cognitive_account_project_connection_custom_keys" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_project_id = azurerm_cognitive_account_project.test.id
  category             = "CustomKeys"
  target               = azurerm_cognitive_account.openai2.endpoint

  metadata = {
    apiType    = "Azure"
    resourceId = azurerm_cognitive_account.openai2.id
    location   = azurerm_cognitive_account.openai2.location
  }

  custom_keys = {
    primaryKey   = azurerm_cognitive_account.openai2.primary_access_key
    secondaryKey = azurerm_cognitive_account.openai2.secondary_access_key
  }
}
`, r.template(data), data.RandomInteger)
}
