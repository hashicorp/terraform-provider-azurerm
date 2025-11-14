package iotoperations_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/brokerauthentication"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type BrokerAuthenticationResource struct{}

func TestAccBrokerAuthentication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_broker_authentication", "test")
	r := BrokerAuthenticationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("test-broker-auth"),
				check.That(data.ResourceName).Key("authentication_methods.0.method").HasValue("ServiceAccountToken"),
			),
		},
		data.ImportStep(),
	})
}

func (r BrokerAuthenticationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := brokerauthentication.ParseAuthenticationID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.IoTOperations.BrokerAuthenticationClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("IoT Operations Broker Authentication %q was not found", id.AuthenticationName)
		}
		return nil, fmt.Errorf("retrieving IoT Operations Broker Authentication %q: %+v", id.AuthenticationName, err)
	}

	exists := resp.Model != nil
	return &exists, nil
}

func (r BrokerAuthenticationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iot-%d"
  location = "%s"
}

resource "azurerm_iotoperations_broker_authentication" "test" {
  name                = "test-broker-auth"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = "test-instance"
  broker_name         = "test-broker"
  location            = azurerm_resource_group.test.location

  authentication_methods {
    method = "ServiceAccountToken"
    service_account_token_settings {
      audiences = ["test-audience"]
    }
  }

  extended_location {
    name = "test-custom-location"
    type = "CustomLocation"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
