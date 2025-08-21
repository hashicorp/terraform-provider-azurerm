// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AIProject struct{}

func TestAccAIProject_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ai_project", "test")
	r := AIProject{}

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

func TestAccAIProject_userIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ai_project", "test")
	r := AIProject{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.userIdentityUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAIProject_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ai_project", "test")
	r := AIProject{}

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

func TestAccAIProject_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ai_project", "test")
	r := AIProject{}

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

func TestAccAIProject_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ai_project", "test")
	r := AIProject{}

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

func (AIProject) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := workspaces.ParseWorkspaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MachineLearning.Workspaces.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r AIProject) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_ai_project" "test" {
  name      = "acctestaip-%[2]d"
  location  = azurerm_ai_hub.test.location
  ai_hub_id = azurerm_ai_hub.test.id

  identity {
    type = "SystemAssigned"
  }
}
`, AIHub{}.basic(data), data.RandomInteger)
}

func (r AIProject) userIdentityTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctestuai-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_resource_group.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_key_vault_access_policy" "test2" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_user_assigned_identity.test.tenant_id
  object_id    = azurerm_user_assigned_identity.test.client_id

  key_permissions = [
    "Create",
    "Get",
  ]
}
`, AIHub{}.basic(data), data.RandomInteger)
}

func (r AIProject) userIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_ai_project" "test" {
  name                           = "acctestaip-%[2]d"
  location                       = azurerm_ai_hub.test.location
  ai_hub_id                      = azurerm_ai_hub.test.id
  primary_user_assigned_identity = azurerm_user_assigned_identity.test.id

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  depends_on = [azurerm_key_vault_access_policy.test2, azurerm_role_assignment.test]
}
`, r.userIdentityTemplate(data), data.RandomInteger)
}

func (r AIProject) userIdentityUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test2" {
  location            = azurerm_resource_group.test.location
  name                = "acctestuai2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_role_assignment" "test2" {
  scope                = azurerm_resource_group.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_user_assigned_identity.test2.principal_id
}

resource "azurerm_key_vault_access_policy" "test3" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_user_assigned_identity.test2.tenant_id
  object_id    = azurerm_user_assigned_identity.test2.client_id

  key_permissions = [
    "Create",
    "Get",
  ]
}

resource "azurerm_ai_project" "test" {
  name                           = "acctestaip-%[2]d"
  location                       = azurerm_ai_hub.test.location
  ai_hub_id                      = azurerm_ai_hub.test.id
  primary_user_assigned_identity = azurerm_user_assigned_identity.test2.id

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id, azurerm_user_assigned_identity.test2.id]
  }

  depends_on = [azurerm_key_vault_access_policy.test2, azurerm_role_assignment.test2]
}
`, r.userIdentityTemplate(data), data.RandomInteger)
}

func (r AIProject) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_ai_project" "test" {
  name      = "acctestaip-%[2]d"
  location  = azurerm_ai_hub.test.location
  ai_hub_id = azurerm_ai_hub.test.id

  description                  = "AI Project created by Terraform"
  friendly_name                = "AI Project"
  high_business_impact_enabled = false

  identity {
    type = "SystemAssigned"
  }

  tags = {
    model = "regression"
  }
}
`, AIHub{}.complete(data), data.RandomInteger)
}

func (r AIProject) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test_project" {
  name                = "acctestuaip-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_ai_project" "test" {
  name      = "acctestaip-%[2]d"
  location  = azurerm_ai_hub.test.location
  ai_hub_id = azurerm_ai_hub.test.id

  description                  = "AI Project updated by Terraform"
  friendly_name                = "AI Project for OS models"
  high_business_impact_enabled = false

  identity {
    type = "SystemAssigned"
  }

  tags = {
    model = "regression"
    env   = "test"
  }
}
`, AIHub{}.complete(data), data.RandomInteger)
}

func (AIProject) requiresImport(data acceptance.TestData) string {
	template := AIProject{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_ai_project" "import" {
  name      = azurerm_ai_project.test.name
  location  = azurerm_ai_project.test.location
  ai_hub_id = azurerm_ai_project.test.ai_hub_id

  identity {
    type = "SystemAssigned"
  }
}
`, template)
}
