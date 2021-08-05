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

type SynapseSparkPoolResource struct{}

func TestAccSynapseSparkPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_spark_pool", "test")
	r := SynapseSparkPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// not returned by service
		data.ImportStep("spark_events_folder", "spark_log_folder"),
	})
}

func TestAccSynapseSparkPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_spark_pool", "test")
	r := SynapseSparkPoolResource{}

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

func TestAccSynapseSparkPool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_spark_pool", "test")
	r := SynapseSparkPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "2.4"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// not returned by service
		data.ImportStep("spark_events_folder", "spark_log_folder"),
	})
}

func TestAccSynapseSparkPool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_spark_pool", "test")
	r := SynapseSparkPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("spark_events_folder", "spark_log_folder"),
		{
			Config: r.complete(data, "2.4"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("spark_events_folder", "spark_log_folder"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("spark_events_folder", "spark_log_folder"),
	})
}

func TestAccSynapseSpark3Pool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_spark_pool", "test")
	r := SynapseSparkPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "3.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// not returned by service
		data.ImportStep("spark_events_folder", "spark_log_folder"),
	})
}

func TestAccSynapseSpark3Pool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_spark_pool", "test")
	r := SynapseSparkPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("spark_events_folder", "spark_log_folder"),
		{
			Config: r.complete(data, "3.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("spark_events_folder", "spark_log_folder"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("spark_events_folder", "spark_log_folder"),
	})
}

func (r SynapseSparkPoolResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SparkPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Synapse.SparkPoolClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.BigDataPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Synapse Spark Pool %q (Workspace %q / Resource Group %q): %+v", id.BigDataPoolName, id.WorkspaceName, id.ResourceGroup, err)
	}

	return utils.Bool(true), nil
}

func (r SynapseSparkPoolResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_spark_pool" "test" {
  name                 = "acctestSSP%s"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  node_size_family     = "MemoryOptimized"
  node_size            = "Small"
  node_count           = 3
}
`, template, data.RandomString)
}

func (r SynapseSparkPoolResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_spark_pool" "import" {
  name                 = azurerm_synapse_spark_pool.test.name
  synapse_workspace_id = azurerm_synapse_spark_pool.test.synapse_workspace_id
  node_size_family     = azurerm_synapse_spark_pool.test.node_size_family
  node_size            = azurerm_synapse_spark_pool.test.node_size
  node_count           = azurerm_synapse_spark_pool.test.node_count
}
`, config)
}

func (r SynapseSparkPoolResource) complete(data acceptance.TestData, sparkVersion string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_spark_pool" "test" {
  name                 = "acctestSSP%s"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  node_size_family     = "MemoryOptimized"
  node_size            = "Medium"

  auto_pause {
    delay_in_minutes = 15
  }

  auto_scale {
    max_node_count = 50
    min_node_count = 3
  }

  library_requirement {
    content  = <<EOF
appnope==0.1.0
beautifulsoup4==4.6.3
EOF
    filename = "requirements.txt"
  }

  spark_log_folder    = "/logs"
  spark_events_folder = "/events"
  spark_version       = "%s"

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomString, sparkVersion)
}

func (r SynapseSparkPoolResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
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
