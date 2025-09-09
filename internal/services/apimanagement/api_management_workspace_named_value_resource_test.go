// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/namedvalue"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementWorkspaceNamedValueTestResource struct{}

func TestAccApiManagementWorkspaceNamedValue_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_named_value", "test")
	r := ApiManagementWorkspaceNamedValueTestResource{}

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

func TestAccApiManagementWorkspaceNamedValue_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_named_value", "test")
	r := ApiManagementWorkspaceNamedValueTestResource{}

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

func TestAccApiManagementWorkspaceNamedValue_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_named_value", "test")
	r := ApiManagementWorkspaceNamedValueTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("value"),
	})
}

func TestAccApiManagementWorkspaceNamedValue_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_named_value", "test")
	r := ApiManagementWorkspaceNamedValueTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("value"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("value"),
	})
}

func TestAccApiManagementWorkspaceNamedValue_keyVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_named_value", "test")
	r := ApiManagementWorkspaceNamedValueTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementWorkspaceNamedValue_keyVaultWithIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_named_value", "test")
	r := ApiManagementWorkspaceNamedValueTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVaultWithIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultWithIdentityUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultWithIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ApiManagementWorkspaceNamedValueTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := namedvalue.ParseWorkspaceNamedValueID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.NamedValueClient_v2024_05_01.WorkspaceNamedValueGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ApiManagementWorkspaceNamedValueTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Premium_1"
}

resource "azurerm_api_management_workspace" "test" {
  name              = "acctestWS-%d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "Test Workspace"
  description       = "Test workspace description"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementWorkspaceNamedValueTestResource) templateWithIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Premium_1"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_api_management_workspace" "test" {
  name              = "acctestAMWorkspace-%d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "Test Workspace"
  description       = "Test workspace description"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementWorkspaceNamedValueTestResource) templateWithUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-api-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctestUAI2-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Premium_1"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
      azurerm_user_assigned_identity.test2.id,
    ]
  }
}

resource "azurerm_api_management_workspace" "test" {
  name              = "acctestAMWorkspace-%[1]d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "Test Workspace"
  description       = "Test workspace description"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ApiManagementWorkspaceNamedValueTestResource) keyVaultTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {
}

resource "azurerm_key_vault" "test" {
  name                = "acct%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Get",
      "Set",
      "Delete",
      "Purge",
    ]
  }

  access_policy {
    tenant_id = azurerm_api_management.test.identity.0.tenant_id
    object_id = azurerm_api_management.test.identity.0.principal_id

    secret_permissions = [
      "Get",
      "Set",
      "Delete",
      "Purge",
    ]
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%[1]d"
  value        = "szechuan"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_key_vault" "test2" {
  name                = "acct2%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Get",
      "Set",
      "Delete",
      "Purge",
    ]
  }

  access_policy {
    tenant_id = azurerm_api_management.test.identity.0.tenant_id
    object_id = azurerm_api_management.test.identity.0.principal_id

    secret_permissions = [
      "Get",
      "Set",
      "Delete",
      "Purge",
    ]
  }
}

resource "azurerm_key_vault_secret" "test2" {
  name         = "secret2-%[1]d"
  value        = "szechuan"
  key_vault_id = azurerm_key_vault.test.id
}
`, data.RandomInteger)
}

func (r ApiManagementWorkspaceNamedValueTestResource) keyVaultTemplateWithUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {
}

resource "azurerm_key_vault" "test" {
  name                = "acct%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Get",
      "Set",
      "Delete",
      "Purge",
    ]
  }

  access_policy {
    tenant_id = azurerm_user_assigned_identity.test.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    secret_permissions = [
      "Get",
    ]
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%[1]d"
  value        = "szechuan"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_key_vault" "test2" {
  name                = "acct2%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Get",
      "Set",
      "Delete",
      "Purge",
    ]
  }

  access_policy {
    tenant_id = azurerm_user_assigned_identity.test2.tenant_id
    object_id = azurerm_user_assigned_identity.test2.principal_id

    secret_permissions = [
      "Get",
    ]
  }
}

resource "azurerm_key_vault_secret" "test2" {
  name         = "secret-%[1]d"
  value        = "szechuan2"
  key_vault_id = azurerm_key_vault.test2.id
}
`, data.RandomInteger)
}

func (r ApiManagementWorkspaceNamedValueTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_api_management_workspace_named_value" "test" {
  name                        = "acctestAMNamedValue-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "TestDisplayName"
  value                       = "Example Value"
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceNamedValueTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_api_management_workspace_named_value" "test" {
  name                        = "acctestAMNamedValue-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "TestDisplayName"
  secret_enabled              = true
  value                       = "Secret Value"
  tags                        = ["tag1", "tag2"]
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceNamedValueTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_workspace_named_value" "import" {
  name                        = azurerm_api_management_workspace_named_value.test.name
  api_management_workspace_id = azurerm_api_management_workspace_named_value.test.api_management_workspace_id
  display_name                = azurerm_api_management_workspace_named_value.test.display_name
  value                       = azurerm_api_management_workspace_named_value.test.value
}
`, r.basic(data))
}

func (r ApiManagementWorkspaceNamedValueTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_named_value" "test" {
  name                        = "acctestAMNamedValue-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "TestDisplayNameUpdate"
  secret_enabled              = false
  value                       = "Secret Value update"
  tags                        = ["tag3", "tag4"]
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceNamedValueTestResource) keyVault(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = true
    }
  }
}

%s

%s

resource "azurerm_api_management_workspace_named_value" "test" {
  name                        = "acctestAMNamedValue-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "TestKeyVaultNamedValue"
  secret_enabled              = true

  value_from_key_vault {
    key_vault_secret_id = azurerm_key_vault_secret.test.id
  }
}
`, r.templateWithIdentity(data), r.keyVaultTemplate(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceNamedValueTestResource) keyVaultUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = true
    }
  }
}

%s

%s

resource "azurerm_api_management_workspace_named_value" "test" {
  name                        = "acctestAMNamedValue-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "TestKeyVaultNamedValue"
  secret_enabled              = true

  value_from_key_vault {
    key_vault_secret_id = azurerm_key_vault_secret.test2.id
  }
}
`, r.templateWithIdentity(data), r.keyVaultTemplate(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceNamedValueTestResource) keyVaultWithIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = true
    }
  }
}

%s

%s

resource "azurerm_api_management_workspace_named_value" "test" {
  name                        = "acctestAMNamedValue-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "TestKeyVaultNamedValuewithIdentity"
  secret_enabled              = true

  value_from_key_vault {
    key_vault_secret_id              = azurerm_key_vault_secret.test.id
    user_assigned_identity_client_id = azurerm_user_assigned_identity.test.client_id
  }
}
`, r.templateWithUserAssignedIdentity(data), r.keyVaultTemplateWithUserAssignedIdentity(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceNamedValueTestResource) keyVaultWithIdentityUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = true
    }
  }
}

%s

%s

resource "azurerm_api_management_workspace_named_value" "test" {
  name                        = "acctestAMNamedValue-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "TestKeyVaultNamedValuewithIdentity"
  secret_enabled              = true

  value_from_key_vault {
    key_vault_secret_id              = azurerm_key_vault_secret.test2.id
    user_assigned_identity_client_id = azurerm_user_assigned_identity.test2.client_id
  }
}
`, r.templateWithUserAssignedIdentity(data), r.keyVaultTemplateWithUserAssignedIdentity(data), data.RandomInteger)
}
