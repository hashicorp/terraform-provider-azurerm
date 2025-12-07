// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nginx_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxapikey"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/nginx"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type APIKeyResource struct{}

func (a APIKeyResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := nginxapikey.ParseApiKeyID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Nginx.NginxApiKey.ApiKeysGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestAccAPIKeyResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, nginx.APIKeyResource{}.ResourceType(), "test")
	secretText := uuid.NewString()
	endDateTime := getEndDateTime(3)
	r := APIKeyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, secretText, endDateTime),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("end_date_time").HasValue(endDateTime),
				check.That(data.ResourceName).Key("hint").HasValue(secretText[:3]),
			),
		},
		// The secret_text is provided by the user and never returned by the API.
		data.ImportStep("secret_text"),
	})
}

func TestAccAPIKeyResource_update(t *testing.T) {
	data := acceptance.BuildTestData(t, nginx.APIKeyResource{}.ResourceType(), "test")
	secretTextCreate := uuid.NewString()
	secretTextUpdate := uuid.NewString()
	endDateTimeCreate := getEndDateTime(3)
	endDateTimeUpdate := getEndDateTime(6)
	r := APIKeyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, secretTextCreate, endDateTimeCreate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("end_date_time").HasValue(endDateTimeCreate),
				check.That(data.ResourceName).Key("hint").HasValue(secretTextCreate[:3]),
			),
		},
		// The secret_text is provided by the user and never returned by the API.
		data.ImportStep("secret_text"),
		{
			Config: r.complete(data, secretTextUpdate, endDateTimeUpdate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("end_date_time").HasValue(endDateTimeUpdate),
				check.That(data.ResourceName).Key("hint").HasValue(secretTextUpdate[:3]),
			),
		},
		// The secret_text is provided by the user and never returned by the API.
		data.ImportStep("secret_text"),
	})
}

func TestAccAPIKeyResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, nginx.APIKeyResource{}.ResourceType(), "test")
	secretText := uuid.NewString()
	endDateTime := getEndDateTime(3)
	r := APIKeyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, secretText, endDateTime),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func getEndDateTime(monthsFromNow int) string {
	return time.Now().AddDate(0, monthsFromNow, 0).UTC().Format("2006-01-02T15:04:05+00:00")
}

func (a APIKeyResource) complete(data acceptance.TestData, secretText, endDateTime string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_nginx_api_key" "test" {
  name                = "acctest-%d"
  nginx_deployment_id = azurerm_nginx_deployment.test.id
  end_date_time       = "%s"
  secret_text         = "%s"
}
`, a.template(data), data.RandomInteger, endDateTime, secretText)
}

func (a APIKeyResource) requiresImport(data acceptance.TestData) string {
	secretText := uuid.NewString()
	endDateTime := getEndDateTime(3)
	return fmt.Sprintf(`
%s

resource "azurerm_nginx_api_key" "import" {
  name                = azurerm_nginx_api_key.test.name
  nginx_deployment_id = azurerm_nginx_api_key.test.nginx_deployment_id
  end_date_time       = azurerm_nginx_api_key.test.end_date_time
  secret_text         = azurerm_nginx_api_key.test.secret_text
}
`, a.complete(data, secretText, endDateTime))
}

func (a APIKeyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctest%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  delegation {
    name = "delegation"

    service_delegation {
      name = "NGINX.NGINXPLUS/nginxDeployments"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_nginx_deployment" "test" {
  name                     = "acctest-%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  sku                      = "standardv2_Monthly"
  capacity                 = 10
  location                 = azurerm_resource_group.test.location
  diagnose_support_enabled = false

  frontend_public {
    ip_address = [azurerm_public_ip.test.id]
  }

  network_interface {
    subnet_id = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
