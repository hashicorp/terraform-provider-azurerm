// nolint: megacheck
// entire automation SDK has been depreciated in v21.3 in favor of logic apps, an entirely different service.
package tests

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/scheduler/mgmt/2016-03-01/scheduler"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSchedulerJobCollection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_scheduler_job_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJobCollection_basic(data, ""),
				Check:  checkAccAzureRMSchedulerJobCollection_basic(data.ResourceName),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSchedulerJobCollection_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_scheduler_job_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJobCollection_basic(data, ""),
				Check:  checkAccAzureRMSchedulerJobCollection_basic(data.ResourceName),
			},
			data.RequiresImportErrorStep(func(data acceptance.TestData) string {
				return testAccAzureRMSchedulerJobCollection_requiresImport(data, "")
			}),
		},
	})
}

func TestAccAzureRMSchedulerJobCollection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_scheduler_job_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSchedulerJobCollection_basic(data, ""),
				Check:  checkAccAzureRMSchedulerJobCollection_basic(data.ResourceName),
			},
			{
				Config: testAccAzureRMSchedulerJobCollection_complete(data),
				Check:  checkAccAzureRMSchedulerJobCollection_complete(data.ResourceName),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMSchedulerJobCollectionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Scheduler.JobCollectionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_scheduler_job_collection" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

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

func testCheckAzureRMSchedulerJobCollectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Scheduler.JobCollectionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Scheduler Job Collection: %q", name)
		}

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Scheduler Job Collection %q (resource group: %q) was not found: %+v", name, resourceGroup, err)
			}

			return fmt.Errorf("Bad: Get on schedulerJobCollectionsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMSchedulerJobCollection_basic(data acceptance.TestData, additional string) string {
	return fmt.Sprintf(` 
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_scheduler_job_collection" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  %s
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, additional)
}

func testAccAzureRMSchedulerJobCollection_requiresImport(data acceptance.TestData, additional string) string {
	template := testAccAzureRMSchedulerJobCollection_basic(data, additional)
	return fmt.Sprintf(`
%s

resource "azurerm_scheduler_job_collection" "import" {
  name                = "${azurerm_scheduler_job_collection.test.name}"
  location            = "${azurerm_scheduler_job_collection.test.location}"
  resource_group_name = "${azurerm_scheduler_job_collection.test.resource_group_name}"
  sku                 = "${azurerm_scheduler_job_collection.test.sku}"

  %s
}
`, template, additional)
}

func testAccAzureRMSchedulerJobCollection_complete(data acceptance.TestData) string {
	return testAccAzureRMSchedulerJobCollection_basic(data, ` 
  state = "disabled" 
  quota { 
    max_recurrence_frequency = "Hour" 
    max_recurrence_interval  = 10  
    max_job_count            = 10 
  } 
`)
}

func checkAccAzureRMSchedulerJobCollection_basic(resourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testCheckAzureRMSchedulerJobCollectionExists(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "name"),
		resource.TestCheckResourceAttrSet(resourceName, "location"),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
		resource.TestCheckResourceAttr(resourceName, "state", string(scheduler.Enabled)),
	)
}

func checkAccAzureRMSchedulerJobCollection_complete(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMSchedulerJobCollectionExists(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "name"),
		resource.TestCheckResourceAttrSet(resourceName, "location"),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
		resource.TestCheckResourceAttr(resourceName, "state", string(scheduler.Disabled)),
		resource.TestCheckResourceAttr(resourceName, "quota.0.max_job_count", "10"),
		resource.TestCheckResourceAttr(resourceName, "quota.0.max_recurrence_interval", "10"),
		resource.TestCheckResourceAttr(resourceName, "quota.0.max_recurrence_frequency", "hour"),
	)
}
