// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apigatewayconfigconnection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementGatewayConfigConnectionResource struct{}

func TestAccApiManagementGatewayConfigConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_gateway_config_connection", "test")
	r := ApiManagementGatewayConfigConnectionResource{}

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

func TestAccApiManagementGatewayConfigConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_gateway_config_connection", "test")
	r := ApiManagementGatewayConfigConnectionResource{}

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

func TestAccApiManagementGatewayConfigConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_gateway_config_connection", "test")
	r := ApiManagementGatewayConfigConnectionResource{}

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

func (r ApiManagementGatewayConfigConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := apigatewayconfigconnection.ParseConfigConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.ApiGatewayConfigConnectionClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ApiManagementGatewayConfigConnectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_gateway_config_connection" "test" {
  name                      = "acctest-gc-%d"
  api_management_gateway_id = azurerm_api_management_standalone_gateway.test.id
  source_id                 = azurerm_api_management_workspace.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementGatewayConfigConnectionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_gateway_config_connection" "import" {
  name                      = azurerm_api_management_gateway_config_connection.test.name
  api_management_gateway_id = azurerm_api_management_gateway_config_connection.test.api_management_gateway_id
  source_id                 = azurerm_api_management_gateway_config_connection.test.source_id
}
`, r.basic(data))
}

func (r ApiManagementGatewayConfigConnectionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_gateway_config_connection" "test" {
  name                      = "acctest-gc-%d"
  api_management_gateway_id = azurerm_api_management_standalone_gateway.test.id
  source_id                 = azurerm_api_management_workspace.test.id

  hostnames = [
    "example.com",
    "api.example.com",
    "gateway.example.com"
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementGatewayConfigConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestapim-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Premium_1"
}

resource "azurerm_api_management_workspace" "test" {
  name              = "acctest-workspace-%[1]d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "my workspace"
}

resource "azurerm_api_management_standalone_gateway" "test" {
  name                = "acctest-gateway-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "WorkspaceGatewayPremium"
    capacity = 1
  }

  depends_on = [azurerm_api_management_workspace.test]
}
`, data.RandomInteger, data.Locations.Primary)
}
