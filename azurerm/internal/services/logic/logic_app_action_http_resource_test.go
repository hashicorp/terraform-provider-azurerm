package logic_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccLogicAppActionHttp_basic(t *testing.T) {
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

func TestAccLogicAppActionHttp_requiresImport(t *testing.T) {
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

func TestAccLogicAppActionHttp_headers(t *testing.T) {
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

func TestAccLogicAppActionHttp_disappears(t *testing.T) {
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
					testCheckAzureRMLogicAppWorkflowExists("azurerm_logic_app_workflow.test"),
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

func TestAccLogicAppActionHttp_runAfter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_action_http", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppActionHttp_runAfterCondition(data, "Succeeded"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppActionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogicAppActionHttp_runAfterCondition(data, "Failed"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppActionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMLogicAppActionHttp_basic(data acceptance.TestData) string {
	template := testAccAzureRMLogicAppActionHttp_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_http" "test" {
  name         = "action%d"
  logic_app_id = azurerm_logic_app_workflow.test.id
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
  name         = azurerm_logic_app_action_http.test.name
  logic_app_id = azurerm_logic_app_action_http.test.logic_app_id
  method       = azurerm_logic_app_action_http.test.method
  uri          = azurerm_logic_app_action_http.test.uri
}
`, template)
}

func testAccAzureRMLogicAppActionHttp_headers(data acceptance.TestData) string {
	template := testAccAzureRMLogicAppActionHttp_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_http" "test" {
  name         = "action%d"
  logic_app_id = azurerm_logic_app_workflow.test.id
  method       = "GET"
  uri          = "http://example.com/hello"

  headers = {
    "Hello"     = "World"
    "Something" = "New"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogicAppActionHttp_runAfterCondition(data acceptance.TestData, condition string) string {
	template := testAccAzureRMLogicAppActionHttp_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_http" "testp1" {
  name         = "action%dp1"
  logic_app_id = azurerm_logic_app_workflow.test.id
  method       = "GET"
  uri          = "http://example.com/hello"
}

resource "azurerm_logic_app_action_http" "testp2" {
  name         = "action%dp2"
  logic_app_id = azurerm_logic_app_workflow.test.id
  method       = "GET"
  uri          = "http://example.com/hello"
}

resource "azurerm_logic_app_action_http" "test" {
  name         = "action%d"
  logic_app_id = azurerm_logic_app_workflow.test.id
  method       = "GET"
  uri          = "http://example.com/hello"
  run_after {
    action_name   = azurerm_logic_app_action_http.testp1.name
    action_result = "%s"
  }
  run_after {
    action_name   = azurerm_logic_app_action_http.testp2.name
    action_result = "%s"
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, condition, condition)
}

func testAccAzureRMLogicAppActionHttp_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
