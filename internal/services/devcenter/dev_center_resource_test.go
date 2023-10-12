package devcenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/devcenters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DevCenterResource struct{}

func TestAccDevCenter_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center", "test")
	r := DevCenterResource{}

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

func TestAccDevCenter_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center", "test")
	r := DevCenterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r DevCenterResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := devcenters.ParseDevCenterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.DevCenter.DevCenterClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r DevCenterResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_dev_center" "test" {
  name                = "testdevcenter-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomIntOfLength(8))
}

func (r DevCenterResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "testuserassignedidentity-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_dev_center" "test" {
  name                = "testdevcenter-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, template, data.RandomIntOfLength(8))
}

func (r DevCenterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRg-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
