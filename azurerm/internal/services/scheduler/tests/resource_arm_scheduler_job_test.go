// nolint: megacheck
// entire automation SDK has been depreciated in v21.3 in favor of logic apps, an entirely different service.
package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSchedulerJob_web_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_scheduler_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_basic(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.url", "http://example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.method", "get"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSchedulerJob_web_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_scheduler_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_basic(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.url", "http://example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.method", "get"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMSchedulerJob_web_requiresImport),
		},
	})
}

func TestAccAzureRMSchedulerJob_storageQueue(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_scheduler_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_storageQueue(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "action_storage_queue.0.storage_account_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "action_storage_queue.0.storage_queue_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "action_storage_queue.0.sas_token"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_storage_queue.0.message", "storage message"),
				),
			},
			data.ImportStep("action_storage_queue.0.sas_token"),
		},
	})
}

func TestAccAzureRMSchedulerJob_storageQueue_errorAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_scheduler_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_storageQueue_errorAction(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.url", "http://example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "error_action_storage_queue.0.storage_account_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "error_action_storage_queue.0.storage_queue_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "error_action_storage_queue.0.sas_token"),
					resource.TestCheckResourceAttr(data.ResourceName, "error_action_storage_queue.0.message", "storage message"),
				),
			},
			data.ImportStep("error_action_storage_queue.0.sas_token"),
		},
	})
}

func TestAccAzureRMSchedulerJob_web_put(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_scheduler_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_put(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.url", "http://example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.method", "put"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.body", "this is some text"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSchedulerJob_web_authBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_scheduler_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_authBasic(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.authentication_basic.0.username", "login"),
				),
			},
			data.ImportStep("action_web.0.authentication_basic.0.password"),
		},
	})
}

func TestAccAzureRMSchedulerJob_web_authCert(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_scheduler_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_authCert(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.authentication_certificate.0.thumbprint", "42C107874FD0E4A9583292A2F1098E8FE4B2EDDA"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.authentication_certificate.0.subject_name", "CN=Terraform App Gateway, OU=Azure, O=Terraform Tests, S=Some-State, C=US"),
				),
			},
			data.ImportStep("action_web.0.authentication_certificate.0.pfx",
				"action_web.0.authentication_certificate.0.password"),
		},
	})
}

func TestAccAzureRMSchedulerJob_web_authAd(t *testing.T) {
	if os.Getenv(resource.TestEnvVar) == "" {
		t.Skipf("Skipping since %q isn't set", resource.TestEnvVar)
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_scheduler_job", "test")

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
				Config: testAccAzureRMSchedulerJob_web_authAd(data, tenantId, clientId, secret, audience),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.authentication_active_directory.0.tenant_id", tenantId),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.authentication_active_directory.0.client_id", clientId),
					resource.TestCheckResourceAttrSet(data.ResourceName, "action_web.0.authentication_active_directory.0.audience"),
				),
			},
			data.ImportStep("action_web.0.authentication_active_directory.0.secret"),
		},
	})
}

