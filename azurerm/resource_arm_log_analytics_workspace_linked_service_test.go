package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMLogAnalyticsWorkspaceLinkedService_requiredOnly(t *testing.T) {
	resourceName := "azurerm_log_analytics_workspace_linked_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMLogAnalyticsWorkspaceLinkedServiceRequiredOnly(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceLinkedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceLinkedServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctestworkspace-%d/Automation", ri)),
					resource.TestCheckResourceAttr(resourceName, "workspace_name", fmt.Sprintf("acctestworkspace-%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "linked_service_name", "automation"),
				),
			},
		},
	})
}

func TestAccAzureRMLogAnalyticsWorkspaceLinkedService_optionalArguments(t *testing.T) {
	resourceName := "azurerm_log_analytics_workspace_linked_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMLogAnalyticsWorkspaceLinkedServiceOptionalArguments(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceLinkedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceLinkedServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "linked_service_name", "automation"),
				),
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

func testCheckAzureRMLogAnalyticsWorkspaceLinkedServiceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
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

func TestAccAzureRMLogAnalyticsWorkspaceLinkedService_importRequiredOnly(t *testing.T) {
	resourceName := "azurerm_log_analytics_workspace_linked_service.test"

	ri := acctest.RandInt()
	config := testAccAzureRMLogAnalyticsWorkspaceLinkedServiceRequiredOnly(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceLinkedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMLogAnalyticsWorkspaceLinkedService_importOptionalArguments(t *testing.T) {
	resourceName := "azurerm_log_analytics_workspace_linked_service.test"

	ri := acctest.RandInt()
	config := testAccAzureRMLogAnalyticsWorkspaceLinkedServiceOptionalArguments(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceLinkedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMLogAnalyticsWorkspaceLinkedServiceRequiredOnly(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestresourcegroup-%d"
  location = "%v"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestautomation-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "Basic"
  }

  tags {
    environment = "development"
  }
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestworkspace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_log_analytics_workspace_linked_service" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  workspace_name      = "${azurerm_log_analytics_workspace.test.name}"

  linked_service_properties {
    resource_id = "${azurerm_automation_account.test.id}"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMLogAnalyticsWorkspaceLinkedServiceOptionalArguments(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestresourcegroup-%d"
  location = "%v"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestautomation-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "Basic"
  }

  tags {
    environment = "development"
  }
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestworkspace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_log_analytics_workspace_linked_service" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  workspace_name      = "${azurerm_log_analytics_workspace.test.name}"
  linked_service_name = "automation"

  linked_service_properties {
    resource_id = "${azurerm_automation_account.test.id}"
  }
}
`, rInt, location, rInt, rInt)
}
