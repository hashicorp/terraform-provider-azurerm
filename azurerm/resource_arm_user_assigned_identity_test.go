package azurerm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMUserAssignedIdentity_create(t *testing.T) {
	resourceIdRegex := "/subscriptions/.+/resourcegroups/.+/providers/Microsoft.ManagedIdentity/userAssignedIdentities/test"
	principalIdRegex := "^[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}$"
	resourceName := "azurerm_user_assigned_identity.test"
	ri := acctest.RandInt()
	config := testAccAzureRMUserAssignedIdentityCreateTemplate(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(resourceIdRegex)),
					resource.TestMatchResourceAttr(resourceName, "principal_id", regexp.MustCompile(principalIdRegex)),
				),
			},
		},
	})
}

func testAccAzureRMUserAssignedIdentityCreateTemplate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name = "acctestRG-%d"
	location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
	resource_group_name = "${azurerm_resource_group.test.name}"
	location = "${azurerm_resource_group.test.location}"

	name = "test"
}
`, rInt, location)
}
