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

type SqlScriptResource struct {
}

func TestAccSynapseSqlScript_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_script", "test")
	r := SqlScriptResource{}

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

func TestAccSynapseSqlScript_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_script", "test")
	r := SqlScriptResource{}

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

func TestAccSynapseSqlScript_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_script", "test")
	r := SqlScriptResource{}

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

func TestAccSynapseSqlScript_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_script", "test")
	r := SqlScriptResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
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
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SqlScriptResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SqlScriptID(state.ID)
	if err != nil {
		return nil, err
	}

	environment := clients.Account.Environment
	client, err := clients.Synapse.SQLScriptClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return nil, err
	}

	resp, err := client.GetSQLScript(ctx, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r SqlScriptResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_sql_script" "test" {
  name                 = "acctestls%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id

  depends_on = [
    azurerm_synapse_firewall_rule.test,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r SqlScriptResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_sql_script" "import" {
  name                 = azurerm_synapse_sql_script.test.name
  synapse_workspace_id = azurerm_synapse_sql_script.test.synapse_workspace_id
}
`, r.basic(data))
}

func (r SqlScriptResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_sql_script" "test" {
  name                 = "acctestls%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  description          = "test"
  language             = "sql"
  query                = "SELECT TOP 100 * FROM example_table_name;"
  sql_connection {
    name = "master"
    type = "SqlOnDemand"
  }

  depends_on = [
    azurerm_synapse_firewall_rule.test,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (SqlScriptResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-synapse-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
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

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
  managed_virtual_network_enabled      = true
}

resource "azurerm_synapse_firewall_rule" "test" {
  name                 = "allowAll"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  start_ip_address     = "0.0.0.0"
  end_ip_address       = "255.255.255.255"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}
