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

type SparkJobDefinitionResource struct {
}

func TestAccSynapseSparkJobDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_spark_job_definition", "test")
	r := SparkJobDefinitionResource{}

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

func TestAccSynapseSparkJobDefinition_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_spark_job_definition", "test")
	r := SparkJobDefinitionResource{}

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

func TestAccSynapseSparkJobDefinition_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_spark_job_definition", "test")
	r := SparkJobDefinitionResource{}

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

func TestAccSynapseSparkJobDefinition_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_spark_job_definition", "test")
	r := SparkJobDefinitionResource{}

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

func (r SparkJobDefinitionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SparkJobDefinitionID(state.ID)
	if err != nil {
		return nil, err
	}

	environment := clients.Account.Environment
	client, err := clients.Synapse.SparkJobDefinitionClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return nil, err
	}

	resp, err := client.GetSparkJobDefinition(ctx, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r SparkJobDefinitionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_spark_job_definition" "test" {
  name                 = "acctestls%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id

  depends_on = [
    azurerm_synapse_firewall_rule.test,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r SparkJobDefinitionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_spark_job_definition" "import" {
  name                 = azurerm_synapse_spark_job_definition.test.name
  synapse_workspace_id = azurerm_synapse_spark_job_definition.test.synapse_workspace_id
}
`, r.basic(data))
}

func (r SparkJobDefinitionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_spark_job_definition" "test" {
  name                 = "acctestls%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  description          = "test"
  language             = "Java"
  spark_pool_name      = azurerm_synapse_spark_pool.test.name
  spark_version        = "2.4"
  job {
    file       = "abfss://test@test.dfs.core.windows.net/artefacts/sample.jar"
    class_name = "dev.test.tools.sample.Main"
    arguments = [
      "exampleArg"
    ]
    jars            = []
    files           = []
    archives        = []
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

func (SparkJobDefinitionResource) template(data acceptance.TestData) string {
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
