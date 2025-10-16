// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/resourceanchors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ResourceAnchorResource struct{}

func (a ResourceAnchorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := resourceanchors.ParseResourceAnchorID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Oracle.OracleClient09.ResourceAnchors.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestResourceAnchorResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ResourceAnchorResource{}.ResourceType(), "test")
	r := ResourceAnchorResource{}
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

func TestResourceAnchorResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ResourceAnchorResource{}.ResourceType(), "test")
	r := ResourceAnchorResource{}
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

func TestResourceAnchorResource_update(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ResourceAnchorResource{}.ResourceType(), "test")
	r := ResourceAnchorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{

		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),

		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestResourceAnchorResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ResourceAnchorResource{}.ResourceType(), "test")
	r := ResourceAnchorResource{}
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

func (a ResourceAnchorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_oracle_database_resource_anchor" "test" {
  name                = "ra1-%[2]d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name

}`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a ResourceAnchorResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_oracle_database_resource_anchor" "test" {
  name                = "ra1-%[2]d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  tags = {
    key236 = "wbucrnidikivbujndfk"
  }
}`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a ResourceAnchorResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_oracle_database_resource_anchor" "test" {
  location            = "global"
  name                = "ra1-%[2]d"
  resource_group_name = azurerm_resource_group.test.name

  # Updated Tags
  tags = {
    newtag = "newvalue"
  }
}`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a ResourceAnchorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_database_resource_anchor" "import" {
  location              = azurerm_oracle_database_resource_anchor.test.location
  name                  = azurerm_oracle_database_resource_anchor.test.name
  resource_group_name   = azurerm_oracle_database_resource_anchor.test.resource_group_name
  linked_compartment_id = azurerm_oracle_database_resource_anchor.test.linked_compartment_id
}
`, a.basic(data))
}

func (a ResourceAnchorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-anchor-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary)
}
