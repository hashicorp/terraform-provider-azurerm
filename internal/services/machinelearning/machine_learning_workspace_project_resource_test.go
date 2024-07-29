// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkspaceProjectResource struct{}

func TestAccMachineLearningWorkspaceProject_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace_project", "test")
	r := WorkspaceProjectResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningWorkspaceProject_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace_project", "test")
	r := WorkspaceProjectResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMachineLearningWorkspaceProject_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace_project", "test")
	r := WorkspaceProjectResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningWorkspaceProject_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace_project", "test")
	r := WorkspaceProjectResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningWorkspaceProject_systemAssignedUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace_project", "test")
	r := WorkspaceProjectResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemAssignedUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningWorkspaceProject_purgeSoftDelete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace_project", "test")
	r := WorkspaceProjectResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.purgeSoftDelete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (r WorkspaceProjectResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	workspacesClient := client.MachineLearning.Workspaces
	id, err := workspaces.ParseWorkspaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := workspacesClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Machine Learning Workspace Projet %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.Model.Properties != nil), nil
}

func (r WorkspaceProjectResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ml-%[1]d"
  location = "%[2]s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestai-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestvault%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Get",
    "Delete",
    "Purge",
    "GetRotationPolicy",
  ]
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[4]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_container_registry" "test" {
  name                = "acctestacr%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  admin_enabled       = true

  lifecycle {
    ignore_changes = [network_rule_set]
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctest-kv-key-%[1]d"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
  depends_on = [azurerm_key_vault.test, azurerm_key_vault_access_policy.test]
}

resource "azurerm_machine_learning_workspace_hub" "test" {
  name                    = "acctest-MLW-%[1]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id
  container_registry_id   = azurerm_container_registry.test.id
  public_network_access   = "Disabled"
  friendly_name           = "test-workspace"

  identity {
    type = "SystemAssigned"
  }

  encryption {
    key_vault_id = azurerm_key_vault.test.id
    key_id       = azurerm_key_vault_key.test.id
  }

  tags = {
    ENV = "Test"
  }
}`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomIntOfLength(15))
}

func (r WorkspaceProjectResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

%s

resource "azurerm_machine_learning_workspace_project" "test" {
  name                = "acctest-MLW-project-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_hub_id    = azurerm_machine_learning_workspace_hub.test.id
  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomStringOfLength(8))
}

func (r WorkspaceProjectResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

%[1]s


resource "azurerm_machine_learning_workspace_project" "test" {
  name                = "acctest-MLW-project-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_hub_id    = azurerm_machine_learning_workspace_hub.test.id
  friendly_name       = "test-workspace"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomIntOfLength(8))
}

func (r WorkspaceProjectResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

%[1]s

resource "azurerm_machine_learning_workspace_project" "test" {
  name                = "acctest-MLW-project-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_hub_id    = azurerm_machine_learning_workspace_hub.test.id
  friendly_name       = "test-workspace-updated"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    ENV = "Test"
    FOO = "Updated"
  }
}
`, template, data.RandomIntOfLength(8))
}

func (r WorkspaceProjectResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_workspace_project" "import" {
  name                = azurerm_machine_learning_workspace_project.test.name
  location            = azurerm_machine_learning_workspace_project.test.location
  resource_group_name = azurerm_machine_learning_workspace_project.test.resource_group_name
  workspace_hub_id    = azurerm_machine_learning_workspace_hub.test.id
  identity {
    type = "SystemAssigned"
  }
}
`, template)
}

func (r WorkspaceProjectResource) systemAssignedUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_machine_learning_workspace_project" "test" {
  name                = "acctest-MLW-project-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_hub_id    = azurerm_machine_learning_workspace_hub.test.id

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, r.template(data), data.RandomIntOfLength(8))
}

func (r WorkspaceProjectResource) purgeSoftDelete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
    machine_learning {
      purge_soft_deleted_workspace_on_destroy = true
    }
  }
}

%s

resource "azurerm_machine_learning_workspace_project" "test" {
  name                = "acctest-MLW-project-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_hub_id    = azurerm_machine_learning_workspace_hub.test.id

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomIntOfLength(8))
}
