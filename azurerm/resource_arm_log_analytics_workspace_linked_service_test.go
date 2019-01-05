package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMLogAnalyticsWorkspaceLinkedService_basic(t *testing.T) {
	resourceName := "azurerm_log_analytics_workspace_linked_service.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceLinkedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsWorkspaceLinkedService_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceLinkedServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctestlaw-%d/Automation", ri)),
					resource.TestCheckResourceAttr(resourceName, "workspace_name", fmt.Sprintf("acctestlaw-%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "linked_service_name", "automation"),
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

func TestAccAzureRMLogAnalyticsWorkspaceLinkedService_complete(t *testing.T) {
	resourceName := "azurerm_log_analytics_workspace_linked_service.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceLinkedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsWorkspaceLinkedService_complete(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceLinkedServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "linked_service_name", "automation"),
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

func testCheckAzureRMLogAnalyticsWorkspaceLinkedServiceDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).linkedServicesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_analytics_workspace_linked_service" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		workspaceName := rs.Primary.Attributes["workspace_name"]
		lsName := rs.Primary.Attributes["linked_service_name"]

		resp, err := conn.Get(ctx, resourceGroup, workspaceName, lsName)
		if err != nil {
			return nil
		}
		if resp.ID == nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Log Analytics Linked Service still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMLogAnalyticsWorkspaceLinkedServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		workspaceName := rs.Primary.Attributes["workspace_name"]
		lsName := rs.Primary.Attributes["linked_service_name"]
		name := rs.Primary.Attributes["name"]

		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Log Analytics Linked Service: '%s'", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).linkedServicesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, workspaceName, lsName)
		if err != nil {
			return fmt.Errorf("Bad: Get on Log Analytics Linked Service Client: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Log Analytics Linked Service '%s' (resource group: '%s') does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMLogAnalyticsWorkspaceLinkedService_prereqs(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]v"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAutomation-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "Basic"
  }

  tags {
    Environment = "Test"
  }
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "PerGB2018"
  retention_in_days   = 30
}
`, rInt, location)
}

func testAccAzureRMLogAnalyticsWorkspaceLinkedService_basic(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace_linked_service" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  workspace_name      = "${azurerm_log_analytics_workspace.test.name}"

  linked_service_properties {
    resource_id = "${azurerm_automation_account.test.id}"
  }
}
`, testAccAzureRMLogAnalyticsWorkspaceLinkedService_prereqs(rInt, location))
}

func testAccAzureRMLogAnalyticsWorkspaceLinkedService_complete(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace_linked_service" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  workspace_name      = "${azurerm_log_analytics_workspace.test.name}"
  linked_service_name = "automation"

  linked_service_properties {
    resource_id = "${azurerm_automation_account.test.id}"
  }
}
`, testAccAzureRMLogAnalyticsWorkspaceLinkedService_prereqs(rInt, location))
}
