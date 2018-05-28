package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMPolicyAssignment_importComplete(t *testing.T) {
	rInt := acctest.RandInt()
	location := testLocation()
	resourceName := "azurerm_policy_assignment.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPolicyAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyAssignment_complete(rInt, location),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
