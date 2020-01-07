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

func TestAccAzureRMLogicAppWorkflow_empty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_workflow", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppWorkflow_empty(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppWorkflowExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "parameters.%", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogicAppWorkflow_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_logic_app_workflow", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppWorkflow_empty(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppWorkflowExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMLogicAppWorkflow_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_logic_app_workflow"),
			},
		},
	})
}

func TestAccAzureRMLogicAppWorkflow_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_workflow", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppWorkflow_empty(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppWorkflowExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "parameters.%", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogicAppWorkflow_tags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppWorkflowExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "parameters.%", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Source", "AcceptanceTests"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMLogicAppWorkflowExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Logic.WorkflowsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		workflowName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Logic App Workflow: %s", workflowName)
		}

		resp, err := client.Get(ctx, resourceGroup, workflowName)
		if err != nil {
			return fmt.Errorf("Bad: Get on logicWorkflowsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Logic App Workflow %q (resource group %q) does not exist", workflowName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMLogicAppWorkflowDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Logic.WorkflowsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_logic_app_workflow" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Logic App Workflow still exists: \n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMLogicAppWorkflow_empty(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMLogicAppWorkflow_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLogicAppWorkflow_empty(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_workflow" "import" {
  name                = "${azurerm_logic_app_workflow.test.name}"
  location            = "${azurerm_logic_app_workflow.test.location}"
  resource_group_name = "${azurerm_logic_app_workflow.test.resource_group_name}"
}
`, template)
}

func testAccAzureRMLogicAppWorkflow_tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    "Source" = "AcceptanceTests"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
