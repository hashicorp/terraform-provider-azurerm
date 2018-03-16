package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"os"
)

// check recurring basic -> recurring each type?
// check: basic everything -> complete everything?
// check : base + web + error action

func TestAccAzureRMSchedulerJob_web_basic(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_scheduler_job.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("action_web"),
					"", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_basic(resourceName, "action_web"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_startTime(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_scheduler_job.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("action_web"),
					"", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_basic(resourceName, "action_web"),
				),
			},
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("action_web"),
					testAccAzureRMSchedulerJob_block_startTime("2019-07-07T07:07:07-07:00"),
					"", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_startTime(resourceName, "2019-07-07T14:07:07Z"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_web_retry(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_scheduler_job.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("action_web"),
					testAccAzureRMSchedulerJob_block_retry_empty(),
					"", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_basic(resourceName, "action_web"),
					checkAccAzureRMSchedulerJob_retry_empty(resourceName),
				),
			},
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("action_web"),
					testAccAzureRMSchedulerJob_block_retry_complete(),
					"", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_basic(resourceName, "action_web"),
					checkAccAzureRMSchedulerJob_retry_complete(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_web_recurring(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_scheduler_job.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("action_web"),
					testAccAzureRMSchedulerJob_block_recurrence_basic(),
					"", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_basic(resourceName, "action_web"),
					checkAccAzureRMSchedulerJob_recurrence_basic(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_web_recurringDaily(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_scheduler_job.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("action_web"),
					testAccAzureRMSchedulerJob_block_recurrence_daily(),
					"", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_basic(resourceName, "action_web"),
					checkAccAzureRMSchedulerJob_recurrence_daily(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_web_recurringWeekly(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_scheduler_job.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("action_web"),
					testAccAzureRMSchedulerJob_block_recurrence_weekly(),
					"", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_basic(resourceName, "action_web"),
					checkAccAzureRMSchedulerJob_recurrence_weekly(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_web_recurringMonthly(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_scheduler_job.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("action_web"),
					testAccAzureRMSchedulerJob_block_recurrence_monthly(),
					"", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_basic(resourceName, "action_web"),
					checkAccAzureRMSchedulerJob_recurrence_monthly(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_web_recurringMonthlyOccurrences(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_scheduler_job.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("action_web"),
					testAccAzureRMSchedulerJob_block_recurrence_monthlyOccurrences(),
					"", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_basic(resourceName, "action_web"),
					checkAccAzureRMSchedulerJob_recurrence_monthlyOccurrences(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_web_recurringAll(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_scheduler_job.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("action_web"),
					testAccAzureRMSchedulerJob_block_recurrence_basic(),
					"", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_basic(resourceName, "action_web"),
					checkAccAzureRMSchedulerJob_recurrence_basic(resourceName),
				),
			},
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("action_web"),
					testAccAzureRMSchedulerJob_block_recurrence_daily(),
					"", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_basic(resourceName, "action_web"),
					checkAccAzureRMSchedulerJob_recurrence_daily(resourceName),
				),
			},
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("action_web"),
					testAccAzureRMSchedulerJob_block_recurrence_weekly(),
					"", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_basic(resourceName, "action_web"),
					checkAccAzureRMSchedulerJob_recurrence_weekly(resourceName),
				),
			},
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("action_web"),
					testAccAzureRMSchedulerJob_block_recurrence_monthly(),
					"", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_basic(resourceName, "action_web"),
					checkAccAzureRMSchedulerJob_recurrence_monthly(resourceName),
				),
			},
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("action_web"),
					testAccAzureRMSchedulerJob_block_recurrence_monthlyOccurrences(),
					"", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_basic(resourceName, "action_web"),
					checkAccAzureRMSchedulerJob_recurrence_monthlyOccurrences(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_web_put(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_scheduler_job.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("action_web"),
					"", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_basic(resourceName, "action_web"),
				),
			},
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_put("action_web"),
					"", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_put(resourceName, "action_web"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_web_authBasic(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_scheduler_job.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_authBasic("action_web"),
					"", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_authBasic(resourceName, "action_web"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"action_web.0.authentication_basic.0.password"},
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_web_authCert(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_scheduler_job.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_authCert(),
					"", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_authCert(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"action_web.0.authentication_certificate.0.pfx",
					"action_web.0.authentication_certificate.0.password",
				},
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_web_authAd(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_scheduler_job.test"

	clientId := os.Getenv("ARM_CLIENT_ID")
	tenantId := os.Getenv("ARM_TENANT_ID")
	secret := os.Getenv("ARM_CLIENT_SECRET")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_authAd(tenantId, clientId, secret),
					"", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_authAd(resourceName, tenantId, clientId, secret),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"action_web.0.authentication_active_directory.0.secret",
				},
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_web_basic_onceToRecurring(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_scheduler_job.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("action_web"),
					"", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_basic(resourceName, "action_web"),
				),
			},
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("action_web"),
					"", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_basic(resourceName, "action_web"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_errrActionWeb_basic(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_scheduler_job.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("action_web"),
					"", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_basic(resourceName, "action_web"),
				),
			},
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("action_web"),
					testAccAzureRMSchedulerJob_block_actionWeb_basic("error_action_web"),
					"", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_basic(resourceName, "action_web"),
					checkAccAzureRMSchedulerJob_web_basic(resourceName, "error_action_web"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_errorActionWeb_put(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_scheduler_job.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_put("action_web"),
					testAccAzureRMSchedulerJob_block_actionWeb_put("error_action_web"),
					"", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_put(resourceName, "action_web"),
					checkAccAzureRMSchedulerJob_web_put(resourceName, "error_action_web"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_errorActionWeb_authBasic(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_scheduler_job.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_authBasic("action_web"),
					testAccAzureRMSchedulerJob_block_actionWeb_authBasic("error_action_web"),
					"", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_authBasic(resourceName, "action_web"),
					checkAccAzureRMSchedulerJob_web_authBasic(resourceName, "error_action_web"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"action_web.0.authentication_basic.0.password",
					"error_action_web.0.authentication_basic.0.password",
				},
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_errorAndActionWebAuth(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_scheduler_job.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_base(ri, testLocation(),
					testAccAzureRMSchedulerJob_block_actionWeb_authCert(),
					testAccAzureRMSchedulerJob_block_actionWeb_authBasic("error_action_web"),
					"", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMSchedulerJob_base(resourceName),
					checkAccAzureRMSchedulerJob_web_authCert(resourceName),
					checkAccAzureRMSchedulerJob_web_authBasic(resourceName, "error_action_web"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{
					"action_web.0.authentication_certificate.0.pfx",
					"action_web.0.authentication_certificate.0.password",
					"error_action_web.0.authentication_basic.0.password",
				},
			},
		},
	})
}

func testCheckAzureRMSchedulerJobDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_scheduler_job.test" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		jobCollection := rs.Primary.Attributes["job_collection_name"]

		client := testAccProvider.Meta().(*ArmClient).schedulerJobsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, jobCollection, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Scheduler Job Collection still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMSchedulerJobExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %q", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		jobCollection := rs.Primary.Attributes["job_collection_name"]

		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Scheduler Job: %q", name)
		}

		client := testAccProvider.Meta().(*ArmClient).schedulerJobsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, jobCollection, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Scheduler Job  %q (resource group: %q) was not found: %+v", name, resourceGroup, err)
			}

			return fmt.Errorf("Bad: Get on schedulerJobsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMSchedulerJob_base(rInt int, location, block1, block2, block3, block4 string) string {
	return fmt.Sprintf(` 
resource "azurerm_resource_group" "rg" { 
  name     = "acctestRG-%[1]d" 
  location = "%[2]s" 
} 
 
resource "azurerm_scheduler_job_collection" "jc" {
    name                = "acctestRG-%[1]d-job_collection"
    location            = "${azurerm_resource_group.rg.location}"
    resource_group_name = "${azurerm_resource_group.rg.name}"
    sku                 = "standard"
}

resource "azurerm_scheduler_job" "test" {
    name                = "acctestRG-%[1]d-job"
    resource_group_name = "${azurerm_resource_group.rg.name}"
    job_collection_name = "${azurerm_scheduler_job_collection.jc.name}"

    %[3]s

    %[4]s

    %[5]s

    %[6]s
} 
`, rInt, location, block1, block2, block3, block4)
}
func checkAccAzureRMSchedulerJob_base(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMSchedulerJobExists(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "name"),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttrSet(resourceName, "job_collection_name"),
		resource.TestCheckResourceAttr(resourceName, "state", "enabled"),
	)
}

func testAccAzureRMSchedulerJob_block_startTime(time string) string {
	return fmt.Sprintf(` 
  start_time = "%s"
`, time)
}

func checkAccAzureRMSchedulerJob_startTime(resourceName, time string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "start_time", time),
	)
}

func testAccAzureRMSchedulerJob_block_recurrence_basic() string {
	return ` 
  recurrence {
    frequency  = "minute"
    interval   = 5
    count      = 10
  //end_time  = "2019-07-17T07:07:07-07:00"
  } 
`
}

func checkAccAzureRMSchedulerJob_recurrence_basic(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.frequency", "minute"),
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.interval", "5"),
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.count", "10"),
	)
}

func testAccAzureRMSchedulerJob_block_recurrence_daily() string {
	return ` 
  recurrence {
    frequency = "day"
    count     = 100 
    hours     = [0,12]
    minutes   = [0,15,30,45] 
  } 
`
}

func checkAccAzureRMSchedulerJob_recurrence_daily(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.frequency", "day"),
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.count", "100"),
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.hours.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.minutes.#", "4"),
	)
}

func testAccAzureRMSchedulerJob_block_recurrence_weekly() string {
	return ` 
  recurrence {
     frequency    = "week"
     count        = 100 
     week_days = ["Sunday", "Saturday"] 
  } 
`
}

func checkAccAzureRMSchedulerJob_recurrence_weekly(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.frequency", "week"),
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.count", "100"),
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.week_days.#", "2"),
	)
}

func testAccAzureRMSchedulerJob_block_recurrence_monthly() string {
	return ` 
  recurrence {
    frequency  = "month"
    count      = 100 
    month_days = [-11,-1,1,11]
  } 
`
}

func checkAccAzureRMSchedulerJob_recurrence_monthly(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.frequency", "month"),
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.count", "100"),
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.month_days.#", "4"),
	)
}

