package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMLogAnalyticsWorkspace_importRequiredOnly(t *testing.T) {
	resourceName := "azurerm_log_analytics_workspace.test"

	ri := acctest.RandInt()
	config := testAccAzureRMLogAnalyticsWorkspace_requiredOnly(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceDestroy,
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

func TestAccAzureRMLogAnalyticsWorkspace_importRetentionInDaysComplete(t *testing.T) {
	resourceName := "azurerm_log_analytics_workspace.test"

	ri := acctest.RandInt()
	config := testAccAzureRMLogAnalyticsWorkspace_retentionInDaysComplete(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceDestroy,
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
