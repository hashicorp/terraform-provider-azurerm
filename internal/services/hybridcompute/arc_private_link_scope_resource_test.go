// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hybridcompute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/privatelinkscopes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ArcPrivateLinkScopeResource struct{}

func TestAccArcPrivateLinkScope_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_private_link_scope", "test")
	r := ArcPrivateLinkScopeResource{}
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

func TestAccArcPrivateLinkScope_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_private_link_scope", "test")
	r := ArcPrivateLinkScopeResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("tags.Environment").HasValue("Production"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccArcPrivateLinkScope_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_private_link_scope", "test")
	r := ArcPrivateLinkScopeResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("tags.Environment").HasValue("Production"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("tags.Environment").HasValue("Staging"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccArcPrivateLinkScope_requireImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_private_link_scope", "test")
	r := ArcPrivateLinkScopeResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requireImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_arc_private_link_scope"),
		},
	})
}

func (r ArcPrivateLinkScopeResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := privatelinkscopes.ParseProviderPrivateLinkScopeID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.HybridCompute.PrivateLinkScopesClient

	resp, err := client.Get(ctx, *id)
	exists := false
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &exists, nil
		}
		return nil, fmt.Errorf("retrieving: %+v", err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r ArcPrivateLinkScopeResource) template(data acceptance.TestData) interface{} {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ArcPrivateLinkScopeResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_arc_private_link_scope" "test" {
  name                = "acctestPLS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, r.template(data), data.RandomInteger)
}

func (r ArcPrivateLinkScopeResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_arc_private_link_scope" "test" {
  name                = "acctestPLS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  public_network_access_enabled = true

  tags = {
    "Environment" = "Production"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ArcPrivateLinkScopeResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_arc_private_link_scope" "test" {
  name                = "acctestPLS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  public_network_access_enabled = false

  tags = {
    "Environment" = "Staging"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ArcPrivateLinkScopeResource) requireImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_arc_private_link_scope" "import" {
  name                = azurerm_arc_private_link_scope.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, r.basic(data))
}
