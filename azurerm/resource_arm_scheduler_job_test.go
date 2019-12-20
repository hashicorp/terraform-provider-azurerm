// nolint: megacheck
// entire automation SDK has been depreciated in v21.3 in favor of logic apps, an entirely different service.
package azurerm

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSchedulerJob_web_basic(t *testing.T) {
	resourceName := "azurerm_scheduler_job.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_basic(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.url", "http://example.com"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.method", "get"),
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

func TestAccAzureRMSchedulerJob_web_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_scheduler_job.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_basic(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.url", "http://example.com"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.method", "get"),
				),
			},
			{
				Config:      testAccAzureRMSchedulerJob_web_requiresImport(ri, acceptance.Location()),
				ExpectError: acceptance.RequiresImportError("azurerm_scheduler_job"),
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_storageQueue(t *testing.T) {
	resourceName := "azurerm_scheduler_job.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_storageQueue(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "action_storage_queue.0.storage_account_name"),
					resource.TestCheckResourceAttrSet(resourceName, "action_storage_queue.0.storage_queue_name"),
					resource.TestCheckResourceAttrSet(resourceName, "action_storage_queue.0.sas_token"),
					resource.TestCheckResourceAttr(resourceName, "action_storage_queue.0.message", "storage message"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"action_storage_queue.0.sas_token"},
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_storageQueue_errorAction(t *testing.T) {
	resourceName := "azurerm_scheduler_job.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_storageQueue_errorAction(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.url", "http://example.com"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttrSet(resourceName, "error_action_storage_queue.0.storage_account_name"),
					resource.TestCheckResourceAttrSet(resourceName, "error_action_storage_queue.0.storage_queue_name"),
					resource.TestCheckResourceAttrSet(resourceName, "error_action_storage_queue.0.sas_token"),
					resource.TestCheckResourceAttr(resourceName, "error_action_storage_queue.0.message", "storage message"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"error_action_storage_queue.0.sas_token"},
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_web_put(t *testing.T) {
	resourceName := "azurerm_scheduler_job.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_put(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.url", "http://example.com"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.method", "put"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.body", "this is some text"),
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
	resourceName := "azurerm_scheduler_job.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_authBasic(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.authentication_basic.0.username", "login"),
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
	resourceName := "azurerm_scheduler_job.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_authCert(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.authentication_certificate.0.thumbprint", "42C107874FD0E4A9583292A2F1098E8FE4B2EDDA"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.authentication_certificate.0.subject_name", "CN=Terraform App Gateway, OU=Azure, O=Terraform Tests, S=Some-State, C=US"),
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
	if os.Getenv(resource.TestEnvVar) == "" {
		t.Skipf("Skipping since %q isn't set", resource.TestEnvVar)
		return
	}

	resourceName := "azurerm_scheduler_job.test"
	ri := tf.AccRandTimeInt()

	clientId := os.Getenv("ARM_CLIENT_ID")
	tenantId := os.Getenv("ARM_TENANT_ID")
	secret := os.Getenv("ARM_CLIENT_SECRET")

	env, err := acceptance.Environment()
	if err != nil {
		t.Fatalf("Error loading Environment: %+v", err)
		return
	}

	audience := env.ServiceManagementEndpoint
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_authAd(ri, acceptance.Location(), tenantId, clientId, secret, audience),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.authentication_active_directory.0.tenant_id", tenantId),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.authentication_active_directory.0.client_id", clientId),
					resource.TestCheckResourceAttrSet(resourceName, "action_web.0.authentication_active_directory.0.audience"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"action_web.0.authentication_active_directory.0.secret"},
			},
		},
	})
}

