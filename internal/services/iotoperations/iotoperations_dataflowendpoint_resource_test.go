// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iotoperations_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/dataflowendpoint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// IotOperationsDataflowEndpointResource is a test harness for azurerm_iotoperations_dataflow_endpoint acceptance tests.
type IotOperationsDataflowEndpointResource struct{}

func TestAccIotOperationsDataflowEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_dataflow_endpoint", "test")
	r := IotOperationsDataflowEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctest-de-%s", data.RandomString)),
				check.That(data.ResourceName).Key("properties.0.endpoint_type").HasValue("DataExplorer"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotOperationsDataflowEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_dataflow_endpoint", "test")
	r := IotOperationsDataflowEndpointResource{}

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

func TestAccIotOperationsDataflowEndpoint_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_dataflow_endpoint", "test")
	r := IotOperationsDataflowEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctest-de-%s", data.RandomString)),
				check.That(data.ResourceName).Key("properties.0.endpoint_type").HasValue("DataLakeStorage"),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotOperationsDataflowEndpoint_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_dataflow_endpoint", "test")
	r := IotOperationsDataflowEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("properties.0.endpoint_type").HasValue("DataExplorer"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("properties.0.endpoint_type").HasValue("DataLakeStorage"),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func (r IotOperationsDataflowEndpointResource) Exists(ctx context.Context, c *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dataflowendpoint.ParseDataflowEndpointID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing ID %q: %w", state.ID, err)
	}

	resp, err := c.IoTOperations.DataflowEndpointClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id.ID(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r IotOperationsDataflowEndpointResource) template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (r IotOperationsDataflowEndpointResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_dataflow_endpoint" "test" {
  name                = "acctest-de-%s"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = azurerm_iotoperations_instance.test.name
  location            = azurerm_resource_group.test.location

  properties {
    endpoint_type = "DataExplorer"
    data_explorer_settings {
      authentication {
        method = "SystemAssignedManagedIdentity"
        system_assigned_managed_identity_settings {
          audience = "https://help.kusto.windows.net"
        }
      }
      database = "testdb-%s"
      host     = "testcluster-%s.region.kusto.windows.net"
      batching {
        latency_seconds = 5
        max_messages    = 100
      }
    }
  }

  extended_location {
    name = azurerm_iotoperations_instance.test.extended_location[0].name
    type = azurerm_iotoperations_instance.test.extended_location[0].type
  }
}
`, r.template(data), data.RandomString, data.RandomString, data.RandomString)
}

func (r IotOperationsDataflowEndpointResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_dataflow_endpoint" "import" {
  name                = azurerm_iotoperations_dataflow_endpoint.test.name
  resource_group_name = azurerm_iotoperations_dataflow_endpoint.test.resource_group_name
  instance_name       = azurerm_iotoperations_dataflow_endpoint.test.instance_name
  location            = azurerm_iotoperations_dataflow_endpoint.test.location

  properties {
    endpoint_type = azurerm_iotoperations_dataflow_endpoint.test.properties[0].endpoint_type
    data_explorer_settings {
      authentication {
        method = azurerm_iotoperations_dataflow_endpoint.test.properties[0].data_explorer_settings[0].authentication[0].method
        system_assigned_managed_identity_settings {
          audience = azurerm_iotoperations_dataflow_endpoint.test.properties[0].data_explorer_settings[0].authentication[0].system_assigned_managed_identity_settings[0].audience
        }
      }
      database = azurerm_iotoperations_dataflow_endpoint.test.properties[0].data_explorer_settings[0].database
      host     = azurerm_iotoperations_dataflow_endpoint.test.properties[0].data_explorer_settings[0].host
      batching {
        latency_seconds = azurerm_iotoperations_dataflow_endpoint.test.properties[0].data_explorer_settings[0].batching[0].latency_seconds
        max_messages    = azurerm_iotoperations_dataflow_endpoint.test.properties[0].data_explorer_settings[0].batching[0].max_messages
      }
    }
  }

  extended_location {
    name = azurerm_iotoperations_dataflow_endpoint.test.extended_location[0].name
    type = azurerm_iotoperations_dataflow_endpoint.test.extended_location[0].type
  }
}
`, r.basic(data))
}

func (r IotOperationsDataflowEndpointResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_dataflow_endpoint" "test" {
  name                = "acctest-de-%s"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = azurerm_iotoperations_instance.test.name
  location            = azurerm_resource_group.test.location

  properties {
    endpoint_type = "DataLakeStorage"
    data_lake_storage_settings {
      authentication {
        method = "SystemAssignedManagedIdentity"
        system_assigned_managed_identity_settings {
          audience = "https://storage.azure.com/"
        }
      }
      host = "testaccount-%s.dfs.core.windows.net"
      batching {
        latency_seconds = 10
        max_messages    = 1000
      }
    }
  }

  extended_location {
    name = azurerm_iotoperations_instance.test.extended_location[0].name
    type = azurerm_iotoperations_instance.test.extended_location[0].type
  }

  tags = {
    environment = "testing"
    purpose     = "dataflow-endpoint-acceptance-test"
  }
}
`, r.template(data), data.RandomString, data.RandomString)
}
