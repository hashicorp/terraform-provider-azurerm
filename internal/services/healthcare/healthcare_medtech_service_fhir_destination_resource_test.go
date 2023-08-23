// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package healthcare_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/iotconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type HealthCareMedTechServiceFhirDestinationResource struct{}

func TestAccHealthCareMedTechServiceFhirDestination_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_medtech_service_fhir_destination", "test")
	r := HealthCareMedTechServiceFhirDestinationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccHealthCareMedTechServiceFhirDestination_updateTemplate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_medtech_service_fhir_destination", "test")
	r := HealthCareMedTechServiceFhirDestinationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
		{
			Config: r.updateTemplate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccHealthCareMedTechServiceFhirDestination_updateFhir(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_medtech_service_fhir_destination", "test")
	r := HealthCareMedTechServiceFhirDestinationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
		{
			Config: r.updateFhirId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccHealthCareMedTechServiceFhirDestination_updateResolutionType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_medtech_service_fhir_destination", "test")
	r := HealthCareMedTechServiceFhirDestinationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
		{
			Config: r.updateResolutionType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func (r HealthCareMedTechServiceFhirDestinationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := iotconnectors.ParseFhirDestinationID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.HealthCare.HealthcareWorkspaceIotConnectorsClient.IotConnectorFhirDestinationGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s, %+v", *id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r HealthCareMedTechServiceFhirDestinationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_healthcare_medtech_service_fhir_destination" "test" {
  name                                 = "des%d"
  location                             = azurerm_resource_group.test.location
  medtech_service_id                   = azurerm_healthcare_medtech_service.test.id
  destination_fhir_service_id          = azurerm_healthcare_fhir_service.test.id
  destination_identity_resolution_type = "Create"

  destination_fhir_mapping_json = <<JSON
{
"templateType": "CollectionFhirTemplate",
"template": []
}
JSON
}
`, r.template(data), data.RandomInteger)
}

func (r HealthCareMedTechServiceFhirDestinationResource) updateTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_healthcare_medtech_service_fhir_destination" "test" {
  name                                 = "des%d"
  location                             = azurerm_resource_group.test.location
  medtech_service_id                   = azurerm_healthcare_medtech_service.test.id
  destination_fhir_service_id          = azurerm_healthcare_fhir_service.test.id
  destination_identity_resolution_type = "Create"

  destination_fhir_mapping_json = <<JSON
{
  "templateType": "CollectionFhirTemplate",
            "template": [
              {
                "templateType": "CodeValueFhir",
                "template": {
                  "codes": [
                    {
                      "code": "8867-4",
                      "system": "http://loinc.org",
                      "display": "Heart rate"
                    }
                  ],
                  "periodInterval": 60,
                  "typeName": "heartrate",
                  "value": {
                    "defaultPeriod": 5000,
                    "unit": "count/min",
                    "valueName": "hr",
                    "valueType": "SampledData"
                  }
                }
              }
            ]
}
JSON
}
`, r.template(data), data.RandomInteger)
}

func (r HealthCareMedTechServiceFhirDestinationResource) updateFhirId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_healthcare_medtech_service_fhir_destination" "test" {
  name                                 = "des%d"
  location                             = azurerm_resource_group.test.location
  medtech_service_id                   = azurerm_healthcare_medtech_service.test.id
  destination_fhir_service_id          = azurerm_healthcare_fhir_service.test1.id
  destination_identity_resolution_type = "Create"

  destination_fhir_mapping_json = <<JSON
{
"templateType": "CollectionFhirTemplate",
"template": []
}
JSON
}`, r.template(data), data.RandomInteger)
}

func (r HealthCareMedTechServiceFhirDestinationResource) updateResolutionType(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_healthcare_medtech_service_fhir_destination" "test" {
  name                                 = "des%d"
  location                             = azurerm_resource_group.test.location
  medtech_service_id                   = azurerm_healthcare_medtech_service.test.id
  destination_fhir_service_id          = azurerm_healthcare_fhir_service.test.id
  destination_identity_resolution_type = "Lookup"

  destination_fhir_mapping_json = <<JSON
{
"templateType": "CollectionFhirTemplate",
"template": []
}
JSON
}`, r.template(data), data.RandomInteger)
}

func (r HealthCareMedTechServiceFhirDestinationResource) template(data acceptance.TestData) string {
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

resource "azurerm_healthcare_fhir_service" "test1" {
  name                = "fhir1%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_healthcare_workspace.test.id
  kind                = "fhir-R4"

  authentication {
    authority = "https://login.microsoftonline.com/72f988bf-86f1-41af-91ab-2d7cd011db47"
    audience  = "https://acctestfhir.fhir.azurehealthcareapis.com"
  }
}
`, HealthCareWorkspaceMedTechServiceResource{}.basic(data), data.RandomInteger, data.RandomInteger)
}
