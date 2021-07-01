package databricks_test

import (
	"context"
	"fmt"
	"testing"

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
	r := DatabricksWorkspaceCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "premium"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabricksWorkspaceCustomerManagedKey_delete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace_customer_managed_key", "test")
	r := DatabricksWorkspaceCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "premium"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.removed(data, "premium"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabricksWorkspaceCustomerManagedKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace_customer_managed_key", "test")
	r := DatabricksWorkspaceCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "premium"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDatabricksWorkspaceCustomerManagedKey_machineLearning(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceCustomerManagedKeyResource{}
	cmkTemplate := r.cmkTemplate(data)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.machineLearning(data, "premium", cmkTemplate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.machineLearning(data, "premium", ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (DatabricksWorkspaceCustomerManagedKeyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.CustomerManagedKeyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataBricks.WorkspacesClient.Get(ctx, id.ResourceGroup, id.CustomerMangagedKeyName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Databricks Workspace Customer Mangaged Key %q (resource group: %q): %+v", id.CustomerMangagedKeyName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.WorkspaceProperties != nil), nil
}

func (DatabricksWorkspaceCustomerManagedKeyResource) basic(data acceptance.TestData, sku string) string {
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
  sku                 = "%[4]s"

  custom_parameters {
    customer_managed_key_enabled = true
  }
}

resource "azurerm_databricks_workspace_customer_managed_key" "test" {
  workspace_id     = azurerm_databricks_workspace.test.id
  key_vault_key_id = azurerm_key_vault_key.test.id
}
`, data.RandomInteger, data.Locations.Primary, keyVault, sku, data.RandomString)
}

func (DatabricksWorkspaceCustomerManagedKeyResource) removed(data acceptance.TestData, sku string) string {
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
  sku                 = "%[4]s"

  custom_parameters {
    customer_managed_key_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, keyVault, sku, data.RandomString)
}

func (DatabricksWorkspaceCustomerManagedKeyResource) requiresImport(data acceptance.TestData) string {
	template := DatabricksWorkspaceCustomerManagedKeyResource{}.basic(data, "premium")
	return fmt.Sprintf(`
%s

resource "azurerm_databricks_workspace_customer_managed_key" "test2" {
  workspace_id     = azurerm_databricks_workspace.test.id
  key_vault_key_id = azurerm_key_vault_key.test.id
}
`, template)
}

func (DatabricksWorkspaceCustomerManagedKeyResource) machineLearning(data acceptance.TestData, sku string, cmk string) string {
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

resource "azurerm_application_insights" "test" {
  name                = "acctest-ai-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[5]s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"

  static_website {
    error_404_document = "error.html"
    index_document     = "index.html"
  }
}

resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acctest-mlws-%[1]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id

  identity {
    type = "SystemAssigned"
  }
}

%[3]s

resource "azurerm_databricks_workspace" "test" {
  name                = "acctestDBW-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "%[4]s"

  custom_parameters {
    customer_managed_key_enabled  = true
    machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  }
}

%[6]s
`, data.RandomInteger, data.Locations.Primary, keyVault, sku, data.RandomString, cmk)
}

func (DatabricksWorkspaceCustomerManagedKeyResource) cmkTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_databricks_workspace_customer_managed_key" "test" {
  workspace_id     = azurerm_databricks_workspace.test.id
  key_vault_key_id = azurerm_key_vault_key.test.id
}
`)
}

func (DatabricksWorkspaceCustomerManagedKeyResource) keyVaultTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_key_vault" "test" {
  name                = "acctest-kv-%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"
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
