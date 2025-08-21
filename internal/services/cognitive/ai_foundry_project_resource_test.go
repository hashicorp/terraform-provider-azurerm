// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitive_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/cognitiveservicesprojects"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AIFoundryProjectTestResource struct{}

func TestAccCognitiveAIFoundryProject_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ai_foundry_project", "test")
	r := AIFoundryProjectTestResource{}

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

func TestAccCognitiveAIFoundryProject_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ai_foundry_project", "test")
	r := AIFoundryProjectTestResource{}

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

func TestAccCognitiveAIFoundryProject_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ai_foundry_project", "test")
	r := AIFoundryProjectTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveAIFoundryProject_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ai_foundry_project", "test")
	r := AIFoundryProjectTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveAIFoundryProject_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ai_foundry_project", "test")
	r := AIFoundryProjectTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
			),
		},
		data.ImportStep(),
		{
			Config: r.identityUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identitySystemAssignedUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
			),
		},
		data.ImportStep(),
	})
}

func (r AIFoundryProjectTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := cognitiveservicesprojects.ParseProjectID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cognitive.ProjectsClient
	resp, err := client.ProjectsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r AIFoundryProjectTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
resource "azurerm_ai_foundry" "test" {
  name                       = "acctest-cogacc-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "acctest-cogacc-%d"
  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomInteger, data.RandomInteger)
}

func (r AIFoundryProjectTestResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_ai_foundry_project" "test" {
  name          = "acctest-cogproj-%d"
  location      = azurerm_resource_group.test.location
  ai_foundry_id = azurerm_ai_foundry.test.id
  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func (AIFoundryProjectTestResource) identityUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
resource "azurerm_ai_foundry" "test" {
  name                       = "acctest-cogacc-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "acctest-cogacc-%d"
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_ai_foundry_project" "test" {
  name          = "acctest-cogproj-%d"
  location      = azurerm_resource_group.test.location
  ai_foundry_id = azurerm_ai_foundry.test.id
  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (AIFoundryProjectTestResource) identitySystemAssignedUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
resource "azurerm_ai_foundry" "test" {
  name                       = "acctest-cogacc-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "acctest-cogacc-%d"
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_ai_foundry_project" "test" {
  name          = "acctest-cogproj-%d"
  location      = azurerm_resource_group.test.location
  ai_foundry_id = azurerm_ai_foundry.test.id
  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r AIFoundryProjectTestResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_ai_foundry_project" "import" {
  name          = azurerm_ai_foundry_project.test.name
  location      = azurerm_ai_foundry_project.test.location
  ai_foundry_id = azurerm_ai_foundry.test.id
  identity {
    type = "SystemAssigned"
  }
}
`, config)
}

func (r AIFoundryProjectTestResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_ai_foundry_project" "test" {
  name          = "acctest-cogproj-%d"
  location      = azurerm_resource_group.test.location
  ai_foundry_id = azurerm_ai_foundry.test.id
  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func (r AIFoundryProjectTestResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_ai_foundry_project" "test" {
  name          = "acctest-cogproj-%d"
  location      = azurerm_resource_group.test.location
  ai_foundry_id = azurerm_ai_foundry.test.id
  description   = "Updated Description"
  display_name  = "Test Project Updated"
  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}
