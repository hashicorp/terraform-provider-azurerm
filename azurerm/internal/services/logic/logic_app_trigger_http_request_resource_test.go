package logic_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccLogicAppTriggerHttpRequest_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_http_request", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerHttpRequest_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "schema", "{}"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLogicAppTriggerHttpRequest_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_http_request", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerHttpRequest_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMLogicAppTriggerHttpRequest_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_logic_app_trigger_http_request"),
			},
		},
	})
}

func TestAccLogicAppTriggerHttpRequest_fullSchema(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_http_request", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerHttpRequest_fullSchema(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "schema"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLogicAppTriggerHttpRequest_method(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_http_request", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerHttpRequest_method(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "method", "PUT"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLogicAppTriggerHttpRequest_relativePath(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_http_request", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerHttpRequest_relativePath(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "method", "POST"),
					resource.TestCheckResourceAttr(data.ResourceName, "relative_path", "customers/{id}"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLogicAppTriggerHttpRequest_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_http_request", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerHttpRequest_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(data.ResourceName),
				),
			},
			{
				// delete it
				Config: testAccAzureRMLogicAppTriggerHttpRequest_template(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppWorkflowExists("azurerm_logic_app_workflow.test"),
				),
			},
			{
				Config:             testAccAzureRMLogicAppTriggerHttpRequest_basic(data),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccAzureRMLogicAppTriggerHttpRequest_basic(data acceptance.TestData) string {
	template := testAccAzureRMLogicAppTriggerHttpRequest_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_http_request" "test" {
  name         = "some-http-trigger"
  logic_app_id = azurerm_logic_app_workflow.test.id
  schema       = "{}"
}
`, template)
}

func testAccAzureRMLogicAppTriggerHttpRequest_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLogicAppTriggerHttpRequest_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_http_request" "import" {
  name         = azurerm_logic_app_trigger_http_request.test.name
  logic_app_id = azurerm_logic_app_trigger_http_request.test.logic_app_id
  schema       = azurerm_logic_app_trigger_http_request.test.schema
}
`, template)
}

func testAccAzureRMLogicAppTriggerHttpRequest_fullSchema(data acceptance.TestData) string {
	template := testAccAzureRMLogicAppTriggerHttpRequest_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_http_request" "test" {
  name         = "some-http-trigger"
  logic_app_id = azurerm_logic_app_workflow.test.id

  schema = <<SCHEMA
{
    "type": "object",
    "properties": {
        "hello": {
            "type": "string"
        }
    }
}
SCHEMA

}
`, template)
}

func testAccAzureRMLogicAppTriggerHttpRequest_method(data acceptance.TestData) string {
	template := testAccAzureRMLogicAppTriggerHttpRequest_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_http_request" "test" {
  name         = "some-http-trigger"
  logic_app_id = azurerm_logic_app_workflow.test.id
  schema       = "{}"
  method       = "PUT"
}
`, template)
}

func testAccAzureRMLogicAppTriggerHttpRequest_relativePath(data acceptance.TestData) string {
	template := testAccAzureRMLogicAppTriggerHttpRequest_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_http_request" "test" {
  name          = "some-http-trigger"
  logic_app_id  = azurerm_logic_app_workflow.test.id
  schema        = "{}"
  method        = "POST"
  relative_path = "customers/{id}"
}
`, template)
}

func testAccAzureRMLogicAppTriggerHttpRequest_template(data acceptance.TestData) string {
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
