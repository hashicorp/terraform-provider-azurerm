// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iotoperations_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/broker"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func TestAccIotOperationsBroker_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_broker", "test")
	r := BrokerResource{}

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

func TestAccIotOperationsBroker_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_broker", "test")
	r := BrokerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_iotoperations_broker"),
		},
	})
}

func TestAccIotOperationsBroker_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_broker", "test")
	r := BrokerResource{}

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

func TestAccIotOperationsBroker_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_broker", "test")
	r := BrokerResource{}

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

type BrokerResource struct{}

func (BrokerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := broker.ParseBrokerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.IoTOperations.BrokerClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

// Template function to create common resources
func (BrokerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iot-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r BrokerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_broker" "test" {
  name                = "acctest-br-%s"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = "test-instance-%d"

  properties {
    memory_profile = "Tiny"
  }

  extended_location {
    name = "/subscriptions/%s/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.ExtendedLocation/customLocations/testlocation-%s"
    type = "CustomLocation"
  }
}
`, r.template(data), data.RandomString, data.RandomInteger, data.Client().SubscriptionID, data.RandomString)
}

func (r BrokerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_broker" "import" {
  name                = azurerm_iotoperations_broker.test.name
  resource_group_name = azurerm_iotoperations_broker.test.resource_group_name
  instance_name       = azurerm_iotoperations_broker.test.instance_name

  properties {
    memory_profile = "Tiny"
  }

  extended_location {
    name = azurerm_iotoperations_broker.test.extended_location[0].name
    type = azurerm_iotoperations_broker.test.extended_location[0].type
  }
}
`, r.basic(data))
}

func (r BrokerResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_broker" "test" {
  name                = "acctest-br-%s"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = "test-instance-%d"
  location            = azurerm_resource_group.test.location

  tags = {
    ENV = "Test"
  }

  properties {
    memory_profile = "Large"
    
    cardinality {
      backend_chain {
        partitions       = 2
        redundancy_factor = 2
        workers          = 2
      }
      frontend {
        replicas = 2
        workers  = 2
      }
    }
  }

  extended_location {
    name = "/subscriptions/%s/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.ExtendedLocation/customLocations/testlocation-%s"
    type = "CustomLocation"
  }
}
`, r.template(data), data.RandomString, data.RandomInteger, data.Client().SubscriptionID, data.RandomString)
}
