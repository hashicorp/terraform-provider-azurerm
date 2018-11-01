package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMAppServicePlan_importBasicWindows(t *testing.T) {
	resourceName := "azurerm_app_service_plan.test"

	ri := acctest.RandInt()
	config := testAccAzureRMAppServicePlan_basicWindows(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
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

func TestAccAzureRMAppServicePlan_importWindowsStandard(t *testing.T) {
	resourceName := "azurerm_app_service_plan.test"

	ri := acctest.RandInt()
	config := testAccAzureRMAppServicePlan_standardWindows(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
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

func TestAccAzureRMAppServicePlan_importPremiumWindows(t *testing.T) {
	resourceName := "azurerm_app_service_plan.test"

	ri := acctest.RandInt()
	config := testAccAzureRMAppServicePlan_premiumWindows(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
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

func TestAccAzureRMAppServicePlan_importCompleteWindows(t *testing.T) {
	resourceName := "azurerm_app_service_plan.test"

	ri := acctest.RandInt()
	config := testAccAzureRMAppServicePlan_completeWindows(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
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
