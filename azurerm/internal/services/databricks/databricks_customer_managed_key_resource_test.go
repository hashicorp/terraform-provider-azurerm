package databricks_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/databricks/mgmt/2018-04-01/databricks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databricks/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DatabricksWorkspaceCustomerManagedKeyResource struct {
}

func TestAccDatabricksWorkspaceCustomerManagedKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace_customer_managed_key", "test")
	parent := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceCustomerManagedKeyResource{}
	cmkTemplate := r.cmkTemplate()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, cmkTemplate),
			Check: acceptance.ComposeTestCheckFunc(
				// You must look for the parent resource (e.g. Databricks Workspace)
				// and then derive if the CMK object has been set or not...
				check.That(parent.ResourceName).ExistsInAzure(r),
			),
		},
		parent.ImportStep(),
	})
}

func TestAccDatabricksWorkspaceCustomerManagedKey_remove(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace_customer_managed_key", "test")
	parent := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceCustomerManagedKeyResource{}
	cmkTemplate := r.cmkTemplate()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, cmkTemplate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(parent.ResourceName).ExistsInAzure(r),
			),
		},
		parent.ImportStep(),
		{
			Config: r.basic(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				// Then ensure the encryption settings on the Databricks Workspace
				// have been reverted to their default state
				check.That(parent.ResourceName).DoesNotExistInAzure(r),
			),
		},
		parent.ImportStep(),
	})
}

func TestAccDatabricksWorkspaceCustomerManagedKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace_customer_managed_key", "test")
	parent := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceCustomerManagedKeyResource{}
	cmkTemplate := r.cmkTemplate()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, cmkTemplate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(parent.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDatabricksWorkspaceCustomerManagedKey_noIp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace_customer_managed_key", "test")
	parent := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceCustomerManagedKeyResource{}
	cmkTemplate := r.cmkTemplate()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.noip(data, cmkTemplate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(parent.ResourceName).ExistsInAzure(r),
			),
		},
		parent.ImportStep(),
		{
			Config: r.noip(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(parent.ResourceName).DoesNotExistInAzure(r),
				check.That(parent.ResourceName).Key("custom_parameters.0.no_public_ip").IsSet(),
			),
		},
		parent.ImportStep(),
	})
}

func (DatabricksWorkspaceCustomerManagedKeyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.WorkspaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataBricks.WorkspacesClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Databricks Workspace Customer Mangaged Key %q (resource group: %q): %+v", id.Name, id.ResourceGroup, err)
	}

	// This is the only way we can tell if the CMK has actually been provisioned or not...
	if resp.WorkspaceProperties.Parameters != nil && resp.WorkspaceProperties.Parameters.Encryption != nil {
		if resp.WorkspaceProperties.Parameters.Encryption.Value.KeySource == databricks.MicrosoftKeyvault {
			return utils.Bool(true), nil
		}
	}

	return utils.Bool(false), nil
}

func (DatabricksWorkspaceCustomerManagedKeyResource) basic(data acceptance.TestData, cmk string) string {
	keyVault := DatabricksWorkspaceCustomerManagedKeyResource{}.keyVaultTemplate(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-db-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_databricks_workspace" "test" {
  name                = "acctestDBW-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "premium"

  customer_managed_key_enabled = true
}

%[4]s
`, data.RandomInteger, data.Locations.Primary, keyVault, cmk)
}

func (DatabricksWorkspaceCustomerManagedKeyResource) requiresImport(data acceptance.TestData) string {
	cmkTemplate := DatabricksWorkspaceCustomerManagedKeyResource{}.cmkTemplate()
	template := DatabricksWorkspaceCustomerManagedKeyResource{}.basic(data, cmkTemplate)
	return fmt.Sprintf(`
%s

resource "azurerm_databricks_workspace_customer_managed_key" "import" {
  workspace_id     = azurerm_databricks_workspace.test.id
  key_vault_key_id = azurerm_key_vault_key.test.id
}
`, template)
}

func (DatabricksWorkspaceCustomerManagedKeyResource) noip(data acceptance.TestData, cmk string) string {
	keyVault := DatabricksWorkspaceCustomerManagedKeyResource{}.keyVaultTemplate(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-db-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_databricks_workspace" "test" {
  name                = "acctestDBW-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "premium"

  customer_managed_key_enabled = true

  custom_parameters {
    no_public_ip = true
  }
}

%[5]s
`, data.RandomInteger, data.Locations.Primary, keyVault, data.RandomString, cmk)
}

func (DatabricksWorkspaceCustomerManagedKeyResource) cmkTemplate() string {
	return `
resource "azurerm_databricks_workspace_customer_managed_key" "test" {
  depends_on = [azurerm_key_vault_access_policy.databricks]

  workspace_id     = azurerm_databricks_workspace.test.id
  key_vault_key_id = azurerm_key_vault_key.test.id
}
`
}

func (DatabricksWorkspaceCustomerManagedKeyResource) keyVaultTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_key_vault" "test" {
  name                = "acctest-kv-%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_key" "test" {
  depends_on = [azurerm_key_vault_access_policy.test]

  name         = "acctest-key-%[1]d"
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
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_key_vault.test.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "get",
    "list",
    "create",
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
    "delete",
    "restore",
    "recover",
    "update",
    "purge",
  ]
}

resource "azurerm_key_vault_access_policy" "databricks" {
  depends_on = [azurerm_databricks_workspace.test]

  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_databricks_workspace.test.storage_account_identity.0.tenant_id
  object_id    = azurerm_databricks_workspace.test.storage_account_identity.0.principal_id

  key_permissions = [
    "get",
    "unwrapKey",
    "wrapKey",
    "delete",
    "purge",
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
