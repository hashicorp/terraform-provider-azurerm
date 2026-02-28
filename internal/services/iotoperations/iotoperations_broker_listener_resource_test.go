// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iotoperations_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/brokerlistener"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// IotOperationsBrokerListenerResource is a test harness for azurerm_iotoperations_broker_listener acceptance tests.
type IotOperationsBrokerListenerResource struct{}

func TestAccIotOperationsBrokerListener_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_broker_listener", "test")
	r := IotOperationsBrokerListenerResource{}

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

func TestAccIotOperationsBrokerListener_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_broker_listener", "test")
	r := IotOperationsBrokerListenerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_iotoperations_broker_listener"),
		},
	})
}

func TestAccIotOperationsBrokerListener_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_broker_listener", "test")
	r := IotOperationsBrokerListenerResource{}

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

func TestAccIotOperationsBrokerListener_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_broker_listener", "test")
	r := IotOperationsBrokerListenerResource{}

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

func (IotOperationsBrokerListenerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	// Parse the ID to get the structured ID object, not individual strings
	id, err := brokerlistener.ParseListenerID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %+v", state.ID, err)
	}

	// Use the parsed ID object in the Get call
	resp, err := clients.IoTOperations.BrokerListenerClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", state.ID, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

// template builds the minimal provider + resource_group; NOTE: you must create an IoT Operations instance and broker
func (IotOperationsBrokerListenerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iotops-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r IotOperationsBrokerListenerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# NOTE: These values should be replaced with actual IoT Operations instance and broker names
# that exist in your test environment. You can either:
# 1. Set environment variables: IOT_OPERATIONS_INSTANCE_NAME, IOT_OPERATIONS_BROKER_NAME
# 2. Create the instance and broker resources in this template
# 3. Reference existing resources in your test subscription

resource "azurerm_iotoperations_broker_listener" "test" {
  name                = "acctest-bl-%s"
  resource_group_name = azurerm_resource_group.test.name

  # TODO: Replace these with actual values or environment variables
  instance_name = "test-instance-%d"  # or use: os.Getenv("IOT_OPERATIONS_INSTANCE_NAME")
  broker_name   = "test-broker-%d"    # or use: os.Getenv("IOT_OPERATIONS_BROKER_NAME")

  properties {
    ports {
      port = 1883
    }
  }
}
`, r.template(data), data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r IotOperationsBrokerListenerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_broker_listener" "import" {
  name                = azurerm_iotoperations_broker_listener.test.name
  resource_group_name = azurerm_iotoperations_broker_listener.test.resource_group_name
  instance_name       = azurerm_iotoperations_broker_listener.test.instance_name
  broker_name         = azurerm_iotoperations_broker_listener.test.broker_name

  properties {
    ports {
      port = 1883
    }
  }
}
`, r.basic(data))
}

func (r IotOperationsBrokerListenerResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# NOTE: Same as basic template - replace with actual IoT Operations instance and broker names

resource "azurerm_iotoperations_broker_listener" "test" {
  name                = "acctest-bl-%s"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = "test-instance-%s"  # TODO: Replace with actual value or env var
  broker_name         = "test-broker-%s"    # TODO: Replace with actual value or env var
  location            = azurerm_resource_group.test.location

  tags = {
    ENV = "Test"
  }

  properties {
    service_type = "LoadBalancer"
    ports {
      port = 8080
      protocol = "WebSockets"
      authentication_ref = "example-auth"
    }
    ports {
      port = 8443
      protocol = "WebSockets"
      authentication_ref = "example-auth"
      tls {
        mode = "Automatic"
        cert_manager_certificate_spec {
          issuer_ref {
            group = "example-group"
            name  = "example-issuer"
            kind  = "Issuer"
          }
        }
      }
    }
    ports {
      port = 1883
      authentication_ref = "example-auth"
    }
    ports {
      port = 8883
      authentication_ref = "example-auth"
      tls {
        mode = "Manual"
        manual {
          secret_ref = "example-secret"
        }
      }
    }
  }
}
`, r.template(data), data.RandomString, data.RandomString, data.RandomString)
}


