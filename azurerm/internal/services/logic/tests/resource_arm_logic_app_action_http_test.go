package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMLogicAppActionHttp_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_action_http", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppActionHttp_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppActionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogicAppActionHttp_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_logic_app_action_http", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppActionHttp_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppActionExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMLogicAppActionHttp_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_logic_app_action_http"),
			},
		},
	})
}

func TestAccAzureRMLogicAppActionHttp_headers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_action_http", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppActionHttp_headers(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppActionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogicAppActionHttp_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_action_http", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppActionHttp_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppActionExists(data.ResourceName),
				),
			},
			{
				// delete it
				Config: testAccAzureRMLogicAppActionHttp_template(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppWorkflowExists("azurerm_logic_app_workflow"),
				),
			},
			{
				Config:             testAccAzureRMLogicAppActionHttp_basic(data),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccAzureRMLogicAppActionHttp_basic(data acceptance.TestData) string {
	template := testAccAzureRMLogicAppActionHttp_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_http" "test" {
  name         = "action%d"
  logic_app_id = "${azurerm_logic_app_workflow.test.id}"
  method       = "GET"
  uri          = "http://example.com/hello"
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogicAppActionHttp_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLogicAppActionHttp_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_http" "import" {
  name         = "${azurerm_logic_app_action_http.test.name}"
  logic_app_id = "${azurerm_logic_app_action_http.test.logic_app_id}"
  method       = "${azurerm_logic_app_action_http.test.method}"
  uri          = "${azurerm_logic_app_action_http.test.uri}"
}
`, template)
}

func testAccAzureRMLogicAppActionHttp_headers(data acceptance.TestData) string {
	template := testAccAzureRMLogicAppActionHttp_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_http" "test" {
  name         = "action%d"
  logic_app_id = "${azurerm_logic_app_workflow.test.id}"
  method       = "GET"
  uri          = "http://example.com/hello"

  headers = {
    "Hello"     = "World"
    "Something" = "New"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogicAppActionHttp_template(data acceptance.TestData) string {
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
