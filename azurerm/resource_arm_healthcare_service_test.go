package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/healthcareapis/mgmt/2018-08-20-preview/healthcareapis"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMHealthcareService(t *testing.T) {
	var healthcareServiceDescription healthcareapis.ServicesDescription
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMHealthcareServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHealthcareService_basic(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHealthcareServiceExists("azurerm_healthcare_service.test", &healthcareServiceDescription),
				),
			},
			{
				ResourceName: "azurerm_healthcare_service.test",
				ImportState:  true,
				ImportStateVerifyIgnore: []string{
					// since these are read from the existing state
					"access_policy_object_ids",
					"cosmosdb_throughput",
					"kind",
				},
			},
		},
	})
}

func testCheckAzureRMHealthcareServiceExists(resourceName string, healthcareServiceDescription *healthcareapis.ServicesDescription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		healthcareServiceName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for healthcare service: %s", healthcareServiceName)
		}

		client := testAccProvider.Meta().(*ArmClient).healthcare.HealthcareServiceClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, healthcareServiceName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Bad: Healthcare service %q (resource group: %q) does not exist", healthcareServiceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on healthcareServiceClient: %+v", err)
		}

		*healthcareServiceDescription = resp

		return nil
	}
}

func testCheckAzureRMHealthcareServiceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).healthcare.HealthcareServiceClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_healthcare_service" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("HealthCare Service still exists:\n%#v", resp.Status)
		}
	}

	return nil
}

func testAccAzureRMHealthcareService_basic(rInt int) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "westus2"
}

resource "azurerm_healthcare_service" "test" {
  name                = "accfa-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    environment = "production"
    purpose     = "AcceptanceTests"
  }

  access_policy_object_ids {
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"
  }
}
`, rInt, rInt)
}
