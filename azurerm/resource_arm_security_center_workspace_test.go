package azurerm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSecurityCenterWorkspace_basic(t *testing.T) {
	resourceName := "azurerm_security_center_workspace.test"
	ri := acctest.RandInt()

	scope := fmt.Sprintf("/subscriptions/%s", os.Getenv("ARM_SUBSCRIPTION_ID"))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSecurityCenterWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSecurityCenterWorkspace_basic(ri, testLocation(), scope),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSecurityCenterWorkspaceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "scope", scope),
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

func TestAccAzureRMSecurityCenterWorkspace_update(t *testing.T) {
	resourceName := "azurerm_security_center_workspace.test"
	ri := acctest.RandInt()

	scope := fmt.Sprintf("/subscriptions/%s", os.Getenv("ARM_SUBSCRIPTION_ID"))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSecurityCenterWorkspace_basic(ri, testLocation(), scope),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSecurityCenterWorkspaceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "scope", scope),
				),
			},
			{
				Config: testAccAzureRMSecurityCenterWorkspace_differentWorkspace(ri, testLocation(), scope),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSecurityCenterWorkspaceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "scope", scope),
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

func testCheckAzureRMSecurityCenterWorkspaceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ArmClient).securityCenterWorkspaceClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		contactName := rs.Primary.Attributes["workspaceSettings"]

		resp, err := client.Get(ctx, contactName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Security Center Subscription Workspace %q was not found: %+v", contactName, err)
			}

			return fmt.Errorf("Bad: GetWorkspace: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMSecurityCenterWorkspaceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).securityCenterWorkspaceClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, res := range s.RootModule().Resources {
		if res.Type != "azurerm_security_center_workspace" {
			continue
		}

		resp, err := client.Get(ctx, resourceArmSecurityCenterWorkspaceName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("security center worspace settings still exists")
	}

	return nil
}

func testAccAzureRMSecurityCenterWorkspace_basic(rInt int, location, scope string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_log_analytics_workspace" "test1" {
  name                = "acctest-%[1]d-1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "PerGB2018"
}

resource "azurerm_security_center_workspace" "test" {
    scope        = "%[3]s"
    workspace_id = "${azurerm_log_analytics_workspace.test1.id}"
}
`, rInt, location, scope)
}

func testAccAzureRMSecurityCenterWorkspace_differentWorkspace(rInt int, location, scope string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_log_analytics_workspace" "test2" {
  name                = "acctest-%[1]d-2"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "PerGB2018"
}

resource "azurerm_security_center_workspace" "test" {
    scope        = "%[3]s"
    workspace_id = "${azurerm_log_analytics_workspace.test2.id}"
}
`, rInt, location, scope)
}
