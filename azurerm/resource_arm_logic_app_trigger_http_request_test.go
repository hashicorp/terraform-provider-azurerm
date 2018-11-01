package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMLogicAppTriggerHttpRequest_basic(t *testing.T) {
	resourceName := "azurerm_logic_app_trigger_http_request.test"
	ri := acctest.RandInt()
	location := testLocation()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerHttpRequest_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "schema", "{}"),
				),
			},
		},
	})
}

func TestAccAzureRMLogicAppTriggerHttpRequest_fullSchema(t *testing.T) {
	resourceName := "azurerm_logic_app_trigger_http_request.test"
	ri := acctest.RandInt()
	location := testLocation()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerHttpRequest_fullSchema(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "schema"),
				),
			},
		},
	})
}

func TestAccAzureRMLogicAppTriggerHttpRequest_method(t *testing.T) {
	resourceName := "azurerm_logic_app_trigger_http_request.test"
	ri := acctest.RandInt()
	location := testLocation()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerHttpRequest_method(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "method", "PUT"),
				),
			},
		},
	})
}

func TestAccAzureRMLogicAppTriggerHttpRequest_relativePath(t *testing.T) {
	resourceName := "azurerm_logic_app_trigger_http_request.test"
	ri := acctest.RandInt()
	location := testLocation()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerHttpRequest_relativePath(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "method", "POST"),
					resource.TestCheckResourceAttr(resourceName, "relative_path", "customers/{id}"),
				),
			},
		},
	})
}

func TestAccAzureRMLogicAppTriggerHttpRequest_disappears(t *testing.T) {
	ri := acctest.RandInt()
	location := testLocation()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerHttpRequest_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists("azurerm_logic_app_trigger_http_request.test"),
				),
			},
			{
				// delete it
				Config: testAccAzureRMLogicAppTriggerHttpRequest_template(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppWorkflowExists("azurerm_logic_app_workflow.test"),
				),
			},
			{
				Config:             testAccAzureRMLogicAppTriggerHttpRequest_basic(ri, location),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccAzureRMLogicAppTriggerHttpRequest_basic(rInt int, location string) string {
	template := testAccAzureRMLogicAppTriggerHttpRequest_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_http_request" "test" {
  name         = "some-http-trigger"
  logic_app_id = "${azurerm_logic_app_workflow.test.id}"
  schema       = "{}"
}
`, template)
}

func testAccAzureRMLogicAppTriggerHttpRequest_fullSchema(rInt int, location string) string {
	template := testAccAzureRMLogicAppTriggerHttpRequest_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_http_request" "test" {
  name         = "some-http-trigger"
  logic_app_id = "${azurerm_logic_app_workflow.test.id}"
  schema       = <<SCHEMA
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

func testAccAzureRMLogicAppTriggerHttpRequest_method(rInt int, location string) string {
	template := testAccAzureRMLogicAppTriggerHttpRequest_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_http_request" "test" {
  name         = "some-http-trigger"
  logic_app_id = "${azurerm_logic_app_workflow.test.id}"
  schema       = "{}"
  method       = "PUT"
}
`, template)
}

func testAccAzureRMLogicAppTriggerHttpRequest_relativePath(rInt int, location string) string {
	template := testAccAzureRMLogicAppTriggerHttpRequest_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_http_request" "test" {
  name          = "some-http-trigger"
  logic_app_id  = "${azurerm_logic_app_workflow.test.id}"
  schema        = "{}"
  method        = "POST"
  relative_path = "customers/{id}"
}
`, template)
}

func testAccAzureRMLogicAppTriggerHttpRequest_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name = "acctestlaw-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}
