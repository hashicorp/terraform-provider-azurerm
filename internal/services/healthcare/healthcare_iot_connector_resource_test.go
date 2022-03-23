package healthcare_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type HealthCareWorkspaceIotConnectorResource struct{}

func TestAccHealthCareIotConnector_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_iot_connector", "test")
	r := HealthCareWorkspaceIotConnectorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccHealthCareIotConnector_updateIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_iot_connector", "test")
	r := HealthCareWorkspaceIotConnectorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
		{
			Config: r.identity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccHealthCareIotConnector_updateTemplate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_iot_connector", "test")
	r := HealthCareWorkspaceIotConnectorResource{}

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

func TestAccHealthCareIotConnector_updateEventhubs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_iot_connector", "test")
	r := HealthCareWorkspaceIotConnectorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
		{
			Config: r.updateEventhubs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccHealthCareIotConnector_updateConsumerGroups(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_iot_connector", "test")
	r := HealthCareWorkspaceIotConnectorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
		{
			Config: r.updateConsumerGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func (r HealthCareWorkspaceIotConnectorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.IotConnectorID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.HealthCare.HealthcareWorkspaceIotConnectorClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s, %+v", *id, err)
	}
	return utils.Bool(resp.IotConnectorProperties != nil), nil
}

func (r HealthCareWorkspaceIotConnectorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_healthcare_iot_connector" "test" {
  name         = "iot%d"
  workspace_id = azurerm_healthcare_workspace.test.id
  location     = azurerm_resource_group.test.location

  eventhub_namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name                = azurerm_eventhub.test.name
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.test.name
  device_mapping               = <<JSON
{
"templateType": "CollectionContent",
"template": []
}
JSON
}
`, r.template(data), data.RandomInteger)
}

func (r HealthCareWorkspaceIotConnectorResource) identity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_healthcare_iot_connector" "test" {
  name         = "iot%d"
  workspace_id = azurerm_healthcare_workspace.test.id
  location     = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
  eventhub_namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name                = azurerm_eventhub.test.name
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.test.name
  device_mapping               = <<JSON
{
            "templateType": "CollectionContent",
            "template": [
              {
                "templateType": "JsonPathContent",
                "template": {
                  "typeName": "heartrate",
                  "typeMatchExpression": "$..[?(@heartrate)]",
                  "deviceIdExpression": "$.deviceid",
                  "timestampExpression": "$.measurementdatetime",
                  "values": [
                    {
                      "required": "true",
                      "valueExpression": "$.heartrate",
                      "valueName": "hr"
                    }
                  ]
                }
              }
            ]
}
JSON
}`, r.template(data), data.RandomInteger)
}

func (r HealthCareWorkspaceIotConnectorResource) updateEventhubs(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_healthcare_iot_connector" "test" {
  name         = "iot%d"
  workspace_id = azurerm_healthcare_workspace.test.id
  location     = azurerm_resource_group.test.location

  eventhub_namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name                = azurerm_eventhub.test1.name
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.test.name
  device_mapping               = <<JSON
{
"templateType": "CollectionContent",
"template": []
}
JSON
}
`, r.template(data), data.RandomInteger)
}

func (r HealthCareWorkspaceIotConnectorResource) updateConsumerGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_healthcare_iot_connector" "test" {
  name         = "iot%d"
  workspace_id = azurerm_healthcare_workspace.test.id
  location     = azurerm_resource_group.test.location

  eventhub_namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name                = azurerm_eventhub.test.name
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.test1.name
  device_mapping               = <<JSON
{
"templateType": "CollectionContent",
"template": []
}
JSON
}
`, r.template(data), data.RandomInteger)
}

func (r HealthCareWorkspaceIotConnectorResource) updateTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_healthcare_iot_connector" "test" {
  name         = "iot%d"
  workspace_id = azurerm_healthcare_workspace.test.id
  location     = azurerm_resource_group.test.location

  eventhub_namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name                = azurerm_eventhub.test.name
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.test.name
  device_mapping               = <<JSON
{
   "templateType": "CollectionContent",
            "template": [
              {
                "templateType": "JsonPathContent",
                "template": {
                  "typeName": "heartrate",
                  "typeMatchExpression": "$..[?(@heartrate)]",
                  "deviceIdExpression": "$.deviceid",
                  "timestampExpression": "$.measurementdatetime",
                  "values": [
                    {
                      "required": "true",
                      "valueExpression": "$.heartrate",
                      "valueName": "hr"
                    }
                  ]
                }
              }
            ]
}
JSON
}
`, r.template(data), data.RandomInteger)
}

func (HealthCareWorkspaceIotConnectorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iot-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-ehn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest-eh-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 1
  message_retention   = 1
}

resource "azurerm_eventhub" "test1" {
  name                = "acctest-eh-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 1
  message_retention   = 1
}

resource "azurerm_eventhub_consumer_group" "test" {
  name                = "acctestCG-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_eventhub_consumer_group" "test1" {
  name                = "acctestCG-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_healthcare_workspace" "test" {
  name                = "wks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
