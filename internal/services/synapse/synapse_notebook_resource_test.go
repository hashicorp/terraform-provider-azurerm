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

type SynapseNotebookResource struct {
}

func TestAccSynapseNotebook_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_notebook", "test")
	r := SynapseNotebookResource{}

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

func TestAccSynapseNotebook_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_notebook", "test")
	r := SynapseNotebookResource{}

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

func TestAccSynapseNotebook_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_notebook", "test")
	r := SynapseNotebookResource{}

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

func TestAccSynapseNotebook_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_notebook", "test")
	r := SynapseNotebookResource{}

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

func (r SynapseNotebookResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.NotebookID(state.ID)
	if err != nil {
		return nil, err
	}

	environment := clients.Account.Environment
	client, err := clients.Synapse.NotebookClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return nil, err
	}

	resp, err := client.GetNotebook(ctx, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r SynapseNotebookResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_notebook" "test" {
  name                 = "acctestls%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id

  cells = <<BODY
[
            {
              "cell_type": "code",
              "metadata": {},
              "source": [
                "def my_function():\n",
                " print(\"Hello from a function\")\n",
                "\n",
                "my_function()"
              ],
              "attachments": {},
              "outputs": [
                {
                  "execution_count": 3,
                  "output_type": "execute_result",
                  "data": {
                    "text/plain": "Hello from a function"
                  },
                  "metadata": {}
                }
              ]
            }
          ]
BODY

  depends_on = [
    azurerm_synapse_firewall_rule.test,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r SynapseNotebookResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_notebook" "import" {
  name                 = azurerm_synapse_notebook.test.name
  synapse_workspace_id = azurerm_synapse_notebook.test.synapse_workspace_id
  cells                = azurerm_synapse_notebook.test.cells
}
`, r.basic(data))
}

func (r SynapseNotebookResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_notebook" "test" {
  name                 = "acctestls%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  cells                = <<BODY
[
            {
              "cell_type": "code",
              "metadata": {},
              "source": [
                "def my_function():\n",
                " print(\"Hello from a function\")\n",
                "\n",
                "my_function()"
              ],
              "attachments": {},
              "outputs": [
                {
                  "execution_count": 3,
                  "output_type": "execute_result",
                  "data": {
                    "text/plain": "Hello from a function"
                  },
                  "metadata": {}
                }
              ]
            }
          ]
BODY

  description     = "test"
  display_name    = "notebook test"
  language        = "python"
  major_version   = 4
  minor_version   = 2
  spark_pool_name = azurerm_synapse_spark_pool.test.name
  session_config {
    driver_memory   = "28g"
    driver_cores    = 4
    executor_memory = "28g"
    executor_cores  = 4
    num_executors   = 2
  }

  depends_on = [
    azurerm_synapse_firewall_rule.test,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (SynapseNotebookResource) template(data acceptance.TestData) string {
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

resource "azurerm_synapse_spark_pool" "test" {
  name                 = "acctestSSP%s"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  node_size_family     = "MemoryOptimized"
  node_size            = "Small"
  node_count           = 3
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomString)
}
