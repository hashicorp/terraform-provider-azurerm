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

func TestAccAzureRMSynapseSqlPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSynapseSqlPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSynapseSqlPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseSqlPoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSynapseSqlPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSynapseSqlPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSynapseSqlPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseSqlPoolExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMSynapseSqlPool_requiresImport),
		},
	})
}

func TestAccAzureRMSynapseSqlPool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSynapseSqlPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSynapseSqlPool_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseSqlPoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSynapseSqlPool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSynapseSqlPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSynapseSqlPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseSqlPoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSynapseSqlPool_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseSqlPoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSynapseSqlPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseSqlPoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMSynapseSqlPoolExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Synapse.SqlPoolClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("synapse SqlPool not found: %s", resourceName)
		}
		id, err := parse.SqlPoolID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.Workspace.ResourceGroup, id.Workspace.Name, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Synapse SqlPool %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on Synapse.SqlPoolClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMSynapseSqlPoolDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Synapse.SqlPoolClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_synapse_sql_pool" {
			continue
		}
		id, err := parse.SqlPoolID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.Workspace.ResourceGroup, id.Workspace.Name, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Synapse.SqlPoolClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMSynapseSqlPool_basic(data acceptance.TestData) string {
	template := testAccAzureRMSynapseSqlPool_template(data)
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

func testAccAzureRMSynapseSqlPool_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMSynapseSqlPool_basic(data)
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

func testAccAzureRMSynapseSqlPool_complete(data acceptance.TestData) string {
	template := testAccAzureRMSynapseSqlPool_template(data)
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

func testAccAzureRMSynapseSqlPool_template(data acceptance.TestData) string {
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