func TestAccAzureRMSchedulerJob_web_retry(t *testing.T) {
	resourceName := "azurerm_scheduler_job.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_retry(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttr(resourceName, "retry.0.interval", "00:05:00"),
					resource.TestCheckResourceAttr(resourceName, "retry.0.count", "10"),
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
	resourceName := "azurerm_scheduler_job.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_recurring(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.frequency", "minute"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.count", "10"),
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
	resourceName := "azurerm_scheduler_job.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_recurringDaily(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.frequency", "day"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.count", "100"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.hours.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.minutes.#", "4"),
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
	resourceName := "azurerm_scheduler_job.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_recurringWeekly(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.frequency", "week"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.count", "100"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.week_days.#", "2"),
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
	resourceName := "azurerm_scheduler_job.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_recurringMonthly(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.frequency", "month"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.count", "100"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.month_days.#", "4"),
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
	resourceName := "azurerm_scheduler_job.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_recurringMonthlyOccurrences(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.frequency", "month"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.count", "100"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.monthly_occurrences.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.monthly_occurrences.2181640481.day", "sunday"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.monthly_occurrences.2181640481.occurrence", "1"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.monthly_occurrences.2956940195.day", "sunday"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.monthly_occurrences.2956940195.occurrence", "3"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.monthly_occurrences.679325150.day", "sunday"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.monthly_occurrences.679325150.occurrence", "-1"),
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

func TestAccAzureRMSchedulerJob_web_errorAction(t *testing.T) {
	resourceName := "azurerm_scheduler_job.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_errorAction(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttr(resourceName, "error_action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "error_action_web.0.method", "get"),
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

func TestAccAzureRMSchedulerJob_web_complete(t *testing.T) {
	resourceName := "azurerm_scheduler_job.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_complete(ri, acceptance.Location(), "2019-07-07T07:07:07-07:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.url", "http://example.com"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.method", "put"),
					resource.TestCheckResourceAttr(resourceName, "action_web.0.body", "this is some text"),
					resource.TestCheckResourceAttr(resourceName, "retry.0.interval", "00:05:00"),
					resource.TestCheckResourceAttr(resourceName, "retry.0.count", "10"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.frequency", "month"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.count", "100"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.month_days.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "start_time", "2019-07-07T14:07:07Z"),
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

func testCheckAzureRMSchedulerJobDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_scheduler_job.test" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		jobCollection := rs.Primary.Attributes["job_collection_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).Scheduler.JobsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testCheckAzureRMSchedulerJobExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		jobCollection := rs.Primary.Attributes["job_collection_name"]

		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Scheduler Job: %q", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Scheduler.JobsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMSchedulerJob_template(rInt int, location string) string {
	return fmt.Sprintf(` 
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_scheduler_job_collection" "test" {
  name                = "acctest-%[1]d-job_collection"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}
`, rInt, location)
}

func testAccAzureRMSchedulerJob_web_basic(rInt int, location string) string {
	//need a valid URL here otherwise on a slow connection job might fault before the test check
	return fmt.Sprintf(`
%s

resource "azurerm_scheduler_job" "test" {
  name                = "acctest-%d-job"
  resource_group_name = "${azurerm_resource_group.test.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.test.name}"

  action_web {
    url    = "http://example.com"
    method = "get"
  }
}
`, testAccAzureRMSchedulerJob_template(rInt, location), rInt)
}

func testAccAzureRMSchedulerJob_web_requiresImport(rInt int, location string) string {
	//need a valid URL here otherwise on a slow connection job might fault before the test check
	return fmt.Sprintf(`
%s

resource "azurerm_scheduler_job" "import" {
  name                = "${azurerm_scheduler_job.test.name}"
  resource_group_name = "${azurerm_scheduler_job.test.resource_group_name}"
  job_collection_name = "${azurerm_scheduler_job.test.job_collection_name}"

  action_web {
    url    = "http://example.com"
    method = "get"
  }
}
`, testAccAzureRMSchedulerJob_web_basic(rInt, location))
}

func testAccAzureRMSchedulerJob_web_put(rInt int, location string) string {
	return fmt.Sprintf(`%s 
resource "azurerm_scheduler_job" "test" {
  name                = "acctest-%d-job"
  resource_group_name = "${azurerm_resource_group.test.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.test.name}"

  action_web {
    url    = "http://example.com"
    method = "put"
    body   = "this is some text"

    headers = {
      "Content-Type" = "text"
    }
  }
}
`, testAccAzureRMSchedulerJob_template(rInt, location), rInt)
}

