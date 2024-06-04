// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressrouteportauthorizations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ExpressRoutePortAuthorizationResource struct{}

func TestAccExpressRoutePortAuthorization_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_port_authorization", "test")
	r := ExpressRoutePortAuthorizationResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("authorization_key").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccExpressRoutePortAuthorization_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_port_authorization", "test")
	r := ExpressRoutePortAuthorizationResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("authorization_key").Exists(),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_express_route_port_authorization"),
		},
	})
}

func TestAccExpressRoutePortAuthorization_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_port_authorization", "test1")
	r := ExpressRoutePortAuthorizationResource{}
	secondResourceName := "azurerm_express_route_port_authorization.test2"

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("authorization_key").Exists(),
				acceptance.TestCheckResourceAttrSet(secondResourceName, "authorization_key"),
			),
		},
	})
}

func (r ExpressRoutePortAuthorizationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := expressrouteportauthorizations.ParseExpressRoutePortAuthorizationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.ExpressRoutePortAuthorizations.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ExpressRoutePortAuthorizationResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_port_authorization" "test" {
  name                    = "acctestauth%[2]d"
  express_route_port_name = azurerm_express_route_port.test.name
  resource_group_name     = azurerm_resource_group.test.name
}
`, template, data.RandomInteger)
}

func (r ExpressRoutePortAuthorizationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_port_authorization" "import" {
  name                    = azurerm_express_route_port_authorization.test.name
  express_route_port_name = azurerm_express_route_port_authorization.test.express_route_port_name
  resource_group_name     = azurerm_express_route_port_authorization.test.resource_group_name
}
`, r.basic(data))
}

func (r ExpressRoutePortAuthorizationResource) multiple(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_port_authorization" "test1" {
  name                    = "acctestauth1%[2]d"
  express_route_port_name = azurerm_express_route_port.test.name
  resource_group_name     = azurerm_resource_group.test.name
}

resource "azurerm_express_route_port_authorization" "test2" {
  name                    = "acctestauth2%[2]d"
  express_route_port_name = azurerm_express_route_port.test.name
  resource_group_name     = azurerm_resource_group.test.name
}
`, template, data.RandomInteger)
}

func (ExpressRoutePortAuthorizationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_express_route_port" "test" {
  name                = "acctestERP-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  peering_location    = "Airtel-Chennai2-CLS"
  bandwidth_in_gbps   = 10
  encapsulation       = "Dot1Q"
  billing_type        = "MeteredData"
  tags = {
    ENV = "Test"
  }
}`, data.RandomInteger, data.Locations.Primary)
}
