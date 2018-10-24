package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMLogAnalyticsWorkspaceLinkedService_importRequiredOnly(t *testing.T) {
	resourceName := "azurerm_log_analytics_workspace_linked_service.test"

	ri := acctest.RandInt()
	config := testAccAzureRMLogAnalyticsWorkspaceLinkedServiceRequiredOnly(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceLinkedServiceDestroy,
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

func TestAccAzureRMLogAnalyticsWorkspaceLinkedService_importOptionalArguments(t *testing.T) {
	resourceName := "azurerm_log_analytics_workspace_linked_service.test"

	ri := acctest.RandInt()
	config := testAccAzureRMLogAnalyticsWorkspaceLinkedServiceOptionalArguments(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceLinkedServiceDestroy,
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