func testAccAzureRMSchedulerJob_web_authBasic(rInt int, location string) string {
	return fmt.Sprintf(`%s 
resource "azurerm_scheduler_job" "test" {
  name                = "acctest-%d-job"
  resource_group_name = "${azurerm_resource_group.test.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.test.name}"

  action_web {
    url    = "https://example.com"
    method = "get"

    authentication_basic {
      username = "login"
      password = "apassword"
    }
  }
}
`, testAccAzureRMSchedulerJob_template(rInt, location), rInt)
}

func testAccAzureRMSchedulerJob_web_authCert(rInt int, location string) string {
	return fmt.Sprintf(`%s 
resource "azurerm_scheduler_job" "test" {
  name                = "acctest-%d-job"
  resource_group_name = "${azurerm_resource_group.test.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.test.name}"

  action_web {
    url    = "https://example.com"
    method = "get"

    authentication_certificate {
      pfx      = "${filebase64("testdata/application_gateway_test.pfx")}"
      password = "terraform"
    }
  }
}
`, testAccAzureRMSchedulerJob_template(rInt, location), rInt)
}

func testAccAzureRMSchedulerJob_web_authAd(rInt int, location, tenantId, clientId, secret, audience string) string {
	return fmt.Sprintf(`%s 
resource "azurerm_scheduler_job" "test" {
  name                = "acctest-%d-job"
  resource_group_name = "${azurerm_resource_group.test.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.test.name}"

  action_web {
    url    = "https://example.com"
    method = "get"

    authentication_active_directory {
      tenant_id = "%s"
      client_id = "%s"
      secret    = "%s"
      audience  = "%s"
    }
  }
}
`, testAccAzureRMSchedulerJob_template(rInt, location), rInt, tenantId, clientId, secret, audience)
}

func testAccAzureRMSchedulerJob_web_retry(rInt int, location string) string {
	return fmt.Sprintf(`%s 
resource "azurerm_scheduler_job" "test" {
  name                = "acctest-%d-job"
  resource_group_name = "${azurerm_resource_group.test.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.test.name}"

  action_web {
    url    = "https://example.com"
    method = "get"
  }

  retry {
    interval = "00:05:00" //5 min
    count    = 10
  }
}
`, testAccAzureRMSchedulerJob_template(rInt, location), rInt)
}

func testAccAzureRMSchedulerJob_web_recurring(rInt int, location string) string {
	return fmt.Sprintf(`%s 
resource "azurerm_scheduler_job" "test" {
  name                = "acctest-%d-job"
  resource_group_name = "${azurerm_resource_group.test.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.test.name}"

  action_web {
    url    = "https://example.com"
    method = "get"
  }

  recurrence {
    frequency = "minute"
    interval  = 5
    count     = 10
  }
}
`, testAccAzureRMSchedulerJob_template(rInt, location), rInt)
}

func testAccAzureRMSchedulerJob_web_recurringDaily(rInt int, location string) string {
	return fmt.Sprintf(`%s 
resource "azurerm_scheduler_job" "test" {
  name                = "acctest-%d-job"
  resource_group_name = "${azurerm_resource_group.test.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.test.name}"

  action_web {
    url    = "https://example.com"
    method = "get"
  }

  recurrence {
    frequency = "day"
    count     = 100
    hours     = [0, 12]
    minutes   = [0, 15, 30, 45]
  }
}
`, testAccAzureRMSchedulerJob_template(rInt, location), rInt)
}

func testAccAzureRMSchedulerJob_web_recurringWeekly(rInt int, location string) string {
	return fmt.Sprintf(`%s 
resource "azurerm_scheduler_job" "test" {
  name                = "acctest-%d-job"
  resource_group_name = "${azurerm_resource_group.test.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.test.name}"

  action_web {
    url    = "https://example.com"
    method = "get"
  }

  recurrence {
    frequency = "week"
    count     = 100
    week_days = ["Sunday", "Saturday"]
  }
}
`, testAccAzureRMSchedulerJob_template(rInt, location), rInt)
}

func testAccAzureRMSchedulerJob_web_recurringMonthly(rInt int, location string) string {
	return fmt.Sprintf(`%s 
resource "azurerm_scheduler_job" "test" {
  name                = "acctest-%d-job"
  resource_group_name = "${azurerm_resource_group.test.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.test.name}"

  action_web {
    url    = "https://example.com"
    method = "get"
  }

  recurrence {
    frequency  = "month"
    count      = 100
    month_days = [-11, -1, 1, 11]
  }
}
`, testAccAzureRMSchedulerJob_template(rInt, location), rInt)
}

