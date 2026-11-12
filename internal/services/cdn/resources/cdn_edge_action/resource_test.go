// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_edge_action_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CdnEdgeActionResource struct{}

func TestAccCdnEdgeAction_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_edge_action", "test")
	r := CdnEdgeActionResource{}

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

func TestAccCdnEdgeAction_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_edge_action", "test")
	r := CdnEdgeActionResource{}

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

func TestAccCdnEdgeAction_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_edge_action", "test")
	r := CdnEdgeActionResource{}

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

func TestAccCdnEdgeAction_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_edge_action", "test")
	r := CdnEdgeActionResource{}

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

func (CdnEdgeActionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := edgeactions.ParseEdgeActionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Cdn.EdgeActionsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (CdnEdgeActionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-ea-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CdnEdgeActionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_edge_action" "test" {
  name                = "acctestea%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard"
}
`, r.template(data), data.RandomInteger)
}

func (r CdnEdgeActionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_edge_action" "import" {
  name                = azurerm_cdn_edge_action.test.name
  resource_group_name = azurerm_cdn_edge_action.test.resource_group_name
  sku_name            = azurerm_cdn_edge_action.test.sku_name
}
`, r.basic(data))
}

func (r CdnEdgeActionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_edge_action" "test" {
  name                = "acctestea%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard"

  tags = {
    Environment = "Test"
    Label       = "EdgeAction"
  }
}
`, r.template(data), data.RandomInteger)
}
