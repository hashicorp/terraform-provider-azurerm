// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package purview_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/purview/2021-07-01/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PurviewAccountResourceTest struct{}

func TestAccPurviewAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_purview_account", "test")
	r := PurviewAccountResourceTest{}

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

func TestAccPurviewAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_purview_account", "test")
	r := PurviewAccountResourceTest{}

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

func TestAccPurviewAccount_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_purview_account", "test")
	r := PurviewAccountResourceTest{}

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
	})
}

func TestAccPurviewAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_purview_account", "test")
	r := PurviewAccountResourceTest{}

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

func TestAccPurviewAccount_withManagedResourceGroupName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_purview_account", "test")
	r := PurviewAccountResourceTest{}
	managedResourceGroupName := fmt.Sprintf("acctestRG-purview-managed-%d", data.RandomInteger)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withManagedResourceGroupName(data, managedResourceGroupName),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("managed_resource_group_name").HasValue(managedResourceGroupName),
			),
		},
		data.ImportStep(),
	})
}

func (r PurviewAccountResourceTest) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := account.ParseAccountID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Purview.AccountsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(true), nil
}

func (r PurviewAccountResourceTest) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_purview_account" "test" {
  name                = "acctestsw%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PurviewAccountResourceTest) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_purview_account" "test" {
  name                   = "acctestsw%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  public_network_enabled = false

  identity {
    type = "SystemAssigned"
  }

  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PurviewAccountResourceTest) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_purview_account" "import" {
  name                = azurerm_purview_account.test.name
  resource_group_name = azurerm_purview_account.test.resource_group_name
  location            = azurerm_purview_account.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, r.basic(data))
}

func (r PurviewAccountResourceTest) withManagedResourceGroupName(data acceptance.TestData, managedResourceGroupName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_purview_account" "test" {
  name                        = "acctestsw%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  managed_resource_group_name = %q

  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomInteger, managedResourceGroupName)
}

func (r PurviewAccountResourceTest) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-purview-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Ternary)
}
