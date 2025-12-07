// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package healthcare_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/fhirservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type HealthcareApiFhirServiceResource struct{}

func TestAccHealthcareApiFhirService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_fhir_service", "test")
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
	data := acceptance.BuildTestData(t, "azurerm_healthcare_fhir_service", "test")
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

func TestAccHealthcareApiFhirService_updateIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_fhir_service", "test")
	r := HealthcareApiFhirServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateIdentitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateIdentityUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateIdentitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHealthcareApiFhirService_updateAcrLoginServer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_fhir_service", "test")
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
	data := acceptance.BuildTestData(t, "azurerm_healthcare_fhir_service", "test")
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

func TestAccHealthcareApiFhirService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_fhir_service", "test")
	r := HealthcareApiFhirServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
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
	return pointer.To(resp.Model != nil), nil
}

func (r HealthcareApiFhirServiceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_healthcare_fhir_service" "test" {
  name                = "fhir%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_healthcare_workspace.test.id
  kind                = "fhir-R4"

  authentication {
    authority = "https://login.microsoftonline.com/72f988bf-86f1-41af-91ab-2d7cd011db47"
    audience  = "https://acctestfhir.fhir.azurehealthcareapis.com"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HealthcareApiFhirServiceResource) updateIdentitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_healthcare_fhir_service" "test" {
  name                = "fhir%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_healthcare_workspace.test.id
  kind                = "fhir-R4"

  authentication {
    authority = "https://login.microsoftonline.com/72f988bf-86f1-41af-91ab-2d7cd011db47"
    audience  = "https://acctestfhir.fhir.azurehealthcareapis.com"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HealthcareApiFhirServiceResource) updateIdentityUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_healthcare_fhir_service" "test" {
  name                = "fhir%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_healthcare_workspace.test.id
  kind                = "fhir-R4"

  authentication {
    authority = "https://login.microsoftonline.com/72f988bf-86f1-41af-91ab-2d7cd011db47"
    audience  = "https://acctestfhir.fhir.azurehealthcareapis.com"
  }

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r HealthcareApiFhirServiceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_healthcare_fhir_service" "import" {
  name                = azurerm_healthcare_fhir_service.test.name
  location            = azurerm_healthcare_fhir_service.test.location
  resource_group_name = azurerm_healthcare_fhir_service.test.resource_group_name
  workspace_id        = azurerm_healthcare_fhir_service.test.workspace_id
  kind                = azurerm_healthcare_fhir_service.test.kind

  authentication {
    authority = azurerm_healthcare_fhir_service.test.authentication[0].authority
    audience  = azurerm_healthcare_fhir_service.test.authentication[0].audience
  }
}
`, r.basic(data))
}

func (r HealthcareApiFhirServiceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {
}
%s

resource "azurerm_container_registry" "test" {
  name                = "acc%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
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

resource "azurerm_healthcare_fhir_service" "test" {
  name                = "fhir%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_healthcare_workspace.test.id
  kind                = "fhir-R4"

  authentication {
    authority = "https://login.microsoftonline.com/72f988bf-86f1-41af-91ab-2d7cd011db47"
    audience  = "https://acctestfhir.fhir.azurehealthcareapis.com"
  }

  access_policy_object_ids = []

  identity {
    type = "SystemAssigned"
  }

  container_registry_login_server_url = [azurerm_container_registry.test.login_server]

  oci_artifact {
    login_server = azurerm_container_registry.test.login_server
  }

  cors {
    allowed_origins     = ["https://acctest.com:123", "https://acctest1.com:3389"]
    allowed_headers     = ["*"]
    allowed_methods     = ["GET", "DELETE", "PUT"]
    max_age_in_seconds  = 3600
    credentials_allowed = true
  }

  configuration_export_storage_account_name = azurerm_storage_account.test.name
}
`, r.template(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary, data.RandomInteger, data.RandomInteger)
}

func (r HealthcareApiFhirServiceResource) updateAcrLoginServer(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {
}
%s

resource "azurerm_container_registry" "test" {
  name                = "acc%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
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

resource "azurerm_healthcare_fhir_service" "test" {
  name                = "fhir%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_healthcare_workspace.test.id
  kind                = "fhir-R4"

  authentication {
    authority = "https://login.microsoftonline.com/72f988bf-86f1-41af-91ab-2d7cd011db47"
    audience  = "https://acctestfhir.fhir.azurehealthcareapis.com"
  }

  access_policy_object_ids = []

  identity {
    type = "SystemAssigned"
  }

  container_registry_login_server_url = []

  cors {
    allowed_origins     = ["https://acctest.com:123", "https://acctest1.com:3389"]
    allowed_headers     = ["*"]
    allowed_methods     = ["GET", "DELETE", "PUT"]
    max_age_in_seconds  = 3600
    credentials_allowed = true
  }

  configuration_export_storage_account_name = azurerm_storage_account.test.name
}
`, r.template(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary, data.RandomInteger, data.RandomInteger)
}

func (r HealthcareApiFhirServiceResource) updateCors(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {
}
%s

resource "azurerm_container_registry" "test" {
  name                = "acc%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
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

resource "azurerm_healthcare_fhir_service" "test" {
  name                = "fhir%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_healthcare_workspace.test.id
  kind                = "fhir-R4"

  authentication {
    authority = "https://login.microsoftonline.com/72f988bf-86f1-41af-91ab-2d7cd011db47"
    audience  = "https://acctestfhir.fhir.azurehealthcareapis.com"
  }

  access_policy_object_ids = []

  identity {
    type = "SystemAssigned"
  }

  container_registry_login_server_url = [azurerm_container_registry.test.login_server]

  cors {
    allowed_origins     = ["https://acctest.com:123"]
    allowed_headers     = ["*"]
    allowed_methods     = ["GET", "DELETE"]
    max_age_in_seconds  = 0
    credentials_allowed = false
  }

  configuration_export_storage_account_name = azurerm_storage_account.test.name
}
`, r.template(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary, data.RandomInteger, data.RandomInteger)
}

func (HealthcareApiFhirServiceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-healthcareapi-%d"
  location = "%s"
}

resource "azurerm_healthcare_workspace" "test" {
  name                = "acc%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
