package synapse_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/synapse/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SynapseSqlPoolResource struct{}

func TestAccSynapseSqlPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool", "test")
	r := SynapseSqlPoolResource{}

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

func TestAccSynapseSqlPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool", "test")
	r := SynapseSqlPoolResource{}

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

func TestAccSynapseSqlPool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool", "test")
	r := SynapseSqlPoolResource{}

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

func TestAccSynapseSqlPool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool", "test")
	r := SynapseSqlPoolResource{}

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

func (r SynapseSqlPoolResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SqlPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Synapse.SqlPoolClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Synapse Sql Pool %q (Workspace %q / Resource Group %q): %+v", id.Name, id.WorkspaceName, id.ResourceGroup, err)
	}

	return utils.Bool(true), nil
}

func (r SynapseSqlPoolResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_synapse_sql_pool" "test" {
  name                 = "acctestSP%s"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  sku_name             = "DW100c"
  create_mode          = "Default"
}
`, template, data.RandomString)
}

func (r SynapseSqlPoolResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_sql_pool" "import" {
  name                 = azurerm_synapse_sql_pool.test.name
  synapse_workspace_id = azurerm_synapse_sql_pool.test.synapse_workspace_id
  sku_name             = azurerm_synapse_sql_pool.test.sku_name
  create_mode          = azurerm_synapse_sql_pool.test.create_mode
}
`, config)
}

func (r SynapseSqlPoolResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_synapse_sql_pool" "test" {
  name                 = "acctestSP%s"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  sku_name             = "DW500c"
  create_mode          = "Default"
  collation            = "SQL_Latin1_General_CP1_CI_AS"
  data_encrypted       = true

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomString)
}

func (r SynapseSqlPoolResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
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

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}
