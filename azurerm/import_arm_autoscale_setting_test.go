package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMAutoscaleSetting_importBasic(t *testing.T) {
	resourceName := "azurerm_autoscale_setting.test"

	ri := acctest.RandInt()
	config := testAccAzureRMAutoscaleSetting_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutoscaleSettingDestroy,
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

func TestAccAzureRMAutoscaleSetting_importRecurrence(t *testing.T) {
	resourceName := "azurerm_autoscale_setting.test"

	ri := acctest.RandInt()
	config := testAccAzureRMAutoscaleSetting_recurrence(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutoscaleSettingDestroy,
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

func TestAccAzureRMAutoscaleSetting_importFixedDate(t *testing.T) {
	resourceName := "azurerm_autoscale_setting.test"

	ri := acctest.RandInt()
	config := testAccAzureRMAutoscaleSetting_fixedDate(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutoscaleSettingDestroy,
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
