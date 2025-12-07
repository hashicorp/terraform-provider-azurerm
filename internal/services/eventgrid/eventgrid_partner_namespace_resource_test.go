// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/partnernamespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type EventGridPartnerNamespaceTestResource struct{}

func TestAccEventGridPartnerNamespace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_partner_namespace", "test")

	r := EventGridPartnerNamespaceTestResource{}
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

func TestAccEventGridPartnerNamespace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_partner_namespace", "test")
	r := EventGridPartnerNamespaceTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_eventgrid_partner_namespace"),
		},
	})
}

func TestAccEventGridPartnerNamespace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_partner_namespace", "test")
	r := EventGridPartnerNamespaceTestResource{}
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

func TestAccEventGridPartnerNamespace_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_partner_namespace", "test")
	r := EventGridPartnerNamespaceTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (EventGridPartnerNamespaceTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := partnernamespaces.ParsePartnerNamespaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.EventGrid.PartnerNamespaces.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf(("retrieving %s: %+v"), *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r EventGridPartnerNamespaceTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventgrid_partner_registration" "test" {
  name                = "acctest-egpr-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_eventgrid_partner_namespace" "test" {
  name                    = "acctest-egpn-%[1]d"
  location                = "%[2]s"
  resource_group_name     = azurerm_resource_group.test.name
  partner_registration_id = azurerm_eventgrid_partner_registration.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (EventGridPartnerNamespaceTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_eventgrid_partner_namespace" "import" {
  name                    = azurerm_eventgrid_partner_namespace.test.name
  location                = azurerm_eventgrid_partner_namespace.test.location
  resource_group_name     = azurerm_eventgrid_partner_namespace.test.resource_group_name
  partner_registration_id = azurerm_eventgrid_partner_namespace.test.partner_registration_id
}
`, EventGridPartnerNamespaceTestResource{}.basic(data))
}

func (EventGridPartnerNamespaceTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventgrid_partner_registration" "test" {
  name                = "acctest-egpr-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_eventgrid_partner_namespace" "test" {
  name                         = "acctest-egpn-%[1]d"
  location                     = "%[2]s"
  resource_group_name          = azurerm_resource_group.test.name
  partner_registration_id      = azurerm_eventgrid_partner_registration.test.id
  local_authentication_enabled = false
  inbound_ip_rule {
    ip_mask = "10.0.0.0/16"
    action  = "Allow"
  }

  inbound_ip_rule {
    ip_mask = "10.1.0.0/16"
    action  = "Allow"
  }

  partner_topic_routing_mode = "ChannelNameHeader"
  public_network_access      = "Enabled"

  tags = {
    "foo" = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (EventGridPartnerNamespaceTestResource) completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventgrid_partner_registration" "test" {
  name                = "acctest-egpr-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_eventgrid_partner_namespace" "test" {
  name                         = "acctest-egpn-%[1]d"
  location                     = "%[2]s"
  resource_group_name          = azurerm_resource_group.test.name
  partner_registration_id      = azurerm_eventgrid_partner_registration.test.id
  local_authentication_enabled = true
  inbound_ip_rule {
    ip_mask = "10.0.0.0/16"
    action  = "Allow"
  }

  inbound_ip_rule {
    ip_mask = "10.10.10.10/16"
    action  = "Allow"
  }

  partner_topic_routing_mode = "ChannelNameHeader"
  public_network_access      = "Disabled"

  tags = {
    "foo"     = "bar"
    "example" = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
