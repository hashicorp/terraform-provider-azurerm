package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMLogicAppActionHttp_basic(t *testing.T) {
	resourceName := "azurerm_logic_app_action_http.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMLogicAppActionHttp_basic(ri, acceptance.Location())
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppActionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMLogicAppActionHttp_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_logic_app_action_http.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppActionHttp_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppActionExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMLogicAppActionHttp_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_logic_app_action_http"),
			},
		},
	})
}

func TestAccAzureRMLogicAppActionHttp_headers(t *testing.T) {
	resourceName := "azurerm_logic_app_action_http.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMLogicAppActionHttp_headers(ri, acceptance.Location())
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppActionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMLogicAppActionHttp_disappears(t *testing.T) {
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppActionHttp_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppActionExists("azurerm_logic_app_action_http.test"),
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

func testAccAzureRMLogicAppActionHttp_requiresImport(rInt int, location string) string {
	template := testAccAzureRMLogicAppActionHttp_basic(rInt, location)
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

func testAccAzureRMLogicAppActionHttp_headers(rInt int, location string) string {
	template := testAccAzureRMLogicAppActionHttp_template(rInt, location)
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
`, template, rInt)
}

func testAccAzureRMLogicAppActionHttp_template(rInt int, location string) string {
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
`, rInt, location, rInt)
}
