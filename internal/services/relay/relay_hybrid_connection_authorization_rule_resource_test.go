// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package relay_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type RelayHybridConnectionAuthorizationRuleResource struct{}

func TestAccRelayHybridConnectionAuthorizationRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_relay_hybrid_connection_authorization_rule", "test")
	r := RelayHybridConnectionAuthorizationRuleResource{}

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

func TestAccRelayHybridConnectionAuthorizationRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_relay_hybrid_connection_authorization_rule", "test")
	r := RelayHybridConnectionAuthorizationRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_relay_hybrid_connection_authorization_rule"),
		},
	})
}

func (t RelayHybridConnectionAuthorizationRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := hybridconnections.ParseHybridConnectionAuthorizationRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Relay.HybridConnectionsClient.GetAuthorizationRule(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (RelayHybridConnectionAuthorizationRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_relay_namespace" "test" {
  name                = "acctestrn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "Standard"
}

resource "azurerm_relay_hybrid_connection" "test" {
  name                 = "acctestrnhc-%d"
  resource_group_name  = azurerm_resource_group.test.name
  relay_namespace_name = azurerm_relay_namespace.test.name
}

resource "azurerm_relay_hybrid_connection_authorization_rule" "test" {
  name                   = "acctestrnak-%d"
  namespace_name         = azurerm_relay_namespace.test.name
  hybrid_connection_name = azurerm_relay_hybrid_connection.test.name
  resource_group_name    = azurerm_resource_group.test.name

  listen = true
  send   = true
  manage = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r RelayHybridConnectionAuthorizationRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_relay_hybrid_connection_authorization_rule" "import" {
  name                   = azurerm_relay_hybrid_connection_authorization_rule.test.name
  namespace_name         = azurerm_relay_hybrid_connection_authorization_rule.test.namespace_name
  hybrid_connection_name = azurerm_relay_hybrid_connection_authorization_rule.test.hybrid_connection_name
  resource_group_name    = azurerm_relay_hybrid_connection_authorization_rule.test.resource_group_name

  listen = azurerm_relay_hybrid_connection_authorization_rule.test.listen
  send   = azurerm_relay_hybrid_connection_authorization_rule.test.send
  manage = azurerm_relay_hybrid_connection_authorization_rule.test.manage
}
`, r.basic(data))
}
