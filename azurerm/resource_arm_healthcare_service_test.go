package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMHealthCareService_basic(t *testing.T) {
	ri := tf.AccRandTimeInt() / 10
	// currently only supported in "ukwest", "northcentralus", "westus2".
	location := "westus2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHealthCareServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHealthCareService_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHealthCareServiceExists("azurerm_healthcare_service.test"),
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

func TestAccAzureRMHealthCareService_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	ri := tf.AccRandTimeInt() / 10
	// currently only supported in "ukwest", "northcentralus", "westus2".
	location := "westus2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHealthCareServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHealthCareService_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHealthCareServiceExists("azurerm_healthcare_service.test"),
				),
			},
			{
				Config:      testAccAzureRMHealthCareService_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_healthcare_service"),
			},
		},
	})
}

func TestAccAzureRMHealthCareService_complete(t *testing.T) {
	ri := tf.AccRandTimeInt() / 10
	// currently only supported in "ukwest", "northcentralus", "westus2".
	location := "westus2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHealthCareServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHealthCareService_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHealthCareServiceExists("azurerm_healthcare_service.test"),
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

func testCheckAzureRMHealthCareServiceExists(resourceName string) resource.TestCheckFunc {
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

		client := acceptance.AzureProvider.Meta().(*clients.Client).HealthCare.HealthcareServiceClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testCheckAzureRMHealthCareServiceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).HealthCare.HealthcareServiceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMHealthCareService_basic(rInt int, location string) string {
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

  access_policy_object_ids = [
    "${data.azurerm_client_config.current.service_principal_object_id}",
  ]
}
`, rInt, location, rInt)
}

func testAccAzureRMHealthCareService_requiresImport(rInt int, location string) string {
	template := testAccAzureRMHealthCareService_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_healthcare_service" "import" {
  name                = azurerm_healthcare_service.test.name
  location            = azurerm_healthcare_service.test.location
  resource_group_name = azurerm_healthcare_service.test.resource_group_name

  access_policy_object_ids = [
    "${data.azurerm_client_config.current.service_principal_object_id}",
  ]
}
`, template)
}

func testAccAzureRMHealthCareService_complete(rInt int, location string) string {
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
    smart_proxy_enabled = true
  }

  cors_configuration {
    allowed_origins    = ["http://www.example.com", "http://www.example2.com"]
    allowed_headers    = ["*"]
    allowed_methods    = ["GET", "PUT"]
    max_age_in_seconds = 500
    allow_credentials  = true
  }
}
`, rInt, location, rInt)
}
