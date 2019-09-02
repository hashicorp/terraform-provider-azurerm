package azurerm

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/preview/healthcareapis/mgmt/2018-08-20-preview/healthcareapis"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMFhirApi(t *testing.T) {
	var fhir_api healthcareapis.ServicesDescription
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFhirApiDestroy,
		Steps: []resource.TestStep{
			{
				//Config: testAccAzureRMFhirApi_basic(ri, testLocation()),
				Config: testAccAzureRMFhirApi_basic(ri),
				Check: resource.ComposeTestCheckFunc(
					//testCheckAzureRMFhirApiExists("azurerm_fhir_api_service.test", &fhir_api),
					testCheckAzureRMFhirApiExists("azurerm_fhir_api_service.test", &fhir_api),
				),
			},
			{
				ResourceName: "azurerm_fhir_api_service.test",
				ImportState:  true,
				ImportStateVerifyIgnore: []string{
					// since these are read from the existing state
					"access_policy_object_ids",
					"cosmodb_throughput",
					"kind",
				},
			},
		},
	})
}

func testCheckAzureRMFhirApiExists(resourceName string, fhir_api *healthcareapis.ServicesDescription) resource.TestCheckFunc {
	//func testCheckAzureRMFhirApiExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		fhirApiName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for loadbalancer: %s", fhirApiName)
		}

		client := testAccProvider.Meta().(*ArmClient).fhirApiServiceClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, fhirApiName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Bad: Fhir Api service %q (resource group: %q) does not exist", fhirApiName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on loadBalancerClient: %+v", err)
		}

		*fhir_api = resp

		return nil
	}
}

func testCheckAzureRMFhirApiDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).fhirApiServiceClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_fhir_api_service" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("LoadBalancer still exists:\n%#v", resp.Status)
		}
	}

	return nil
}

//func testAccAzureRMFhirApi_basic(rInt int, location string) string {
func testAccAzureRMFhirApi_basic(rInt int) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "westus2"
}

resource "azurerm_fhir_api_service" "test" {
  name                = "accfa-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    environment = "production"
    purpose     = "AcceptanceTests"
  }

  access_policy_object_ids {
    object_id          = "${data.azurerm_client_config.current.service_principal_object_id}"
  }
}
`, rInt, rInt)
}
