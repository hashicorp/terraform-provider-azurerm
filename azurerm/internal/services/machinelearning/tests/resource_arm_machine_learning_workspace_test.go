package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/machinelearning/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMachineLearningWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMachineLearningWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMachineLearningWorkspace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMachineLearningWorkspaceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "key_vault_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "application_insights_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMachineLearningWorkspace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMachineLearningWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMachineLearningWorkspace_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMachineLearningWorkspaceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "key_vault_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "application_insights_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "friendly_name", "test-workspace"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Test machine learning workspace"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMachineLearningWorkspace_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMachineLearningWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMachineLearningWorkspace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMachineLearningWorkspaceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "key_vault_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "application_insights_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.#", "1"),
				),
			},
			{
				Config: testAccAzureRMMachineLearningWorkspace_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMachineLearningWorkspaceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "key_vault_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "application_insights_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "friendly_name", "test-workspace"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Test machine learning workspace"),
				),
			},
			{
				Config: testAccAzureRMMachineLearningWorkspace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMachineLearningWorkspaceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "key_vault_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "application_insights_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMachineLearningWorkspace_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMachineLearningWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMachineLearningWorkspace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMachineLearningWorkspaceExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMMachineLearningWorkspace_requiresImport),
		},
	})
}

func testCheckAzureRMMachineLearningWorkspaceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Machine Learning Workspace not found: %s", resourceName)
		}

		id, err := parse.WorkspaceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).MachineLearning.WorkspacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Machine Learning Workspace %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Bad: Get on machinelearningservices.WorkspacesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMMachineLearningWorkspaceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MachineLearning.WorkspacesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_machine_learning_workspace" {
			continue
		}

		id, err := parse.WorkspaceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on machinelearningservices.WorkspacesClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMMachineLearningWorkspace_basic(data acceptance.TestData) string {
	template := testAccAzureRMMachineLearningWorkspace_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acctestworkspace-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomIntOfLength(16))
}

func testAccAzureRMMachineLearningWorkspace_complete(data acceptance.TestData) string {
	template := testAccAzureRMMachineLearningWorkspace_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acctestworkspace-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  friendly_name           = "test-workspace"
  description             = "Test machine learning workspace"
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomIntOfLength(16))
}

func testAccAzureRMMachineLearningWorkspace_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMachineLearningWorkspace_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_workspace" "import" {
  name                    = azurerm_machine_learning_workspace.test.name
  location                = azurerm_machine_learning_workspace.test.location
  resource_group_name     = azurerm_machine_learning_workspace.test.resource_group_name
  application_insights_id = azurerm_machine_learning_workspace.test.application_insights_id
  key_vault_id            = azurerm_machine_learning_workspace.test.key_vault_id
  storage_account_id      = azurerm_machine_learning_workspace.test.storage_account_id

  identity {
    type = "SystemAssigned"
  }
}
`, template)
}

func testAccAzureRMMachineLearningWorkspace_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ml-%[1]d"
  location = "%[2]s"
}

resource "azurerm_application_insights" "test" {
  name                 = "acctestai-%[1]d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  application_type     = "web"
}

resource "azurerm_key_vault" "test" {
  name                 = "acctestvault%[3]d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  tenant_id            = data.azurerm_client_config.current.tenant_id

  sku_name = "premium"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[4]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(12), data.RandomIntOfLength(15))
}
