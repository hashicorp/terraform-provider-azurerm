// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package synapse_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SynapseWorkspaceDataSource struct{}

func TestAccDataSourceSynapseWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_synapse_workspace", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SynapseWorkspaceDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("storage_data_lake_gen2_filesystem_id").MatchesOtherKey(check.That("azurerm_storage_data_lake_gen2_filesystem.test").Key("id")),
				check.That(data.ResourceName).Key("sql_administrator_login").HasValue("sqladminuser"),
				check.That(data.ResourceName).Key("connectivity_endpoints.%").Exists(),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
			),
		},
	})
}

func TestAccDataSourceSynapseWorkspace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_synapse_workspace", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SynapseWorkspaceDataSource{}.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("storage_data_lake_gen2_filesystem_id").MatchesOtherKey(check.That("azurerm_storage_data_lake_gen2_filesystem.test").Key("id")),
				check.That(data.ResourceName).Key("sql_administrator_login").HasValue("sqladminuser"),
				check.That(data.ResourceName).Key("linking_allowed_for_aad_tenant_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("compute_subnet_id").MatchesOtherKey(check.That("azurerm_subnet.test").Key("id")),
				check.That(data.ResourceName).Key("data_exfiltration_protection_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("managed_virtual_network_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("managed_resource_group_name").HasValue(fmt.Sprintf("acctest-ManagedSynapse-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("purview_id").MatchesOtherKey(check.That("azurerm_purview_account.test").Key("id")),
				check.That(data.ResourceName).Key("sql_identity_control_enabled").HasValue("true"),
			),
		},
	})
}

func TestAccDataSourceSynapseWorkspace_azureDevOps(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_synapse_workspace", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SynapseWorkspaceDataSource{}.azureDevOps(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("azure_devops_repo.#").HasValue("1"),
			),
		},
	})
}

func TestAccDataSourceSynapseWorkspace_github(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_synapse_workspace", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SynapseWorkspaceDataSource{}.github(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("github_repo.#").HasValue("1"),
			),
		},
	})
}

func TestAccDataSourceSynapseWorkspace_customerManagedKeyAndAzureAdOnlyAuthentication(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_synapse_workspace", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SynapseWorkspaceDataSource{}.customerManagedKeyAndAzureAdOnlyAuthentication(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("customer_managed_key.#").HasValue("1"),
				check.That(data.ResourceName).Key("azuread_authentication_only").HasValue("true"),
			),
		},
	})
}

func (d SynapseWorkspaceDataSource) basic(data acceptance.TestData) string {
	config := SynapseWorkspaceResource{}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_synapse_workspace" "test" {
  name                = azurerm_synapse_workspace.test.name
  resource_group_name = azurerm_synapse_workspace.test.resource_group_name
}
`, config)
}

func (d SynapseWorkspaceDataSource) complete(data acceptance.TestData) string {
	config := SynapseWorkspaceResource{}.complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_synapse_workspace" "test" {
  name                = azurerm_synapse_workspace.test.name
  resource_group_name = azurerm_synapse_workspace.test.resource_group_name
}
`, config)
}

func (d SynapseWorkspaceDataSource) azureDevOps(data acceptance.TestData) string {
	config := SynapseWorkspaceResource{}.azureDevOps(data)
	return fmt.Sprintf(`
%s

data "azurerm_synapse_workspace" "test" {
  name                = azurerm_synapse_workspace.test.name
  resource_group_name = azurerm_synapse_workspace.test.resource_group_name
}
`, config)
}

func (d SynapseWorkspaceDataSource) github(data acceptance.TestData) string {
	config := SynapseWorkspaceResource{}.github(data)
	return fmt.Sprintf(`
%s

data "azurerm_synapse_workspace" "test" {
  name                = azurerm_synapse_workspace.test.name
  resource_group_name = azurerm_synapse_workspace.test.resource_group_name
}
`, config)
}

func (d SynapseWorkspaceDataSource) customerManagedKeyAndAzureAdOnlyAuthentication(data acceptance.TestData) string {
	config := SynapseWorkspaceResource{}.cmkWithAADAdmin(data)
	return fmt.Sprintf(`
%s

data "azurerm_synapse_workspace" "test" {
  name                = azurerm_synapse_workspace.test.name
  resource_group_name = azurerm_synapse_workspace.test.resource_group_name
}
`, config)
}