func testAccAzureRMSchedulerJob_web_recurringMonthlyOccurrences(rInt int, location string) string {
	return fmt.Sprintf(`%s 
resource "azurerm_scheduler_job" "test" {
  name                = "acctest-%d-job"
  resource_group_name = "${azurerm_resource_group.test.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.test.name}"

  action_web {
    url    = "https://example.com"
    method = "get"
  }

  recurrence {
    frequency = "month"
    count     = 100

    monthly_occurrences {
      day        = "sunday"
      occurrence = 1
	}

    monthly_occurrences {
      day        = "sunday"
      occurrence = 3
    }

    monthly_occurrences {
      day        = "sunday"
      occurrence = -1
    }
  }
}
`, testAccAzureRMSchedulerJob_template(rInt, location), rInt)
}

func testAccAzureRMSchedulerJob_web_errorAction(rInt int, location string) string {
	return fmt.Sprintf(`%s 
resource "azurerm_scheduler_job" "test" {
  name                = "acctest-%d-job"
  resource_group_name = "${azurerm_resource_group.test.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.test.name}"

  action_web {
    url    = "https://example.com"
    method = "get"
  }

  error_action_web {
    url    = "https://example.com"
    method = "get"
  }
}
`, testAccAzureRMSchedulerJob_template(rInt, location), rInt)
}

func testAccAzureRMSchedulerJob_web_complete(rInt int, location, time string) string {
	return fmt.Sprintf(`%s 
resource "azurerm_scheduler_job" "test" {
  name                = "acctest-%d-job"
  resource_group_name = "${azurerm_resource_group.test.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.test.name}"

  action_web {
    url    = "http://example.com"
    method = "put"
    body   = "this is some text"

    headers = {
      "Content-Type" = "text"
    }
  }

  retry {
    interval = "00:05:00" //5 min
    count    = 10
  }

  recurrence {
    frequency  = "month"
    count      = 100
    month_days = [-11, -1, 1, 11]
  }

  start_time = "%s"
}
`, testAccAzureRMSchedulerJob_template(rInt, location), rInt, time)
}

func testAccAzureRMSchedulerJob_storageQueue(rInt int, location string) string {
	//need a valid URL here otherwise on a slow connection job might fault before the test check
	return fmt.Sprintf(`%[1]s
resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_queue" "test" {
  name                 = "acctest-%[3]d-job"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
}

resource "azurerm_scheduler_job" "test" {
  name                = "acctest-%[3]d-job"
  resource_group_name = "${azurerm_resource_group.test.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.test.name}"

  action_storage_queue {
    storage_account_name = "${azurerm_storage_account.test.name}"
    storage_queue_name   = "${azurerm_storage_queue.test.name}"
    sas_token            = "${azurerm_storage_account.test.primary_access_key}"
    message              = "storage message"
  }
}
`, testAccAzureRMSchedulerJob_template(rInt, location), strconv.Itoa(rInt)[12:17], rInt)
}

func testAccAzureRMSchedulerJob_storageQueue_errorAction(rInt int, location string) string {
	//need a valid URL here otherwise on a slow connection job might fault before the test check
	return fmt.Sprintf(`%[1]s
resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_queue" "test" {
  name                 = "acctest-%[3]d-job"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
}

resource "azurerm_scheduler_job" "test" {
  name                = "acctest-%[3]d-job"
  resource_group_name = "${azurerm_resource_group.test.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.test.name}"

  action_web {
    url    = "http://example.com"
    method = "get"
  }

  error_action_storage_queue {
    storage_account_name = "${azurerm_storage_account.test.name}"
    storage_queue_name   = "${azurerm_storage_queue.test.name}"
    sas_token            = "${azurerm_storage_account.test.primary_access_key}"
    message              = "storage message"
  }
}
`, testAccAzureRMSchedulerJob_template(rInt, location), strconv.Itoa(rInt)[12:17], rInt)
}
