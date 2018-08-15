package azurerm

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func testAccAzureRMRoleAssignment_importBasic(t *testing.T) {
	resourceName := "azurerm_role_assignment.test"

	roleDefinitionId := uuid.New().String()
	roleAssignmentId := uuid.New().String()
	ri := acctest.RandInt()
	config := testAccAzureRMRoleAssignment_customConfig(roleDefinitionId, roleAssignmentId, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
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

func testAccAzureRMRoleAssignment_importCustom(t *testing.T) {
	resourceName := "azurerm_role_assignment.test"

	roleDefinitionId := uuid.New().String()
	roleAssignmentId := uuid.New().String()
	ri := acctest.RandInt()
	config := testAccAzureRMRoleAssignment_customConfig(roleDefinitionId, roleAssignmentId, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
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
