package purview_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/purview/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type PurviewAccountResource struct{}

func TestAccPurviewAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_purview_account", "test")
	r := PurviewAccountResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPurviewAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_purview_account", "test")
	r := PurviewAccountResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r PurviewAccountResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.AccountID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Purview.AccountsClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Purview Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(true), nil
}

func (r PurviewAccountResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_purview_account" "test" {
  name                = "acctestsw%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Standard_4"
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
  sku_name            = azurerm_purview_account.test.sku_name
}
`, template)
}

func (r PurviewAccountResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-purview-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
