// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gatewayapi"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementGatewayAPIResource struct{}

func TestAccApiManagementGatewayApi_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_gateway_api", "test")
	r := ApiManagementGatewayAPIResource{}

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

func TestAccApiManagementGatewayApi_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_gateway_api", "test")
	r := ApiManagementGatewayAPIResource{}

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

func (ApiManagementGatewayAPIResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := gatewayapi.ParseGatewayApiID(state.ID)
	if err != nil {
		return nil, err
	}
	if resp, err := clients.ApiManagement.GatewayApisClient.GetEntityTag(ctx, *id); err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if !response.WasStatusCode(resp.HttpResponse, http.StatusNoContent) {
			return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
		}
	}

	return pointer.To(true), nil
}

func (ApiManagementGatewayAPIResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_gateway" "test" {
  name              = "acctestAMGateway-%d"
  api_management_id = azurerm_api_management.test.id
  description       = "this is a test gateway"

  location_data {
    name     = "old world"
    city     = "test city"
    district = "test district"
    region   = "test region"
  }
}

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"
}

resource "azurerm_api_management_gateway_api" "test" {
  gateway_id = azurerm_api_management_gateway.test.id
  api_id     = azurerm_api_management_api.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementGatewayAPIResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_gateway_api" "import" {
  gateway_id = azurerm_api_management_gateway_api.test.gateway_id
  api_id     = azurerm_api_management_gateway_api.test.api_id
}
`, r.basic(data))
}
