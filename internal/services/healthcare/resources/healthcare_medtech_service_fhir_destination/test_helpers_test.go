// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package healthcare_medtech_service_fhir_destination_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type HealthCareWorkspaceMedTechServiceResource struct{}

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
  name              = "acctest-eh-%d"
  namespace_id      = azurerm_eventhub_namespace.test.id
  partition_count   = 1
  message_retention = 1
}

resource "azurerm_eventhub" "test1" {
  name              = "acctest-eh1-%d"
  namespace_id      = azurerm_eventhub_namespace.test.id
  partition_count   = 1
  message_retention = 1
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
