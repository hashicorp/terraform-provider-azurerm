// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package healthcare_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/iotconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type HealthCareWorkspaceMedTechServiceResource struct{}

func TestAccHealthCareMedTechService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_medtech_service", "test")
	r := HealthCareWorkspaceMedTechServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccHealthCareMedTechService_updateIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_medtech_service", "test")
	r := HealthCareWorkspaceMedTechServiceResource{}

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

func TestAccHealthCareMedTechService_updateTemplate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_medtech_service", "test")
	r := HealthCareWorkspaceMedTechServiceResource{}

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

func TestAccHealthCareMedTechService_updateEventhubs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_medtech_service", "test")
	r := HealthCareWorkspaceMedTechServiceResource{}

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

func TestAccHealthCareMedTechService_updateConsumerGroups(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_medtech_service", "test")
	r := HealthCareWorkspaceMedTechServiceResource{}

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

func (r HealthCareWorkspaceMedTechServiceResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := iotconnectors.ParseIotConnectorID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.HealthCare.HealthcareWorkspaceIotConnectorsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s, %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r HealthCareWorkspaceMedTechServiceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_healthcare_medtech_service" "test" {
  name         = "mt%d"
  workspace_id = azurerm_healthcare_workspace.test.id
  location     = azurerm_resource_group.test.location

  eventhub_namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name                = azurerm_eventhub.test.name
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.test.name

  device_mapping_json = <<JSON
{
"templateType": "CollectionContent",
"template": []
}
JSON
}
`, r.template(data), data.RandomInteger)
}

func (r HealthCareWorkspaceMedTechServiceResource) identity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_healthcare_medtech_service" "test" {
  name         = "mt%d"
  workspace_id = azurerm_healthcare_workspace.test.id
  location     = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }

  eventhub_namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name                = azurerm_eventhub.test.name
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.test.name

  device_mapping_json = <<JSON
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

func (r HealthCareWorkspaceMedTechServiceResource) updateEventhubs(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_healthcare_medtech_service" "test" {
  name         = "mt%d"
  workspace_id = azurerm_healthcare_workspace.test.id
  location     = azurerm_resource_group.test.location

  eventhub_namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name                = azurerm_eventhub.test1.name
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.test.name

  device_mapping_json = <<JSON
{
"templateType": "CollectionContent",
"template": []
}
JSON
}
`, r.template(data), data.RandomInteger)
}

func (r HealthCareWorkspaceMedTechServiceResource) updateConsumerGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_healthcare_medtech_service" "test" {
  name         = "mt%d"
  workspace_id = azurerm_healthcare_workspace.test.id
  location     = azurerm_resource_group.test.location

  eventhub_namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name                = azurerm_eventhub.test.name
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.test1.name

  device_mapping_json = <<JSON
{
"templateType": "CollectionContent",
"template": []
}
JSON
}
`, r.template(data), data.RandomInteger)
}

func (r HealthCareWorkspaceMedTechServiceResource) updateTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_healthcare_medtech_service" "test" {
  name         = "mt%d"
  workspace_id = azurerm_healthcare_workspace.test.id
  location     = azurerm_resource_group.test.location

  eventhub_namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name                = azurerm_eventhub.test.name
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.test.name

  device_mapping_json = <<JSON
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

func (HealthCareWorkspaceMedTechServiceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-medTech-%d"
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
  name                = "acctest-eh1-%d"
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
  name                = "acctestCG1-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test1.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_healthcare_workspace" "test" {
  name                = "wks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
