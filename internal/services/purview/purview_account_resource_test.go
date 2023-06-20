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

type PurviewAccountResource struct{}

func TestAccPurviewAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_purview_account", "test")
	r := PurviewAccountResource{}

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
	r := PurviewAccountResource{}

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

func TestAccPurviewAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_purview_account", "test")
	r := PurviewAccountResource{}

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
	r := PurviewAccountResource{}
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

func (r PurviewAccountResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := account.ParseAccountID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Purview.AccountsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (r PurviewAccountResource) basic(data acceptance.TestData) string {
	template := r.template(data)
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
`, template, data.RandomInteger)
}

func (r PurviewAccountResource) complete(data acceptance.TestData) string {
	template := r.template(data)
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
`, template, data.RandomInteger)
}

func (r PurviewAccountResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
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
`, template)
}

func (r PurviewAccountResource) withManagedResourceGroupName(data acceptance.TestData, managedResourceGroupName string) string {
	template := r.template(data)
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
`, template, data.RandomInteger, managedResourceGroupName)
}

func (r PurviewAccountResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-purview-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
