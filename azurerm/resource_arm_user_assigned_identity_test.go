package azurerm

import (
	"fmt"
	"testing"

	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMUserAssignedIdentity_basic(t *testing.T) {
	resourceName := "azurerm_user_assigned_identity.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(14)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMUserAssignedIdentityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMUserAssignedIdentity_basic(ri, testLocation(), rs),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMUserAssignedIdentityExists(resourceName),
					resource.TestMatchResourceAttr(resourceName, "principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(resourceName, "client_id", validate.UUIDRegExp),
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
func TestAccAzureRMUserAssignedIdentity_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_user_assigned_identity.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(14)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMUserAssignedIdentityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMUserAssignedIdentity_basic(ri, testLocation(), rs),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMUserAssignedIdentityExists(resourceName),
					resource.TestMatchResourceAttr(resourceName, "principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(resourceName, "client_id", validate.UUIDRegExp),
				),
			},
			{
				Config:      testAccAzureRMUserAssignedIdentity_requiresImport(ri, testLocation(), rs),
				ExpectError: testRequiresImportError("azurerm_user_assigned_identity"),
			},
		},
	})
}

func testCheckAzureRMUserAssignedIdentityExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for virtual machine: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).MSI.UserAssignedIdentitiesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on userAssignedIdentitiesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: User Assigned Identity %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMUserAssignedIdentityDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).MSI.UserAssignedIdentitiesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_machine" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("User Assigned Identity still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMUserAssignedIdentity_basic(rInt int, location string, rString string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}
`, rInt, location, rString)
}

func testAccAzureRMUserAssignedIdentity_requiresImport(rInt int, location string, rString string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "import" {
  name                = "${azurerm_user_assigned_identity.test.name}"
  resource_group_name = "${azurerm_user_assigned_identity.test.resource_group_name}"
  location            = "${azurerm_user_assigned_identity.test.location}"
}
`, testAccAzureRMUserAssignedIdentity_basic(rInt, location, rString))
}
