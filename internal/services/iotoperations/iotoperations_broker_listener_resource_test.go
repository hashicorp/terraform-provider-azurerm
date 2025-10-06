package iotoperations_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

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
    rg, instance, broker, listener, err := parseBrokerListenerID(state.ID)
    if err != nil {
        return nil, err
    }

    resp, err := clients.IotOperations.BrokerListenerClient.Get(ctx, rg, instance, broker, listener, nil)
    if err != nil {
        return nil, fmt.Errorf("retrieving %s: %+v", state.ID, err)
    }

    return utils.Bool(resp.BrokerListenerResource != nil), nil
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
`, data.RandomInteger, data.Location())
}

func (r IotOperationsBrokerListenerResource) basic(data acceptance.TestData) string {
    return fmt.Sprintf(`
%s

# TODO: create or reference an IoT Operations instance and broker here.

resource "azurerm_iotoperations_broker_listener" "test" {
  name                = "acctest-bl-%s"
  resource_group_name = azurerm_resource_group.test.name

  # Replace with the actual instance and broker names created in this template:
  instance_name = "REPLACE_WITH_INSTANCE_NAME"
  broker_name   = "REPLACE_WITH_BROKER_NAME"

  properties {
    ports {
      port = 1883
    }
  }
}
`, r.template(data), data.RandomString)
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

# See TODO in basic template for instance and broker creation.

resource "azurerm_iotoperations_broker_listener" "test" {
  name                = "acctest-bl-%s"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = "REPLACE_WITH_INSTANCE_NAME"
  broker_name         = "REPLACE_WITH_BROKER_NAME"
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
`, r.template(data), data.RandomString)
}

// parseBrokerListenerID extracts resource group, instance name, broker name, and listener name from a full ARM resource id.
// Expected pattern:
// /subscriptions/.../resourceGroups/{rg}/providers/Microsoft.IoTOperations/instances/{instance}/brokers/{broker}/listeners/{listener}
func parseBrokerListenerID(id string) (rg, instance, broker, listener string, err error) {
    if id == "" {
        return "", "", "", "", fmt.Errorf("empty id")
    }
    parts := strings.Split(id, "/")
    start := 0
    if parts[0] == "" {
        start = 1
    }
    for i := start; i < len(parts)-1; i++ {
        switch strings.ToLower(parts[i]) {
        case "resourcegroups":
            rg = parts[i+1]
        case "instances":
            instance = parts[i+1]
        case "brokers":
            broker = parts[i+1]
        case "listeners":
            listener = parts[i+1]
        }
    }
    if rg == "" || instance == "" || broker == "" || listener == "" {
        return "", "", "", "", fmt.Errorf("failed to parse id: %s", id)
    }
    return rg, instance, broker, listener, nil
}