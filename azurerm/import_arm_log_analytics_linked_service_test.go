package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMLogAnalyticsLinkedService_importRequiredOnly(t *testing.T) {
	resourceName := "azurerm_log_analytics_linked_service.test"

	ri := acctest.RandInt()
	config := testAccAzureRMLogAnalyticsLinkedServiceRequiredOnly(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsLinkedServiceDestroy,
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

func TestAccAzureRMLogAnalyticsLinkedService_importOptionalArguments(t *testing.T) {
	resourceName := "azurerm_log_analytics_linked_service.test"

	ri := acctest.RandInt()
	config := testAccAzureRMLogAnalyticsLinkedServiceOptionalArguments(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsLinkedServiceDestroy,
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
