package healthcare_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"testing"
)

type HealthCareIotConnectorFhirDestinationResource struct{}

func TestAccHealthCareIotConnectorFhirDestination_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_iot_fhir_destination", "test")
	r := HealthCareIotConnectorFhirDestinationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccHealthCareIotConnectorFhirDestination_updateTemplate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_iot_fhir_destination", "test")
	r := HealthCareIotConnectorFhirDestinationResource{}

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

func TestAccHealthCareIotConnectorFhirDestination_updateFhir(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_iot_fhir_destination", "test")
	r := HealthCareIotConnectorFhirDestinationResource{}

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

func TestAccHealthCareIotConnectorFhirDestination_updateResolutionType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_iot_fhir_destination", "test")
	r := HealthCareIotConnectorFhirDestinationResource{}

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

func (r HealthCareIotConnectorFhirDestinationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.IotFhirDestinationID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.HealthCare.HealthcareWorkspaceIotConnectorFhirDestinationClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.IotconnectorName, id.FhirdestinationName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s, %+v", *id, err)
	}
	return utils.Bool(resp.IotFhirDestinationProperties != nil), nil
}

func (r HealthCareIotConnectorFhirDestinationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_healthcare_iot_fhir_destination" "test" {
  name                                 = "des%d"
  location                             = azurerm_resource_group.test.location
  iot_connector_id                     = azurerm_healthcare_iot_connector.test.id
  destination_fhir_service_id          = azurerm_healthcare_fhir_service.test.id
  destination_identity_resolution_type = "Create"
  destination_fhir_mapping             = <<JSON
{
"templateType": "CollectionFhirTemplate",
"template": []
}
JSON
depends_on = [azurerm_healthcare_fhir_service.test]
}                              
`, r.template(data), data.RandomInteger)
}

func (r HealthCareIotConnectorFhirDestinationResource) updateTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_healthcare_iot_fhir_destination" "test" {
  name                                 = "des%d"
  location                             = azurerm_resource_group.test.location
  iot_connector_id                     = azurerm_healthcare_iot_connector.test.id
  destination_fhir_service_id          = azurerm_healthcare_fhir_service.test.id
  destination_identity_resolution_type = "Create"
  destination_fhir_mapping             = <<JSON
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

func (r HealthCareIotConnectorFhirDestinationResource) updateFhirId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_healthcare_iot_fhir_destination" "test" {
  name                                 = "des%d"
  location                             = azurerm_resource_group.test.location
  iot_connector_id                     = azurerm_healthcare_iot_connector.test.id
  destination_fhir_service_id          = azurerm_healthcare_fhir_service.test1.id
  destination_identity_resolution_type = "Create"
  destination_fhir_mapping             = <<JSON
{
"templateType": "CollectionFhirTemplate",
"template": []
}
JSON
}`, r.template(data), data.RandomInteger)
}

func (r HealthCareIotConnectorFhirDestinationResource) updateResolutionType(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_healthcare_iot_fhir_destination" "test" {
  name                                 = "des%d"
  location                             = azurerm_resource_group.test.location
  iot_connector_id                     = azurerm_healthcare_iot_connector.test.id
  destination_fhir_service_id          = azurerm_healthcare_fhir_service.test.id
  destination_identity_resolution_type = "Lookup"
  destination_fhir_mapping             = <<JSON
{
"templateType": "CollectionFhirTemplate",
"template": []
}
JSON
}`, r.template(data), data.RandomInteger)
}

func (r HealthCareIotConnectorFhirDestinationResource) template(data acceptance.TestData) string {
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
  depends_on = [azurerm_healthcare_workspace.test]
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
  depends_on = [azurerm_healthcare_workspace.test]
}

`, HealthCareWorkspaceIotConnectorResource{}.basic(data), data.RandomInteger, data.RandomInteger)
}
