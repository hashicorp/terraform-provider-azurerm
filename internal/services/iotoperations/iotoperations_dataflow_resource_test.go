// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iotoperations_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/dataflow"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// IotOperationsDataflowResource is a test harness for azurerm_iotoperations_dataflow acceptance tests.
type IotOperationsDataflowResource struct{}

func TestAccIotOperationsDataflow_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_dataflow", "test")
	r := IotOperationsDataflowResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctest-dataflow-%s", data.RandomString)),
				check.That(data.ResourceName).Key("properties.0.mode").HasValue("Enabled"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotOperationsDataflow_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_dataflow", "test")
	r := IotOperationsDataflowResource{}

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

func TestAccIotOperationsDataflow_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_dataflow", "test")
	r := IotOperationsDataflowResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctest-dataflow-%s", data.RandomString)),
				check.That(data.ResourceName).Key("properties.0.mode").HasValue("Enabled"),
				check.That(data.ResourceName).Key("properties.0.operations.#").HasValue("2"),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotOperationsDataflow_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_dataflow", "test")
	r := IotOperationsDataflowResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("properties.0.operations.#").HasValue("2"),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func (r IotOperationsDataflowResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dataflow.ParseDataflowID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %+v", state.ID, err)
	}

	resp, err := clients.IoTOperations.DataflowClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", state.ID, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r IotOperationsDataflowResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iot-%d"
  location = "%s"
}

resource "azurerm_iotoperations_instance" "test" {
  name                = "acctest-instance-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  
  extended_location {
    name = "acctest-custom-location-%s"
    type = "CustomLocation"
  }
}

resource "azurerm_iotoperations_dataflow_profile" "test" {
  name                = "acctest-profile-%s"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = azurerm_iotoperations_instance.test.name
  location            = azurerm_resource_group.test.location

  extended_location {
    name = azurerm_iotoperations_instance.test.extended_location[0].name
    type = azurerm_iotoperations_instance.test.extended_location[0].type
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomString)
}

func (r IotOperationsDataflowResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_dataflow" "test" {
  name                   = "acctest-dataflow-%s"
  resource_group_name    = azurerm_resource_group.test.name
  instance_name         = azurerm_iotoperations_instance.test.name
  dataflow_profile_name = azurerm_iotoperations_dataflow_profile.test.name
  location              = azurerm_resource_group.test.location

  properties {
    mode                     = "Enabled"
    request_disk_persistence = "Enabled"
    
    operations {
      operation_type = "Source"
      name          = "temperature-source"
      
      source_settings {
        endpoint_ref = "temperature-endpoint"
        data_sources = ["temperature/*"]
      }
    }
  }

  extended_location {
    name = azurerm_iotoperations_instance.test.extended_location[0].name
    type = azurerm_iotoperations_instance.test.extended_location[0].type
  }
}
`, r.template(data), data.RandomString)
}

func (r IotOperationsDataflowResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_dataflow" "import" {
  name                   = azurerm_iotoperations_dataflow.test.name
  resource_group_name    = azurerm_iotoperations_dataflow.test.resource_group_name
  instance_name          = azurerm_iotoperations_dataflow.test.instance_name
  dataflow_profile_name  = azurerm_iotoperations_dataflow.test.dataflow_profile_name
  location               = azurerm_iotoperations_dataflow.test.location

  properties {
    mode                     = azurerm_iotoperations_dataflow.test.properties[0].mode
    request_disk_persistence = azurerm_iotoperations_dataflow.test.properties[0].request_disk_persistence
    
    operations {
      operation_type = azurerm_iotoperations_dataflow.test.properties[0].operations[0].operation_type
      name          = azurerm_iotoperations_dataflow.test.properties[0].operations[0].name
      
      source_settings {
        endpoint_ref = azurerm_iotoperations_dataflow.test.properties[0].operations[0].source_settings[0].endpoint_ref
        data_sources = azurerm_iotoperations_dataflow.test.properties[0].operations[0].source_settings[0].data_sources
      }
    }
  }

  extended_location {
    name = azurerm_iotoperations_dataflow.test.extended_location[0].name
    type = azurerm_iotoperations_dataflow.test.extended_location[0].type
  }
}
`, r.basic(data))
}

func (r IotOperationsDataflowResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_dataflow" "test" {
  name                   = "acctest-dataflow-%s"
  resource_group_name    = azurerm_resource_group.test.name
  instance_name         = azurerm_iotoperations_instance.test.name
  dataflow_profile_name = azurerm_iotoperations_dataflow_profile.test.name
  location              = azurerm_resource_group.test.location

  properties {
    mode                     = "Enabled"
    request_disk_persistence = "Enabled"
    
    operations {
      operation_type = "Source"
      name          = "temperature-source"
      
      source_settings {
        endpoint_ref   = "temperature-endpoint"
        data_sources   = ["temperature/*", "humidity/*"]
        serialization_format = "Json"
      }
    }

    operations {
      operation_type = "Destination" 
      name          = "adx-destination"
      
      destination_settings {
        endpoint_ref = "adx-endpoint"
        data_destination = "telemetry-table"
      }
    }
  }

  extended_location {
    name = azurerm_iotoperations_instance.test.extended_location[0].name
    type = azurerm_iotoperations_instance.test.extended_location[0].type
  }

  tags = {
    environment = "testing"
    purpose     = "dataflow-acceptance-test"
  }
}
`, r.template(data), data.RandomString)
}