func TestAccAzureRMSchedulerJob_web_retry(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_scheduler_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_retry(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttr(data.ResourceName, "retry.0.interval", "00:05:00"),
					resource.TestCheckResourceAttr(data.ResourceName, "retry.0.count", "10"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSchedulerJob_web_recurring(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_scheduler_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_recurring(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.frequency", "minute"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.interval", "5"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.count", "10"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSchedulerJob_web_recurringDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_scheduler_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_recurringDaily(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.frequency", "day"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.count", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.hours.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.minutes.#", "4"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSchedulerJob_web_recurringWeekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_scheduler_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_recurringWeekly(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.frequency", "week"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.count", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.week_days.#", "2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSchedulerJob_web_recurringMonthly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_scheduler_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_recurringMonthly(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.frequency", "month"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.count", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.month_days.#", "4"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSchedulerJob_web_recurringMonthlyOccurrences(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_scheduler_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_recurringMonthlyOccurrences(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.frequency", "month"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.count", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.monthly_occurrences.#", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.monthly_occurrences.2181640481.day", "sunday"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.monthly_occurrences.2181640481.occurrence", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.monthly_occurrences.2956940195.day", "sunday"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.monthly_occurrences.2956940195.occurrence", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.monthly_occurrences.679325150.day", "sunday"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.monthly_occurrences.679325150.occurrence", "-1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSchedulerJob_web_errorAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_scheduler_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_errorAction(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.method", "get"),
					resource.TestCheckResourceAttr(data.ResourceName, "error_action_web.0.url", "https://example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "error_action_web.0.method", "get"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSchedulerJob_web_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_scheduler_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJob_web_complete(data, "2019-07-07T07:07:07-07:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMSchedulerJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.url", "http://example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.method", "put"),
					resource.TestCheckResourceAttr(data.ResourceName, "action_web.0.body", "this is some text"),
					resource.TestCheckResourceAttr(data.ResourceName, "retry.0.interval", "00:05:00"),
					resource.TestCheckResourceAttr(data.ResourceName, "retry.0.count", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.frequency", "month"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.count", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "recurrence.0.month_days.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "start_time", "2019-07-07T14:07:07Z"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMSchedulerJobDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Scheduler.JobsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_scheduler_job.test" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		jobCollection := rs.Primary.Attributes["job_collection_name"]

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
		client := acceptance.AzureProvider.Meta().(*clients.Client).Scheduler.JobsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMSchedulerJob_template(data acceptance.TestData) string {
	return fmt.Sprintf(` 
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_scheduler_job_collection" "test" {
  name                = "acctest-%d-job_collection"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMSchedulerJob_web_basic(data acceptance.TestData) string {
	template := testAccAzureRMSchedulerJob_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMSchedulerJob_web_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMSchedulerJob_web_basic(data)
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
`, template)
}

func testAccAzureRMSchedulerJob_web_put(data acceptance.TestData) string {
	template := testAccAzureRMSchedulerJob_template(data)
	return fmt.Sprintf(`
%s 

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
`, template, data.RandomInteger)
}

func testAccAzureRMSchedulerJob_web_authBasic(data acceptance.TestData) string {
	template := testAccAzureRMSchedulerJob_template(data)
	return fmt.Sprintf(`
%s

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
`, template, data.RandomInteger)
}

func testAccAzureRMSchedulerJob_web_authCert(data acceptance.TestData) string {
	template := testAccAzureRMSchedulerJob_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMSchedulerJob_web_authAd(data acceptance.TestData, tenantId, clientId, secret, audience string) string {
	template := testAccAzureRMSchedulerJob_template(data)
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
`, template, data.RandomInteger, tenantId, clientId, secret, audience)
}

func testAccAzureRMSchedulerJob_web_retry(data acceptance.TestData) string {
	template := testAccAzureRMSchedulerJob_template(data)
	return fmt.Sprintf(`
%s

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
`, template, data.RandomInteger)
}

func testAccAzureRMSchedulerJob_web_recurring(data acceptance.TestData) string {
	template := testAccAzureRMSchedulerJob_template(data)
	return fmt.Sprintf(`
%s

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
`, template, data.RandomInteger)
}

func testAccAzureRMSchedulerJob_web_recurringDaily(data acceptance.TestData) string {
	template := testAccAzureRMSchedulerJob_template(data)
	return fmt.Sprintf(`
%s 

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
`, template, data.RandomInteger)
}

func testAccAzureRMSchedulerJob_web_recurringWeekly(data acceptance.TestData) string {
	template := testAccAzureRMSchedulerJob_template(data)
	return fmt.Sprintf(`
%s 

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
`, template, data.RandomInteger)
}

func testAccAzureRMSchedulerJob_web_recurringMonthly(data acceptance.TestData) string {
	template := testAccAzureRMSchedulerJob_template(data)
	return fmt.Sprintf(`
%s 

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
`, template, data.RandomInteger)
}

func testAccAzureRMSchedulerJob_web_recurringMonthlyOccurrences(data acceptance.TestData) string {
	template := testAccAzureRMSchedulerJob_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMSchedulerJob_web_errorAction(data acceptance.TestData) string {
	template := testAccAzureRMSchedulerJob_template(data)
	return fmt.Sprintf(`
%s 

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
`, template, data.RandomInteger)
}

func testAccAzureRMSchedulerJob_web_complete(data acceptance.TestData, time string) string {
	template := testAccAzureRMSchedulerJob_template(data)
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
`, template, data.RandomInteger, time)
}

func testAccAzureRMSchedulerJob_storageQueue(data acceptance.TestData) string {
	template := testAccAzureRMSchedulerJob_template(data)
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
`, template, data.RandomString, data.RandomInteger)
}

func testAccAzureRMSchedulerJob_storageQueue_errorAction(data acceptance.TestData) string {
	template := testAccAzureRMSchedulerJob_template(data)
	//need a valid URL here otherwise on a slow connection job might fault before the test check
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_queue" "test" {
  name                 = "acctest-%d-job"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
}

resource "azurerm_scheduler_job" "test" {
  name                = "acctest-%d-job"
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
`, template, data.RandomString, data.RandomInteger, data.RandomInteger)
}
