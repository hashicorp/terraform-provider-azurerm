package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/healthcare/parse"
)

func TestAccAzureRMHealthCareService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHealthCareServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHealthCareService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHealthCareServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMHealthCareService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHealthCareServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHealthCareService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHealthCareServiceExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMHealthCareService_requiresImport),
		},
	})
}

func TestAccAzureRMHealthCareService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHealthCareServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHealthCareService_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHealthCareServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMHealthCareServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).HealthCare.HealthcareServiceClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.ServiceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Bad: Healthcare service %q (resource group: %q) does not exist", id.Name, id.ResourceGroup)
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

		id, err := parse.ServiceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("HealthCare Service still exists:\n%#v", resp.Status)
		}
	}

	return nil
}

func testAccAzureRMHealthCareService_basic(data acceptance.TestData) string {
	// currently only supported in "ukwest", "northcentralus", "westus2".
	location := "westus2"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-health-%d"
  location = "%s"
}

resource "azurerm_healthcare_service" "test" {
  name                = "testacc%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  access_policy_object_ids = [
    data.azurerm_client_config.current.object_id,
  ]
}
`, data.RandomInteger, location, data.RandomIntOfLength(17)) // name can only be 24 chars long
}

func testAccAzureRMHealthCareService_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMHealthCareService_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_healthcare_service" "import" {
  name                = azurerm_healthcare_service.test.name
  location            = azurerm_healthcare_service.test.location
  resource_group_name = azurerm_healthcare_service.test.resource_group_name

  access_policy_object_ids = [
    "${data.azurerm_client_config.current.object_id}",
  ]
}
`, template)
}

func testAccAzureRMHealthCareService_complete(data acceptance.TestData) string {
	// currently only supported in "ukwest", "northcentralus", "westus2".
	location := "westus2"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-health-%d"
  location = "%s"
}

resource "azurerm_healthcare_service" "test" {
  name                = "testacc%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "production"
    purpose     = "AcceptanceTests"
  }

  access_policy_object_ids = [
    data.azurerm_client_config.current.object_id,
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
`, data.RandomInteger, location, data.RandomIntOfLength(17)) // name can only be 24 chars long
}
