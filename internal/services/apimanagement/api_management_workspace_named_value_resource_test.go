package apimanagement_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/namedvalue"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementWorkspaceNamedValueResource struct{}

func TestAccApiManagementWorkspaceNamedValue_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_named_value", "test")
	r := ApiManagementWorkspaceNamedValueResource{}

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
	r := ApiManagementWorkspaceNamedValueResource{}

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

func TestAccApiManagementWorkspaceNamedValue_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_named_value", "test")
	r := ApiManagementWorkspaceNamedValueResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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
		data.ImportStep("value"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementWorkspaceNamedValue_keyVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_named_value", "test")
	r := ApiManagementWorkspaceNamedValueResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementWorkspaceNamedValue_keyVaultInvalidSecretValue(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_named_value", "test")
	r := ApiManagementWorkspaceNamedValueResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.keyVaultWithInvalidSecretValue(data),
			ExpectError: regexp.MustCompile("`secret` must be set to `true` when `value_from_key_vault` is specified"),
		},
	})
}

func TestAccApiManagementWorkspaceNamedValue_keyVaultUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_named_value", "test")
	r := ApiManagementWorkspaceNamedValueResource{}

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
			Config: r.keyVaultUpdateToValue(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("value"),
	})
}

func (ApiManagementWorkspaceNamedValueResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := namedvalue.ParseWorkspaceNamedValueID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.NamedValueClient_v2024_05_01.WorkspaceNamedValueGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (ApiManagementWorkspaceNamedValueResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Premium_1"
}

resource "azurerm_api_management_workspace" "test" {
  name              = "acctestAMW-%[1]d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "Workspace-%[1]d"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ApiManagementWorkspaceNamedValueResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_named_value" "test" {
  name                        = "acctest-nv-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "TestProperty"
  value                       = "Test Value"
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceNamedValueResource) requiresImport(data acceptance.TestData) string {
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

func (r ApiManagementWorkspaceNamedValueResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_named_value" "test" {
  name                        = "acctest-nv-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "TestProperty-Updated"
  value                       = "Test Value Updated"
  secret                      = true
  tags                        = ["tag1", "tag2"]
}
`, r.template(data), data.RandomInteger)
}

func (ApiManagementWorkspaceNamedValueResource) keyVaultTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-apim-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_api_management" "test" {
  name                = "acctestapim-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Premium_1"

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}

resource "azurerm_api_management_workspace" "test" {
  name              = "acctestapimws-%[1]d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "Test-Workspace"
  description       = "Test workspace for named value with key vault"
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "acctestKV-%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id
  certificate_permissions = [
    "Create",
    "Delete",
    "DeleteIssuers",
    "Get",
    "GetIssuers",
    "Import",
    "List",
    "ListIssuers",
    "ManageContacts",
    "ManageIssuers",
    "SetIssuers",
    "Update",
    "Purge",
  ]
  secret_permissions = [
    "Get",
    "Delete",
    "List",
    "Purge",
    "Recover",
    "Set",
  ]
}

resource "azurerm_key_vault_access_policy" "test2" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_user_assigned_identity.test.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id
  secret_permissions = [
    "Get",
    "List",
  ]
}

resource "azurerm_key_vault_access_policy" "test3" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_api_management.test.identity[0].principal_id
  secret_permissions = [
    "Get",
    "List",
  ]
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%[3]s"
  value        = "rick-and-morty"
  key_vault_id = azurerm_key_vault.test.id

  depends_on = [
    azurerm_key_vault_access_policy.test,
    azurerm_key_vault_access_policy.test2,
    azurerm_key_vault_access_policy.test3
  ]
}

resource "azurerm_key_vault_secret" "test2" {
  name         = "secret2-%[3]s"
  value        = "rick-and-morty2"
  key_vault_id = azurerm_key_vault.test.id

  depends_on = [
    azurerm_key_vault_access_policy.test,
    azurerm_key_vault_access_policy.test2,
    azurerm_key_vault_access_policy.test3
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r ApiManagementWorkspaceNamedValueResource) keyVault(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_named_value" "test" {
  name                        = "acctest-nv-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "TestProperty-KeyVault"
  secret                      = true

  value_from_key_vault {
    secret_id = azurerm_key_vault_secret.test.id
  }
}
`, r.keyVaultTemplate(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceNamedValueResource) keyVaultWithInvalidSecretValue(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_named_value" "test" {
  name                        = "acctest-nv-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "TestProperty-KeyVault"
  secret                      = false

  value_from_key_vault {
    secret_id          = azurerm_key_vault_secret.test.id
    identity_client_id = azurerm_user_assigned_identity.test.client_id
  }
}
`, r.keyVaultTemplate(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceNamedValueResource) keyVaultUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_named_value" "test" {
  name                        = "acctest-nv-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "TestProperty-KeyVault-Updated"
  secret                      = true

  value_from_key_vault {
    secret_id          = azurerm_key_vault_secret.test2.id
    identity_client_id = azurerm_user_assigned_identity.test.client_id
  }
}
`, r.keyVaultTemplate(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceNamedValueResource) keyVaultUpdateToValue(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_named_value" "test" {
  name                        = "acctest-nv-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "TestProperty-Value"
  secret                      = true
  value                       = "Test Value"
}
`, r.keyVaultTemplate(data), data.RandomInteger)
}
