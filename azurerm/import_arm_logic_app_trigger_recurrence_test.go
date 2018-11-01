package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMLogicAppTriggerRecurrence_importMonth(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMLogicAppTriggerRecurrence_basic(ri, testLocation(), "Month", 1)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_logic_app_trigger_recurrence.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMLogicAppTriggerRecurrence_importWeek(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMLogicAppTriggerRecurrence_basic(ri, testLocation(), "Week", 3)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_logic_app_trigger_recurrence.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMLogicAppTriggerRecurrence_importDay(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMLogicAppTriggerRecurrence_basic(ri, testLocation(), "Day", 5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_logic_app_trigger_recurrence.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMLogicAppTriggerRecurrence_importHour(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMLogicAppTriggerRecurrence_basic(ri, testLocation(), "Hour", 3)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_logic_app_trigger_recurrence.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMLogicAppTriggerRecurrence_importMinute(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMLogicAppTriggerRecurrence_basic(ri, testLocation(), "Minute", 5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_logic_app_trigger_recurrence.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMLogicAppTriggerRecurrence_importSecond(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMLogicAppTriggerRecurrence_basic(ri, testLocation(), "Second", 30)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_logic_app_trigger_recurrence.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
