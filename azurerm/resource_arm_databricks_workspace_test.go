package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMDatabricksWorkspace_basic(t *testing.T) {
	resourceName := "azurerm_databricks_workspace.test"
	ri := acctest.RandInt()
	config := testAccAzureRMDatabricksWorkspace_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDatabricksWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabricksWorkspaceExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMDatabricksWorkspace_withTags(t *testing.T) {
	resourceName := "azurerm_databricks_workspace.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMDatabricksWorkspace_withTags(ri, location)
	postConfig := testAccAzureRMDatabricksWorkspace_withTagsUpdate(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDatabricksWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabricksWorkspaceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabricksWorkspaceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
				),
			},
		},
	})
}

func testCheckAzureRMDatabricksWorkspaceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Bad: Not found: %s", name)
		}

		name := rs.Primary.Attributes["workspace_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: No resource group found in state for Databricks Workspace: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).databricksWorkspacesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Getting Workspace: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Databricks Workspace %s (resource group: %s) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDatabricksWorkspaceDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).databricksWorkspacesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_databricks_workspace" {
			continue
		}

		name := rs.Primary.Attributes["workspace_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Bad: Databricks Workspace still exists:\n%#v", resp.ID)
		}
	}

	return nil
}

func testAccAzureRMDatabricksWorkspace_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG_%d"
  location = "%s"
}

resource "azurerm_databricks_workspace" "test" {
  workspace_name               = "databricks-test-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location = "%s"
}
`, rInt, location, rInt, location)
}

func testAccAzureRMDatabricksWorkspace_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG_%d"
  location = "%s"
}

resource "azurerm_databricks_workspace" "test" {
	workspace_name               = "databricks-test-%d"
	resource_group_name = "${azurerm_resource_group.test.name}"
	location = "%s"

	tags {
		environment = "Production"
		pricing = "Standard"
	  }
  }
`, rInt, location, rInt, location)
}

func testAccAzureRMDatabricksWorkspace_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG_%d"
  location = "%s"
}

resource "azurerm_databricks_workspace" "test" {
	workspace_name               = "databricks-test-%d"
	resource_group_name = "${azurerm_resource_group.test.name}"
	location = "%s"

	tags {
		pricing = "Premium"
	  }
  }
`, rInt, location, rInt, location)
}