func testAccAzureRMSchedulerJob_block_recurrence_monthlyOccurrences() string {
	return ` 
  recurrence {
    frequency  = "month"
    count      = 100 
    monthly_occurrences = [
            {
                day        = "sunday"
                occurrence = 1
            },
            {
                day        = "sunday"
                occurrence = 3
            },
            {
                day        = "sunday"
                occurrence = -1
            }
        ]
  } 
`
}

func checkAccAzureRMSchedulerJob_recurrence_monthlyOccurrences(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.frequency", "month"),
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.count", "100"),
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.monthly_occurrences.#", "3"),
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.monthly_occurrences.2181640481.day", "sunday"),
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.monthly_occurrences.2181640481.occurrence", "1"),
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.monthly_occurrences.2956940195.day", "sunday"),
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.monthly_occurrences.2956940195.occurrence", "3"),
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.monthly_occurrences.679325150.day", "sunday"),
		resource.TestCheckResourceAttr(resourceName, "recurrence.0.monthly_occurrences.679325150.occurrence", "-1"),
	)
}

func testAccAzureRMSchedulerJob_block_retry_empty() string {
	return ` 
  retry { 
  } 
`
}

func checkAccAzureRMSchedulerJob_retry_empty(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "retry.0.interval", "00:00:30"),
		resource.TestCheckResourceAttr(resourceName, "retry.0.count", "4"),
	)
}

