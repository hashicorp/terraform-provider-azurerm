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

func TestAccAzureRMSynapseWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSynapseWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSynapseWorkspace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep("sql_administrator_login_password"),
		},
	})
}

func TestAccAzureRMSynapseWorkspace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSynapseWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSynapseWorkspace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseWorkspaceExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMSynapseWorkspace_requiresImport),
		},
	})
}

func TestAccAzureRMSynapseWorkspace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace", "test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSynapseWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSynapseWorkspace_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep("sql_administrator_login_password"),
		},
	})
}

func TestAccAzureRMSynapseWorkspace_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSynapseWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSynapseWorkspace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep("sql_administrator_login_password"),
			{
				Config: testAccAzureRMSynapseWorkspace_withUpdateFields(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep("sql_administrator_login_password"),
			{
				Config: testAccAzureRMSynapseWorkspace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseWorkspaceExists(data.ResourceName),
				),
			},
		},
	})
}

func testCheckAzureRMSynapseWorkspaceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Synapse.WorkspaceClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("synapse Workspace not found: %s", resourceName)
		}
		id, err := parse.WorkspaceID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Synapse Workspace %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on Synapse.WorkspaceClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMSynapseWorkspaceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Synapse.WorkspaceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_synapse_workspace" {
			continue
		}
		id, err := parse.WorkspaceID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Synapse.WorkspaceClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMSynapseWorkspace_basic(data acceptance.TestData) string {
	template := testAccAzureRMSynapseWorkspace_template(data)
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

func testAccAzureRMSynapseWorkspace_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMSynapseWorkspace_basic(data)
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

func testAccAzureRMSynapseWorkspace_complete(data acceptance.TestData) string {
	template := testAccAzureRMSynapseWorkspace_template(data)
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

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMSynapseWorkspace_withUpdateFields(data acceptance.TestData) string {
	template := testAccAzureRMSynapseWorkspace_template(data)
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

  aad_admin {
    login     = "AzureAD Admin"
    object_id = data.azurerm_client_config.current.object_id
    tenant_id = data.azurerm_client_config.current.tenant_id
  }

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMSynapseWorkspace_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-Synapse-%d"
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
