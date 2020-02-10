package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMMachineLearningWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMachineLearningWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMachineLearningWorkspaceBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMachineLearningWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMachineLearingWorkspace_requiresImport(t *testing.T) {
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
				Config: testAccAzureRMMachineLearningWorkspaceBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMachineLearningWorkspaceExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMMachineLearningWorkspaceRequiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_machine_learning_workspace"),
			},
		},
	})
}

func TestAccAzureRMMachineLearningWorkspaceWithTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMachineLearningWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMachineLearningWorkspaceWithTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMachineLearningWorkspaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
				),
			},
			{
				Config: testAccAzureRMMachineLearningWorkspaceWithTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMachineLearningWorkspaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMachineLearningWorkspaceWithContainerRegistry(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMachineLearningWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMachineLearningWorkspaceWithContainerRegistry(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMachineLearningWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMachineLearningWorkspaceWithStorageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMachineLearningWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMachineLearningWorkspaceWithStorageAccount(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMachineLearningWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMachineLearningWorkspaceWithKeyVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMachineLearningWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMachineLearningWorkspaceWithKeyVault(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMachineLearningWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMachineLearningWorkspaceWithApplicationInsights(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMachineLearningWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMachineLearningWorkspaceWithApplicationInsights(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMachineLearningWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMachineLearningWorkspaceFull(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMachineLearningWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMachineLearningWorkspaceFull(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMachineLearningWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMMachineLearningWorkspaceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MachineLearning.WorkspacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		wsName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Machine Learning Workspace: %s", wsName)
		}

		resp, err := client.Get(ctx, resourceGroup, wsName)
		if err != nil {
			return fmt.Errorf("Bad: Get Machine Learning Workspace Client: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Machine Learning Workspace %s (resource group: %s) does not exist", wsName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMMachineLearningWorkspaceDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).MachineLearning.WorkspacesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_machine_learning_workspace" {
			continue
		}

		wsName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, wsName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("Machine Learning Workspace still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMMachineLearningWorkspaceRequiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMachineLearningWorkspaceBasic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_workspace" "import" {
  name                = azurerm_machine_learning_workspace.test.name
  resource_group_name = azurerm_machine_learning_workspace.test.resource_group_name
}
`, template)
}

func testAccAzureRMMachineLearningWorkspaceBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
  description = "aml workspace description"
  friendly_name = "aml workspace friendly name"
}

resource "azurerm_machine_learning_workspace" "test" {
  name                = "acctestworkspace%d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMachineLearningWorkspaceWithContainerRegistry(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_registry" "acrtest" {
	name                = "acctestacr%d"
	resource_group_name = azurerm_resource_group.test.name
	location            = "%s"
	sku                 = "Standard"
	admin_enabled       = true
}

resource "azurerm_machine_learning_workspace" "test" {
  name                = "acctestworkspace%d"
  resource_group_name              = azurerm_resource_group.test.name
  container_registry               = azurerm_container_registry.acttest.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMachineLearningWorkspaceWithStorageAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%d"
	location = "%s"
}

resource "azurerm_storage_account" "stortest" {
	name                      = "accteststor%d"
	resource_group_name       = azurerm_resource_group.test.name
	location                  = azurerm_resource_group.test.location
	account_kind              = "StorageV2"
	account_tier              = "Standard"
	account_replication_type  = "LRS"
}

resource "azurerm_machine_learning_workspace" "amltest" {
	name                = "acctestworkspace%d"
	resource_group_name = azurerm_resource_group.test.name
	storage_account     = azurerm_storage_account.stortest.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMMachineLearningWorkspaceWithApplicationInsights(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%d"
	location = "%s"
}

resource "azurerm_storage_account" "aitest" {
	name                 = "acctestai%d"
	resource_group_name  = azurerm_resource_group.test.name
	location             = azurerm_resource_group.test.location
	application_type     = "web"
}

resource "azurerm_machine_learning_workspace" "amltest" {
	name                	= "acctestworkspace%d"
	resource_group_name   = azurerm_resource_group.test.name
	application_insights  = azurerm_storage_account.aitest.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMMachineLearningWorkspaceWithKeyVault(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%d"
	location = "%s"
}

resource "azurerm_storage_account" "kvtest" {
	name                 = "acctestkv%d"
	resource_group_name  = azurerm_resource_group.test.name
	location             = azurerm_resource_group.test.location
	enabled_for_disk_encryption     = true
    tenant_id                       = "%s"
    enabled_for_template_deployment = false
    enabled_for_deployment          = false
    sku_name                        = "Standard"
}

resource "azurerm_machine_learning_workspace" "amltest" {
	name                  = "acctestworkspace%d"
	resource_group_name   = azurerm_resource_group.test.name
	application_insights  = azurerm_storage_account.aitest.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Client().TenantID, data.RandomInteger)
}

func testAccAzureRMMachineLearningWorkspaceWithTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%d"
	location = "%s"
}

resource "azurerm_machine_learning_workspace" "amltest" {
  name                = "acctestworkspace%d"
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMachineLearningWorkspaceWithTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_machine_learning_workspace" "amltest" {
  name                = "acctestworkspace%d"
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMachineLearningWorkspaceFull(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG_%d"
	location = "%s"
}

resource "azurerm_container_registry" "acrtest" {
	name                = "acctestacr%d"
	resource_group_name = azurerm_resource_group.test.name
	location            = azurerm_resource_group.test.location
	sku                 = "Standard"
	admin_enabled       = true
}

resource "azurerm_application_insights" "aitest" {
	name                 = "aitest%d"
	resource_group_name  = azurerm_resource_group.test.name
	location             = azurerm_resource_group.test.location
	application_type     = "web"
}

resource "azurerm_storage_account" "stortest" {
	name                      = "stortest%d"
	resource_group_name       = azurerm_resource_group.test.name
	location                  = azurerm_resource_group.test.location
	account_kind              = "StorageV2"
	account_tier              = "Standard"
	account_replication_type  = "LRS"
}

resource "azurerm_key_vault" "kvtest" {
	name                 = "kvtest%d"
	resource_group_name  = azurerm_resource_group.test.name
	location             = azurerm_resource_group.test.location
	enabled_for_disk_encryption     = true
    tenant_id                       = "%s"
    enabled_for_template_deployment = false
    enabled_for_deployment          = false
    sku_name                        = "Standard"
}

resource "azurerm_machine_learning_workspace" "amltest" {
	name                  = "acctestworkspace%d"
	resource_group_name   = azurerm_resource_group.test.name
	application_insights  = azurerm_storage_account.aitest.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.Client().TenantID, data.RandomInteger)
}
