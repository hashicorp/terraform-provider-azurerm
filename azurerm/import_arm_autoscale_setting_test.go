package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMAutoScaleSetting_importBasic(t *testing.T) {
	resourceName := "azurerm_autoscale_setting.test"

	ri := acctest.RandInt()
	rs := acctest.RandString(6)
	location := testLocation()
	config := testAccAzureRMAutoScaleSetting_basic(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutoScaleSettingDestroy,
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

func TestAccAzureRMAutoScaleSetting_importRecurrence(t *testing.T) {
	resourceName := "azurerm_autoscale_setting.test"

	ri := acctest.RandInt()
	rs := acctest.RandString(6)
	location := testLocation()
	config := testAccAzureRMAutoScaleSetting_recurrence(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutoScaleSettingDestroy,
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

func TestAccAzureRMAutoScaleSetting_importFixedDate(t *testing.T) {
	resourceName := "azurerm_autoscale_setting.test"

	ri := acctest.RandInt()
	rs := acctest.RandString(6)
	location := testLocation()
	config := testAccAzureRMAutoScaleSetting_fixedDate(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutoScaleSettingDestroy,
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
