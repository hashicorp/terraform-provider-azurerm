package healthcare_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type HealthCareApisDicomResource struct{}

func (HealthCareApisDicomResource) basic (data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-hcwk-%d"
  location = "%s"
}

resource "azurerm_healthcare_apis_dicom_service" "test"{
  name = "acctest-dicom%d"
  workspace_id = "/subscriptions/85b3dbca-5974-4067-9669-67a141095a76/resourceGroups/xiaxintestrg-healthcare/providers/Microsoft.HealthcareApis/workspaces/xiaxintestworkspace"
  location = "east us"
  public_network_access_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(8))
}


