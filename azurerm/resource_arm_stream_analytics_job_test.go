package azurerm

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMStreamAnalyticsJob_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsJob_basic("westus", "integration"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsExists("azurerm_stream_analytics_job.test"),
				),
			},
		},
	})
}

func testAccAzureRMStreamAnalyticsJob_basic(location, name string) string {

	return fmt.Sprintf(`
		resource "azurerm_stream_analytics_job" "test" {
		name = "acctest-%s"
		sku = "Standard"
		resource_group_name = "girishsandbox"
		location = "%s"
		tags {
			"Purpose" = "test"
		}
	}`, name, location)
}

func testCheckAzureRMStreamAnalyticsExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		jobName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for storage queue: %s", name)
		}
		armClient := testAccProvider.Meta().(*ArmClient)
		result, err := armClient.streamingJobClient.Get(resourceGroup, jobName, "")
		if err != nil {
			return err
		}
		if *result.JobState != rs.Primary.Attributes["job_state"] {
			return errors.New("remote state and local state have diveraged")
		}

		return nil
	}
}
