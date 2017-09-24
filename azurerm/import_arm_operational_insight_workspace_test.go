package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMOperationalInsightWorkspace_importRequiredOnly(t *testing.T) {
	resourceName := "azurerm_log_analytics.test"

	ri := acctest.RandInt()
	config := testAccAzureRMOperationalInsightWorkspace_requiredOnly(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMOperationalInsightWorkspaceDestroy,
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

func TestAccAzureRMOperationalInsightWorkspace_importRetentionInDaysComplete(t *testing.T) {
	resourceName := "azurerm_log_analytics.test"

	ri := acctest.RandInt()
	config := testAccAzureRMOperationalInsightWorkspace_retentionInDaysComplete(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMOperationalInsightWorkspaceDestroy,
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
