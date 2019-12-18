package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMHealthcareService_basic(t *testing.T) {
	ri := tf.AccRandTimeInt() / 10

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMHealthcareServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHealthcareService_basic(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHealthcareServiceExists("azurerm_healthcare_service.test"),
				),
			},
			{
				ResourceName:      "azurerm_healthcare_service.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMHealthcareServiceExists(resourceName string) resource.TestCheckFunc {
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

		client := testAccProvider.Meta().(*ArmClient).HealthCare.HealthcareServiceClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, healthcareServiceName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Bad: Healthcare service %q (resource group: %q) does not exist", healthcareServiceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on healthcareServiceClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMHealthcareServiceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).HealthCare.HealthcareServiceClient
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
	// currently only supported in "ukwest", "northcentralus", "westus2".
	location := "westus2"
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-health-%d"
  location = "%s"
}

resource "azurerm_healthcare_service" "test" {
  name                = "testacc%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    environment = "production"
    purpose     = "AcceptanceTests"
  }

  access_policy_object_ids = [
    "${data.azurerm_client_config.current.service_principal_object_id}",
  ]

  authentication_configuration {
    authority           = "https://login.microsoftonline.com/${data.azurerm_client_config.current.tenant_id}"
    audience            = "https://azurehealthcareapis.com"
    smart_proxy_enabled = "true"
  }

  cors_configuration {
    allowed_origins    = ["http://www.example.com", "http://www.example2.com"]
    allowed_headers    = ["x-tempo-*", "x-tempo2-*"]
    allowed_methods    = ["GET", "PUT"]
    max_age_in_seconds = "500"
    allow_credentials  = "true"
  }
}
`, rInt, location, rInt)
}
