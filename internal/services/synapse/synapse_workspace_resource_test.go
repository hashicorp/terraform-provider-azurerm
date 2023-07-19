// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SynapseWorkspaceResource struct{}

func TestAccSynapseWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace", "test")
	r := SynapseWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
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

func TestAccSynapseWorkspace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace", "test")
	r := SynapseWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("managed_resource_group_name").HasValue(fmt.Sprintf("acctest-ManagedSynapse-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("sql_identity_control_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("data_exfiltration_protection_enabled").HasValue("true"),
			),
		},
		data.ImportStep("sql_administrator_login_password"),
	})
}

func TestAccSynapseWorkspace_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace", "test")
	r := SynapseWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sql_administrator_login_password"),
		{
			Config: r.withAadAdmin(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sql_administrator_login_password"),
		{
			Config: r.withAadAdmins(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sql_administrator_login_password"),
		{
			Config: r.withSqlAadAdmin(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sql_administrator_login_password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccSynapseWorkspace_azdo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace", "test")
	r := SynapseWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureDevOps(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("azure_devops_repo.0.account_name").HasValue("myorg"),
				check.That(data.ResourceName).Key("azure_devops_repo.0.project_name").HasValue("myproj"),
				check.That(data.ResourceName).Key("azure_devops_repo.0.repository_name").HasValue("myrepo"),
				check.That(data.ResourceName).Key("azure_devops_repo.0.branch_name").HasValue("dev"),
				check.That(data.ResourceName).Key("azure_devops_repo.0.root_folder").HasValue("/"),
				check.That(data.ResourceName).Key("azure_devops_repo.0.tenant_id").IsEmpty(),
			),
		},
	})
}

func TestAccSynapseWorkspace_azdoTenant(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace", "test")
	r := SynapseWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureDevOpsTenant(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("azure_devops_repo.0.account_name").HasValue("myorg"),
				check.That(data.ResourceName).Key("azure_devops_repo.0.project_name").HasValue("myproj"),
				check.That(data.ResourceName).Key("azure_devops_repo.0.repository_name").HasValue("myrepo"),
				check.That(data.ResourceName).Key("azure_devops_repo.0.branch_name").HasValue("dev"),
				check.That(data.ResourceName).Key("azure_devops_repo.0.root_folder").HasValue("/"),
				check.That(data.ResourceName).Key("azure_devops_repo.0.tenant_id").Exists(),
			),
		},
	})
}

func TestAccSynapseWorkspace_github(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace", "test")
	r := SynapseWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.github(data),
			Check: acceptance.ComposeTestCheckFunc(
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

func TestAccSynapseWorkspace_customerManagedKeyActivation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace", "test")
	r := SynapseWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customerManagedKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("customer_managed_key.0.key_versionless_id").Exists(),
			),
		},
		data.ImportStep("sql_administrator_login_password"),
	})
}

func (r SynapseWorkspaceResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

data "azurerm_client_config" "current" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuaid%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, template, data.RandomInteger, data.RandomInteger)
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

  identity {
    type = "SystemAssigned"
  }
}
`, config)
}

func (r SynapseWorkspaceResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}


resource "azurerm_purview_account" "test" {
  name                = "acctestacc%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "%s-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "%s-subnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
  data_exfiltration_protection_enabled = true
  managed_virtual_network_enabled      = true
  managed_resource_group_name          = "acctest-ManagedSynapse-%d"
  sql_identity_control_enabled         = true
  public_network_access_enabled        = false
  linking_allowed_for_aad_tenant_ids   = [data.azurerm_client_config.current.tenant_id]
  purview_id                           = azurerm_purview_account.test.id
  compute_subnet_id                    = azurerm_subnet.test.id

  identity {
    type = "SystemAssigned"
  }

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomString, data.Locations.Secondary, data.RandomString, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r SynapseWorkspaceResource) withAadAdmin(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuaid%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
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

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    ENV = "Test2"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r SynapseWorkspaceResource) withSqlAadAdmin(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuaid%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR4!"
  sql_identity_control_enabled         = true

  sql_aad_admin {
    login     = "AzureAD Admin"
    object_id = data.azurerm_client_config.current.object_id
    tenant_id = data.azurerm_client_config.current.tenant_id
  }

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    ENV = "Test2"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r SynapseWorkspaceResource) withAadAdmins(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuaid%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
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

  sql_aad_admin {
    login     = "AzureAD Admin"
    object_id = data.azurerm_client_config.current.object_id
    tenant_id = data.azurerm_client_config.current.tenant_id
  }

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    ENV = "Test2"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r SynapseWorkspaceResource) azureDevOps(data acceptance.TestData) string {
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
    last_commit_id  = "1592393b38543d51feb12714cbd39501d697610c"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func (r SynapseWorkspaceResource) azureDevOpsTenant(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

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
    tenant_id       = data.azurerm_client_config.current.tenant_id
  }

  identity {
    type = "SystemAssigned"
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
    last_commit_id  = "1592393b38543d51feb12714cbd39501d697610c"
  }

  identity {
    type = "SystemAssigned"
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
      "Create",
      "Get",
      "Delete",
      "Purge",
      "GetRotationPolicy",
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

  customer_managed_key {
    key_versionless_id = azurerm_key_vault_key.test.versionless_id
  }

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r SynapseWorkspaceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
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
