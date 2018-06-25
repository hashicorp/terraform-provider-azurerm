package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMUserAssignedIdentity_create(t *testing.T) {
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
				Check:  resource.TestCheckResourceAttrSet(resourceName, "id"),
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
