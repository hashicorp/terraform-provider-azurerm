package healthcare_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/sdk/2021-06-01-preview/fhirservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type HealthcareApiFhirServiceResource struct{}

func TestAccHealthcareApiFhirService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcareapis_fhir_service", "test")
	r := HealthcareApiFhirServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHealthcareApiFhirService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcareapis_fhir_service", "test")
	r := HealthcareApiFhirServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHealthcareApiFhirService_updateAcrLoginServer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcareapis_fhir_service", "test")
	r := HealthcareApiFhirServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateAcrLoginServer(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHealthcareApiFhirService_updateCors(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcareapis_fhir_service", "test")
	r := HealthcareApiFhirServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateCors(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHealthcareApisService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcareapis_fhir_service", "test")
	r := HealthcareApiFhirServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (HealthcareApiFhirServiceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fhirservices.ParseFhirServiceID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := clients.HealthCare.HealthcareWorkspaceFhirServiceClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving Healthcare api fhir service %s: %+v", *id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (HealthcareApiFhirServiceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-healthcareapi-%d"
  location = "%s"
}

resource "azurerm_healthcareapis_workspace" "test" {
  name                = "acc%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_healthcareapis_fhir_service" "test" {
  name                = "fhir%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_healthcareapis_workspace.test.id
  kind                = "fhir-R4"
  authentication_configuration {
    authority = "https://login.microsoftonline.com/72f988bf-86f1-41af-91ab-2d7cd011db47"
    audience  = "https://acctestfhir.fhir.azurehealthcareapis.com"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r HealthcareApiFhirServiceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_healthcareapis_fhir_service" "test" {
  name                = azurerm_healthcareapis_fhir_service.test.name
  location            = azurerm_healthcareapis_fhir_service.test.location
  resource_group_name = azurerm_healthcareapis_fhir_service.test.resource_group_name
  workspace_id        = azurerm_healthcareapis_fhir_service.test.workspace_id
  kind                = azurerm_healthcareapis_fhir_service.test.kind
  authentication_configuration {
    authority = azurerm_healthcareapis_fhir_service.test.authentication_configuration[0].authority
    audience  = azurerm_healthcareapis_fhir_service.test.authentication_configuration[0].audience
  }
}
`, r.basic(data))
}

func (HealthcareApiFhirServiceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-healthcareapi-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "acc%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  admin_enabled       = false
  georeplications {
    location                = "%s"
    zone_redundancy_enabled = true
    tags                    = {}
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acc%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_healthcareapis_workspace" "test" {
  name                = "acc%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_healthcareapis_fhir_service" "test" {
  name                = "fhir%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_healthcareapis_workspace.test.id
  kind                = "fhir-R4"
  authentication_configuration {
    authority = "https://login.microsoftonline.com/72f988bf-86f1-41af-91ab-2d7cd011db47"
    audience  = "https://acctestfhir.fhir.azurehealthcareapis.com"
  }
  identity {
    type = "SystemAssigned"
  }

  acr_login_servers = [azurerm_container_registry.test.login_server]

  cors_configuration {
    allowed_origins    = ["https://acctest.com:123", "https://acctest1.com:3389"]
    allowed_headers    = ["*"]
    allowed_methods    = ["GET", "DELETE", "PUT"]
    max_age_in_seconds = 3600
    allow_credentials  = true
  }
  export_storage_account_name = azurerm_storage_account.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Secondary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (HealthcareApiFhirServiceResource) updateAcrLoginServer(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-healthcareapi-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acc%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_healthcareapis_workspace" "test" {
  name                = "acc%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_healthcareapis_fhir_service" "test" {
  name                = "fhir%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_healthcareapis_workspace.test.id
  kind                = "fhir-R4"
  authentication_configuration {
    authority = "https://login.microsoftonline.com/72f988bf-86f1-41af-91ab-2d7cd011db47"
    audience  = "https://acctestfhir.fhir.azurehealthcareapis.com"
  }
  identity {
    type = "SystemAssigned"
  }

  acr_login_servers = []

  cors_configuration {
    allowed_origins    = ["https://acctest.com:123", "https://acctest1.com:3389"]
    allowed_headers    = ["*"]
    allowed_methods    = ["GET", "DELETE", "PUT"]
    max_age_in_seconds = 3600
    allow_credentials  = true
  }

  export_storage_account_name = azurerm_storage_account.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (HealthcareApiFhirServiceResource) updateCors(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-healthcareapi-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "acc%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  admin_enabled       = false
  georeplications {
    location                = "%s"
    zone_redundancy_enabled = true
    tags                    = {}
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acc%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_healthcareapis_workspace" "test" {
  name                = "acc%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_healthcareapis_fhir_service" "test" {
  name                = "fhir%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_healthcareapis_workspace.test.id
  kind                = "fhir-R4"
  authentication_configuration {
    authority = "https://login.microsoftonline.com/72f988bf-86f1-41af-91ab-2d7cd011db47"
    audience  = "https://acctestfhir.fhir.azurehealthcareapis.com"
  }
  identity {
    type = "SystemAssigned"
  }

  acr_login_servers = []

  cors_configuration {
    allowed_origins    = ["https://acctest.com:123"]
    allowed_headers    = ["*"]
    allowed_methods    = ["GET", "DELETE"]
    max_age_in_seconds = 0
    allow_credentials  = false
  }

  export_storage_account_name = azurerm_storage_account.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Secondary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
