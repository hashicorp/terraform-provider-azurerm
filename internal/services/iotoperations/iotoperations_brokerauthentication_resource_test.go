// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type BrokerAuthenticationResource struct{}

func TestAccIotOperationsBrokerAuthentication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_broker_authentication", "test")
	r := BrokerAuthenticationResource{}

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

func TestAccIotOperationsBrokerAuthentication_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_broker_authentication", "test")
	r := BrokerAuthenticationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_iotoperations_broker_authentication"),
		},
	})
}

func TestAccIotOperationsBrokerAuthentication_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_broker_authentication", "test")
	r := BrokerAuthenticationResource{}

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

func TestAccIotOperationsBrokerAuthentication_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_broker_authentication", "test")
	r := BrokerAuthenticationResource{}

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

func (r BrokerAuthenticationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := brokerauthentication.ParseAuthenticationID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.IoTOperations.BrokerAuthenticationClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving IoT Operations Broker Authentication %q: %+v", id.AuthenticationName, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

// Template function to create common resources
func (BrokerAuthenticationResource) template(data acceptance.TestData) string {
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

func (r BrokerAuthenticationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_broker_authentication" "test" {
  name                = "acctest-ba-%s"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = "test-instance-%d"
  broker_name         = "test-broker-%s"
  location            = azurerm_resource_group.test.location

  authentication_methods {
    method = "ServiceAccountToken"
    service_account_token_settings {
      audiences = ["test-audience-%s"]
    }
  }

  extended_location {
    name = "/subscriptions/%s/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.ExtendedLocation/customLocations/testlocation-%s"
    type = "CustomLocation"
  }
}
`, r.template(data), data.RandomString, data.RandomInteger, data.RandomString, data.RandomString, data.Client().SubscriptionID, data.RandomString)
}

func (r BrokerAuthenticationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_broker_authentication" "import" {
  name                = azurerm_iotoperations_broker_authentication.test.name
  resource_group_name = azurerm_iotoperations_broker_authentication.test.resource_group_name
  instance_name       = azurerm_iotoperations_broker_authentication.test.instance_name
  broker_name         = azurerm_iotoperations_broker_authentication.test.broker_name
  location            = azurerm_iotoperations_broker_authentication.test.location

  authentication_methods {
    method = "ServiceAccountToken"
    service_account_token_settings {
      audiences = ["test-audience-%s"]
    }
  }

  extended_location {
    name = azurerm_iotoperations_broker_authentication.test.extended_location[0].name
    type = azurerm_iotoperations_broker_authentication.test.extended_location[0].type
  }
}
`, r.basic(data), data.RandomString)
}

func (r BrokerAuthenticationResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_broker_authentication" "test" {
  name                = "acctest-ba-%s"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = "test-instance-%d"
  broker_name         = "test-broker-%s"
  location            = azurerm_resource_group.test.location

  tags = {
    ENV = "Test"
  }

  authentication_methods {
    method = "ServiceAccountToken"
    service_account_token_settings {
      audiences = ["test-audience-%s", "additional-audience-%s"]
    }
  }

  authentication_methods {
    method = "X509Certificate"
    x509_settings {
      trusted_client_ca_cert = "example-cert-%s"
      authorization_attributes = {
        "subject" = "CN=example"
        "issuer"  = "CN=ca"
      }
    }
  }

  extended_location {
    name = "/subscriptions/%s/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.ExtendedLocation/customLocations/testlocation-%s"
    type = "CustomLocation"
  }
}
`, r.template(data), data.RandomString, data.RandomInteger, data.RandomString, data.RandomString, data.RandomString, data.RandomString, data.Client().SubscriptionID, data.RandomString)
}