func testAccAzureRMSchedulerJob_block_retry_complete() string {
	return ` 
  retry { 
    interval = "00:05:00" //5 min
    count    =  10
  } 
`
}

func checkAccAzureRMSchedulerJob_retry_complete(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "retry.0.interval", "00:05:00"),
		resource.TestCheckResourceAttr(resourceName, "retry.0.count", "10"),
	)
}

func testAccAzureRMSchedulerJob_block_actionWeb_basic(blockName string) string {
	return fmt.Sprintf(`
  %s {
    url = "http://this.get.url.fails"
  } 
`, blockName)
}

func checkAccAzureRMSchedulerJob_web_basic(resourceName, blockName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("%s.0.url", blockName), "http://this.get.url.fails"),
	)
}

func testAccAzureRMSchedulerJob_block_actionWeb_put(blockName string) string {
	return fmt.Sprintf(`
  %s {
    url    = "http://this.put.url.fails"
    method = "put"
    body   = "this is some text"
    headers = {
      Content-Type = "text"
	}
  } 
`, blockName)
}

func checkAccAzureRMSchedulerJob_web_put(resourceName, blockName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("%s.0.url", blockName), "http://this.put.url.fails"),
		resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("%s.0.method", blockName), "put"),
		resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("%s.0.body", blockName), "this is some text"),
		resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("%s.0.headers.%%", blockName), "1"),
	)
}

