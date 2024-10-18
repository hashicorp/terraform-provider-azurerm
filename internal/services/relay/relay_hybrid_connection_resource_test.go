// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package relay_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type RelayHybridConnectionResource struct{}

func TestAccRelayHybridConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_relay_hybrid_connection", "test")
	r := RelayHybridConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("requires_client_authorization").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRelayHybridConnection_full(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_relay_hybrid_connection", "test")
	r := RelayHybridConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.full(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("requires_client_authorization").HasValue("false"),
				check.That(data.ResourceName).Key("user_metadata").HasValue("metadatatest"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRelayHybridConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_relay_hybrid_connection", "test")
	r := RelayHybridConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("requires_client_authorization").Exists(),
			),
		},
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("user_metadata").HasValue("metadataupdated"),
			),
		},
	})
}

func TestAccRelayHybridConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_relay_hybrid_connection", "test")
	r := RelayHybridConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("requires_client_authorization").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t RelayHybridConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := hybridconnections.ParseHybridConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Relay.HybridConnectionsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (RelayHybridConnectionResource) basic(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (RelayHybridConnectionResource) full(data acceptance.TestData) string {
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
  name                          = "acctestrnhc-%d"
  resource_group_name           = azurerm_resource_group.test.name
  relay_namespace_name          = azurerm_relay_namespace.test.name
  requires_client_authorization = false
  user_metadata                 = "metadatatest"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (RelayHybridConnectionResource) update(data acceptance.TestData) string {
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
  user_metadata        = "metadataupdated"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r RelayHybridConnectionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_relay_hybrid_connection" "import" {
  name                 = azurerm_relay_hybrid_connection.test.name
  resource_group_name  = azurerm_relay_hybrid_connection.test.resource_group_name
  relay_namespace_name = azurerm_relay_hybrid_connection.test.relay_namespace_name
}
`, r.basic(data))
}
