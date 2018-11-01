package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMUserAssignedIdentity_importBasic(t *testing.T) {
	resourceName := "azurerm_user_assigned_identity.test"

	ri := acctest.RandInt()
	location := testLocation()
	rs := acctest.RandString(14)
	config := testAccAzureRMUserAssignedIdentity_basic(ri, location, rs)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMUserAssignedIdentityDestroy,
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

//
