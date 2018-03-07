package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMLogAnalyticsSolution(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMLogAnalyticsSolution(ri, testLocation())

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsSolutionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsSolutionExists("azurerm_log_analytics_solution.solution"),
				),
			},
		},
	})
}

func testCheckAzureRMLogAnalyticsSolutionDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).solutionsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_analytics_solution" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Log Analytics solution still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMLogAnalyticsSolutionExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Log Analytics Workspace: '%s'", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).solutionsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on Log Analytics Solutions Client: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Log Analytics Solutions '%s' (resource group: '%s') does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMLogAnalyticsSolution(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "oms-acctestRG-%d"
	location = "%s"
}
  
resource "azurerm_log_analytics_workspace" "workspace" {
	name                = "acctest-dep-%d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
	sku                 = "Free"
}
  
resource "azurerm_log_analytics_solution" "solution" {
	name                  = "acctest-%d"
	location              = "${azurerm_resource_group.test.location}"
	resource_group_name   = "${azurerm_resource_group.test.name}"
	workspace_resource_id = "${azurerm_log_analytics_workspace.workspace.id}"
  
	plan {
	  name           = "Containers"
	  publisher      = "Microsoft"
	  product        = "OMSGallery/Containers"
	}
}
`, rInt, location, rInt, rInt)
}

// func testAccAzureRMLogAnalyticsSolution_Temp(rInt int, location string) string {
// 	return fmt.Sprintf(`
// resource "azurerm_log_analytics_solution" "solution" {
// 	name                  = "acctest"
// 	location              = "westeurope"
// 	resource_group_name   = "lg-terraformtest"
// 	workspace_resource_id = "/subscriptions/5774ad8f-d51e-4456-a72e-0447910568d3/resourcegroups/lg-terraformtest/providers/microsoft.operationalinsights/workspaces/lg-testoms"

// 	plan {
// 	  name           = "Containers"
// 	  publisher      = "Microsoft"
// 	  product        = "OMSGallery/Containers"
// 	}
// }
// `)
// }