func testAccAzureRMSchedulerJob_block_actionWeb_authBasic(blockName string) string {
	return fmt.Sprintf(`
  %s {
    url    = "https://this.url.fails"
    method = "get"

    authentication_basic {
      username = "login"
      password = "apassword"
    }
  }
`, blockName)
}

func checkAccAzureRMSchedulerJob_web_authBasic(resourceName string, blockName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("%s.0.url", blockName), "https://this.url.fails"),
		resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("%s.0.method", blockName), "get"),
		resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("%s.0.authentication_basic.0.username", blockName), "login"),
		resource.TestCheckResourceAttrSet(resourceName, fmt.Sprintf("%s.0.authentication_basic.0.password", blockName)),
	)
}

func testAccAzureRMSchedulerJob_block_actionWeb_authCert() string {
	return `
  action_web {
    url    = "https://this.url.fails"
    method = "get"

    authentication_certificate {
      pfx      = "${base64encode(file("testdata/application_gateway_test.pfx"))}"
      password = "terraform"
    }
  }
`
}

func checkAccAzureRMSchedulerJob_web_authCert(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "action_web.0.url", "https://this.url.fails"),
		resource.TestCheckResourceAttr(resourceName, "action_web.0.method", "get"),
		resource.TestCheckResourceAttrSet(resourceName, "action_web.0.authentication_certificate.0.pfx"),
		resource.TestCheckResourceAttrSet(resourceName, "action_web.0.authentication_certificate.0.password"),
		resource.TestCheckResourceAttr(resourceName, "action_web.0.authentication_certificate.0.thumbprint", "42C107874FD0E4A9583292A2F1098E8FE4B2EDDA"),
		resource.TestCheckResourceAttr(resourceName, "action_web.0.authentication_certificate.0.subject_name", "CN=Terraform App Gateway, OU=Azure, O=Terraform Tests, S=Some-State, C=US"),
	)
}

func testAccAzureRMSchedulerJob_block_actionWeb_authAd(tenantId, clientId, secret string) string {
	return fmt.Sprintf(`
  action_web {
    url    = "https://this.url.fails"
    method = "get"

    authentication_active_directory {
      tenant_id = "%s"
      client_id = "%s"
      secret    = "%s"
      audience  = "https://management.core.windows.net/"
    }
  }
`, tenantId, clientId, secret)
}

func checkAccAzureRMSchedulerJob_web_authAd(resourceName, tenantId, clientId, secret string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "action_web.0.url", "https://this.url.fails"),
		resource.TestCheckResourceAttr(resourceName, "action_web.0.method", "get"),
		resource.TestCheckResourceAttr(resourceName, "action_web.0.authentication_active_directory.0.tenant_id", tenantId),
		resource.TestCheckResourceAttr(resourceName, "action_web.0.authentication_active_directory.0.client_id", clientId),
		resource.TestCheckResourceAttr(resourceName, "action_web.0.authentication_active_directory.0.secret", secret),
		resource.TestCheckResourceAttr(resourceName, "action_web.0.authentication_active_directory.0.audience", "https://management.core.windows.net/"),
	)
}
