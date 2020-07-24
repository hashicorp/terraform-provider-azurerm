package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/synapse/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSynapseBigDataPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_big_data_pool", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSynapseBigDataPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSynapseBigDataPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseBigDataPoolExists(data.ResourceName),
				),
			},
			// not returned by service
			data.ImportStep("spark_events_folder", "spark_log_folder"),
		},
	})
}

func TestAccAzureRMSynapseBigDataPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_big_data_pool", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSynapseBigDataPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSynapseBigDataPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseBigDataPoolExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMSynapseBigDataPool_requiresImport),
		},
	})
}

func TestAccAzureRMSynapseBigDataPool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_big_data_pool", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSynapseBigDataPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSynapseBigDataPool_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseBigDataPoolExists(data.ResourceName),
				),
			},
			// not returned by service
			data.ImportStep("spark_events_folder", "spark_log_folder"),
		},
	})
}

func TestAccAzureRMSynapseBigDataPool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_big_data_pool", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSynapseBigDataPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSynapseBigDataPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseBigDataPoolExists(data.ResourceName),
				),
			},
			data.ImportStep("spark_events_folder", "spark_log_folder"),
			{
				Config: testAccAzureRMSynapseBigDataPool_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseBigDataPoolExists(data.ResourceName),
				),
			},
			data.ImportStep("spark_events_folder", "spark_log_folder"),
			{
				Config: testAccAzureRMSynapseBigDataPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseBigDataPoolExists(data.ResourceName),
				),
			},
			data.ImportStep("spark_events_folder", "spark_log_folder"),
		},
	})
}

func testCheckAzureRMSynapseBigDataPoolExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Synapse.BigDataPoolClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("synapse BigDataPool not found: %s", resourceName)
		}
		id, err := parse.SynapseBigDataPoolID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Synapse BigDataPool %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on Synapse.BigDataPoolClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMSynapseBigDataPoolDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Synapse.BigDataPoolClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_synapse_big_data_pool" {
			continue
		}
		id, err := parse.SynapseBigDataPoolID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Synapse.BigDataPoolClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMSynapseBigDataPool_basic(data acceptance.TestData) string {
	template := testAccAzureRMSynapseBigDataPool_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_big_data_pool" "test" {
  name                 = "acctestSSP%s"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  node_size_family     = "MemoryOptimized"
  node_size            = "Small"
  node_count           = 3
}
`, template, data.RandomString)
}

func testAccAzureRMSynapseBigDataPool_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMSynapseBigDataPool_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_big_data_pool" "import" {
  name                 = azurerm_synapse_big_data_pool.test.name
  synapse_workspace_id = azurerm_synapse_big_data_pool.test.synapse_workspace_id
  node_size_family     = azurerm_synapse_big_data_pool.test.node_size_family
  node_size            = azurerm_synapse_big_data_pool.test.node_size
  node_count           = azurerm_synapse_big_data_pool.test.node_count
}
`, config)
}

func testAccAzureRMSynapseBigDataPool_complete(data acceptance.TestData) string {
	template := testAccAzureRMSynapseBigDataPool_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_big_data_pool" "test" {
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
  spark_version       = "2.4"

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomString)
}

func testAccAzureRMSynapseBigDataPool_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-SW-%d"
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
