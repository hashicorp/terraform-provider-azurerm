package synapse_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/synapse/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SynapseWorkspaceResource struct{}

func TestAccSynapseWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace", "test")
	r := SynapseWorkspaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("managed_resource_group_name").Exists(),
			),
		},
		data.ImportStep("sql_administrator_login_password"),
	})
}

func TestAccSynapseWorkspace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace", "test")
	r := SynapseWorkspaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSynapseWorkspace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace", "test")
	r := SynapseWorkspaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("managed_resource_group_name").HasValue(fmt.Sprintf("acctest-ManagedSynapse-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("sql_identity_control_enabled").HasValue("true"),
			),
		},
		data.ImportStep("sql_administrator_login_password"),
	})
}

func TestAccSynapseWorkspace_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace", "test")
	r := SynapseWorkspaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sql_administrator_login_password"),
		{
			Config: r.withUpdateFields(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sql_administrator_login_password"),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccSynapseWorkspace_azdo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace", "test")
	r := SynapseWorkspaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.azdo(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("azure_devops_repo.0.account_name").HasValue("myorg"),
				check.That(data.ResourceName).Key("azure_devops_repo.0.project_name").HasValue("myproj"),
				check.That(data.ResourceName).Key("azure_devops_repo.0.repository_name").HasValue("myrepo"),
				check.That(data.ResourceName).Key("azure_devops_repo.0.branch_name").HasValue("dev"),
				check.That(data.ResourceName).Key("azure_devops_repo.0.root_folder").HasValue("/"),
			),
		},
	})
}

func TestAccSynapseWorkspace_github(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace", "test")
	r := SynapseWorkspaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.github(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("github_repo.0.account_name").HasValue("myuser"),
				check.That(data.ResourceName).Key("github_repo.0.git_url").HasValue("https://github.mydomain.com"),
				check.That(data.ResourceName).Key("github_repo.0.repository_name").HasValue("myrepo"),
				check.That(data.ResourceName).Key("github_repo.0.branch_name").HasValue("dev"),
				check.That(data.ResourceName).Key("github_repo.0.root_folder").HasValue("/"),
			),
		},
	})
}

func TestAccSynapseWorkspace_customerManagedKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace", "test")
	r := SynapseWorkspaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.customerManagedKey(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("customer_managed_key_versionless_id").Exists(),
			),
		},
		data.ImportStep("sql_administrator_login_password"),
	})
}

func (r SynapseWorkspaceResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.WorkspaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Synapse.WorkspaceClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Synapse Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(true), nil
}

func (r SynapseWorkspaceResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
}
`, template, data.RandomInteger)
}

func (r SynapseWorkspaceResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_workspace" "import" {
  name                                 = azurerm_synapse_workspace.test.name
  resource_group_name                  = azurerm_synapse_workspace.test.resource_group_name
  location                             = azurerm_synapse_workspace.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_synapse_workspace.test.storage_data_lake_gen2_filesystem_id
  sql_administrator_login              = azurerm_synapse_workspace.test.sql_administrator_login
  sql_administrator_login_password     = azurerm_synapse_workspace.test.sql_administrator_login_password
}
`, config)
}

func (r SynapseWorkspaceResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
  managed_virtual_network_enabled      = true
  managed_resource_group_name          = "acctest-ManagedSynapse-%d"
  sql_identity_control_enabled         = true

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r SynapseWorkspaceResource) withUpdateFields(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {
}

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR4!"
  sql_identity_control_enabled         = true

  aad_admin {
    login     = "AzureAD Admin"
    object_id = data.azurerm_client_config.current.object_id
    tenant_id = data.azurerm_client_config.current.tenant_id
  }

  tags = {
    ENV = "Test2"
  }
}
`, template, data.RandomInteger)
}

func (r SynapseWorkspaceResource) azdo(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"

  azure_devops_repo {
    account_name    = "myorg"
    project_name    = "myproj"
    repository_name = "myrepo"
    branch_name     = "dev"
    root_folder     = "/"
  }
}
`, template, data.RandomInteger)
}

func (r SynapseWorkspaceResource) github(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"

  github_repo {
    account_name    = "myuser"
    git_url         = "https://github.mydomain.com"
    repository_name = "myrepo"
    branch_name     = "dev"
    root_folder     = "/"
  }
}
`, template, data.RandomInteger)
}

func (r SynapseWorkspaceResource) customerManagedKey(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                     = "acckv%d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id
    key_permissions = [
      "create",
      "get",
      "delete",
      "purge"
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "key"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts = [
    "unwrapKey",
    "wrapKey"
  ]
}

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
  customer_managed_key_versionless_id  = azurerm_key_vault_key.test.versionless_id
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r SynapseWorkspaceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-synapse-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%d"
  storage_account_id = azurerm_storage_account.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}
