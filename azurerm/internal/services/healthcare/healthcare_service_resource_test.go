package healthcare_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/healthcare/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type HealthCareServiceResource struct {
}

func TestAccHealthCareService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_service", "test")
	r := HealthCareServiceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHealthCareService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_service", "test")
	r := HealthCareServiceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccHealthCareService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_service", "test")
	r := HealthCareServiceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (HealthCareServiceResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.HealthCare.HealthcareServiceClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Healthcare service %q (resource group: %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (HealthCareServiceResource) basic(data acceptance.TestData) string {
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

func (r HealthCareServiceResource) requiresImport(data acceptance.TestData) string {
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
`, r.basic(data))
}

func (HealthCareServiceResource) complete(data acceptance.TestData) string {
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
