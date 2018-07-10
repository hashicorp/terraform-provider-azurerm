package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"strings"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMLogicAppActionHttp_basic(t *testing.T) {
	resourceName := "azurerm_logic_app_action_http.test"
	ri := acctest.RandInt()
	config := testAccAzureRMLogicAppActionHttp_basic(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppActionHttpExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMLogicAppActionHttp_headers(t *testing.T) {
	resourceName := "azurerm_logic_app_action_http.test"
	ri := acctest.RandInt()
	config := testAccAzureRMLogicAppActionHttp_headers(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppActionHttpExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMLogicAppActionHttp_disappears(t *testing.T) {
	ri := acctest.RandInt()
	location := testLocation()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppActionHttp_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppActionHttpExists("azurerm_logic_app_action_http.test"),
				),
			},
			{
				// delete it
				Config: testAccAzureRMLogicAppActionHttp_template(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppWorkflowExists("azurerm_logic_app_workflow.test"),
				),
			},
			{
				Config:             testAccAzureRMLogicAppActionHttp_basic(ri, location),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMLogicAppActionHttpExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		logicAppId := rs.Primary.Attributes["logic_app_id"]
		id, err := parseAzureResourceID(logicAppId)
		if err != nil {
			return err
		}

		actionName := rs.Primary.Attributes["name"]
		workflowName := id.Path["workflows"]
		resourceGroup := id.ResourceGroup

		client := testAccProvider.Meta().(*ArmClient).logicWorkflowsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, workflowName)
		if err != nil {
			return fmt.Errorf("Bad: Get on logicWorkflowsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Logic App Workflow %q (resource group %q) does not exist", workflowName, resourceGroup)
		}

		definition := resp.WorkflowProperties.Definition.(map[string]interface{})
		actions := definition["actions"].(map[string]interface{})

		exists := false
		for k, _ := range actions {
			if strings.EqualFold(k, actionName) {
				exists = true
				break
			}
		}

		if !exists {
			return fmt.Errorf("Action %q was not found on Logic App %q (Resource Group %q)", actionName, workflowName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMLogicAppActionHttp_basic(rInt int, location string) string {
	template := testAccAzureRMLogicAppActionHttp_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_http" "test" {
  name         = "action%d"
  logic_app_id = "${azurerm_logic_app_workflow.test.id}"
  method       = "GET"
  uri          = "http://example.com/hello"
}
`, template, rInt)
}

func testAccAzureRMLogicAppActionHttp_headers(rInt int, location string) string {
	template := testAccAzureRMLogicAppActionHttp_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_http" "test" {
  name         = "action%d"
  logic_app_id = "${azurerm_logic_app_workflow.test.id}"
  method       = "GET"
  uri          = "http://example.com/hello"
  headers {
    "Hello"     = "World"
    "Something" = "New"
  }
}
`, template, rInt)
}

func testAccAzureRMLogicAppActionHttp_template(rInt int, location string) string {
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
