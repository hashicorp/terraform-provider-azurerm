package azurerm

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMRoleDefinition_importBasic(t *testing.T) {
	resourceName := "azurerm_role_definition.test"

	id := uuid.New().String()
	ri := acctest.RandInt()
	config := testAccAzureRMRoleDefinition_basic(id, ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRoleDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"role_definition_id", "scope"},
			},
		},
	})
}

func TestAccAzureRMRoleDefinition_importComplete(t *testing.T) {
	resourceName := "azurerm_role_definition.test"

	id := uuid.New().String()
	ri := acctest.RandInt()
	config := testAccAzureRMRoleDefinition_complete(id, ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRoleDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"role_definition_id", "scope"},
			},
		},
	})
}
