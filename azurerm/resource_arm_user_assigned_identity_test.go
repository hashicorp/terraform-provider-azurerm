package azurerm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMUserAssignedIdentity_create(t *testing.T) {
	principalIdRegex := "^[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}$"
	resourceName := "azurerm_user_assigned_identity.test"
	ri := acctest.RandInt()
	config := testAccAzureRMUserAssignedIdentityCreate(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.TestMatchResourceAttr(resourceName, "principal_id", regexp.MustCompile(principalIdRegex)),
			},
		},
	})
}

func testAccAzureRMUserAssignedIdentityCreate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name = "acctestRG-acctest%[1]d"
	location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
	resource_group_name = "${azurerm_resource_group.test.name}"
	location = "${azurerm_resource_group.test.location}"

	name = "acctest%[1]d"
}
`, rInt, location)
}
